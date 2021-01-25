package task

import "context"

type TaskRunnerContextValue int

const (
    TaskRunnerInstance = 1 << iota
)

type TaskRunner interface {
    SubmitTask(fn Task) chan error
    Tick(ctx context.Context)
}

func IngestTasksToTaskRunner(tasks chan Task, runner TaskRunner) {
    go func() {
        for task := range tasks {
            runner.SubmitTask(task)
        }
    }()
}

type contextTaskRunner struct {
    ctx context.Context
    tasks chan(func(ctx context.Context))
}

func (tr *contextTaskRunner) SubmitTask(fn Task) chan error {
    cb := make(chan error, 1)
    taskFn := func(ctx context.Context) {
        cb <- fn(context.WithValue(ctx, TaskRunnerInstance, tr))
        defer close(cb)
    }
    tr.tasks <-taskFn
    return cb
}

func (tr *contextTaskRunner) Tick(ctx context.Context) {
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

func NewContextTaskRunner(ctx context.Context) TaskRunner {
    return &contextTaskRunner{
        ctx: ctx,
        tasks: make(chan func(context.Context), 64),
    }
}
