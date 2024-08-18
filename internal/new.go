package internal

import (
	"os"
	"strings"
	"test-go/internal/nginx"
	"test-go/internal/xray"

	"github.com/oschwald/maxminddb-golang"
)

func NewXray(key string) Filter {
	return xray.NewXray(key)
}

func NewNginx(key string) Filter {
	return nginx.NewNginx(key)
}

func NewStatistics(logsPath, geolite, geocn string) (*Statistics, error) {
	// 读取日志文件
	logsBytes, err := os.ReadFile(logsPath)
	if err != nil {
		return nil, err
	}

	// 读取IP数据库
	dblite, err := maxminddb.Open(geolite)
	if err != nil {
		return nil, err
	}

	// 读取IP数据库
	dbcn, err := maxminddb.Open(geocn)
	if err != nil {
		return nil, err
	}

	return &Statistics{
		logs: strings.Split(string(logsBytes), "\n"),
		dbs: struct {
			geolite *maxminddb.Reader
			geocn   *maxminddb.Reader
		}{
			geolite: dblite,
			geocn:   dbcn,
		},
	}, nil
}
