package pkg

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/zonesan/clog"
)

type RechargeAgent struct {
	*Agent
	BaseURL *url.URL
}

type Recharge struct {
	Amount  float64 `json:"amount,omitempty"`
	Project string  `json:"namespace,omitempty"`
}

type HongPay apiRechargePayload

func (agent *RechargeAgent) AdminRecharge(r *http.Request, param *RequestParams) (*Balance, error) {

	urlStr := "/charge/v1/couponrecharge"

	balance := new(Balance)

	if err := doRequest(agent, r, "POST", urlStr, param, balance); err != nil {
		clog.Error(err)

		return nil, err
	}
	clog.Debug(balance)
	return balance, nil

}

func (agent *RechargeAgent) Create(r *http.Request, recharge *Recharge) (*HongPay, error) {

	urlStr := "/charge/v1/recharge"

	hongpay := new(HongPay)

	if err := doRequest(agent, r, "POST", urlStr, recharge, hongpay); err != nil {
		clog.Error(err)

		return nil, err
	}
	clog.Debug(hongpay.Payloads)
	return hongpay, nil

}

func (agent *RechargeAgent) Notification(r *http.Request) ([]byte, error) {
	urlStr := "/charge/v1/aipaycallback"

	if r.URL.RawQuery != "" {
		urlStr += "?" + r.URL.RawQuery
	}

	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := agent.BaseURL.ResolveReference(rel)

	reqbody, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		return nil, err
	}

	clog.Debug("Request Body:", string(reqbody))

	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(reqbody))
	if err != nil {
		clog.Error(err)
		return nil, err
	}

	resp, err := agent.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		// Drain up to 512 bytes and close the body to let the Transport reuse the connection
		io.CopyN(ioutil.Discard, resp.Body, 512)
		resp.Body.Close()
	}()

	data, err := ioutil.ReadAll(resp.Body)

	clog.Debugf("%s", data)

	return data, err

}

func (agent *RechargeAgent) Url() *url.URL {
	u := new(url.URL)
	u, _ = url.Parse(agent.BaseURL.String())
	return u
}

func (agent *RechargeAgent) Instance() *Agent {
	return agent.Agent
}

//weixin

type WeixinAgent struct {
	*Agent
	BaseURL *url.URL
}

func (agent *WeixinAgent) CreateWx(r *http.Request, recharge *Recharge) (*Resp, error) {

	urlStr := "/charge/v1/wechat/recharge"

	resp := new(Resp)

	if err := doRequest(agent, r, "POST", urlStr, recharge, resp); err != nil {
		clog.Error(err)

		return nil, err
	}
	clog.Debug(resp)
	return resp, nil

}

func (agent *WeixinAgent) GetStat(r *http.Request, tid string) (*Resp, error) {
	urlStr := fmt.Sprintf("/charge/v1/wechat/order/%v", tid)

	resp := new(Resp)
	if err := doRequest(agent, r, "GET", urlStr, nil, resp); err != nil {
		clog.Error(err)
		return nil, err
	}
	clog.Debug(resp)
	return resp, nil
}

func (agent *WeixinAgent) Url() *url.URL {
	u := new(url.URL)
	u, _ = url.Parse(agent.BaseURL.String())
	return u
}

func (agent *WeixinAgent) Instance() *Agent {
	return agent.Agent
}
