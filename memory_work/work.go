package memory_work

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/zusux/gokit/memory_work/dto"

	"github.com/go-kratos/kratos/v2/log"
)

// WorkerConfig 定义 Worker 的配置项
type WorkerConfig struct {
	CollectBufferSize        int           // 收集数据通道的缓冲区大小
	EventProcessorBufferSize int           // 每个事件处理器的数据通道缓冲区大小
	CleanupInterval          time.Duration // 清理检查的时间间隔
	MaxInactiveTime          time.Duration // 最大不活跃时间
	MaxProcessingTime        time.Duration // 最大处理时间
	MaxErrorRate             float64       // 最大错误率阈值
	CircuitBreakerWindow     time.Duration // 熔断器窗口期
}

// DefaultConfig 返回默认配置
func DefaultConfig() WorkerConfig {
	return WorkerConfig{
		CollectBufferSize:        30000,
		EventProcessorBufferSize: 3000,
		CleanupInterval:          time.Hour * 24,
		MaxInactiveTime:          time.Hour * 24 * 31,
		MaxProcessingTime:        time.Second * 5,
		MaxErrorRate:             0.1, // 10% 错误率阈值
		CircuitBreakerWindow:     time.Minute * 5,
	}
}

type Worker struct {
	sync.Mutex
	config      WorkerConfig
	collectData chan *dto.MarketEvent
	processors  sync.Map
	// 监控指标
	processedEvents uint64
	failedEvents    uint64
	// 熔断器状态
	circuitBreaker     bool
	lastCircuitTripped time.Time
	// 优雅关闭支持
	shutdownChan chan struct{}
	waitGroup    sync.WaitGroup
}

func NewWorker(config ...WorkerConfig) *Worker {
	cfg := DefaultConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	return &Worker{
		config:       cfg,
		collectData:  make(chan *dto.MarketEvent, cfg.CollectBufferSize),
		processors:   sync.Map{},
		shutdownChan: make(chan struct{}),
	}
}

func (w *Worker) Start(ctx context.Context) {
	log.Context(ctx).Infof("Worker starting...")
	for {
		select {
		case <-ctx.Done():
			log.Context(ctx).Infof("Worker received context cancellation")
			w.Stop()
			return
		case <-w.shutdownChan:
			log.Context(ctx).Infof("Worker received shutdown signal")
			return
		case data, ok := <-w.collectData:
			if !ok {
				continue
			}
			w.handleEvent(ctx, data)
		}
	}
}

func (w *Worker) handleEvent(ctx context.Context, data *dto.MarketEvent) {
	defer func() {
		if r := recover(); r != nil {
			log.Context(ctx).Errorf("Panic recovered in handleEvent: %v", r)
			atomic.AddUint64(&w.failedEvents, 1)
		}
	}()

	// 检查熔断器状态
	if w.circuitBreaker {
		if time.Since(w.lastCircuitTripped) < w.config.CircuitBreakerWindow {
			log.Context(ctx).Warnf("Circuit breaker is open, dropping event for key: %s", data.Key)
			return
		}
		// 重置熔断器
		w.circuitBreaker = false
	}

	// 检查错误率
	processed := atomic.LoadUint64(&w.processedEvents)
	failed := atomic.LoadUint64(&w.failedEvents)
	if processed > 1000 && float64(failed)/float64(processed) > w.config.MaxErrorRate {
		w.circuitBreaker = true
		w.lastCircuitTripped = time.Now()
		log.Context(ctx).Errorf("Circuit breaker tripped due to high error rate: %.2f%%", float64(failed)/float64(processed)*100)
		return
	}

	processorObj, ok := w.processors.Load(data.Key)
	if ok {
		processor, _ := processorObj.(*dto.EventProcessor)
		processor.Lock()
		if processor.Opened {
			select {
			case processor.DataChan <- data:
				atomic.AddUint64(&w.processedEvents, 1)
			default:
				log.Context(ctx).Warnf("Channel full for key %s, creating new handler", data.Key)
				w.createEventHandler(ctx, data)
			}
		} else {
			w.createEventHandler(ctx, data)
		}
		processor.Unlock()
	} else {
		w.createEventHandler(ctx, data)
	}
}

func (w *Worker) createEventHandler(ctx context.Context, data *dto.MarketEvent) {
	w.Lock()
	defer w.Unlock()
	processor := &dto.EventProcessor{
		CreateTime: time.Now(),
		DataChan:   make(chan *dto.MarketEvent, w.config.EventProcessorBufferSize),
		Key:        data.Key,
		Opened:     true,
	}
	w.processors.Store(data.Key, processor)
	processor.DataChan <- data
	go w.processEvents(ctx, processor)
}

func (w *Worker) processEvents(ctx context.Context, processor *dto.EventProcessor) {
	updateTime := time.Now()
	tick := time.NewTicker(w.config.CleanupInterval)
	defer tick.Stop()

	w.waitGroup.Add(1)
	defer w.waitGroup.Done()

	for {
		select {
		case <-ctx.Done():
			log.Context(ctx).Infof("Worker matchHandle exit for key: %s", processor.Key)
			w.closeEventProcessor(processor)
			return
		case <-tick.C:
			if time.Since(processor.CreateTime) > w.config.MaxInactiveTime && time.Since(updateTime) > w.config.CleanupInterval {
				log.Context(ctx).Infof("Closing inactive worker for key: %s", processor.Key)
				w.closeEventProcessor(processor)
				return
			}
		case data, ok := <-processor.DataChan:
			if !ok {
				return
			}

			start := time.Now()
			var err error

			// 检查处理时间
			processDone := make(chan struct{})
			go func() {
				defer close(processDone)
				switch data.Event {
				// 在这里处理具体的事件类型
				}
			}()

			// 设置处理超时
			select {
			case <-processDone:
				// 处理完成
			case <-time.After(w.config.MaxProcessingTime):
				err = fmt.Errorf("processing timeout after %v", w.config.MaxProcessingTime)

			}

			// 更新处理状态
			processor.Lock()
			processTime := time.Since(start)
			processor.LastProcessTime = time.Now()
			if err != nil {
				atomic.AddUint64(&processor.FailedCount, 1)
				processor.LastError = err
				log.Context(ctx).Errorf("Error processing event for key %s: %v", processor.Key, err)
			} else {
				atomic.AddUint64(&processor.ProcessedCount, 1)
			}

			// 更新性能指标
			if processor.MaxProcessTime < processTime || processor.MaxProcessTime == 0 {
				processor.MaxProcessTime = processTime
			}
			if processor.MinProcessTime > processTime || processor.MinProcessTime == 0 {
				processor.MinProcessTime = processTime
			}
			// 使用指数移动平均计算平均处理时间
			const alpha = 0.1 // 平滑因子
			if processor.AverageProcessTime == 0 {
				processor.AverageProcessTime = processTime
			} else {
				processor.AverageProcessTime = time.Duration(float64(processor.AverageProcessTime)*(1-alpha) + float64(processTime)*alpha)
			}
			processor.Unlock()

			updateTime = time.Now()
		}
	}
}

func (w *Worker) closeEventProcessor(processor *dto.EventProcessor) {
	// 首先检查输入参数
	if processor == nil {
		log.Errorf("Attempting to close nil EventProcessor")
		return
	}

	// 获取 worker 锁以确保安全访问 groupChan
	w.Lock()
	tmpProcessorObj, ok := w.processors.LoadAndDelete(processor.Key)
	w.Unlock() // 尽早释放锁

	if !ok {
		log.Warnf("WorkerBet with key %s not found in groupChan", processor.Key)
		return
	}

	eventProcessor, ok := tmpProcessorObj.(*dto.EventProcessor)
	if !ok {
		log.Errorf("Invalid EventProcessor type for key %s", processor.Key)
		return
	}

	// 获取 workerBet 锁以确保安全访问
	eventProcessor.Lock()
	defer eventProcessor.Unlock()

	// 如果已经关闭，直接返回
	if !eventProcessor.Opened {
		log.Warnf("EventProcessor %s is already closed", processor.Key)
		return
	}

	// 标记为关闭并关闭通道
	eventProcessor.Opened = false
	close(eventProcessor.DataChan)

	// 记录关闭信息和处理统计
	log.Infof("EventProcessor %s closed. Processed: %d, Failed: %d, Avg Process Time: %v",
		processor.Key,
		atomic.LoadUint64(&eventProcessor.ProcessedCount),
		atomic.LoadUint64(&eventProcessor.FailedCount),
		eventProcessor.AverageProcessTime)
}
func (w *Worker) Stop() {
	log.Infof("Worker stopping, processed events: %d, failed events: %d",
		atomic.LoadUint64(&w.processedEvents),
		atomic.LoadUint64(&w.failedEvents))

	close(w.shutdownChan)
	close(w.collectData)

	// 等待所有处理器完成
	w.waitGroup.Wait()

	// 关闭所有处理器
	w.processors.Range(func(key, value interface{}) bool {
		if processor, ok := value.(*dto.EventProcessor); ok {
			w.closeEventProcessor(processor)
		}
		return true
	})
}

func (w *Worker) GetMetrics() (processed, failed uint64) {
	return atomic.LoadUint64(&w.processedEvents), atomic.LoadUint64(&w.failedEvents)
}

// GetDetailedMetrics 返回详细的监控指标
func (w *Worker) GetDetailedMetrics() map[string]interface{} {
	metrics := make(map[string]interface{})

	// 基础指标
	metrics["processed_events"] = atomic.LoadUint64(&w.processedEvents)
	metrics["failed_events"] = atomic.LoadUint64(&w.failedEvents)
	metrics["error_rate"] = float64(atomic.LoadUint64(&w.failedEvents)) / float64(atomic.LoadUint64(&w.processedEvents)+1) * 100

	// 熔断器状态
	metrics["circuit_breaker_status"] = w.circuitBreaker
	if w.circuitBreaker {
		metrics["circuit_breaker_tripped_time"] = w.lastCircuitTripped
	}

	// 队列状态
	metrics["collect_channel_length"] = len(w.collectData)
	metrics["collect_channel_capacity"] = cap(w.collectData)

	// 处理器状态
	activeHandlers := 0
	totalQueueLength := 0
	w.processors.Range(func(_, value interface{}) bool {
		if processor, ok := value.(*dto.EventProcessor); ok && processor.Opened {
			activeHandlers++
			totalQueueLength += len(processor.DataChan)
		}
		return true
	})
	metrics["active_handlers"] = activeHandlers
	metrics["total_queue_length"] = totalQueueLength

	return metrics
}
