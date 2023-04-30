package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/syumai/workers"
	"github.com/syumai/workers/cloudflare"
)

func task(ctx context.Context, event workers.CronEvent) error {
	fmt.Println(cloudflare.Getenv(ctx, "HELLO"))

	if event.ScheduledTime.Minute()%2 == 0 {
		return errors.New("even numbers cause errors")
	}

	return nil
}

func main() {
	workers.ScheduleTask(task)
}
