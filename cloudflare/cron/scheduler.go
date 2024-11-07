package cron

import (
	"context"
	"fmt"
	"syscall/js"

	"github.com/syumai/workers"
	"github.com/syumai/workers/internal/jsutil"
	"github.com/syumai/workers/internal/runtimecontext"
)

type Task func(ctx context.Context) error

var scheduledTask Task

func runScheduler(eventObj js.Value) error {
	ctx := runtimecontext.New(context.Background(), eventObj)
	if err := scheduledTask(ctx); err != nil {
		return err
	}
	return nil
}

func init() {
	runSchedulerCallback := js.FuncOf(func(_ js.Value, args []js.Value) any {
		if len(args) != 1 {
			panic(fmt.Errorf("invalid number of arguments given to runScheduler: %d", len(args)))
		}
		eventObj := args[0]
		var cb js.Func
		cb = js.FuncOf(func(_ js.Value, pArgs []js.Value) any {
			defer cb.Release()
			resolve := pArgs[0]
			go func() {
				err := runScheduler(eventObj)
				if err != nil {
					panic(err)
				}
				resolve.Invoke(js.Undefined())
			}()
			return js.Undefined()
		})

		return jsutil.NewPromise(cb)
	})
	jsutil.Binding.Set("runScheduler", runSchedulerCallback)
}

// ScheduleTask sets the Task to be executed
func ScheduleTask(task Task) {
	scheduledTask = task
	workers.Ready()
	WaitForCompletion()
}

// ScheduleTaskNonBlock sets the Task to be executed but does not signal readiness or block
// indefinitely. The non-blocking form is meant to be used in conjunction with [workers.Serve].
func ScheduleTaskNonBlock(task Task) {
	scheduledTask = task
}

// WaitForCompletion blocks until the task set by ScheduleTaskWithNonBlock is completed.
// Currently, this function never returns to support cloudflare.WaitUntil feature.
func WaitForCompletion() {
	select {}
}
