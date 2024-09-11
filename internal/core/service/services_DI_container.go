package service

import (
	"context"
	"github-service/config"
	"github-service/internal/adapters/github"
	"github-service/internal/core/domain"

	"github-service/internal/ports"

	"log"
)

func SetupService(ctx context.Context, cfg config.Config, rData domain.RepoData, commitRepo ports.PostgresCommit, repositoryRepo ports.PostgresRepository, bs ports.BadgerImpl) (*CommitService, *RepositoryService, *MonitorService) {
	ghClient := github.NewGithubClient(&cfg, ctx)
	ghService := NewGithubService(&cfg, ctx, ghClient)

	// Initialize service instances
	commitService := NewCommitService(commitRepo, &cfg, ghService)
	repositoryService := NewRepositoryService(repositoryRepo, *commitService, &cfg, bs, ghService)

	// Initialize the commit monitor service
	monitorService := NewMonitorService(commitService, repositoryService, 5, 2, ghService)

	// Seed the database with initial data starting from the defined date
	if err := monitorService.MonitorRepository(ctx, rData); err != nil {
		log.Printf("Failed to add initial repository: %v", err)
	}

	go NewScheduler(monitorService, &cfg, bs).ScheduleMonitoring()

	return commitService, repositoryService, monitorService
}
