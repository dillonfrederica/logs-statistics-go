package internal

import "github.com/oschwald/maxminddb-golang"

type Statistics struct {
	dbs struct {
		geolite *maxminddb.Reader
		geocn   *maxminddb.Reader
	}
	logs  []string
	cache []struct {
		ip      string
		count   int
		origins []struct {
			key  string
			data []string
		}
	}
}
