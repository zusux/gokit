package main

import (
	"context"
	"fmt"
	"github.com/zusux/gokit/notice"
	"time"
)

func example() {
	now := time.Now()
	d := notice.NewDelayNotice()
	go func() {
		time.Sleep(time.Duration(10) * time.Second)
		d.Put("222", 33)
	}()
	v, e := d.Get(context.Background(), "222", time.Duration(20)*time.Second)
	fmt.Println(v, e, time.Since(now))
}
