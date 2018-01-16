FROM golang:1.7.1

MAINTAINER Zonesan <chaizs@asiainfo.com>

ENV TIME_ZONE=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TIME_ZONE /etc/localtime && echo $TIME_ZONE > /etc/timezone

COPY . /go/src/github.com/asiainfoLDP/datafoundry_payment

WORKDIR /go/src/github.com/asiainfoLDP/datafoundry_payment

EXPOSE 8080

RUN go build

CMD ["sh", "-c", "./datafoundry_payment"]
