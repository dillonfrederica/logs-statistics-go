package internal

import (
	"fmt"
	"net"
	"os"
	"slices"
	"sort"
	"strings"
	"test-go/internal/nginx"
	"test-go/internal/xrayaccess"

	"github.com/oschwald/maxminddb-golang"
)

type Filter interface {
	Filter(origin string) (string, string)
}

type (
	filterOrigin struct {
		key  string
		data []string
	}

	filterCache struct {
		ip      string
		count   int
		origins []*filterOrigin
	}

	dbs struct {
		geolite *maxminddb.Reader
		geocn   *maxminddb.Reader
	}

	Statistics struct {
		dbs        dbs
		logsString string
		cache      []*filterCache
	}
)

func NewXrayAccess(key string) Filter {
	return xrayaccess.NewXrayAccess(key)
}

func NewNginx(key string) Filter {
	return nginx.NewNginx(key)
}

func Load(logsPath, geolite, geocn string) (*Statistics, error) {
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
		logsString: string(logsBytes),
		dbs: dbs{
			geolite: dblite,
			geocn:   dbcn,
		},
	}, nil
}

func (f *Statistics) Statistics(origin string, fer Filter) error {
	lines := strings.Split(f.logsString, "\n")
	for _, line := range lines {
		if line != "" {
			ip, key := fer.Filter(line)
			if ip != "" {
				idx := slices.IndexFunc(f.cache, func(item *filterCache) bool { return item.ip == ip })
				if idx <= -1 {
					f.cache = append(f.cache, &filterCache{
						origins: []*filterOrigin{{key: key, data: []string{line}}},
						ip:      ip,
						count:   1,
					})
				} else {
					f.cache[idx].count++

					idx2 := slices.IndexFunc(f.cache[idx].origins, func(item *filterOrigin) bool { return item.key == key })
					if idx2 <= -1 {
						f.cache[idx].origins = append(f.cache[idx].origins, &filterOrigin{key: key, data: []string{line}})
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

type DBCn struct {
	City          string `maxminddb:"city"`
	CityCode      int    `maxminddb:"cityCode"`
	Districts     string `maxminddb:"districts"`
	DistrictsCode int    `maxminddb:"districtsCode"`
	ISP           string `maxminddb:"isp"`
	Net           string `maxminddb:"net"`
	Province      string `maxminddb:"province"`
	ProvinceCode  int    `maxminddb:"provinceCode"`
}

type DBGlobal struct {
	City struct {
		GeonameID int               `maxminddb:"geoname_id"`
		Names     map[string]string `maxminddb:"names"`
	} `maxminddb:"city"`
	Continent struct {
		Code      string            `maxminddb:"code"`
		GeonameID int               `maxminddb:"geoname_id"`
		Names     map[string]string `maxminddb:"names"`
	} `maxminddb:"continent"`
	Country struct {
		GeonameID int               `maxminddb:"geoname_id"`
		IsoCode   string            `maxminddb:"iso_code"`
		Names     map[string]string `maxminddb:"names"`
	} `maxminddb:"country"`
	Location struct {
		AccuracyRadius int     `maxminddb:"accuracy_radius"`
		Latitude       float64 `maxminddb:"latitude"`
		Longitude      float64 `maxminddb:"longitude"`
		TimeZone       string  `maxminddb:"time_zone"`
	} `maxminddb:"location"`
	Postal struct {
		Code string `maxminddb:"code"`
	} `maxminddb:"postal"`
	RegisteredCountry struct {
		GeonameID int               `maxminddb:"geoname_id"`
		IsoCode   string            `maxminddb:"iso_code"`
		Names     map[string]string `maxminddb:"names"`
	} `maxminddb:"registered_country"`
	Subdivisions []struct {
		GeonameID int               `maxminddb:"geoname_id"`
		IsoCode   string            `maxminddb:"iso_code"`
		Names     map[string]string `maxminddb:"names"`
	} `maxminddb:"subdivisions"`
}

func (f *Statistics) findIP(ip net.IP) (string, error) {
	var r1 DBGlobal
	if err := f.dbs.geolite.Lookup(ip, &r1); err != nil {
		return "", err
	}

	if r1.Country.IsoCode != "CN" {
		var address string
		if cn, ok := r1.Country.Names["zh-CN"]; ok {
			address = cn
		} else if en, ok := r1.Country.Names["en"]; ok {
			address = en
		} else {
			address = r1.Country.IsoCode
		}

		if len(r1.Subdivisions) > 0 {
			for _, v := range r1.Subdivisions {
				if cn, ok := v.Names["zh-CN"]; ok {
					address = fmt.Sprintf("%s %s", address, cn)
				} else if en, ok := v.Names["en"]; ok {
					address = fmt.Sprintf("%s %s", address, en)
				}
			}
		}

		if cn, ok := r1.City.Names["zh-CN"]; ok {
			address = fmt.Sprintf("%s %s", address, cn)
		} else if en, ok := r1.City.Names["en"]; ok {
			address = fmt.Sprintf("%s %s", address, en)
		}

		return address, nil
	}

	var r2 DBCn
	if err := f.dbs.geocn.Lookup(ip, &r2); err != nil {
		return "", err
	}

	var address string
	if r2.Province != "" {
		address = r2.Province
	}
	if r2.City != "" {
		address = fmt.Sprintf("%s %s", address, r2.City)
	}
	if r2.Districts != "" {
		address = fmt.Sprintf("%s %s", address, r2.Districts)
	}
	if r2.ISP != "" {
		address = fmt.Sprintf("%s %s", address, r2.ISP)
	}
	return address, nil
}

func (f *Statistics) Print() {
	f.Sort()

	for i, cache := range f.cache {
		address, err := f.findIP(net.ParseIP(cache.ip))
		if err != nil {
			fmt.Printf("错误：%v\n", err)
			continue
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
