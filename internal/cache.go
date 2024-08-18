package internal

import (
	"context"
	"slices"
	"sync"
)

func (stats *Statistics) handlerCache(ks struct {
	ip     string
	key    string
	origin string
}) {
	idx := slices.IndexFunc(stats.cache, func(item struct {
		ip      string
		count   int
		origins []struct {
			key  string
			data []string
		}
	}) bool {
		return item.ip == ks.ip
	})
	if idx <= -1 {
		// 创建新纪录
		stats.cache = append(stats.cache, struct {
			ip      string
			count   int
			origins []struct {
				key  string
				data []string
			}
		}{
			origins: []struct {
				key  string
				data []string
			}{{key: ks.key, data: []string{ks.origin}}},
			ip:    ks.ip,
			count: 1,
		})
		return
	}

	// 更新纪录
	stats.cache[idx].count++

	idx2 := slices.IndexFunc(stats.cache[idx].origins, func(item struct {
		key  string
		data []string
	}) bool {
		return item.key == ks.key
	})
	if idx2 <= -1 {
		stats.cache[idx].origins = append(stats.cache[idx].origins, struct {
			key  string
			data []string
		}{key: ks.key, data: []string{ks.origin}})
		return
	}

	// 更新纪录
	stats.cache[idx].origins[idx2].data = append(stats.cache[idx].origins[idx2].data, ks.origin)

}

func (stats *Statistics) cacheThread(ctx context.Context, wg *sync.WaitGroup, ch chan struct {
	ip     string
	key    string
	origin string
}) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return

		case data := <-ch:
			stats.handlerCache(data)
		}
	}
}
