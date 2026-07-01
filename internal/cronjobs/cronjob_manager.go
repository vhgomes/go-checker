package cronjobs

import (
	"context"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type CronJobManager struct {
	cron    *cron.Cron
	jobs    []CronJobInterface
	context context.Context
}

func NewJobManager(ctx context.Context) *CronJobManager {
	return &CronJobManager{
		cron:    cron.New(cron.WithSeconds()),
		jobs:    []CronJobInterface{},
		context: ctx}
}

type CronJobInterface interface {
	Name() string
	Schedule() string
	Run(ctx context.Context) error
}

func (cm *CronJobManager) RegisterJob(jobInterface CronJobInterface) {
	cm.jobs = append(cm.jobs, jobInterface)
}

func (cm *CronJobManager) StartScheduler() {
	for _, job := range cm.jobs {
		schedule := job.Schedule()
		if _, err := cm.cron.AddFunc(schedule, func() {
			if err := job.Run(cm.context); err != nil {
				zap.L().Error("Error in job", zap.String("job_name", job.Name()), zap.Error(err))
			} else {
				zap.L().Info("Job executed successfully", zap.String("job_name", job.Name()))
			}
		}); err != nil {
			zap.L().Error("Failed to schedule job", zap.String("job_name", job.Name()), zap.Error(err))
		}
	}
	cm.cron.Start()
}

// Por enquanto o codigo está rodando sem redis/rabbitmq então quando o servidor cair ele irá
// resetar todos os jobs, só quando reiniciar o servidor que ele ira voltar.
