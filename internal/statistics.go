package internal

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"sort"
	"strings"
	"test-go/internal/nginx"
	"test-go/internal/xrayaccess"

	"github.com/xiaoqidun/qqwry"
)

type Filter interface {
	Filter(origin string) (string, string)
}

type (
	FilterOrigin struct {
		key  string
		data []string
	}

	FilterCache struct {
		ip      string
		count   int
		origins []*FilterOrigin
	}

	Statistics struct {
		logsPath   string
		ipPath     string
		logsString string
		cache      []*FilterCache
	}
)

func NewXrayAccess(key string) Filter {
	return xrayaccess.NewXrayAccess(key)
}

func NewNginx(key string) Filter {
	return nginx.NewNginx(key)
}

func Load(logsPath, ipPath string) (*Statistics, error) {
	// 读取日志文件
	logsBytes, err := os.ReadFile(logsPath)
	if err != nil {
		return nil, err
	}

	// 读取IP数据库
	if err := qqwry.LoadFile(ipPath); err != nil {
		return nil, err
	}

	return &Statistics{
		logsString: string(logsBytes),
		logsPath:   logsPath,
		ipPath:     ipPath,
	}, nil
}

func (f *Statistics) Statistics(origin string, fer Filter) error {
	lines := strings.Split(f.logsString, "\n")
	for _, line := range lines {
		if line != "" {
			ip, key := fer.Filter(line)
			if ip != "" {
				idx := slices.IndexFunc(f.cache, func(item *FilterCache) bool { return item.ip == ip })
				if idx <= -1 {
					f.cache = append(f.cache, &FilterCache{
						origins: []*FilterOrigin{{key: key, data: []string{line}}},
						ip:      ip,
						count:   1,
					})
				} else {
					f.cache[idx].count++

					idx2 := slices.IndexFunc(f.cache[idx].origins, func(item *FilterOrigin) bool { return item.key == key })
					if idx2 <= -1 {
						f.cache[idx].origins = append(f.cache[idx].origins, &FilterOrigin{key: key, data: []string{line}})
					} else {
						f.cache[idx].origins[idx2].data = append(f.cache[idx].origins[idx2].data, line)
					}
				}
			}
		}
	}

	return nil
}

func (f *Statistics) Sort() {
	sort.Slice(f.cache, func(i, j int) bool {
		return f.cache[i].count > f.cache[j].count
	})
}

// https://raw.githubusercontent.com/FW27623/qqwry/main/qqwry.dat
func DownloadQqwry(output string) error {
	resp, err := http.Get("https://raw.githubusercontent.com/FW27623/qqwry/main/qqwry.dat")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	dataBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := os.WriteFile(output, dataBytes, 0644); err != nil {
		return err
	}

	return nil
}

func (f *Statistics) Print() {
	f.Sort()

	for i, cache := range f.cache {
		// 从内存或缓存查询IP
		location, err := qqwry.QueryIP(cache.ip)
		if err != nil {
			fmt.Printf("错误：%v\n", err)
			continue
		}

		address := ""
		if location.Country != "" {
			address = location.Country
		}
		if location.Province != "" {
			address = fmt.Sprintf("%s %s", address, location.Province)
		}
		if location.City != "" {
			address = fmt.Sprintf("%s %s", address, location.City)
		}
		if location.District != "" {
			address = fmt.Sprintf("%s %s", address, location.District)
		}
		if location.ISP != "" {
			address = fmt.Sprintf("%s %s", address, location.ISP)
		}

		fmt.Printf("No.%d\n  IP: %s 共访问%d次\n  物理地址: %s\n",
			i+1,
			cache.ip,
			cache.count,
			address,
		)

		// cache.origins 排序
		sort.Slice(cache.origins, func(i, j int) bool { return len(cache.origins[j].data) < len(cache.origins[i].data) })

		for _, ov := range cache.origins {
			fmt.Printf("    目标地址: %s 共访问%d次\n", ov.key, len(ov.data))
		}
		fmt.Printf("\n")
	}
}
