package market

import (
	"net/http"

	"github.com/asiainfoLDP/datafoundry_payment/api"
	"github.com/julienschmidt/httprouter"
	"github.com/zonesan/clog"
)

func Market(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	clog.Info("from", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.Proto)

	agent := api.Agent()
	market, err := agent.Market.ListPlan(r)

	if err != nil {
		api.RespError(w, err)
	} else {
		api.RespOK(w, market)
	}

}

func ListRegion(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	clog.Info("from", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.Proto)

	agent := api.Agent()
	market, err := agent.Market.ListRegion(r)

	if err != nil {
		api.RespError(w, err)
	} else {
		api.RespOK(w, market)
	}

}
