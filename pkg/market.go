package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	// "reflect"

	"github.com/zonesan/clog"
)

type MarketAgent struct {
	*Agent
	BaseURL *url.URL
}

type Plan struct {
	PlanId       string                 `json:"plan_id"`
	Name         string                 `json:"plan_name"`
	Level        int                    `json:"plan_level"`
	Type         string                 `json:"type"`
	Price        float32                `json:"price"`
	BillPeriod   string                 `json:"bill_period"`
	Desc         string                 `json:"description"`
	Desc2        string                 `json:"description2,omitempty"`
	Detail       string                 `json:"detail,omitempty"`
	Region       string                 `json:"region,omitempty"`
	RegionID     string                 `json:"region_id,omitempty"`
	CreationTime string                 `json:"creation_time,omitempty"`
	Customize    map[string]interface{} `json:"customize,omitempty"`
}

type Market struct {
	Plans []Plan `json:"plans"`
}

type Region apiRegion

func (agent *MarketAgent) Get(r *http.Request, id string) (*Plan, error) {
	// market := fakeMarket()

	urlStr := fmt.Sprintf("/charge/v1/plans/%v", id)

	if r.URL.RawQuery != "" {
		urlStr += "?" + r.URL.RawQuery
	}

	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := agent.BaseURL.ResolveReference(rel)

	req, err := agent.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	response := new(RemoteResponse)

	if err := agent.Do(req, response); err != nil {
		return nil, err
	}

	apiplan := new(apiPlan)
	plan := new(Plan)

	if err := json.Unmarshal([]byte(response.Data), apiplan); err != nil {
		clog.Error(err)
		return nil, err
	} else {
		clog.Infof("%#v", apiplan)
		plan = convertPlan(apiplan)
	}

	return plan, nil

}

func (agent *MarketAgent) ListPlan(r *http.Request) (*Market, error) {
	urlStr := "/charge/v1/plans"

	if r.URL.RawQuery != "" {
		urlStr += "?" + r.URL.RawQuery
	}

	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := agent.BaseURL.ResolveReference(rel)

	req, err := agent.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	market := new(Market)

	response := new(RemoteListResponse)
	if err := agent.Do(req, response); err != nil {
		clog.Error(err)
		return nil, err
	}

	plans := []apiPlan{}

	if err := json.Unmarshal([]byte(response.Data), &plans); err != nil {
		clog.Error(err)
		return nil, err
	} else {
		clog.Debug(len(plans), "plans listed.")
		for _, apiplan := range plans {
			plan := convertPlan(&apiplan)
			market.Plans = append(market.Plans, *plan)
		}
	}

	return market, nil
}

func (agent *MarketAgent) ListRegion(r *http.Request) (*[]Region, error) {
	urlStr := "/charge/v1/query/plans/region"

	if r.URL.RawQuery != "" {
		urlStr += "?" + r.URL.RawQuery
	}

	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := agent.BaseURL.ResolveReference(rel)

	req, err := agent.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	response := new(RemoteListResponse)
	if err := agent.Do(req, response); err != nil {
		clog.Error(err)
		return nil, err
	}

	regions := new([]Region)

	if err := json.Unmarshal([]byte(response.Data), regions); err != nil {
		clog.Error(err)
		return nil, err
	}

	return regions, nil
}

func convertPlan(apiplan *apiPlan) *Plan {
	plan := new(Plan)

	plan.PlanId = apiplan.PlanId
	plan.Name = apiplan.Name
	plan.Level = apiplan.Level
	plan.Type = apiplan.PlanType
	plan.Price = apiplan.Price
	plan.BillPeriod = apiplan.Cycle
	plan.RegionID = apiplan.Region
	plan.Region = apiplan.RegionDesc
	plan.Desc = apiplan.Spec1
	plan.Desc2 = apiplan.Spec2
	plan.Detail = apiplan.Description
	plan.CreationTime = apiplan.CreateTime
	plan.Customize = apiplan.Customize

	return plan

}
