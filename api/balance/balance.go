package balance

import (
	"net/http"

	"github.com/asiainfoLDP/datafoundry_payment/api"
	"github.com/julienschmidt/httprouter"
	"github.com/zonesan/clog"
)

func Balance(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	clog.Info("from", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.Proto)

	agent := api.Agent()

	if balance, err := agent.Balance.Get(r); err != nil {
		api.RespError(w, err)
	} else {
		api.RespOK(w, balance)
	}

}
