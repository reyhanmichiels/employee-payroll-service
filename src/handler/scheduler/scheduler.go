package scheduler

import (
	"context"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/reyhanmichiels/go-pkg/v2/auth"
	"github.com/reyhanmichiels/go-pkg/v2/log"
	"github.com/reyhanmichies/employee-payroll-service/src/business/usecase"
	"github.com/reyhanmichies/employee-payroll-service/src/utils/config"
)

var (
	once = &sync.Once{}
)

type SchedulerTaskConf struct {
	Name             string
	Enabled          bool
	TimeType         string
	Interval         time.Duration
	Weekday          time.Weekday
	ScheduledTime    string
	MultipleSchedule []string
}

type Interface interface {
	Run()
	TriggerScheduler(name string) error
}

type scheduler struct {
	cron     *gocron.Scheduler
	metaconf config.ApplicationMeta
	log      log.Interface
	auth     auth.Interface
	uc       *usecase.Usecases
}

type InitParam struct {
	MetaConf config.ApplicationMeta
	Log      log.Interface
	Auth     auth.Interface
	Uc       *usecase.Usecases
}

func Init(params InitParam) Interface {
	s := &scheduler{}
	once.Do(func() {
		cron := gocron.NewScheduler(time.UTC)
		cron.TagsUnique()

		s = &scheduler{
			cron:     cron,
			metaconf: params.MetaConf,
			log:      params.Log,
			auth:     params.Auth,
			uc:       params.Uc,
		}

		s.assignScheduledTasks()
	})
	return s
}

func (s *scheduler) assignScheduledTasks() {
	// Declare Scheduler Task
	s.AssignTask(
		SchedulerTaskConf{
			Name:          "ValidateAttendancePeriodScheduler",
			Enabled:       true,
			TimeType:      "daily",
			ScheduledTime: "00:01",
		},
		s.uc.AttendancePeriod.ValidateAttendancePeriodScheduler,
	)
}

func (s *scheduler) Run() {
	s.cron.StartAsync()
	s.log.Info(context.Background(), "Scheduler is running")
}

func (s *scheduler) TriggerScheduler(name string) error {
	return s.cron.RunByTag(name)
}
