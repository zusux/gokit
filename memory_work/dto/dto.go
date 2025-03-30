package dto

import (
	"sync"
	"time"
)

type EventKey string

const (
	//EventKeyTest 水机数据
	EventKeyTest EventKey = "test"
)

type MarketEvent struct {
	Event     EventKey `json:"event"`
	Key       string   `json:"key"`
	ExtraData any
}

type EventProcessor struct {
	sync.Mutex
	DataChan   chan *MarketEvent
	Key        string
	Opened     bool
	CreateTime time.Time
	// 任务处理状态
	LastProcessTime time.Time
	ProcessedCount  uint64
	FailedCount     uint64
	LastError       error
	// 性能指标
	AverageProcessTime time.Duration
	MaxProcessTime     time.Duration
	MinProcessTime     time.Duration
}
