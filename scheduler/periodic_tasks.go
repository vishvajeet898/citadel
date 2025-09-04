package scheduler

import (
	"github.com/Orange-Health/citadel/common/constants"
	workerPeriodicTasks "github.com/Orange-Health/citadel/worker_periodic_tasks"
	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/log"
	"github.com/robfig/cron/v3"
)

func HandlePeriodicTasks(server *machinery.Server) {
	// Setup the periodic tasks scheduler
	c := cron.New()

	RegisterStaleTasksPeriodicTask(server, c)
	RegisterEtsTatBreachedPeriodicTask(server, c)
	RegisterReverseLogisticsSlackAlertPeriodicTask(server, c)

	c.Start()
}

func RegisterStaleTasksPeriodicTask(server *machinery.Server, c *cron.Cron) {
	log.INFO.Printf("RegisterStaleTasksPeriodicTask")

	signature := workerPeriodicTasks.StaleTasksPeriodicTaskSignature()
	_, err := c.AddFunc(constants.StaleTasksPeriodicTaskFrequency, func() {
		_, err := server.SendTask(signature)
		if err != nil {
			log.ERROR.Printf("error while sending periodic task: %v", err.Error())
		}
	})

	if err != nil {
		log.ERROR.Printf("error while adding stale tasks periodic task: %v", err.Error())
	}
}

func RegisterEtsTatBreachedPeriodicTask(server *machinery.Server, c *cron.Cron) {
	log.INFO.Printf("RegisterEtsTatBreachedPeriodicTask")

	signature := workerPeriodicTasks.EtsTatBreachedPeriodicTaskSignature()
	_, err := c.AddFunc(constants.EtsTatBreachedPeriodicTaskFrequency, func() {
		log.INFO.Println("Sending EtsTatBreachedPeriodicTask")
		_, err := server.SendTask(signature)
		if err != nil {
			log.ERROR.Printf("error while sending periodic task: %v", err.Error())
		}
	})

	if err != nil {
		log.ERROR.Printf("error while adding ets tat breached periodic task: %v", err.Error())
	}
}

func RegisterReverseLogisticsSlackAlertPeriodicTask(server *machinery.Server, c *cron.Cron) {
	log.INFO.Printf("RegisterReverseLogisticsSlackAlertPeriodicTask")

	signature := workerPeriodicTasks.SlackAlertForDelayedReverseLogisticsPeriodicTaskSignature()
	_, err := c.AddFunc(constants.SlackAlertForReverseLogisticsPeriodicTaskFrequency, func() {
		log.INFO.Println("Sending SlackAlertForReverseLogisticsPeriodicTask")
		_, err := server.SendTask(signature)
		if err != nil {
			log.ERROR.Printf("error while sending periodic task: %v", err.Error())
		}
	})

	if err != nil {
		log.ERROR.Printf("error while adding reverse logistics slack alert periodic task: %v", err.Error())
	}
}
