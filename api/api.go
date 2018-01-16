package api

import (
	"fmt"
	"net/http"

	"github.com/asiainfoLDP/datafoundry_payment/pkg"
	"github.com/julienschmidt/httprouter"
	"github.com/zonesan/clog"
)

type Mux struct{}

var PaymentAgent *pkg.Agent

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	clog.Info("from", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.Proto)
	//RespError(w, ErrorNew(ErrCodeNotFound), http.StatusNotFound)
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("not found"))
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	clog.Info("from", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.Proto)
	fmt.Fprint(w, "Welcome!\n")
}

func Agent() *pkg.Agent {
	if PaymentAgent == nil {
		PaymentAgent = pkg.NewAgent(nil)
	}
	return PaymentAgent
}
