package cronjobs

import (
	"context"
	"encoding/json"
	"fmt"
	"go-checker/internal/repository"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type DashboardCronJob struct {
	name     string
	schedule string
	siteRepo *repository.SiteRepo
	userRepo *repository.UserRepo
	userID   uint
	redis    *redis.Client
}

func NewDashboardCronJob(siteRepo *repository.SiteRepo, userRepo *repository.UserRepo, redisClient *redis.Client, schedule string) *DashboardCronJob {
	return &DashboardCronJob{
		name:     "DashboardCronJob",
		schedule: schedule,
		siteRepo: siteRepo,
		userRepo: userRepo,
		redis:    redisClient,
	}
}

func (d DashboardCronJob) Name() string {
	return d.name
}

func (d DashboardCronJob) Schedule() string {
	return d.schedule
}

func (d DashboardCronJob) Run(ctx context.Context) error {
	start := time.Now()

	users, err := d.userRepo.GetAllUsersId(ctx)
	if err != nil {
		zap.L().Error("error to get all the users", zap.Error(err))
		return err
	}

	for _, user := range users {
		select {
		case <-ctx.Done():
			zap.L().Warn("execution canceled by context!")
			return ctx.Err()
		default:
			info, err := d.siteRepo.GetAllSiteInfoByUserId(ctx, user)
			if err != nil {
				zap.L().Error("error from compute the dashboard of the user", zap.Uint("user_id", user), zap.Error(err))
				continue
			}

			data, err := json.Marshal(info)
			if err != nil {
				zap.L().Error("error to serialize infos by user", zap.Uint("user_id", user), zap.Error(err))
				continue
			}

			key := fmt.Sprintf("dashboard:user:%d", user)

			err = d.redis.Set(ctx, key, data, 1*time.Minute).Err()
			if err != nil {
				zap.L().Error("error on save on redis of user", zap.Uint("user_id", user), zap.Error(err))
				continue
			}

			zap.L().Info("user dashboard updated on redis", zap.Uint("user_id", user))
		}
	}

	zap.L().Info("AllUsersDashboardJob ended", zap.Duration("duration", time.Since(start)))
	return nil
}
