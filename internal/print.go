package internal

import (
	"fmt"
	"net"
	"sort"
)

func (stats *Statistics) Print() {
	for i, cache := range stats.cache {
		address, err := stats.findIP(net.ParseIP(cache.ip))
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
