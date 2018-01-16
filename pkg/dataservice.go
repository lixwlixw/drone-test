package pkg

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/zonesan/clog"
)

type DataServiceAgent struct {
	*Agent
	BaseURL *url.URL
}

type InstanceCredential struct {
	Uri      string `json:"uri"`
	Hostname string `json:"hostname"`
	Port     string `json:"port"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type DataServiceInfo struct {
	ServiceId    string `json:"service_id"`
	Name         string `json:"service_name"`
	Alias        string `json:"display_name"`
	Class        string `json:"class"`
	Provider     string `json:"provider"`
	InstanceData string `json:"instance_data,omitempty"`
	DataName     string `json:"data_name,omitempty"`
	Desc         string `json:"description"`
	ImageUrl     string `json:"image_url"`
}

func (agent *DataServiceAgent) ListServices(r *http.Request) (*[]DataServiceInfo, error) {

	urlStr := "/integration/v1/services"

	services := []DataServiceInfo{}

	if err := doRequestList(agent, r, "GET", urlStr, nil, &services); err != nil {
		clog.Error(err)
		return nil, err
	} else {
		for idx, service := range services {
			services[idx].Name = service.InstanceData
			services[idx].InstanceData = ""
			services[idx].Alias = service.DataName
			services[idx].DataName = ""

		}
	}

	clog.Infof("%v service(s) listed.", len(services))
	return &services, nil

}

func (agent *DataServiceAgent) CreateServiceInstance(r *http.Request, instance_id string) (*InstanceCredential, error) {

	urlStr := fmt.Sprintf("/integration/v1/instance/%v", instance_id)

	credential := new(InstanceCredential)

	if err := doRequest(agent, r, "POST", urlStr, nil, credential); err != nil {
		clog.Error(err)
		return nil, err
	}

	clog.Infof("service instance created.\n%#v", credential)
	return credential, nil

}

func (agent *DataServiceAgent) Url() *url.URL {
	u := new(url.URL)
	u, _ = url.Parse(agent.BaseURL.String())
	return u
}

func (agent *DataServiceAgent) Instance() *Agent {
	return agent.Agent
}
