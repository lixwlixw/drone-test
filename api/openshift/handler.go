package openshift

import (
	"net/http"

	"github.com/asiainfoLDP/datafoundry_payment/api"
	"github.com/julienschmidt/httprouter"
	"github.com/zonesan/clog"
)

func CreateProject(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	clog.Info("from", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.Proto)

	org := new(Orgnazition)

	if err := api.ParseRequestBody(r, org); err != nil {
		clog.Error("read request body error.", err)
		api.RespError(w, err)
		return
	}

	oc, err := NewClient(r, true)

	if err != nil {
		clog.Error("NewClient", err)
		api.RespError(w, err)
		return
	}

	if proj, err := oc.CreateProject(r, org.Name); err != nil {
		api.RespError(w, err)
	} else {
		api.RespOK(w, proj)
	}

}

func ListMembers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	clog.Info("from", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.Proto)

	project := ps.ByName("project")

	oc, err := NewAdminClient(r, project)

	if err != nil {
		clog.Error("NewAdminClient", err)
		api.RespError(w, err)
		return
	}

	if roles, err := oc.ListRoles(r, project); err != nil {
		api.RespError(w, err)
	} else {
		api.RespOK(w, roles)
	}

}

func InviteMember(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	clog.Info("from", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.Proto)

	member := new(OrgMember)

	if err := api.ParseRequestBody(r, member); err != nil {
		clog.Error("read request body error.", err)
		api.RespError(w, err)
		return
	}

	project := ps.ByName("project")

	oc, err := NewClient(r, false)

	if err != nil {
		clog.Error("NewClient", err)
		api.RespError(w, err)
		return
	}

	if role, err := oc.RoleAdd(r, project, member.Name, member.IsAdmin); err != nil {
		api.RespError(w, err)
	} else {
		api.RespOK(w, role)
	}

}

func RemoveMember(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	clog.Info("from", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.Proto)

	member := new(OrgMember)

	if err := api.ParseRequestBody(r, member); err != nil {
		clog.Error("read request body error.", err)
		api.RespError(w, err)
		return
	}

	project := ps.ByName("project")

	oc, err := NewClient(r, true)

	if err != nil {
		clog.Error("NewClient", err)
		api.RespError(w, err)
		return
	}

	if err := oc.RoleRemove(r, project, member.Name); err != nil {
		api.RespError(w, err)
	} else {
		api.RespOK(w, nil)
	}

}
