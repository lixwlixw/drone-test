GET /payment/v1/account?namespace=hello-org&region=cn-north-1

```json
{
  "purchased": false,
  "notification": false,
  "plans": [
    {
      "plan_id": "2736123def7232",
      "type": "normal",
      "price": 10,
      "bill_period": "monthly",
      "description": "1 CPU Core, 512M Memory"
    }
  ],
  "status": "",
  "balance": {
    "balance": 0
  }
}
```

GET /payment/v1/market?region=cn-north-1&type=c

```json
{
  "plans": [
    {
      "plan_id": "1d3452ea-7f14-11e6-9fe0-2344dd5557c3",
      "type": "C",
      "price": 20,
      "level":"S",
      "name":"asdwe",
      "bill_period": "M",
      "description": "1 CPU Core, 512M Memory",
      "creation_time": "2016-09-26T17:27:04+08:00"
    },
    {
      "plan_id": "1d3452ea-7f14-11e6-9fe0-2344dd5557c3",
      "type": "C",
      "price": 40.88,
      "bill_period": "M",
      "description": "2 CPU Cores, 1G Memory",
      "creation_time": "2016-09-26T17:27:04+08:00"
    },
    {
      "plan_id": "1d3452ea-7f14-11e6-9fe0-2344dd5557c3",
      "type": "C",
      "price": 88.88,
      "bill_period": "M",
      "description": "4 CPU Cores, 2G Memory",
      "creation_time": "2016-09-26T17:27:04+08:00"
    }
  ]
}
```

POST /payment/v1/checkout -d '{"plan_id":"zwwqe","namespace":"chaizs","region":"cn-north-1"}'

```json
{
  "balance": 3000.01,
  "status": "active"
}
```
POST /payment/v1/recharge -d '{"amount":1234.34,namespace:"chaizs"}'

```json
{
  "balance": 6000.89,
  "status": "active"
}
```

GET /payment/v1/amounts?namespace=hello

```json
[
  {
    "trans_id": "03Fwerqe2",
    "creation_time": "2016-09-26T17:30:03+08:00",
    "amount": 12.23,
    "description": "Plan A",
    "payment_method": "balance",
    "status": "finish"
  },
  {
    "trans_id": "qwer238DKJ",
    "creation_time": "2016-09-26T17:30:03+08:00",
    "amount": 12.23,
    "description": "Plan A",
    "payment_method": "balance",
    "status": "finish"
  },
  {
    "trans_id": "03F232X238DKJ",
    "creation_time": "2016-09-26T17:30:03+08:00",
    "amount": 2.34,
    "description": "Plan A",
    "payment_method": "balance",
    "status": "refunded"
  }
]
```
GET /payment/v1/amounts/:tid?namespace=hello

```json
{
  "trans_id": "03F232X238DKJ",
  "creation_time": "2016-09-26T17:30:04+08:00",
  "amount": 12.23,
  "description": "Plan A",
  "payment_method": "balance",
  "status": "finish"
}
```
GET /payment/v1/balance?namespace=hello

```json
{
  "balance": 50000.89,
  "status": "active"
}
```
GET /payment/v1/regions

```json
[
  {
    "identification": "cn-north-1",
    "region_describe": "铸造一区"
  },
  {
    "identification": "cn-north-2",
    "region_describe": "铸造二区"
  }
]
```

GET /payment/v1/coupon/:couponcode

```json
{
  "serial": "xxeefd",
  "amount": 10,
  "expire_on": "2017-01-01T00:00:00Z",
  "status": "available"
}
```

POST /payment/v1/redeem -d  '{"serial":"xxeefd","code":"ssa","namespace":"chaizs","region":"cn-north-1"}'

```json
{
  "amount": 10,
  "namespace": "chaizs",
  "region": "cn-north-1"
}
```

GET /integration/v1/repos

```json
[
  {
    "repo_name": "repo_aaa",
    "class": "Class A",
    "label": "Label 1",
    "description": "description of this repo.",
    "image_url": "www.example.com"
  },
  {
    "repo_name": "repo_bbb",
    "class": "Class B",
    "label": "Label 2",
    "description": "description of this repo",
    "image_url": "www.example.com"
  }
]
```


GET /integration/v1/repos/:repo

```json
{
  "repo_name": "repo_aaa",
  "description": "description of repo_aaa",
  "owner": "owner info of repo_aaa",
  "items": [
    {
      "item_name": "item_1",
      "url": "www.example.com"
    },
    {
      "item_name": "item_2",
      "url": "www.example.com"
    }
  ]
}
```

GET /integration/v1/repos/:repo/items/:item

```json
{
  "item_name": "item_1",
  "url": "www.example.com",
  "update_at": "2016-11-04T07:45:59Z",
  "owner": "own info of item.",
  "attrs": [
    {
      "attr_name": "attr1",
      "comment": "comment of attr1",
      "example": "a",
      "order": 1
    },
    {
      "attr_name": "attr2",
      "comment": "comment of attr2",
      "example": "2",
      "order": 2
    },
    {
      "attr_name": "attr3",
      "comment": "comment of attr3",
      "example": "ok",
      "order": 3
    }
  ]
}
```

GET /integration/v1/services

```json
[
  {
    "service_id": "c13d4-3123-11e6-8ffb-0323c75f5c",
    "class": "数据库",
    "provider": "Asiainfo",
    "service_name": "Meteorological",
    "description": "国气象科学数据共享服务网是提供气象资料共享的公益性网站，由一个主节点和分布在国家级和省级气象部门的若干个分节点网站组成。国家气象信息中心负责对中国气象科学数据共享服务网的建设和管理。",
    "image_url": "www.example.com"
  }
]
```

POST /integration/v1/instance/:instance_id

```json
{
  "uri": "mongo://username:password@mysqlhost:3306/database",
  "hostname": "mysqlhost",
  "port": "3306",
  "name": "database",
  "username": "username",
  "password": "password"
}
```


POST /payment/v1/coupon?region=cn-north-1  -d '{"amount":100,"expire_on":30,"kind":"recharge"}'

```json
{
  "serial": "DF075754236988704R",
  "code": "WLLK720DVPYVUD3B",
  "amount": 100,
  "expire_on": "2017-01-26"
}
```

POST /payment/v1/directrecharge?region=cn-north-2 -d '{"namespace":"xxx","amount":1000,"reason":" byadmin"}'

```json
{
  "namespace": "xxx",
  "create_at": "",
  "update_at": "",
  "balance": 5207.5
}
```

  ** identification  jd

  ** region cn-north-1

  ** desc taocan 1

  ** comment

  ** status available
