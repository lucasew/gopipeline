package task

import (
	"context"
	"time"
)

// ExecuteJobs runs the tasks of a TaskRunner until the context is cancelled
func ExecuteJobs(ctx context.Context, runner TaskRunner) {
    for {
        select {
        case <-ctx.Done():
            return
        default:
            runner.Tick(ctx)
        }
    }
}

// ExecuteJobsTimeout runs the task of a TaskRunner until the context is cancelled, each task have a constant timeout
func ExecuteJobsTimeout(ctx context.Context, runner TaskRunner, timeout time.Duration) {
    for {
        select {
        case <-ctx.Done():
            return
        default:
            toutctx, cancel := context.WithTimeout(ctx, timeout)
            runner.Tick(toutctx)
            cancel()
        }
    }
}

// ScheduleJobs spawn ExecuteJobs loops using replicas goroutines until the context is cancelled
func ScheduleJobs(ctx context.Context, runner TaskRunner, replicas int) {
    spawn := func (ctx context.Context) {
        begin:
        select {
        case <-ctx.Done():
            return
        default:
            ExecuteJobs(ctx, runner)
            goto begin
        }
    }
    for i := 0; i < replicas; i++ {
        go spawn(ctx)
    }
}
