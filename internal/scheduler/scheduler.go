package scheduler

import (
	"context"
	"fmt"
	"github-service/config"
	"github-service/internal/core/domain"
	"github-service/internal/core/service"
	"github-service/internal/ports"
	"github-service/pkg/logger"
	"time"

	"github.com/go-co-op/gocron"
)

type Scheduler struct {
	monitorService *service.MonitorService
	cfg            *config.Config
	b              ports.BadgerImpl
	schedulers     map[string]*gocron.Scheduler // Map to track schedulers by repo ID
}

func NewScheduler(monitorService *service.MonitorService, cfg *config.Config) *Scheduler {
	return &Scheduler{
		monitorService: monitorService,
		cfg:            cfg,
		schedulers:     make(map[string]*gocron.Scheduler),
	}
}

func (s *Scheduler) ScheduleMonitoring(r domain.RepoData) {
	repoDataArray, _ := s.b.GetRepoArray("repos")
	for _, repo := range repoDataArray {
		repoKey := repo.RepoName // Use RepoName as the key for the schedulers map
		logger.LogInfo(fmt.Sprintf("Monitoring scheduled for repository: %s", repoKey))

		if _, exists := s.schedulers[repoKey]; !exists {
			scheduler := gocron.NewScheduler(time.UTC)
			s.schedulerJob(scheduler, repo)
			s.schedulerStart(scheduler, repo)
		}
	}
}

func (s *Scheduler) schedulerJob(scheduler *gocron.Scheduler, r domain.RepoData) {
	scheduler.Every(s.cfg.POLL_INTERVAL).Do(func() {
		s.monitorRepository(r)
	})
}

func (s *Scheduler) schedulerStart(scheduler *gocron.Scheduler, r domain.RepoData) {
	s.schedulers[r.RepoName] = scheduler
	scheduler.StartAsync()
}

func (s *Scheduler) monitorRepository(r domain.RepoData) {
	ctx := context.Background()
	if err := s.monitorService.MonitorRepository(ctx, r); err != nil {
		logger.LogError(fmt.Errorf("monitoring failed for repository name %s: %w", r.RepoName, err))
	}
}
