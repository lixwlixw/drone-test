package pkg

import (
	"net/http"

	"github.com/zonesan/clog"
)

type AccountAgent service

type Account struct {
	Purchased bool     `json:"purchased"`
	Notify    bool     `json:"notification"`
	Plans     []Plan   `json:"subscriptions,omitempty"`
	Status    string   `json:"status"`
	Balance   *Balance `json:"balance"`
}

func (agent *AccountAgent) Get(r *http.Request) (*Account, error) {
	r.ParseForm()

	project := r.FormValue("namespace")

	clog.Debug(project)

	_, err := getToken(r)
	if err != nil {
		clog.Error(err)
		return nil, err
	}

	account := new(Account)

	c := make(chan bool)

	go func() {
		if account.Balance, err = agent.Balance.Get(r); err != nil {
			clog.Error(err)
		}
		close(c)
	}()

	var plans *Market
	c1 := make(chan bool)

	go func() {
		if plans, err = agent.Market.ListPlan(r); err != nil {
			clog.Error(err)
		}
		close(c1)
	}()

	if orders, err := agent.Checkout.ListOrders(r); err != nil {
		clog.Error(err)
		return nil, err
	} else {
		//clog.Debugf("%#v", orders)

		if len(*orders) > 0 {
			account.Purchased = true

			<-c1

			if plans != nil {
				func() {
					for _, order := range *orders {
						found := false
						for _, plan := range plans.Plans {
							if order.Order.Plan_id == plan.PlanId {
								account.Plans = append(account.Plans, plan)
								clog.Debug(order.Order.Plan_id, "found in plan list.")
								found = true
								break
							}
						}
						if !found {
							clog.Warnf("order with plan id '%v' not found in market.", order.Order.Plan_id)
						}
					}

				}()
			}

		}
	}

	<-c

	//account := fakeAccount(r)
	return account, nil
}
