package pkg

import (
	// "encoding/json"
	"net/http"
	"net/url"
	// "reflect"

	"github.com/zonesan/clog"
)

//type BalanceAgent service

type BalanceAgent struct {
	*Agent
	BaseURL *url.URL
}

type Balance apiBalance

func (agent *BalanceAgent) Get(r *http.Request) (*Balance, error) {

	urlStr := "/charge/v1/balance"

	balance := new(Balance)

	if err := doRequest(agent, r, "GET", urlStr, nil, balance); err != nil {
		clog.Error(err)
		return nil, err
	}
	clog.Debug(balance)

	return balance, nil

}

func (agent *BalanceAgent) Url() *url.URL {
	u := new(url.URL)
	u, _ = url.Parse(agent.BaseURL.String())
	return u
}

func (agent *BalanceAgent) Instance() *Agent {
	return agent.Agent
}
