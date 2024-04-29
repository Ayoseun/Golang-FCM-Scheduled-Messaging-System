
package main

import (
	app "firebase-fcm-cron-job/app/services"
	"fmt"

	"github.com/robfig/cron"
)

func main() {

	fmt.Println("Started Scheduler")
	c := cron.New()
	c.AddFunc("*/15 * * * *", app.FetchScheduledMessages)
	c.Start()

	select {}
}
