package recharge

import (
	"net/http"

	"github.com/asiainfoLDP/datafoundry_payment/api"
	"github.com/asiainfoLDP/datafoundry_payment/pkg"
	"github.com/julienschmidt/httprouter"
	"github.com/zonesan/clog"
)

func AiPayRecharge(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	clog.Info("from", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.Proto)

	recharge := new(pkg.Recharge)

	if err := api.ParseRequestBody(r, recharge); err != nil {
		clog.Error("read request body error.", err)
		api.RespError(w, err)
		return
	}

	agent := api.Agent()
	if hongpay, err := agent.Recharge.Create(r, recharge); err != nil {
		api.RespError(w, err)
	} else {
		api.RespOK(w, hongpay)
	}
}

func WxRecharge(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	clog.Info("from", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.Proto)

	recharge := new(pkg.Recharge)

	if err := api.ParseRequestBody(r, recharge); err != nil {
		clog.Error("read request body error.", err)
		api.RespError(w, err)
		return
	}

	agent := api.Agent()
	if hongpay, err := agent.Weixin.CreateWx(r, recharge); err != nil {
		api.RespError(w, err)
	} else {
		api.RespOK(w, hongpay)
	}
}

func WxGetRechargeStat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	clog.Info("from", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.Proto)

	tid := ps.ByName("tid")

	agent := api.Agent()

	if result, err := agent.Weixin.GetStat(r, tid); err != nil {
		api.RespError(w, err)
	} else {
		api.RespOK(w, result)
	}
}

func DirectRecharge(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	clog.Info("from", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.Proto)

	params := new(pkg.RequestParams)

	if err := api.ParseRequestBody(r, params); err != nil {
		clog.Error("read request body error.", err)
		api.RespError(w, err)
		return
	}

	agent := api.Agent()
	if hongpay, err := agent.Recharge.AdminRecharge(r, params); err != nil {
		api.RespError(w, err)
	} else {
		api.RespOK(w, hongpay)
	}
}

func Notification(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	clog.Info("from", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.Proto)

	agent := api.Agent()
	if resp, err := agent.Recharge.Notification(r); err != nil {
		api.RespError(w, err)
	} else {
		api.RespOK(w, resp)
	}
}
