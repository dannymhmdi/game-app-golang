package scheduler

import (
	"context"
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"mymodule/entity"
	"mymodule/params"
	"mymodule/service/matchmakingService"
	"time"
)

type Scheduler struct {
	matchMakingSvc matchmakingService.Service
	scheduler      gocron.Scheduler
}

func New(matchmakingSvc matchmakingService.Service) *Scheduler {
	scheduler, nErr := gocron.NewScheduler()
	if nErr != nil {
		fmt.Println("failed to init scheduler", nErr)

		return nil
	}

	return &Scheduler{
		matchMakingSvc: matchmakingSvc,
		scheduler:      scheduler,
	}
}

func (s Scheduler) Start(ch <-chan bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	_, err := s.scheduler.NewJob(
		gocron.DurationJob(
			3*time.Second,
		),
		gocron.NewTask(
			func() {
				res, err := s.matchMakingSvc.MatchMaking(ctx, []entity.Category{entity.SoccorCategory, entity.HistoryCategory}, params.MatchMakingRequest{})
				if err != nil {
					fmt.Println("error in matchmaking service in scheduler", err)

					return
				}

				fmt.Println("matchmaking called successfully", res.MatchedUsers)
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
