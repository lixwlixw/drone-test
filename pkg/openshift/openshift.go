package openshift

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"
	"time"

	project "github.com/openshift/origin/pkg/project/api/v1"
	rolebinding "github.com/openshift/origin/pkg/rolebinding/api/v1"
	user "github.com/openshift/origin/pkg/user/api/v1"
	"github.com/zonesan/clog"
	kapi "k8s.io/kubernetes/pkg/api/v1"
)

type OpenshiftClient struct {
	name string

	host string
	//authUrl string
	oapiUrl string
	kapiUrl string

	namespace string
	username  string
	password  string
	//bearerToken string
	bearerToken atomic.Value
}

func Hello() {
	proj := project.Project{}
	usr := user.User{}
	role := rolebinding.RoleBinding{}

	fmt.Printf("%#v\n%#v\n%#v", proj, usr, role)

}

//==============================================================
//
//==============================================================

func CreateOpenshiftClientFromUserToken(host, token string) *OpenshiftClient {
	host = httpsAddr(host)
	oc := &OpenshiftClient{
		host:    host,
		oapiUrl: host + "/oapi/v1",
		kapiUrl: host + "/api/v1",
	}

	oc.setBearerToken(token)

	return oc
}

func NewOpenshiftClient(name, host, username, password string, durPhase time.Duration) *OpenshiftClient {
	host = httpsAddr(host)
	oc := &OpenshiftClient{
		name: name,

		host: host,
		//authUrl: host + "/oauth/authorize?response_type=token&client_id=openshift-challenging-client",
		oapiUrl: host + "/oapi/v1",
		kapiUrl: host + "/api/v1",

		username: username,
		password: password,
	}
	oc.bearerToken.Store("")

	go oc.updateBearerToken(durPhase)

	return oc
}

func NewOpenshiftTokenClient(host, bearerToken string) *OpenshiftClient {
	host = httpsAddr(host)
	oc := &OpenshiftClient{
		host:    host,
		oapiUrl: host + "/oapi/v1",
		kapiUrl: host + "/api/v1",
	}

	oc.setBearerToken(bearerToken)

	return oc
}

// // for general user
// // the token must contains "Bearer "
// func (baseOC *OpenshiftClient) NewOpenshiftClient(token string) *OpenshiftClient {
// 	oc := &OpenshiftClient{
// 		host:    baseOC.host,
// 		oapiUrl: baseOC.oapiUrl,
// 		kapiUrl: baseOC.kapiUrl,
// 	}

// 	oc.setBearerToken(token)

// 	return oc
// }

func (oc *OpenshiftClient) BearerToken() string {
	//return oc.bearerToken
	return oc.bearerToken.Load().(string)
}

func (oc *OpenshiftClient) setBearerToken(token string) {
	oc.bearerToken.Store(token)
}

func (oc *OpenshiftClient) updateBearerToken(durPhase time.Duration) {
	for {

		// clog.Debugf("Request bearer token from: %v(%v) ", oc.name, oc.host)

		token, err := RequestToken(oc.host, oc.username, oc.password)
		if err != nil {
			clog.Error("RequestToken error, try in 15 seconds. error detail: ", err)

			time.Sleep(15 * time.Second)
		} else {

			oc.setBearerToken("Bearer " + token)

			clog.Infof("[%v] [%v] [%v]", oc.name, oc.host, token)

			// durPhase is to avoid mulitple OCs updating tokens at the same time
			time.Sleep(3*time.Hour + durPhase)
			durPhase = 0
		}
	}
}

func RequestToken(host, username, password string) (token string, err error) {

	tr := &http.Transport{
		DisableKeepAlives: true,
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		//RoundTrip:       roundTrip,
	}

	var DefaultTransport http.RoundTripper = tr

	oauthUrl := httpsAddr(host) + "/oauth/authorize?client_id=openshift-challenging-client&response_type=token"

	req, _ := http.NewRequest("HEAD", oauthUrl, nil)
	req.SetBasicAuth(username, password)

	resp, err := DefaultTransport.RoundTrip(req)

	//resp, err := client.Do(req)
	if err != nil {
		clog.Error(err)
		return "", err
	} else {
		defer resp.Body.Close()
		location, err := resp.Location()
		if err == nil {
			//fmt.Println("resp", url.Fragment)
			fragments := strings.Split(location.Fragment, "&")
			//n := proc(m)
			n := func(s []string) map[string]string {
				m := map[string]string{}
				for _, v := range s {
					n := strings.Split(v, "=")
					m[n[0]] = n[1]
				}
				return m
			}(fragments)

			//r, _ := json.Marshal(n)

			// return string(r), nil
			return n["access_token"], nil
		}
	}
	return token, err
}

// func proc(s []string) (m map[string]string) {
// 	m = map[string]string{}
// 	for _, v := range s {
// 		n := strings.Split(v, "=")
// 		m[n[0]] = n[1]
// 	}
// 	return
// }

func (oc *OpenshiftClient) request(method string, url string, body []byte, timeout time.Duration) (*http.Response, error) {
	token := oc.BearerToken()
	if token == "" {
		return nil, errors.New("token is blank")
	}

	var req *http.Request
	var err error
	if len(body) == 0 {
		req, err = http.NewRequest(method, url, nil)
	} else {
		req, err = http.NewRequest(method, url, bytes.NewReader(body))
	}

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	transCfg := &http.Transport{
		DisableKeepAlives: true,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: transCfg,
		Timeout:   timeout,
	}
	return client.Do(req)
}

type WatchStatus struct {
	Info []byte
	Err  error
}

func (oc *OpenshiftClient) doWatch(url string) (<-chan WatchStatus, chan<- struct{}, error) {
	res, err := oc.request("GET", url, nil, 0)
	if err != nil {
		return nil, nil, err
	}
	//if res.Body == nil {
	//	return nil, nil, errors.New("response.body is nil")
	//}

	statuses := make(chan WatchStatus, 5)
	canceled := make(chan struct{}, 1)

	go func() {
		defer func() {
			close(statuses)
			res.Body.Close()
		}()

		reader := bufio.NewReader(res.Body)
		for {
			select {
			case <-canceled:
				return
			default:
			}

			line, err := reader.ReadBytes('\n')
			if err != nil {
				statuses <- WatchStatus{nil, err}
				return
			}

			statuses <- WatchStatus{line, nil}
		}
	}()

	return statuses, canceled, nil
}

func (oc *OpenshiftClient) OWatch(uri string) (<-chan WatchStatus, chan<- struct{}, error) {
	return oc.doWatch(oc.oapiUrl + "/watch" + uri)
}

func (oc *OpenshiftClient) KWatch(uri string) (<-chan WatchStatus, chan<- struct{}, error) {
	return oc.doWatch(oc.kapiUrl + "/watch" + uri)
}

const GeneralRequestTimeout = time.Duration(60) * time.Second

/*
func (oc *OpenshiftClient) doRequest (method, url string, body []byte) ([]byte, error) {
	res, err := oc.request(method, url, body, GeneralRequestTimeout)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}

func (oc *OpenshiftClient) ORequest (method, uri string, body []byte) ([]byte, error) {
	return oc.doRequest(method, oc.oapiUrl + uri, body)
}

func (oc *OpenshiftClient) KRequest (method, uri string, body []byte) ([]byte, error) {
	return oc.doRequest(method, oc.kapiUrl + uri, body)
}
*/

type OpenshiftREST struct {
	oc         *OpenshiftClient
	Err        error
	StatusCode int
}

// client can't be nil now!
func NewOpenshiftREST(client *OpenshiftClient) *OpenshiftREST {
	//if client == nil {
	//	return &OpenshiftREST{oc: adminClient()}
	//}

	return &OpenshiftREST{oc: client}
}

// func (osr *OpenshiftREST) doRequest(method, url string, bodyParams interface{}, into interface{}) *OpenshiftREST {
// 	if osr.Err != nil {
// 		return osr
// 	}

// 	var body []byte
// 	if bodyParams != nil {
// 		body, osr.Err = json.Marshal(bodyParams)
// 		if osr.Err != nil {
// 			return osr
// 		}
// 	}

// 	//res, osr.Err := oc.request(method, url, body, GeneralRequestTimeout) // non-name error
// 	res, err := osr.oc.request(method, url, body, GeneralRequestTimeout)

// 	if err != nil {
// 		osr.Err = err
// 		return osr
// 	}
// 	defer res.Body.Close()

// 	osr.StatusCode = res.StatusCode

// 	var data []byte
// 	data, osr.Err = ioutil.ReadAll(res.Body)
// 	if osr.Err != nil {
// 		return osr
// 	}

// 	if res.StatusCode < 200 || res.StatusCode >= 400 {
// 		osr.Err = errors.New(string(data))
// 	} else {
// 		if into != nil {
// 			osr.Err = json.Unmarshal(data, into)
// 		}
// 	}

// 	return osr
// }

func (osr *OpenshiftREST) doRequest(method, url string, bodyParams interface{}, v interface{}) *OpenshiftREST {

	clog.Debugf("%v %v %#v", method, url, bodyParams)

	if osr.Err != nil {
		return osr
	}

	var body []byte
	if bodyParams != nil {
		body, osr.Err = json.Marshal(bodyParams)
		if osr.Err != nil {
			return osr
		}
	}

	//res, osr.Err := oc.request(method, url, body, GeneralRequestTimeout) // non-name error
	resp, err := osr.oc.request(method, url, body, GeneralRequestTimeout)

	if err != nil {
		osr.Err = err
		return osr
	}

	defer func() {
		// Drain up to 512 bytes and close the body to let the Transport reuse the connection
		io.CopyN(ioutil.Discard, resp.Body, 512)
		resp.Body.Close()
	}()

	////////////////
	defer resp.Body.Close()

	osr.StatusCode = resp.StatusCode

	err = CheckApiStatus(resp)
	if err != nil {
		// even though there was an error, we still return the response
		// in case the caller wants to inspect it further
		osr.Err = err
		clog.Error(err, osr.StatusCode)
		return osr
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err == io.EOF {
				err = nil // ignore EOF errors caused by empty response body
			}
			clog.Tracef("%#v", v)
			osr.Err = err
		}
	}

	return osr
}

func buildUriWithSelector(uri string, selector map[string]string) string {
	var buf bytes.Buffer
	for k, v := range selector {
		if buf.Len() > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(k)
		buf.WriteByte('=')
		buf.WriteString(v)
	}

	if buf.Len() == 0 {
		return uri
	}

	values := url.Values{}
	values.Set("labelSelector", buf.String())

	if strings.IndexByte(uri, '?') < 0 {
		uri = uri + "?"
	}

	println("\n uri=", uri+values.Encode(), "\n")

	return uri + values.Encode()
}

// o

func (osr *OpenshiftREST) OList(uri string, selector map[string]string, into interface{}) *OpenshiftREST {

	return osr.doRequest("GET", osr.oc.oapiUrl+buildUriWithSelector(uri, selector), nil, into)
}

func (osr *OpenshiftREST) OGet(uri string, into interface{}) *OpenshiftREST {
	return osr.doRequest("GET", osr.oc.oapiUrl+uri, nil, into)
}

func (osr *OpenshiftREST) ODelete(uri string, into interface{}) *OpenshiftREST {
	return osr.doRequest("DELETE", osr.oc.oapiUrl+uri, &kapi.DeleteOptions{}, into)
}

func (osr *OpenshiftREST) OPost(uri string, body interface{}, into interface{}) *OpenshiftREST {
	return osr.doRequest("POST", osr.oc.oapiUrl+uri, body, into)
}

func (osr *OpenshiftREST) OPut(uri string, body interface{}, into interface{}) *OpenshiftREST {
	return osr.doRequest("PUT", osr.oc.oapiUrl+uri, body, into)
}

// k

func (osr *OpenshiftREST) KList(uri string, selector map[string]string, into interface{}) *OpenshiftREST {
	return osr.doRequest("GET", osr.oc.kapiUrl+buildUriWithSelector(uri, selector), nil, into)
}

func (osr *OpenshiftREST) KGet(uri string, into interface{}) *OpenshiftREST {
	return osr.doRequest("GET", osr.oc.kapiUrl+uri, nil, into)
}

func (osr *OpenshiftREST) KDelete(uri string, into interface{}) *OpenshiftREST {
	return osr.doRequest("DELETE", osr.oc.kapiUrl+uri, &kapi.DeleteOptions{}, into)
}

func (osr *OpenshiftREST) KPost(uri string, body interface{}, into interface{}) *OpenshiftREST {
	return osr.doRequest("POST", osr.oc.kapiUrl+uri, body, into)
}

func (osr *OpenshiftREST) KPut(uri string, body interface{}, into interface{}) *OpenshiftREST {
	return osr.doRequest("PUT", osr.oc.kapiUrl+uri, body, into)
}
