package openshift

type Orgnazition struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	CreateBy    string         `json:"create_by"`
	CreateTime  string         `json:"creation_time"`
	MemberList  []OrgMember    `json:"members"`
	Status      OrgStatusPhase `json:"status"`
	RoleBinding bool           `json:"rolebinding"`
	Reason      string         `json:"reason,omitempty"`
}

type OrgnazitionList struct {
	Orgnazitions []Orgnazition `json:"orgnazitions"`
}

type OrgMember struct {
	Name         string            `json:"member_name"`
	IsAdmin      bool              `json:"admin"`
	PrivilegedBy string            `json:"privileged_by"`
	JoinedAt     string            `json:"joined_at"`
	Status       MemberStatusPhase `json:"status"`
}

type MemberStatusPhase string

const (
	OrgMemberStatusInvited MemberStatusPhase = "invited"
	OrgMemberStatusjoined  MemberStatusPhase = "joined"
	OrgMemberStatusNone    MemberStatusPhase = "none"
)

type OrgStatusPhase string

const (
	OrgStatusCreated OrgStatusPhase = "created"
	OrgStatusPending OrgStatusPhase = "creating"
	OrgStatusError   OrgStatusPhase = "failed"
)
