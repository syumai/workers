package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/syumai/workers/cloudflare"
	"github.com/syumai/workers/cloudflare/cron"
)

func task(ctx context.Context, event cron.Event) error {
	fmt.Println(cloudflare.Getenv(ctx, "HELLO"))

	if event.ScheduledTime.Minute()%2 == 0 {
		return errors.New("even numbers cause errors")
	}

	return nil
}

func main() {
	cron.ScheduleTask(task)
}
