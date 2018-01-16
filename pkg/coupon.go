package pkg

import (
	"fmt"
	"net/http"
	"net/url"

	apierrors "github.com/asiainfoLDP/datafoundry_payment/pkg/errors"
	"github.com/zonesan/clog"
)

type CouponAgent struct {
	*Agent
	BaseURL *url.URL
}

type Coupon apiCoupon

type Redeem struct {
	Serial    string `json:"serial,omitempty"`
	Code      string `json:"code"`
	Namespace string `json:"namespace"`
	Region    string `json:"region"`
}

func (agent *CouponAgent) Get(r *http.Request, code string) (*Coupon, error) {

	urlStr := fmt.Sprintf("/charge/v1/coupons/%v", code)

	coupon := new(Coupon)
	if err := doRequest(agent, r, "GET", urlStr, nil, coupon); err != nil {
		clog.Error(err)
		return nil, err
	}
	clog.Debug(coupon)
	coupon.Region = ""
	return coupon, nil
}

func (agent *CouponAgent) AdminCreate(r *http.Request, param *RequestParams) (*Coupon, error) {

	urlStr := "/charge/v1/coupons"

	coupon := new(Coupon)
	if err := doRequest(agent, r, "POST", urlStr, param, coupon); err != nil {
		clog.Error(err)
		return nil, err
	}
	clog.Debug(coupon)
	coupon.Region = ""
	return coupon, nil
}

func (agent *CouponAgent) Redeem(r *http.Request, redeem *Redeem) (*Coupon, error) {
	if redeem.Serial == "" || redeem.Code == "" {
		return nil, apierrors.ErrorNew(apierrors.ErrCodeBadRequest)
	}

	urlStr := fmt.Sprintf("/charge/v1/coupons/use/%s", redeem.Serial)

	coupon := new(Coupon)
	if err := doRequest(agent, r, "PUT", urlStr, redeem, coupon); err != nil {
		clog.Error(err)
		return nil, err
	}

	clog.Debug(coupon)
	return coupon, nil
}

func (agent *CouponAgent) Url() *url.URL {
	u := new(url.URL)
	u, _ = url.Parse(agent.BaseURL.String())
	return u
}

func (agent *CouponAgent) Instance() *Agent {
	return agent.Agent
}
