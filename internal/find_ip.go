package internal

import (
	"fmt"
	"net"
)

func (stats *Statistics) findIP(ip net.IP) (string, error) {
	var r1 struct {
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
	if err := stats.dbs.geolite.Lookup(ip, &r1); err != nil {
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

	var r2 struct {
		City          string `maxminddb:"city"`
		CityCode      int    `maxminddb:"cityCode"`
		Districts     string `maxminddb:"districts"`
		DistrictsCode int    `maxminddb:"districtsCode"`
		ISP           string `maxminddb:"isp"`
		Net           string `maxminddb:"net"`
		Province      string `maxminddb:"province"`
		ProvinceCode  int    `maxminddb:"provinceCode"`
	}
	if err := stats.dbs.geocn.Lookup(ip, &r2); err != nil {
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
