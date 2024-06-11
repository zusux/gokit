package notice

import (
	"context"
	"sync"
	"time"
)

type oChan struct {
	sync.Once
	ch chan struct{}
}

func newOChan() *oChan {
	return &oChan{
		ch: make(chan struct{}),
	}
}
func (o *oChan) close() {
	o.Do(func() {
		close(o.ch)
	})
}

type DelayNotice struct {
	sync.Mutex
	data    map[string]interface{}
	keyChan map[string]*oChan
}

func NewDelayNotice() *DelayNotice {
	return &DelayNotice{
		data:    make(map[string]interface{}),
		keyChan: make(map[string]*oChan),
	}
}

func (d *DelayNotice) Put(key string, v interface{}) {
	d.Lock()
	defer d.Unlock()
	d.data[key] = v
	ch, ok := d.keyChan[key]
	if !ok {
		return
	}
	ch.close()
}

func (d *DelayNotice) Get(ctx context.Context, key string, maxWaitTime time.Duration) (interface{}, error) {
	d.Lock()
	v, ok := d.data[key]
	if ok {
		d.Unlock()
		return v, nil
	}
	ch, ok := d.keyChan[key]
	if !ok {
		ch = newOChan()
		d.keyChan[key] = ch
	}
	tCtx, cancel := context.WithTimeout(ctx, maxWaitTime)
	defer cancel()
	d.Unlock()
	select {
	case <-tCtx.Done():
		return nil, tCtx.Err()
	case <-ch.ch:
	}
	d.Lock()
	res := d.data[key]
	d.Unlock()
	return res, nil
}

func (d *DelayNotice) Del(key string) {
	d.Lock()
	defer d.Unlock()
	ch, ok := d.keyChan[key]
	delete(d.data, key)
	if !ok {
		return
	}
	delete(d.keyChan, key)
	ch.close()
}
