package pkg

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/zonesan/clog"
)

var (
	// defaultBaseURL = "http://localhost:7071"

	// BaseURL should always be specified WITHOUT a trailing slash.
	defaultCouponBaseURL      = "https://datafoundry.coupon.app.dataos.io"
	defaultMarketBaseURL      = "https://datafoundry.plan.app.dataos.io"
	defaultCheckoutBaseURL    = "https://datafoundry.serviceusage.app.dataos.io"
	defaultBalanceBaseURL     = "https://datafoundry.recharge.app.dataos.io"
	defaultRechargeBaseURL    = defaultBalanceBaseURL
	defaultAmountBaseURL      = defaultBalanceBaseURL
	defaultIntegrationBaseURL = "https://datafoundry.integration.app.dataos.io"
	defaultDataServiceBaseURL = "https://datafoundry.datainstance.app.dataos.io"
	defaultWeixinBaseURL      = "https://datafoundry.wechat.app.dataos.io"
)

func InitBaseUrls() {
	couponurl := combainHostPort("ENV_COUPON_HOST", "ENV_COUPON_PORT")
	if len(couponurl) > 0 {
		defaultCouponBaseURL = httpAddr(couponurl)
	}

	marketurl := combainHostPort("ENV_MARKET_HOST", "ENV_MARKET_PORT")
	if len(marketurl) > 0 {
		defaultMarketBaseURL = httpAddr(marketurl)
	}

	checkouturl := combainHostPort("ENV_CHECKOUT_HOST", "ENV_CHECKOUT_PORT")
	if len(checkouturl) > 0 {
		defaultCheckoutBaseURL = httpAddr(checkouturl)
	}

	balanceurl := combainHostPort("ENV_BALANCE_HOST", "ENV_BALANCE_PORT")
	if len(balanceurl) > 0 {
		defaultBalanceBaseURL = httpAddr(balanceurl)
		defaultRechargeBaseURL = defaultBalanceBaseURL
		defaultAmountBaseURL = defaultBalanceBaseURL
	}

	integrationurl := combainHostPort("ENV_INTEGRATION_HOST", "ENV_INTEGRATION_PORT")
	if len(integrationurl) > 0 {
		defaultIntegrationBaseURL = httpAddr(integrationurl)
	}

	dataserviceurl := combainHostPort("ENV_DATAINSTANCE_HOST", "ENV_DATAINSTANCE_PORT")
	if len(dataserviceurl) > 0 {
		defaultDataServiceBaseURL = httpAddr(dataserviceurl)
	}

	weixinurl := combainHostPort("ENV_WEIXIN_HOST", "ENV_WEIXIN_PORT")
	if len(weixinurl) > 0 {
		defaultWeixinBaseURL = httpAddr(weixinurl)
	}

	clog.Info("couponurl", defaultCouponBaseURL)
	clog.Info("marketurl", defaultMarketBaseURL)
	clog.Info("checkouturl", defaultCheckoutBaseURL)
	clog.Info("balanceurl", defaultBalanceBaseURL)
	clog.Info("rechargeurl", defaultRechargeBaseURL)
	clog.Info("amounturl", defaultAmountBaseURL)
	clog.Info("integrationurl", defaultIntegrationBaseURL)
	clog.Info("dataserviceurl", defaultDataServiceBaseURL)
	clog.Info("weixinurl", defaultWeixinBaseURL)

}

func combainHostPort(hostenv, portenv string) string {
	host := os.Getenv(os.Getenv(hostenv))
	port := os.Getenv(os.Getenv(portenv))
	if port == "" {
		return host
	}
	return host + ":" + port
}

// An Agent manages communication with the payment components API.
type Agent struct {
	clientMu sync.Mutex   // clientMu protects the client during calls that modify the CheckRedirect func.
	client   *http.Client // HTTP client used to communicate with the API.

	// Base URL for API requests.  Defaults to the public GitHub API, BaseURL should
	// always be specified with a trailing slash.

	// BaseURL *url.URL

	common *service // Reuse a single struct instead of allocating one for each service on the heap.

	// Agents used for talking to different parts of the payment components API.
	Recharge    *RechargeAgent
	Checkout    *CheckoutAgent
	Balance     *BalanceAgent
	Market      *MarketAgent
	Amount      *AmountAgent
	Account     *AccountAgent
	Coupon      *CouponAgent
	Integration *IntegrationAgent
	DataService *DataServiceAgent
	Weixin      *WeixinAgent
}

type service struct {
	*Agent
}

func NewAgent(httpClient *http.Client) *Agent {

	if httpClient == nil {
		tr := &http.Transport{
			//DisableKeepAlives: true,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		httpClient = &http.Client{Transport: tr}
	}

	// baseURL, _ := url.Parse(defaultBaseURL)

	// agent := &Agent{client: httpClient, BaseURL: baseURL}
	agent := &Agent{client: httpClient}

	marketBaseURL, _ := url.Parse(defaultMarketBaseURL)
	couponBaseURL, _ := url.Parse(defaultCouponBaseURL)
	checkoutBaseURL, _ := url.Parse(defaultCheckoutBaseURL)
	balanceBaseURL, _ := url.Parse(defaultBalanceBaseURL)
	rechargeBaseURL, _ := url.Parse(defaultRechargeBaseURL)
	amountBaseURL, _ := url.Parse(defaultAmountBaseURL)
	integrationBaseURL, _ := url.Parse(defaultIntegrationBaseURL)
	dataServiceBaseURL, _ := url.Parse(defaultDataServiceBaseURL)
	weixinBaseURL, _ := url.Parse(defaultWeixinBaseURL)

	service := &service{agent}

	agent.common = service
	agent.Account = (*AccountAgent)(agent.common)
	agent.Amount = &AmountAgent{Agent: agent.common.Agent, BaseURL: amountBaseURL}
	agent.Balance = &BalanceAgent{Agent: agent.common.Agent, BaseURL: balanceBaseURL}
	agent.Checkout = &CheckoutAgent{Agent: agent.common.Agent, BaseURL: checkoutBaseURL}
	agent.Coupon = &CouponAgent{Agent: agent.common.Agent, BaseURL: couponBaseURL}
	agent.Market = &MarketAgent{Agent: agent.common.Agent, BaseURL: marketBaseURL}
	agent.Recharge = &RechargeAgent{Agent: agent.common.Agent, BaseURL: rechargeBaseURL}
	agent.Integration = &IntegrationAgent{Agent: agent.common.Agent, BaseURL: integrationBaseURL}
	agent.DataService = &DataServiceAgent{Agent: agent.common.Agent, BaseURL: dataServiceBaseURL}
	agent.Weixin = &WeixinAgent{Agent: agent.common.Agent, BaseURL: weixinBaseURL}

	return agent
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash.  If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (*Agent) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {

	clog.Info(method, urlStr, body)

	if _, err := url.Parse(urlStr); err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(strings.ToUpper(method), urlStr, buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	// req.Header.Set("Accept", mediaTypeV3)
	// if c.UserAgent != "" {
	// 	req.Header.Set("User-Agent", c.UserAgent)
	// }
	return req, nil
}

// Do sends an API request and returns the API response.  The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred.  If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.  If rate limit is exceeded and reset time is in the future,
// Do returns *RateLimitError immediately without making a network API call.
func (c *Agent) Do(req *http.Request, v interface{}) error {

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		// Drain up to 512 bytes and close the body to let the Transport reuse the connection
		io.CopyN(ioutil.Discard, resp.Body, 512)
		resp.Body.Close()
	}()

	err = CheckResponse(resp)
	if err != nil {
		// even though there was an error, we still return the response
		// in case the caller wants to inspect it further
		return err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			// var data []byte
			// data, err = ioutil.ReadAll(resp.Body)
			// if err == nil && data != nil {
			// 	clog.Debugf("%s", data)
			// 	err = json.Unmarshal(data, v)
			// }

			err = json.NewDecoder(resp.Body).Decode(v)
			if err == io.EOF {
				err = nil // ignore EOF errors caused by empty response body
			}
		}
	}

	return err
}

type Error struct {
	Resource string `json:"resource"` // resource on which the error occurred
	Field    string `json:"field"`    // field on which the error occurred
	Code     string `json:"code"`     // validation error code
	Message  string `json:"message"`  // Message describing the error. Errors with Code == "custom" will always have this set.
}

func (e *Error) Error() string {
	return fmt.Sprintf("%v error caused by %v field on %v resource",
		e.Code, e.Field, e.Resource)
}

type ErrorResponse struct {
	Response *http.Response // HTTP response that caused this error
	Message  string         `json:"msg"`  // error message
	Code     int            `json:"code"` // more detail on individual errors

}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v %+v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Message, r.Code)
}

// CheckResponse checks the API response for errors, and returns them if
// present.  A response is considered an error if it has a status code outside
// the 200 range.  API error responses are expected to have either no response
// body, or a JSON response body that maps to ErrorResponse.  Any other
// response body will be silently ignored.
//
// The error type will be *RateLimitError for rate limit exceeded errors,
// and *TwoFactorAuthError for two-factor authentication errors.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		//clog.Errorf("%s", data)
		json.Unmarshal(data, errorResponse)
	}
	return errorResponse
}

func init() {
	InitBaseUrls()
}
