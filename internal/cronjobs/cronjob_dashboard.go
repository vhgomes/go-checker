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

	users, err := d.userRepo.GetAllUsersId()
	if err != nil {
		log.Println("error to get all the users")
		return err
	}

	for _, user := range users {
		select {
		case <-ctx.Done():
			log.Println("execution canceled by context!")
			return ctx.Err()
		default:
			info, err := d.siteRepo.GetAllSiteInfoByUserId(ctx, user)
			if err != nil {
				log.Printf("error from compute the dashboard of the user: %d: %v", user, err)
				continue
			}

			data, err := json.Marshal(info)
			if err != nil {
				log.Printf("error to serialize infos by user: %d: %v", user, err)
				continue
			}

			key := "dashboard:user:" + string(rune(user))

			err = d.redis.Set(ctx, key, data, 5*time.Minute).Err()
			if err != nil {
				log.Printf("error on save on redis of user: %d: %v", user, err)
				continue
			}

			log.Printf("user dashboard %d updated on redis", user)
		}
	}

	log.Printf("AllUsersDashboardJob ended in %s", time.Since(start))
	return nil
}
