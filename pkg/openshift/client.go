package openshift

import (
	"fmt"
	"net/http"
	"time"

	apierrors "github.com/asiainfoLDP/datafoundry_payment/pkg/errors"
	projectapi "github.com/openshift/origin/pkg/project/api/v1"
	rolebindingapi "github.com/openshift/origin/pkg/rolebinding/api/v1"
	"github.com/zonesan/clog"
	kapi "k8s.io/kubernetes/pkg/api/v1"
)

type OClient struct {
	client *OpenshiftREST
	user   string
}

const (
	JoinedTimePrefix = "joinedTime/"
)

func NewOClient(host, token, username string) *OClient {

	clog.Debugf("%v:(%v)@%v", username, token, host)

	client := NewOpenshiftREST(NewOpenshiftTokenClient(host, token))
	return &OClient{client: client, user: username}
}

func NewAdminOClient(adminClient *OpenshiftClient) *OClient {
	adminRESTClient := NewOpenshiftREST(adminClient)
	client := &OClient{client: adminRESTClient}
	return client
}

func (oc *OClient) CreateProject(r *http.Request, name string) (*projectapi.Project, error) {

	uri := "/projectrequests"

	projRequest := new(projectapi.ProjectRequest)
	{
		projRequest.DisplayName = name
		projRequest.Name = oc.user + "-org-" + genRandomName(8)
		// projRequest.Annotations = make(map[string]string)
		// projRequest.Annotations["datafoundry.io/requester"] = oc.user
	}

	proj := new(projectapi.Project)

	oc.client.OPost(uri, projRequest, proj)
	if oc.client.Err != nil {
		clog.Error(oc.client.Err)
		return nil, oc.client.Err
	}

	return proj, nil
}

func (oc *OClient) ListRoles(r *http.Request, project string) (*rolebindingapi.RoleBindingList, error) {
	uri := fmt.Sprintf("/namespaces/%v/rolebindings", project)

	roles := new(rolebindingapi.RoleBindingList)

	oc.client.OGet(uri, roles)

	if oc.client.Err != nil {
		clog.Error(oc.client.Err)
		return nil, oc.client.Err
	}
	//clog.Debug(roles)

	rolesResult := new(rolebindingapi.RoleBindingList)

	for _, role := range roles.Items {
		if role.Name == "view" || role.Name == "admin" || role.Name == "edit" {
			rolesResult.Items = append(rolesResult.Items, role)
		} else {
			for _, subject := range role.Subjects {
				if subject.Kind == "User" {
					if role.RoleRef.Name == "view" || role.RoleRef.Name == "admin" ||
						role.RoleRef.Name == "edit" {
						//clog.Debugf("%#v", role)
						rolesResult.Items = append(rolesResult.Items, role)
						break
					}
				}
			}
		}
	}
	return rolesResult, nil
}

func (oc *OClient) GetProject(r *http.Request, project string) (*projectapi.Project, error) {
	uri := fmt.Sprintf("/projects/%v/", project)

	proj := new(projectapi.Project)

	oc.client.OGet(uri, proj)

	if oc.client.Err != nil {
		clog.Error(oc.client.Err)
		return nil, oc.client.Err
	}

	return proj, nil
}

func (oc *OClient) RoleAdd(r *http.Request, project, name string, admin bool) (*rolebindingapi.RoleBinding, error) {

	if name == "" || project == "" {
		return nil, apierrors.ErrorNew(apierrors.ErrCodeInvalidParam)
	}

	uri := fmt.Sprintf("/namespaces/%v/rolebindings", project)

	roleList, err := oc.ListRoles(r, project)
	if err != nil {
		clog.Error(err)
		return nil, err
	}

	if exist := findUserInRoles(roleList, name); exist != nil {
		clog.Warnf("duplicate user: %v, role: %v, project: %v", name, exist.RoleRef.Name, project)
		return nil, apierrors.ErrorNew(apierrors.ErrCodeConflict)
	}

	roleRef := "edit"
	if admin {
		roleRef = "admin"
	}

	role := findRole(roleList, roleRef)
	create := false

	if role == nil { //post else put
		clog.Infof("role '%v' not exist in project '%v', will be created.", roleRef, project)

		create = true
		role = new(rolebindingapi.RoleBinding)
		role.Name = roleRef
		role.RoleRef.Name = roleRef
	}

	subject := kapi.ObjectReference{Kind: "User", Name: name}
	role.Subjects = append(role.Subjects, subject)
	role.UserNames = append(role.UserNames, name)

	if role.Annotations == nil {
		role.Annotations = make(map[string]string)
	}
	role.Annotations[JoinedTimePrefix+name] = time.Now().Format(time.RFC3339)

	if create {
		oc.client.OPost(uri, role, role)
	} else {
		uri += "/" + roleRef
		oc.client.OPut(uri, role, role)
	}

	return role, oc.client.Err
}

func (oc *OClient) RoleRemove(r *http.Request, project, name string) error {
	if name == "" || project == "" {
		return apierrors.ErrorNew(apierrors.ErrCodeInvalidParam)
	}

	if name == oc.user {
		return apierrors.ErrorNew(apierrors.ErrCodeActionNotSupport)
	}

	uri := fmt.Sprintf("/namespaces/%v/rolebindings", project)

	roleList, err := oc.ListRoles(r, project)
	if err != nil {
		clog.Error(err)
		return err
	}

	role := findUserInRoles(roleList, name)
	if role == nil {
		clog.Errorf("can't find user '%v' from roles in project '%v'", name, project)
		return apierrors.ErrorNew(apierrors.ErrCodeUserNotFound)
	} else {
		role = removeUserInRole(role, name)
		uri += "/" + role.Name
		oc.client.OPut(uri, role, role)
	}

	if oc.client.Err != nil {
		clog.Error(oc.client.Err)
	}

	return oc.client.Err
}

func findRole(roles *rolebindingapi.RoleBindingList, roleRef string) *rolebindingapi.RoleBinding {
	for _, role := range roles.Items {
		if role.Name == roleRef {
			return &role
		}
	}
	return nil
}

// func findUserInRole(users []string, user string) bool {
// 	for _, v := range users {
// 		if user == v {
// 			return true
// 		}
// 	}
// 	return false
// }

func findUserInRoles(roles *rolebindingapi.RoleBindingList, username string) *rolebindingapi.RoleBinding {
	for _, role := range roles.Items {
		// if ok := findUserInRole(role.UserNames, username); ok {
		// 	return &role
		// }
		for _, v := range role.UserNames {
			if username == v {
				return &role
			}
		}
	}
	return nil
}

func removeUserInRole(role *rolebindingapi.RoleBinding, user string) *rolebindingapi.RoleBinding {
	for idx, userName := range role.UserNames {
		if userName == user {
			role.UserNames = append(role.UserNames[:idx], role.UserNames[idx+1:]...)
		}
	}
	for idx, subject := range role.Subjects {
		if subject.Name == user {
			role.Subjects = append(role.Subjects[:idx], role.Subjects[idx+1:]...)
		}
	}

	delete(role.Annotations, JoinedTimePrefix+user)

	return role
}
