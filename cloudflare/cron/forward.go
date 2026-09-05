// Deprecated: use github.com/syumai/workers-go/cloudflare/cron instead.
package cron

import (
	src "github.com/syumai/workers-go/cloudflare/cron"
)

var Done = src.Done

type Event = src.Event

var NewEvent = src.NewEvent
var ScheduleTask = src.ScheduleTask
var ScheduleTaskNonBlock = src.ScheduleTaskNonBlock

type Task = src.Task
