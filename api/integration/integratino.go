package integration

import (
	"net/http"

	"github.com/asiainfoLDP/datafoundry_payment/api"
	"github.com/julienschmidt/httprouter"
	"github.com/zonesan/clog"
)

func ListRepos(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	clog.Info("from", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.Proto)

	agent := api.Agent()
	if result, err := agent.Integration.ListRepos(r); err != nil {
		api.RespError(w, err)
	} else {
		api.RespOK(w, result)
	}
}

func GetRepo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	clog.Info("from", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.Proto)

	repo := ps.ByName("repo")

	agent := api.Agent()
	if result, err := agent.Integration.GetRepo(r, repo); err != nil {
		api.RespError(w, err)
	} else {
		api.RespOK(w, result)
	}
}

func ListItems(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	clog.Info("from", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.Proto)

	agent := api.Agent()
	if result, err := agent.Recharge.Notification(r); err != nil {
		api.RespError(w, err)
	} else {
		api.RespOK(w, result)
	}
}

func GetItem(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	clog.Info("from", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.Proto)

	repo := ps.ByName("repo")
	item := ps.ByName("item")

	agent := api.Agent()
	if result, err := agent.Integration.GetItem(r, repo, item); err != nil {
		api.RespError(w, err)
	} else {
		api.RespOK(w, result)
	}
}

func ListDataServices(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	clog.Info("from", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.Proto)

	agent := api.Agent()
	if result, err := agent.DataService.ListServices(r); err != nil {
		api.RespError(w, err)
	} else {
		api.RespOK(w, result)
	}
}

func DataServiceInstance(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	clog.Info("from", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.Proto)

	instance_id := ps.ByName("instance_id")

	agent := api.Agent()
	if result, err := agent.DataService.CreateServiceInstance(r, instance_id); err != nil {
		api.RespError(w, err)
	} else {
		api.RespOK(w, result)
	}
}
