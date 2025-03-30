package memory_work

import (
	"sync"
)

var workerManagerOnce sync.Once
var workerManager *WorkManager

type WorkManager struct {
	workers sync.Map
}

func GetWorkManager() *WorkManager {
	workerManagerOnce.Do(func() {
		workerManager = &WorkManager{}
	})
	return workerManager
}

func (w *WorkManager) AddWorker(key string, worker *Worker) {
	w.workers.Store(key, worker)
}

func (w *WorkManager) GetWorker(key string) (*Worker, bool) {
	worker, exists := w.workers.Load(key)
	if !exists {
		return nil, false
	}
	return worker.(*Worker), true
}

func (w *WorkManager) RemoveWorker(key string) {
	worker, exists := w.workers.LoadAndDelete(key)
	if w, ok := worker.(*Worker); exists && ok {
		w.Stop()
	}
}

func (w *WorkManager) GetAllWorkers() map[string]*Worker {
	result := make(map[string]*Worker)
	w.workers.Range(func(key, value interface{}) bool {
		result[key.(string)] = value.(*Worker)
		return true
	})
	return result
}
