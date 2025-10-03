package cronjobs

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"go-checker/internal/repository"
	"log"
	"time"
)

type DashboardCronJob struct {
	name     string
	schedule string
	siteRepo *repository.SiteRepo
	userID   uint
	redis    *redis.Client
}

func NewDashboardCronJob(siteRepo *repository.SiteRepo, userID uint, schedule string, redisClient *redis.Client) *DashboardCronJob {
	return &DashboardCronJob{
		name:     "DashboardCronJob",
		schedule: schedule,
		siteRepo: siteRepo,
		userID:   userID,
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

	info, err := d.siteRepo.GetAllSiteInfoByUserId(ctx, d.userID)
	if err != nil {
		log.Printf("error to get the dashboard data from: %d: %v", d.userID, err)
		return err
	}

	data, err := json.Marshal(info)
	if err != nil {
		log.Printf("error to serialize the data from:  %d: %v", d.userID, err)
		return err
	}

	key := "dashboard:user:" + string(rune(d.userID))

	err = d.redis.Set(ctx, key, data, 5*time.Minute).Err()
	if err != nil {
		log.Printf("error to save on redis to: %d: %v", d.userID, err)
		return err
	}

	log.Printf("user dashboard %d updated on redis (time: %s)", d.userID, time.Since(start))
	return nil
}
