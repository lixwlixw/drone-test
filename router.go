package main

import (
	"github.com/asiainfoLDP/datafoundry_payment/api"
	"github.com/asiainfoLDP/datafoundry_payment/api/account"
	"github.com/asiainfoLDP/datafoundry_payment/api/amount"
	"github.com/asiainfoLDP/datafoundry_payment/api/balance"
	"github.com/asiainfoLDP/datafoundry_payment/api/checkout"
	"github.com/asiainfoLDP/datafoundry_payment/api/coupon"
	"github.com/asiainfoLDP/datafoundry_payment/api/integration"
	"github.com/asiainfoLDP/datafoundry_payment/api/market"
	"github.com/asiainfoLDP/datafoundry_payment/api/openshift"
	"github.com/asiainfoLDP/datafoundry_payment/api/recharge"
	"github.com/julienschmidt/httprouter"
)

const (
	PAYMENT_API_PREFIX     = "/payment/v1"
	INTEGRATION_API_PREFIX = "/integration/v1"
	LDP_API_PREFIX         = "/lapi/v1"
)

func createRouter() *httprouter.Router {
	router := httprouter.New()

	router.GET("/", api.Index)

	router.POST(PAYMENT_API_PREFIX+"/recharge", recharge.AiPayRecharge)
	router.POST(PAYMENT_API_PREFIX+"/recharge/weixin", recharge.WxRecharge)
	router.GET(PAYMENT_API_PREFIX+"/recharge/weixin/:tid", recharge.WxGetRechargeStat)
	router.POST(PAYMENT_API_PREFIX+"/directrecharge", recharge.DirectRecharge) //admin permission
	router.POST(PAYMENT_API_PREFIX+"/checkout", checkout.Checkout)
	router.GET(PAYMENT_API_PREFIX+"/balance", balance.Balance)
	router.GET(PAYMENT_API_PREFIX+"/market", market.Market)
	router.GET(PAYMENT_API_PREFIX+"/amounts", amount.AmountList)
	router.GET(PAYMENT_API_PREFIX+"/amounts/:tid", amount.Amount)
	router.GET(PAYMENT_API_PREFIX+"/account", account.Account)
	router.GET(PAYMENT_API_PREFIX+"/coupon/:serial", coupon.Coupon)
	router.POST(PAYMENT_API_PREFIX+"/coupon", coupon.CreateCoupon) //admin permission
	router.POST(PAYMENT_API_PREFIX+"/redeem", coupon.Redeem)
	router.GET(PAYMENT_API_PREFIX+"/orders", checkout.Order)
	router.DELETE(PAYMENT_API_PREFIX+"/orders/:orderid", checkout.Unsubscribe)
	router.POST(PAYMENT_API_PREFIX+"/notification", recharge.Notification)
	router.GET(PAYMENT_API_PREFIX+"/regions", market.ListRegion)

	router.GET(INTEGRATION_API_PREFIX+"/repos", integration.ListRepos)
	router.GET(INTEGRATION_API_PREFIX+"/repos/:repo", integration.GetRepo)
	router.GET(INTEGRATION_API_PREFIX+"/repos/:repo/items/:item", integration.GetItem)

	router.GET(INTEGRATION_API_PREFIX+"/services", integration.ListDataServices)
	router.POST(INTEGRATION_API_PREFIX+"/instance/:instance_id", integration.DataServiceInstance)

	router.POST(LDP_API_PREFIX+"/orgs", openshift.CreateProject)
	router.GET(LDP_API_PREFIX+"/orgs/:project/roles", openshift.ListMembers)
	router.PUT(LDP_API_PREFIX+"/orgs/:project/invite", openshift.InviteMember)
	router.PUT(LDP_API_PREFIX+"/orgs/:project/remove", openshift.RemoveMember)

	router.NotFound = &api.Mux{}

	return router
}
