package checkout

import (
	"net/http"

	"github.com/asiainfoLDP/datafoundry_payment/api"
	"github.com/asiainfoLDP/datafoundry_payment/pkg"
	"github.com/julienschmidt/httprouter"
	"github.com/zonesan/clog"
)

func Checkout(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	clog.Info("from", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.Proto)

	checkout := new(pkg.Checkout)

	if err := api.ParseRequestBody(r, checkout); err != nil {
		clog.Error("read request body error.", err)
		api.RespError(w, err)
		return
	}

	agent := api.Agent()

	checkoutResult, err := agent.Checkout.Create(r, checkout)

	if err != nil {
		api.RespError(w, err)
	} else {
		api.RespOK(w, checkoutResult)
	}
	//http.Redirect(w, r, "http://www.google.com", http.StatusMovedPermanently)
}

func Order(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	clog.Info("from", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.Proto)

	agent := api.Agent()

	orders, err := agent.Checkout.ListOrders(r)

	if err != nil {
		api.RespError(w, err)
	} else {
		api.RespOK(w, orders)
	}
	//http.Redirect(w, r, "http://www.google.com", http.StatusMovedPermanently)
}

func Unsubscribe(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	clog.Info("from", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.Proto)

	orderid := ps.ByName("orderid")

	agent := api.Agent()
	result, err := agent.Checkout.Unsubscribe(r, orderid)

	if err != nil {
		api.RespError(w, err)
	} else {
		api.RespOK(w, result)
	}
	//http.Redirect(w, r, "http://www.google.com", http.StatusMovedPermanently)
}
