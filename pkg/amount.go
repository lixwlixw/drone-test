package pkg

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/zonesan/clog"
)

type AmountAgent struct {
	*Agent
	BaseURL *url.URL
}

type Amount struct {
	Id            string  `json:"trans_id"`
	CreationTime  string  `json:"creation_time"`
	Amount        float64 `json:"amount"`
	User          string  `json:"user"`
	Desc          string  `json:"description"`
	PaymentMethod string  `json:"payment_method"`
	Reason        string  `json:"reason"`
	Namespace     string  `json:"namespace"`
	Status        string  `json:"status"`
}

type Amounts struct {
	Amounts []Amount `json:"amounts"`
}

func (agent *AmountAgent) Get(r *http.Request, tid string) (*Amount, error) {

	urlStr := fmt.Sprintf("/charge/v1/recharge/%v", tid)

	amount := new(Amount)
	if err := doRequest(agent, r, "GET", urlStr, nil, amount); err != nil {
		clog.Error(err)
		return nil, err
	}
	clog.Debug(amount)
	return amount, nil
}

func (agent *AmountAgent) List(r *http.Request) (*Amounts, error) {

	urlStr := "/charge/v1/recharge"

	transactions := []apiTransaction{}
	amounts := new(Amounts)

	if err := doRequestList(agent, r, "GET", urlStr, nil, &transactions); err != nil {
		clog.Error(err)
		return nil, err
	} else {
		for _, transaction := range transactions {
			amount := Amount{
				Id:            transaction.TransactionId,
				CreationTime:  transaction.CreateTime,
				Amount:        transaction.Amount,
				Desc:          transaction.Type,
				User:          transaction.User,
				Reason:        transaction.Reason,
				Namespace:     transaction.Namespace,
				PaymentMethod: transaction.PaymentMethod,
				Status:        transaction.Status,
			}
			amounts.Amounts = append(amounts.Amounts, amount)
		}
	}

	return amounts, nil
}

func (agent *AmountAgent) Url() *url.URL {
	u := new(url.URL)
	u, _ = url.Parse(agent.BaseURL.String())
	return u
}

func (agent *AmountAgent) Instance() *Agent {
	return agent.Agent
}
