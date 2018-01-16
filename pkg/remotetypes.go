package pkg

import (
	"encoding/json"
)

type RequestParams map[string]interface{}

type RemoteResponse struct {
	Code    uint            `json:"code"`
	Message string          `json:"msg"`
	Data    json.RawMessage `json:"data"`
}

type RemoteListResponse struct {
	Code            uint   `json:"code"`
	Msg             string `json:"msg"`
	QueryListResult `json:"data,omitempty"`
}

type QueryListResult struct {
	Total int64           `json:"total"`
	Data  json.RawMessage `json:"results"`
}

type apiPlan struct {
	id          int
	PlanId      string                 `json:"plan_id,omitempty"`
	Name        string                 `json:"plan_name,omitempty"`
	PlanType    string                 `json:"plan_type,omitempty"`
	Level       int                    `json:"plan_level,omitempty"`
	Spec1       string                 `json:"specification1,omitempty"`
	Spec2       string                 `json:"specification2,omitempty"`
	Description string                 `json:"description,omitempty"`
	Price       float32                `json:"price,omitempty"`
	Cycle       string                 `json:"cycle,omitempty"`
	Region      string                 `json:"region,omitempty"`
	RegionDesc  string                 `json:"region_describe,omitempty"`
	CreateTime  string                 `json:"creation_time,omitempty"`
	Status      string                 `json:"status,omitempty"`
	Customize   map[string]interface{} `json:"customize,omitempty"`
}

type apiPurchaseOrder struct {
	Money string `json:"money"`
	Order struct {
		Id              int64  `json:"id,omitempty"`
		Order_id        string `json:"order_id,omitempty"`
		Account_id      string `json:"namespace,omitempty"` // accountId
		Region          string `json:"region,omitempty"`
		Plan_id         string `json:"plan_id,omitempty"`
		Plan_type       string `json:"_,omitempty"`
		Start_time      string `json:"start_time,omitempty"`
		End_time        string `json:"_,omitempty"`        // po
		EndTime         string `json:"end_time,omitempty"` // vo
		Deadline_time   string `json:"deadline,omitempty"`
		Last_consume_id int    `json:"_,omitempty"`
		Status          string `json:"status,omitempty"`
		Creator         string `json:"creator,omitempty"`
	} `json:"order"`
}

type apiBalance struct {
	Namespace string  `json:"namespace"`
	CreateAt  string  `json:"create_at"`
	UpdateAt  string  `json:"update_at"`
	Balance   float64 `json:"balance"`
	Status    string  `json:"state,omitempty"`
}

type apiTransaction struct {
	TransactionId string  `json:"transactionId"`
	Type          string  `json:"type"`
	Amount        float64 `json:"amount"`
	Namespace     string  `json:"namespace"`
	User          string  `json:"user,omitempty"`
	Reason        string  `json:"reason,omitempty"`
	CreateTime    string  `json:"createtime,omitempty"`
	PaymentMethod string  `json:"paymode"`
	Status        string  `json:"status,omitempty"`
	StatusTime    string  `json:"statustime,omitempty"`
}

type apiRechargePayload struct {
	//aihongpay
	AiPayUrl string    `json:"aiurl,omitempty"`
	Method   string    `json:"method,omitempty"`
	Payloads []Payload `json:"payloads,omitempty"`
	Packet   string    `json:"requestpacket,omitempty"`
	// //weixin pay
	// TradeId   string `json:"out_trade_no,omitempty"`
	// QRDataUrl string `json:"code_url,omitempty"`
	// TradeType string `json:"trade_type,omitempty"`
}

type Payload struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type apiRegion struct {
	Id   string `json:"identification"`
	Desc string `json:"region_describe"`
}

type apiCoupon struct {
	SerialNumber string  `json:"serial,omitempty"`
	Code         string  `json:"code,omitempty"`
	Amount       float64 `json:"amount,omitempty"`
	ExpireOn     string  `json:"expire_on,omitempty"`
	Status       string  `json:"status,omitempty"`
	Region       string  `json:"region,omitempty"`
	Namespace    string  `json:"namespace,omitempty"`
}

type apiDataRepo struct {
	RepoID      int           `json:"repoId"`
	Name        string        `json:"repoName"`
	DisplayName string        `json:"chRepoName"`
	Class       string        `json:"class"`
	Label       string        `json:"label"`
	Desc        string        `json:"description"`
	Owner       string        `json:"createUser"`
	ImageUrl    string        `json:"imageUrl"`
	Items       []apiDataItem `json:"items"`
}

type apiDataItem struct {
	ItemID   int               `json:"itemId"`
	Name     string            `json:"itemName"`
	Url      string            `json:"url"`
	UpdateAt string            `json:"updateTime"`
	Sameple  string            `json:"simple"`
	Owner    string            `json:"createUser"`
	Attrs    []apiDataItemAttr `json:"attrs"`
}

type apiDataItemAttr struct {
	Name    string `json:"attrName"`
	Comment string `json:"instruction"`
	Example string `json:"example"`
	Order   int    `json:"orderId"`
}

type UndefinedResp map[string]interface{}

type Resp map[string]interface{}
