package task

import (
	"context"
	"testing"
)

func TestTaskIsRunning(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    tr := NewContextTaskRunner(ctx)
    i := 0
    increment := func () {
        i++
    }
    await := tr.SubmitTask(NewTask(increment))
    if i == 1 {
        t.Errorf("job started without a executor")
    }
    ScheduleJobs(ctx, tr, 1)
    <-await
    if i != 1 {
        t.Errorf("job didn't run")
    }
}
