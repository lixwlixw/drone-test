package main

import (
	"net/http"

	"github.com/asiainfoLDP/datafoundry_payment/api/openshift"
	"github.com/zonesan/clog"
)

func main() {

	initOpenshift()

	router := createRouter()

	//clog.SetLogLevel(clog.LOG_LEVEL_DEBUG)
	clog.Info("listening on port 8080...")
	clog.Fatal(http.ListenAndServe(":8080", router))
}

func initOpenshift() {
	openshift.Init()
}
