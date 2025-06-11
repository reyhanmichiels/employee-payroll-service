package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/reyhanmichiels/go-pkg/v2/appcontext"
	"github.com/reyhanmichiels/go-pkg/v2/codes"
	"github.com/reyhanmichiels/go-pkg/v2/errors"
)

const (
	schedulerUserAgent string = "Cron Scheduler : %s"

	schedulerAssignError string = "Assigning Scheduler %s error: %s"

	schedulerRunning       string = "Scheduler %s is running"
	schedulerDoneError     string = "Scheduler %s is error: %v"
	schedulerDoneSuccess   string = "Scheduler %s is success"
	schedulerTimeExecution string = "Scheduler %s done in %v"

	schedulerTimeTypeExact           string = "daily"
	schedulerTimeTypeInterval        string = "interval"
	schedulerTimeTypeWeekly          string = "weekly"
	schedulerTimeTypeMultipleInDaily string = "multiple"
)

type handlerFunc func(ctx context.Context) error

func (s *scheduler) createContext(conf SchedulerTaskConf) context.Context {
	ctx := context.Background()
	ctx = appcontext.SetUserAgent(ctx, fmt.Sprintf(schedulerUserAgent, conf.Name))
	ctx = appcontext.SetRequestId(ctx, uuid.New().String())
	ctx = appcontext.SetRequestStartTime(ctx, time.Now())
	ctx = appcontext.SetServiceVersion(ctx, s.metaconf.Version)

	return ctx
}

func (s *scheduler) AssignTask(conf SchedulerTaskConf, task handlerFunc) {
	if conf.Enabled {
		var err error
		ctx := context.Background()
		schedulerFunc := s.taskWrapper(conf, task)

		switch conf.TimeType {
		case schedulerTimeTypeInterval:
			_, err = s.cron.Every(conf.Interval).Tag(conf.Name).Do(schedulerFunc)
		case schedulerTimeTypeExact:
			_, err = s.cron.Every(1).Day().Tag(conf.Name).At(conf.ScheduledTime).Do(schedulerFunc)
		case schedulerTimeTypeWeekly:
			_, err = s.cron.Every(1).Weekday(conf.Weekday).Tag(conf.Name).At(conf.ScheduledTime).Do(schedulerFunc)
		case schedulerTimeTypeMultipleInDaily:
			for _, scheduleString := range conf.MultipleSchedule {
				_, err = s.cron.Every(1).Day().Tag(fmt.Sprintf("%s-%s", conf.Name, scheduleString)).At(scheduleString).Do(schedulerFunc)
			}
		default:
			err = errors.NewWithCode(codes.CodeInternalServerError, "Unknown Scheduler Task Time Type")
		}

		if err != nil {
			s.log.Fatal(ctx, fmt.Sprintf(schedulerAssignError, conf.Name, err.Error()))
		}

	}
}

func (s *scheduler) taskWrapper(conf SchedulerTaskConf, task handlerFunc) func() {
	return func() {
		ctx := s.createContext(conf)
		s.log.Info(ctx, fmt.Sprintf(schedulerRunning, conf.Name))
		if err := task(ctx); err != nil {
			s.log.Error(ctx, fmt.Sprintf(schedulerDoneError, conf.Name, err))
		} else {
			s.log.Info(ctx, fmt.Sprintf(schedulerDoneSuccess, conf.Name))
		}

		startTime := appcontext.GetRequestStartTime(ctx)
		s.log.Info(ctx, fmt.Sprintf(schedulerTimeExecution, conf.Name, time.Since(startTime)))
	}
}
