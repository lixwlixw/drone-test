package plan

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/zonesan/clog"
	"net/http"
)

func Plans(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	clog.Info("from", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.Proto)
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}
