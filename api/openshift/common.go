package openshift

import (
	"fmt"
	"net/http"
	"os"
	"time"

	apierrors "github.com/asiainfoLDP/datafoundry_payment/pkg/errors"
	"github.com/asiainfoLDP/datafoundry_payment/pkg/openshift"
	userapi "github.com/openshift/origin/pkg/user/api/v1"
	"github.com/zonesan/clog"
)

const (
	Region01 = "cn-north-1"
	Region02 = "cn-north-2"

	TotalRegionsCnt = 2

	RegionDefault = Region01
)

var (
	regions      map[string]string
	adminClients map[string]*openshift.OpenshiftClient
)

func init_regions() {
	if regions == nil {
		regions = make(map[string]string, TotalRegionsCnt)
	}

	regions[Region01] = os.Getenv("APISERVER_CN_NORTH_01")
	regions[Region02] = os.Getenv("APISERVER_CN_NORTH_02")

	if len(regions[Region01]) == 0 || len(regions[Region02]) == 0 {
		clog.Fatal("env 'APISERVER_CN_NORTH_01' and 'APISERVER_CN_NORTH_02' must be specified.")
	}

	clog.Infof("%#v", regions)

}

func init_authinfo() (user, pass string) {

	auth := make(map[string]string)

	auth["user"] = os.Getenv("DATAFOUNDRY_ADMIN_USER")
	auth["pass"] = os.Getenv("DATAFOUNDRY_ADMIN_PASS")

	if len(auth["user"]) == 0 || len(auth["pass"]) == 0 {
		clog.Fatal("env 'DATAFOUNDRY_ADMIN_USER' and 'DATAFOUNDRY_ADMIN_PASS' must be specified.")
	}

	clog.Infof("%#v", auth)
	return auth["user"], auth["pass"]
}

func init_admin() {
	var durPhase time.Duration
	phaseStep := time.Hour / TotalRegionsCnt

	init_regions()

	adminuser, adminpassword := init_authinfo()

	adminClients = make(map[string]*openshift.OpenshiftClient, TotalRegionsCnt)

	adminClients[Region01] = openshift.NewOpenshiftClient(Region01, regions[Region01], adminuser, adminpassword, durPhase)
	durPhase += phaseStep
	adminClients[Region02] = openshift.NewOpenshiftClient(Region02, regions[Region02], adminuser, adminpassword, durPhase)

}

func Init() {
	init_admin()
}

func authDF(region, bearertoken string) (*userapi.User, error) {

	ocRestClient, err := NewRestClient(region, bearertoken)
	if err != nil {
		clog.Error(err)
		return nil, err
	}

	u := &userapi.User{}
	uri := "/users/~"
	ocRestClient.OGet(uri, u)
	if ocRestClient.Err != nil {
		clog.Errorf("authDF, region(%s), uri(%s) error: %s", region, uri, ocRestClient.Err)
		//Logger.Infof("authDF, region(%s), token(%s), uri(%s) error: %s", region, userToken, uri, osRest.Err)
		return nil, ocRestClient.Err
	}

	return u, nil
}

func dfUser(user *userapi.User) string {
	return user.Name
}

func getDFUserame(region, token string) (string, error) {

	user, err := authDF(region, token)
	if err != nil {
		return "", err
	}
	return dfUser(user), nil
}

func NewRestClient(region, bearertoken string) (*openshift.OpenshiftREST, error) {
	if regions[region] == "" {
		return nil, fmt.Errorf("user noud found @ region (%s).", region)
	}

	client := openshift.NewOpenshiftREST(openshift.NewOpenshiftTokenClient(regions[region], bearertoken))
	return client, nil
}

func RegionHostname(region string) string {
	return regions[region]
}

func NewClient(r *http.Request, validation bool) (*openshift.OClient, error) {
	r.ParseForm()

	region := r.FormValue("region")
	host := RegionHostname(region)
	token := r.Header.Get("Authorization")

	clog.Debugf("[%v] [%v] [%v]", region, host, token)

	if token == "" {
		return nil, apierrors.ErrorNew(apierrors.ErrCodeUnauthorized)
	}

	if host == "" {
		return nil, apierrors.ErrorNew(apierrors.ErrCodeRegionNotFound)
	}

	user := ""
	var err error = nil
	if validation {
		user, err = getDFUserame(region, token)

		if err != nil {
			return nil, err
		}
	}
	client := openshift.NewOClient(host, token, user)
	return client, nil
}

func NewAdminClient(r *http.Request, project string) (*openshift.OClient, error) {
	r.ParseForm()

	if project == "" {
		return nil, apierrors.ErrorNew(apierrors.ErrCodeInvalidParam)
	}

	region := r.FormValue("region")
	host := RegionHostname(region)
	token := r.Header.Get("Authorization")

	clog.Debugf("[%v] [%v] [%v]", region, host, token)

	if token == "" {
		return nil, apierrors.ErrorNew(apierrors.ErrCodeUnauthorized)
	}

	if host == "" {
		return nil, apierrors.ErrorNew(apierrors.ErrCodeRegionNotFound)
	}

	user, err := getDFUserame(region, token)

	if err != nil {
		return nil, err
	}

	if err := func() error {
		oc, err := NewClient(r, false)
		if err != nil {
			return err
		}
		_, err = oc.GetProject(r, project)

		return err

	}(); err != nil {
		clog.Error(err)
		return nil, err
	}

	clog.Debug(user, "is reuqesting admin permission.")

	if adminClients[region] == nil || adminClients[region].BearerToken() == "" {
		return nil, apierrors.ErrorNew(apierrors.ErrCodeAdminNotPresented)
	}

	client := openshift.NewAdminOClient(adminClients[region])

	return client, nil
}
