package workers

import (
	"context"
	"errors"
	"fmt"
	"syscall/js"
	"time"

	"github.com/syumai/workers/internal/jsutil"
	"github.com/syumai/workers/internal/runtimecontext"
)

// CronEvent represents information about the Cron that invoked this worker.
type CronEvent struct {
	// Type string
	Cron          string
	ScheduledTime time.Time
}

// toGo converts JS Object to CronEvent
func (ce *CronEvent) toGo(obj js.Value) error {
	if obj.IsUndefined() {
		return errors.New("event is null")
	}
	cronVal := obj.Get("cron").String()
	ce.Cron = cronVal
	scheduledTimeVal := obj.Get("scheduledTime").Float()
	ce.ScheduledTime = time.Unix(int64(scheduledTimeVal)/1000, 0).UTC()

	return nil
}

type CronFunc func(ctx context.Context, event CronEvent) error

var cronTask CronFunc

// ScheduleTask sets the CronFunc to be executed
func ScheduleTask(task CronFunc) {
	cronTask = task
	jsutil.Global.Call("ready")
	select {}
}

func runScheduler(eventObj js.Value, runtimeCtxObj js.Value) error {
	ctx := runtimecontext.New(context.Background(), runtimeCtxObj)
	event := CronEvent{}
	err := event.toGo(eventObj)
	if err != nil {
		return err
	}
	err = cronTask(ctx, event)
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
	jsutil.Global.Set("runScheduler", runSchedulerCallback)
}
