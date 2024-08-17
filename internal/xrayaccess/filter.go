package xrayaccess

import (
	"regexp"
)

type XrayAccess struct {
	key string
	reg *regexp.Regexp
}

func NewXrayAccess(key string) *XrayAccess {
	if key == "" {
		key = `[^:]+`
	}
	return &XrayAccess{
		key: key,
		reg: regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}).*accepted\s(\w+):(` + key + `):(\d+)\s`),
	}
}

func (xa *XrayAccess) Filter(origin string) (string, string) {
	if origin != "" {
		subMatchResults := xa.reg.FindStringSubmatch(origin)

		if len(subMatchResults) == 5 {
			return subMatchResults[1], subMatchResults[3]
		}
	}

	return "", ""
}
