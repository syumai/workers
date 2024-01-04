package cron

import (
	"context"
	"errors"
	"fmt"
	"syscall/js"
	"time"

	"github.com/syumai/workers/internal/jsutil"
	"github.com/syumai/workers/internal/runtimecontext"
)

// Event represents information about the Cron that invoked this worker.
type Event struct {
	Cron          string
	ScheduledTime time.Time
}

// toEvent converts JS Object to Go Event struct
func toEvent(obj js.Value) (*Event, error) {
	if obj.IsUndefined() {
		return nil, errors.New("event is null")
	}
	cronVal := obj.Get("cron").String()
	scheduledTimeVal := obj.Get("scheduledTime").Float()
	return &Event{
		Cron:          cronVal,
		ScheduledTime: time.Unix(int64(scheduledTimeVal)/1000, 0).UTC(),
	}, nil
}

type Task func(ctx context.Context, event *Event) error

var scheduledTask Task

// ScheduleTask sets the Task to be executed
func ScheduleTask(task Task) {
	scheduledTask = task
	js.Global().Call("ready")
	select {}
}

func runScheduler(eventObj js.Value, runtimeCtxObj js.Value) error {
	ctx := runtimecontext.New(context.Background(), runtimeCtxObj)
	event, err := toEvent(eventObj)
	if err != nil {
		return err
	}
	err = scheduledTask(ctx, event)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	runSchedulerCallback := js.FuncOf(func(_ js.Value, args []js.Value) any {
		if len(args) != 2 {
			panic(fmt.Errorf("invalid number of arguments given to runScheduler: %d", len(args)))
		}
		event := args[0]
		runtimeCtx := args[1]

		var cb js.Func
		cb = js.FuncOf(func(_ js.Value, pArgs []js.Value) any {
			defer cb.Release()
			resolve := pArgs[0]
			go func() {
				err := runScheduler(event, runtimeCtx)
				if err != nil {
					panic(err)
				}
				resolve.Invoke(js.Undefined())
			}()
			return js.Undefined()
		})

		return jsutil.NewPromise(cb)
	})
	js.Global().Set("runScheduler", runSchedulerCallback)
}
