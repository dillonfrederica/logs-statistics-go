package nginx

import (
	"log"
	"regexp"
)

type Nginx struct {
	key string
	reg *regexp.Regexp
}

func NewNginx(key string) *Nginx {
	if key == "" {
		key = `\/[^\/\?]*`
	}

	return &Nginx{
		key: key,
		reg: regexp.MustCompile(`(?<ip>(?:\d{1,3}\.){3}\d{1,3}).*\[(?<time>.*\s[+-]\d{4})\]\s['"](?<method>\w*?)\s(?<path>(?:` + key + `){1,})(?<params>\?.*){0,1}?\s(?<protocol>.*?)['"]\s(?<status_code>\d+)\s\d+\s['"].*?['"]\s['"](?<user_agent>.*?)['"]`),
	}
}

func (ng *Nginx) Filter(origin string) (string, string) {
	if origin != "" {
		subMatchResults := ng.reg.FindStringSubmatch(origin)
		groupNames := ng.reg.SubexpNames()
		if len(subMatchResults) == 0 || len(groupNames) == 0 {
			// log.Printf("未匹配到结果: %v\n", origin)
			return "", ""
		}

		if len(subMatchResults) != len(groupNames) {
			log.Fatalf("匹配异常: %v\n", origin)
		}

		var (
			ip         string
			statusCode string
			path       string
			method     string
		)
		for i, v := range groupNames {
			switch v {
			case "ip":
				ip = subMatchResults[i]
			case "status_code":
				statusCode = subMatchResults[i]
			case "path":
				path = subMatchResults[i]
			case "method":
				method = subMatchResults[i]
			}
		}

		return ip, statusCode + " " + method + " " + path
	}

	return "", ""
}
