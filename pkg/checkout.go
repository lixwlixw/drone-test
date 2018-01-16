package pkg

import (
	//"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/zonesan/clog"
)

type CheckoutAgent struct {
	*Agent
	BaseURL *url.URL
}

type Checkout struct {
	PlanId  string                 `json:"plan_id"`
	Project string                 `json:"namespace,omitempty"`
	Region  string                 `json:"region"` //need it?
	OrderID string                 `json:"orderid,omitempty"`
	DryTry  int                    `json:"drytry,omitempty"`
	Params  map[string]interface{} `json:"parameters,omitempty"`
}

type PurchasedOrder apiPurchaseOrder

func (agent *CheckoutAgent) Get() *Balance {

	balance := &Balance{
		Balance: 3000.01,
		Status:  "active",
	}
	return balance
}

func (agent *CheckoutAgent) ListOrders(r *http.Request) (*[]PurchasedOrder, error) {
	urlStr := "/usageapi/v1/orders"

	orders := new([]PurchasedOrder)

	if err := doRequestList(agent, r, "GET", urlStr, nil, orders); err != nil {
		clog.Error(err)
		return nil, err
	}

	clog.Infof("%v order(s) listed.", len(*orders))

	return orders, nil

}

func (agent *CheckoutAgent) Create(r *http.Request, checkout *Checkout) (*PurchasedOrder, error) {
	urlStr := "/usageapi/v1/orders"

	order := new(PurchasedOrder)
	if err := doRequest(agent, r, "POST", urlStr, checkout, order); err != nil {
		clog.Error(err)
		return nil, err
	}

	return order, nil
}

func (agent *CheckoutAgent) Unsubscribe(r *http.Request, orderid string) (*UndefinedResp, error) {
	urlStr := fmt.Sprintf("/usageapi/v1/orders/%v", orderid)

	r.ParseForm()

	req := make(map[string]string)

	req["action"] = r.FormValue("action")
	req["namespace"] = r.FormValue("namespace")

	resp := new(UndefinedResp)

	if err := doRequest(agent, r, "PUT", urlStr, req, resp); err != nil {
		clog.Error(err)
		return nil, err
	}

	return resp, nil
}

func (agent *CheckoutAgent) Url() *url.URL {
	u := new(url.URL)
	u, _ = url.Parse(agent.BaseURL.String())
	return u
}

func (agent *CheckoutAgent) Instance() *Agent {
	return agent.Agent
}
