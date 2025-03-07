package tasks

import (
	"log"

	"github.com/go-co-op/gocron/v2"
)

var Scheduler gocron.Scheduler

func Init() {
	var err error

	Scheduler, err := gocron.NewScheduler()
	if err != nil {
		log.Println("Error scheduling tasks")
	}

	Scheduler.Start()
}

func ShutdownTasks() {
	err := Scheduler.Shutdown()
	if err != nil {
		log.Println("Error stopping task scheduler")
	}
}
