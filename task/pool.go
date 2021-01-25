package task

import "context"

type poolContextValue int

const (
    // PoolInstance Acess this Pool instance inside a Task by its context
    PoolInstance = poolContextValue(1 << iota)
)

// Pool Represents a generic task runner
type Pool interface {
    SubmitTask(fn Task) chan error
    Tick(ctx context.Context)
}

// IngestTasksToPool bridge tasks from a channel of tasks to a Pool, runs detached until the channel is closed
func IngestTasksToPool(tasks chan Task, pool Pool) {
    go func() {
        for task := range tasks {
            pool.SubmitTask(task)
        }
    }()
}

type contextPool struct {
    ctx context.Context
    tasks chan(func(ctx context.Context))
}

func (tr *contextPool) SubmitTask(fn Task) chan error {
    cb := make(chan error, 1)
    taskFn := func(ctx context.Context) {
        cb <- fn(context.WithValue(ctx, PoolInstance, tr))
        defer close(cb)
    }
    tr.tasks <-taskFn
    return cb
}

func (tr *contextPool) Tick(ctx context.Context) {
    done := make(chan struct{}, 1)
    mergedContext, cancel := context.WithCancel(ctx)
    go func () {
        select {
        case <-done:
            return
        case <-tr.ctx.Done():
            cancel()
        }
    }()
    task := <-tr.tasks
    task(mergedContext)
    done <-struct{}{}
    close(done)
}

// NewContextPool creates a pool that lives until the context is not cancelled
func NewContextPool(ctx context.Context) Pool {
    return &contextPool{
        ctx: ctx,
        tasks: make(chan func(context.Context), 64),
    }
}
