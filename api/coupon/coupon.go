package coupon

import (
	"net/http"

	"github.com/asiainfoLDP/datafoundry_payment/api"
	"github.com/asiainfoLDP/datafoundry_payment/pkg"
	"github.com/julienschmidt/httprouter"
	"github.com/zonesan/clog"
)

func Coupon(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	clog.Info("from", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.Proto)

	serial := ps.ByName("serial")

	agent := api.Agent()
	coupon, err := agent.Coupon.Get(r, serial)

	if err != nil {
		api.RespError(w, err)
	} else {
		api.RespOK(w, coupon)
	}

	return
}

func CreateCoupon(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	clog.Info("from", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.Proto)

	param := new(pkg.RequestParams)

	if err := api.ParseRequestBody(r, param); err != nil {
		clog.Error("read request body error.", err)
		api.RespError(w, err)
		return
	}

	agent := api.Agent()
	result, err := agent.Coupon.AdminCreate(r, param)

	if err != nil {
		api.RespError(w, err)
	} else {
		api.RespOK(w, result)
	}

	return
}

func Redeem(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	clog.Info("from", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.Proto)

	redeem := new(pkg.Redeem)

	if err := api.ParseRequestBody(r, redeem); err != nil {
		clog.Error("read request body error.", err)
		api.RespError(w, err)
		return
	}

	agent := api.Agent()

	coupon, err := agent.Coupon.Redeem(r, redeem)

	if err != nil {
		api.RespError(w, err)
	} else {
		api.RespOK(w, coupon)
	}
	//http.Redirect(w, r, "http://www.google.com", http.StatusMovedPermanently)
}
