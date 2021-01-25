package task

import (
	"context"
	"time"
)

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

func ExecuteJobsTimeout(ctx context.Context, runner TaskRunner, timeout time.Duration) {
    for {
        select {
        case <-ctx.Done():
            return
        default:
            toutctx, _ := context.WithTimeout(ctx, timeout)
            runner.Tick(toutctx)
        }
    }
}

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
