package internal

import (
	"context"
	"sync"
)

func (stats *Statistics) handlerFilter(wg *sync.WaitGroup, c chan struct{}, ch chan struct {
	ip     string
	key    string
	origin string
}, fer Filter, line string) {
	defer func() {
		<-c
		wg.Done()
	}()

	ip, key := fer.Filter(line)
	if ip == "" {
		return
	}

	ch <- struct {
		ip     string
		key    string
		origin string
	}{
		ip:     ip,
		key:    key,
		origin: line,
	}

}

func (stats *Statistics) filterThread(cancel context.CancelFunc, wg *sync.WaitGroup, ch chan struct {
	ip     string
	key    string
	origin string
}, fer Filter, max int) {
	defer wg.Done()
	defer cancel()

	var (
		c = make(chan struct{}, max-2)
		w sync.WaitGroup
	)

	for _, line := range stats.logs {
		// 过滤空行
		if line != "" {
			c <- struct{}{}
			w.Add(1)
			go stats.handlerFilter(&w, c, ch, fer, line)
		}
	}

	w.Wait()
}
