package amount

import (
	"net/http"

	"github.com/asiainfoLDP/datafoundry_payment/api"
	"github.com/julienschmidt/httprouter"
	"github.com/zonesan/clog"
)

func AmountList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	clog.Info("from", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.Proto)

	agent := api.Agent()
	amounts, err := agent.Amount.List(r)

	if err != nil {
		api.RespError(w, err)
	} else {
		api.RespOK(w, amounts)
	}

	return
}

func Amount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	clog.Info("from", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.Proto)

	tid := ps.ByName("tid")

	agent := api.Agent()
	amount, err := agent.Amount.Get(r, tid)

	if err != nil {
		api.RespError(w, err)
	} else {
		api.RespOK(w, amount)
	}

	return
}
