package nginx

import (
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
		reg: regexp.MustCompile(`((?:\d{1,3}\.){3}\d{1,3}).*\[(.*)\s[+-]\d{4}\]\s['"](\w*?)\s((` + key + `){1,})(\?(.*)){0,1}?\s(.*?)['"]\s(\d+)\s\d+\s['"].*?['"]\s['"](.*?)['"]`),
	}
}

func (ng *Nginx) Filter(origin string) (string, string) {
	if origin != "" {
		subMatchResults := ng.reg.FindStringSubmatch(origin)

		// for i, v := range subMatchResults {
		// 	log.Println(i, v)
		// }
		// os.Exit(0)
		if len(subMatchResults) == 11 {
			return subMatchResults[1], subMatchResults[9] + " " + subMatchResults[3] + " " + subMatchResults[4]
		}
	}

	return "", ""
}
