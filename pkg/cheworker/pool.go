package cheworker

import (
	"context"
	"fmt"
	"sync"
)

// Job represents a unit of work to be processed by the worker pool.
type Job func(context.Context) error

// Pool represents a worker pool that processes jobs concurrently.
type Pool struct {
	workers      int
	jobQueue     chan Job
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
	errors       chan error
	errorsMu     sync.Mutex
	allErrors    []error
	onError      func(error)
	panicHandler func(interface{})
	shutdownOnce sync.Once
	isShutdown   bool
	shutdownMu   sync.RWMutex
}

// Config holds configuration for the worker pool.
type Config struct {
	// Workers is the number of concurrent workers. Default: 10
	Workers int

	// QueueSize is the size of the job queue buffer. Default: 100
	QueueSize int

	// OnError is called when a job returns an error. Optional.
	OnError func(error)

	// PanicHandler is called when a job panics. If nil, panics are converted to errors.
	PanicHandler func(interface{})
}

// New creates a new worker pool with the given configuration.
func New(config *Config) *Pool {
	if config == nil {
		config = &Config{}
	}

	workers := config.Workers
	if workers <= 0 {
		workers = 10
	}

	queueSize := config.QueueSize
	if queueSize <= 0 {
		queueSize = 100
	}

	ctx, cancel := context.WithCancel(context.Background())

	pool := &Pool{
		workers:      workers,
		jobQueue:     make(chan Job, queueSize),
		ctx:          ctx,
		cancel:       cancel,
		errors:       make(chan error, workers),
		allErrors:    make([]error, 0),
		onError:      config.OnError,
		panicHandler: config.PanicHandler,
	}

	return pool
}

// Start starts the worker pool. Must be called before submitting jobs.
func (p *Pool) Start() {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.worker()
	}

	// Error collector goroutine
	go p.collectErrors()
}

// Submit submits a job to the pool.
// Returns an error if the pool is shutting down or the context is cancelled.
func (p *Pool) Submit(job Job) error {
	p.shutdownMu.RLock()
	defer p.shutdownMu.RUnlock()

	if p.isShutdown {
		return fmt.Errorf("pool is shutting down")
	}

	select {
	case <-p.ctx.Done():
		return fmt.Errorf("pool is shutting down")
	case p.jobQueue <- job:
		return nil
	}
}

// SubmitWithContext submits a job with a specific context.
// The job will receive the provided context instead of the pool's context.
func (p *Pool) SubmitWithContext(ctx context.Context, job Job) error {
	wrappedJob := func(poolCtx context.Context) error {
		// Use the provided context, but also check pool context
		select {
		case <-poolCtx.Done():
			return poolCtx.Err()
		case <-ctx.Done():
			return ctx.Err()
		default:
			return job(ctx)
		}
	}

	return p.Submit(wrappedJob)
}

// Shutdown gracefully shuts down the worker pool.
// It stops accepting new jobs and waits for all submitted jobs to complete.
func (p *Pool) Shutdown() {
	p.shutdownOnce.Do(func() {
		p.shutdownMu.Lock()
		p.isShutdown = true
		p.shutdownMu.Unlock()

		close(p.jobQueue)
		p.wg.Wait()
		close(p.errors)
	})
}

// ShutdownWithContext gracefully shuts down the worker pool with a context.
// Returns an error if the context is cancelled before all jobs complete.
func (p *Pool) ShutdownWithContext(ctx context.Context) error {
	done := make(chan struct{})

	go func() {
		p.Shutdown()
		close(done)
	}()

	select {
	case <-ctx.Done():
		p.cancel() // Cancel all running jobs
		<-done     // Wait for shutdown to complete
		return ctx.Err()
	case <-done:
		return nil
	}
}

// Stop immediately stops the worker pool.
// Running jobs will be cancelled via context cancellation.
func (p *Pool) Stop() {
	p.shutdownOnce.Do(func() {
		p.shutdownMu.Lock()
		p.isShutdown = true
		p.shutdownMu.Unlock()

		p.cancel()
		close(p.jobQueue)
		p.wg.Wait()
		close(p.errors)
	})
}

// Errors returns all errors that occurred during job processing.
func (p *Pool) Errors() []error {
	p.errorsMu.Lock()
	defer p.errorsMu.Unlock()

	errorsCopy := make([]error, len(p.allErrors))
	copy(errorsCopy, p.allErrors)
	return errorsCopy
}

// Wait waits for all submitted jobs to complete and returns all errors.
// New jobs can still be submitted while waiting.
func (p *Pool) Wait() []error {
	// This is a best-effort wait - it doesn't guarantee all jobs are complete
	// because new jobs can still be submitted
	p.Shutdown()
	return p.Errors()
}

// worker is the worker goroutine that processes jobs from the queue.
func (p *Pool) worker() {
	defer p.wg.Done()

	for job := range p.jobQueue {
		p.executeJob(job)
	}
}

// executeJob executes a single job with panic recovery.
func (p *Pool) executeJob(job Job) {
	defer func() {
		if r := recover(); r != nil {
			if p.panicHandler != nil {
				p.panicHandler(r)
			} else {
				err := fmt.Errorf("worker panic: %v", r)
				p.errors <- err
			}
		}
	}()

	if err := job(p.ctx); err != nil {
		p.errors <- err
	}
}

// collectErrors collects errors from the error channel.
func (p *Pool) collectErrors() {
	for err := range p.errors {
		p.errorsMu.Lock()
		p.allErrors = append(p.allErrors, err)
		p.errorsMu.Unlock()

		if p.onError != nil {
			p.onError(err)
		}
	}
}

// WorkerCount returns the number of workers in the pool.
func (p *Pool) WorkerCount() int {
	return p.workers
}

// QueueSize returns the size of the job queue buffer.
func (p *Pool) QueueSize() int {
	return cap(p.jobQueue)
}

// PendingJobs returns the number of jobs waiting in the queue.
func (p *Pool) PendingJobs() int {
	return len(p.jobQueue)
}
