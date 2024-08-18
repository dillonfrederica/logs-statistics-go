package xray

import (
	"log"
	"regexp"
)

type XrayAccess struct {
	key string
	reg *regexp.Regexp
}

const IPREG = `(?:(?:\d{1,3}\.){3}\d{1,3})|(?:(?:[a-f0-9]{1,4}(?:\:[a-f0-9]{1,4}){7})|(?:[a-f0-9]{1,4}(?:\:[a-f0-9]{1,4}){0,7}::[a-f0-9]{0,4}(?:\:[a-f0-9]{1,4}){0,7}))`

func NewXray(key string) *XrayAccess {
	if key == "" {
		key = `[^:]+`
	}
	return &XrayAccess{
		key: key,
		reg: regexp.MustCompile(`(?<ip>` + IPREG + `).*accepted\s(\w+):(?<domain>` + key + `):(\d+)\s`),
	}
}

func (xa *XrayAccess) Filter(origin string) (string, string) {
	if origin != "" {
		subMatchResults := xa.reg.FindStringSubmatch(origin)
		groupNames := xa.reg.SubexpNames()
		if len(subMatchResults) == 0 || len(groupNames) == 0 {
			// log.Printf("未匹配到结果: %v\n", origin)
			return "", ""
		}
		if len(subMatchResults) != len(groupNames) {
			log.Fatalf("匹配异常: %v\n", origin)
		}

		var (
			ip     string
			domain string
		)
		for i, v := range groupNames {
			switch v {
			case "ip":
				ip = subMatchResults[i]
			case "domain":
				domain = subMatchResults[i]
			}
		}
		return ip, domain
	}

	return "", ""
}
