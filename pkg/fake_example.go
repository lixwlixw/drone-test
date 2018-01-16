package pkg

import (
	"net/http"
	"time"
)

func fakeAccount(r *http.Request) *Account {

	account := &Account{
		Purchased: false,
		Notify:    false,
		Plans: []Plan{
			{
				PlanId:     "29BA4085-B3B4-8308-2097-7A5340E770B9",
				Type:       "C",
				Price:      10.00,
				BillPeriod: "monthly",
				Desc:       "1 CPU Core, 512M Memory",
			},
		},
	}

	r.ParseForm()
	if len(r.FormValue("n")) > 0 {
		account.Purchased = true
		account.Notify = true
	}

	return account
}

func fakeAmount(r *http.Request) *Amount {

	amount := &Amount{
		Id:            "03F232X238DKJ",
		CreationTime:  time.Now().Format(time.RFC3339),
		Amount:        12.23,
		Desc:          "Plan A",
		PaymentMethod: "balance",
		Status:        "finish",
	}

	r.ParseForm()
	if len(r.FormValue("n")) > 0 {
		amount.Status = "refunded"
	}

	return amount
}

func fakeAmounts(r *http.Request) *[]Amount {

	amounts := &[]Amount{
		{
			Id:            "03Fwerqe2",
			CreationTime:  time.Now().Format(time.RFC3339),
			Amount:        12.23,
			Desc:          "Plan A",
			PaymentMethod: "balance",
			Status:        "finish",
		},
		{
			Id:            "qwer238DKJ",
			CreationTime:  time.Now().Format(time.RFC3339),
			Amount:        12.23,
			Desc:          "Plan A",
			PaymentMethod: "balance",
			Status:        "finish",
		},
		{
			Id:            "03F232X238DKJ",
			CreationTime:  time.Now().Format(time.RFC3339),
			Amount:        2.34,
			Desc:          "Plan A",
			PaymentMethod: "balance",
			Status:        "refunded",
		},
	}

	return amounts
}

func fakeMarket() *Market {
	market := &Market{
		Plans: []Plan{
			{
				PlanId:       "1d3452ea-7f14-11e6-9fe0-2344dd5557c3",
				Type:         "C",
				Price:        20,
				BillPeriod:   "M",
				Desc:         "1 CPU Core, 512M Memory",
				CreationTime: time.Now().Format(time.RFC3339),
			},
			{
				PlanId:       "1d3452ea-7f14-11e6-9fe0-2344dd5557c3",
				Type:         "C",
				Price:        40.88,
				BillPeriod:   "M",
				Desc:         "2 CPU Cores, 1G Memory",
				CreationTime: time.Now().Format(time.RFC3339),
			},
			{
				PlanId:       "1d3452ea-7f14-11e6-9fe0-2344dd5557c3",
				Type:         "C",
				Price:        88.88,
				BillPeriod:   "M",
				Desc:         "4 CPU Cores, 2G Memory",
				CreationTime: time.Now().Format(time.RFC3339),
			},
		},
	}

	return market
}
