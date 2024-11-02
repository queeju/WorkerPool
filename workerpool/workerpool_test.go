package workerpool

import (
	"testing"
)

func TestWorkerPool_AddWorker(t *testing.T) {
	pool := NewWorkerPool(2)
	defer pool.Stop()

	initialWorkers := len(pool.workers)

	if err := pool.AddWorker(); err != nil {
		t.Fatalf("Expected no error when adding worker, got %v", err)
	}

	if len(pool.workers) != initialWorkers+1 {
		t.Fatalf("Expected %d workers, got %d", initialWorkers+1, len(pool.workers))
	}
}

func TestWorkerPool_RemoveWorker(t *testing.T) {
	pool := NewWorkerPool(3)
	defer pool.Stop()

	if err := pool.RemoveWorker(); err != nil {
		t.Fatalf("Expected no error when removing worker, got %v", err)
	}

	if len(pool.workers) != 2 {
		t.Fatalf("Expected 2 workers, got %d", len(pool.workers))
	}
}

func TestWorkerPool_RemoveWorker_EmptyPool(t *testing.T) {
	pool := NewWorkerPool(1)
	defer pool.Stop()

	if err := pool.RemoveWorker(); err != nil {
		t.Fatalf("Expected no error when removing worker, got %v", err)
	}

	// Попробуем удалить еще один воркер, когда в пуле больше нет воркеров
	if err := pool.RemoveWorker(); err == nil {
		t.Fatal("Expected error when removing a worker from an empty pool, got none")
	}
}

func TestWorkerPool_AddJob_AfterStop(t *testing.T) {
	pool := NewWorkerPool(2)
	pool.Stop()

	err := pool.AddJob("test job")
	if err != errPoolStopped {
		t.Fatalf("Expected error %v, got %v", errPoolStopped, err)
	}
}

func TestWorkerPool_Stop(t *testing.T) {
	pool := NewWorkerPool(3)

	if err := pool.Stop(); err != nil {
		t.Fatalf("Expected no error when stopping pool, got %v", err)
	}

	// Проверяем, что канал jobChan закрыт
	select {
	case _, ok := <-pool.jobChan:
		if ok {
			t.Fatal("Expected jobChan to be closed")
		}
	default:
		// Канал закрыт
	}
}

func TestWorkerPool_AddJob_WhenStopped(t *testing.T) {
	pool := NewWorkerPool(2)
	pool.Stop()

	err := pool.AddJob("test job")
	if err != errPoolStopped {
		t.Fatalf("Expected error %v, got %v", errPoolStopped, err)
	}
}

func TestWorkerPool_RemoveWorker_AfterStop(t *testing.T) {
	pool := NewWorkerPool(3)
	pool.Stop()

	err := pool.RemoveWorker()
	if err != errPoolStopped {
		t.Fatalf("Expected error %v, got %v", errPoolStopped, err)
	}
}

func TestWorkerPool_AddWorker_AfterStop(t *testing.T) {
	pool := NewWorkerPool(2)
	pool.Stop()

	err := pool.AddWorker()
	if err != errPoolStopped {
		t.Fatalf("Expected error %v, got %v", errPoolStopped, err)
	}
}
