package main

import (
	"flag"
	"fmt"
	"os"
	"test-go/internal"
)

var (
	logs      string
	geolite   string
	geocn     string
	searchKey string
	app       string
	method    string
	kind      string
)

func init() {
	app = os.Args[0]
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}
	method = os.Args[1]
	os.Args = append([]string{os.Args[0]}, os.Args[2:]...)

	flag.StringVar(&logs, "logs", "", "日志")
	flag.StringVar(&geolite, "geolite", "./GeoLite2-City.mmdb", "IP数据库")
	flag.StringVar(&geocn, "geocn", "./GeoCN.mmdb", "IP数据库")
	flag.StringVar(&searchKey, "search", "", "搜索正则表达式")
	flag.StringVar(&kind, "kind", "", "类型")
	flag.Parse()
}

func printStatisticsHelp() {
	fmt.Printf("Usage: %s statistics --kind [help|xray-access|nginx] --logs [日志] --search [正则表达式]\n", app)
	fmt.Printf("Example(nginx): %s statistics --kind nginx --logs access.log --search \"\\/\"\n", app)
	fmt.Printf("Example(xray-access): %s statistics --kind xray-access --logs access.log --search \"[^\\.]*\\.baidu\\.com\"\n", app)
}

func statistics() {
	if logs == "" || kind == "" {
		printStatisticsHelp()
		os.Exit(0)
	}

	f, err := internal.Load(logs, geolite, geocn)
	if err != nil {
		panic(err)
	}

	switch kind {
	case "xray-access":
		if err := f.Statistics(searchKey, internal.NewXrayAccess(searchKey)); err != nil {
			panic(err)
		}
	case "nginx":
		if err := f.Statistics(searchKey, internal.NewNginx(searchKey)); err != nil {
			panic(err)
		}
	case "help":
		printStatisticsHelp()
		os.Exit(0)
	default:
		printStatisticsHelp()
		os.Exit(1)
	}

	f.Print()
}

func download() {
	if err := internal.DownloadGeoLite(geolite); err != nil {
		panic(err)
	}
	if err := internal.DownloadGeoCN(geocn); err != nil {
		panic(err)
	}
}

func printHelp() {
	fmt.Printf("Usage: %s [help|statistics|download]\n", app)
}

func main() {
	switch method {
	case "statistics":
		statistics()
	case "download":
		download()
	case "help":
		printHelp()
		os.Exit(0)
	default:
		printHelp()
		os.Exit(1)
	}
}
