package shplugexporter

import (
	"fmt"
	"strings"
)

type Plug struct {
	Hostname string
	Device   string
}

func (x Plug) GetBaseUrl() string {
	hostname := strings.TrimPrefix(x.Hostname, "http://")
	return fmt.Sprintf("http://%s/rpc", hostname)
}
