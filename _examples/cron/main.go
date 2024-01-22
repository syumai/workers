package main

import (
	"context"
	"fmt"

	"github.com/syumai/workers/cloudflare/cron"
)

func task(ctx context.Context) error {
	e, err := cron.NewEvent(ctx)
	if err != nil {
		return err
	}

	fmt.Println(e.ScheduledTime.Unix())

	return nil
}

func main() {
	cron.ScheduleTask(task)
}
