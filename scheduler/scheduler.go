package scheduler

import (
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"mymodule/service/matchmakingService"
	"time"
)

type Scheduler struct {
	matchMakingSvc matchmakingService.Service
	scheduler      gocron.Scheduler
}

func New(matchmakingSvc matchmakingService.Service) *Scheduler {
	s, nErr := gocron.NewScheduler()
	if nErr != nil {
		fmt.Println("failed to init scheduler", nErr)

		return nil
	}

	return &Scheduler{
		matchMakingSvc: matchmakingSvc,
		scheduler:      s,
	}
}

func (s Scheduler) Start(ch <-chan bool) {
	_, err := s.scheduler.NewJob(
		gocron.DurationJob(
			3*time.Second,
		),
		gocron.NewTask(
			func() {
				if err := s.matchMakingSvc.MatchMaking(); err != nil {
					fmt.Println("error in service call scheduler", err)

					return
				}

			},
		),
	)
	if err != nil {
		fmt.Println("failed to start job", err)

		return
	}

	s.scheduler.Start()

	<-ch

	if sErr := s.scheduler.Shutdown(); sErr != nil {
		fmt.Println("failed to shutdown scheduler", sErr)
	}

	time.Sleep(2 * time.Second)
	fmt.Println("scheduler shutdown gracefully")
}
