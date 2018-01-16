package openshift

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func httpsAddr(addr string) string {
	if strings.HasSuffix(addr, "/") {
		addr = strings.TrimRight(addr, "/")
	}

	if !strings.HasPrefix(addr, "https://") {
		return fmt.Sprintf("https://%s", addr)
	}

	return addr
}

func genRandomName(strlen int) (name string) {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}
