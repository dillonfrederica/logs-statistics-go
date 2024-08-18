package internal

import (
	"context"
	"runtime"
	"sort"
	"sync"
)

func (stats *Statistics) Statistics(origin string, fer Filter) error {
	var (
		max = runtime.NumCPU() * 2
		ch  = make(chan struct {
			ip     string
			key    string
			origin string
		}, max*2)
		wg          sync.WaitGroup
		ctx, cancel = context.WithCancel(context.Background())
	)

	wg.Add(1)
	go stats.cacheThread(ctx, &wg, ch)

	wg.Add(1)
	go stats.filterThread(cancel, &wg, ch, fer, max)

	wg.Wait()

	// 最后排序
	sort.Slice(stats.cache, func(i, j int) bool {
		return stats.cache[i].count > stats.cache[j].count
	})

	return nil
}
