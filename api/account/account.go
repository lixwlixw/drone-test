package account

import (
	"net/http"

	"github.com/asiainfoLDP/datafoundry_payment/api"
	"github.com/julienschmidt/httprouter"
	"github.com/zonesan/clog"
)

func Account(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	clog.Info("from", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.Proto)
	agent := api.Agent()
	account, err := agent.Account.Get(r)

	if err != nil {
		api.RespError(w, err)
	} else {
		api.RespOK(w, account)
	}

	return
}
