package pkg

import (
	"reflect"
	"testing"
)

func expect(t *testing.T, before, after, result interface{}) {
	if after != result {
		t.Errorf("Expected %v (type %v)---> %v (type %v) - Got %v (type %v)",
			before, reflect.TypeOf(before),
			after, reflect.TypeOf(after),
			result, reflect.TypeOf(result))
	} else {
		t.Logf("Test ok! %v ---> %v", before, result)
	}
}

type addrArray struct {
	before, after string
}

var addrs []addrArray = []addrArray{
	{"abc.com", "http://abc.com"},
	{"HTTp://abc.com", "HTTp://abc.com"},
	{"HTTP://ABC.COM", "HTTP://ABC.COM"},
	{"https://abc.com", "https://abc.com"},
	{"ftp://abc.com", "http://ftp://abc.com"},
}

func TestHttpAddr(t *testing.T) {

	for _, addr := range addrs {
		expect(t, addr.before, addr.after, httpAddr(addr.before))
	}

}
