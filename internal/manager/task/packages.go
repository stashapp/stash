package task

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/pkg"
)

type PackagesJob struct {
	PackageManager *pkg.Manager
	OnComplete     func()
}

func (j *PackagesJob) installPackage(ctx context.Context, p models.PackageSpecInput, progress *job.Progress) error {
	defer progress.Increment()

	if err := j.PackageManager.Install(ctx, p); err != nil {
		return fmt.Errorf("installing package: %w", err)
	}

	return nil
}

type InstallPackagesJob struct {
	PackagesJob
	Packages []*models.PackageSpecInput
}

func (j *InstallPackagesJob) Execute(ctx context.Context, progress *job.Progress) error {
	progress.SetTotal(len(j.Packages))

	for _, p := range j.Packages {
		if job.IsCancelled(ctx) {
			logger.Info("Cancelled installing packages")
			return nil
		}

		logger.Infof("Installing package %s", p.ID)
		taskDesc := fmt.Sprintf("Installing %s", p.ID)
		progress.ExecuteTask(taskDesc, func() {
			if err := j.installPackage(ctx, *p, progress); err != nil {
				logger.Errorf("Error installing package %s from %s: %v", p.ID, p.SourceURL, err)
			}
		})
	}

	if j.OnComplete != nil {
		j.OnComplete()
	}

	logger.Infof("Finished installing packages")
	return nil
}

type UpdatePackagesJob struct {
	PackagesJob
	Packages []*models.PackageSpecInput
}

func (j *UpdatePackagesJob) Execute(ctx context.Context, progress *job.Progress) error {
	// if no packages are specified, update all
	if len(j.Packages) == 0 {
		installed, err := j.PackageManager.InstalledStatus(ctx)
		if err != nil {
			return fmt.Errorf("error getting installed packages: %w", err)
		}

		for _, p := range installed {
			if p.Upgradable() {
				j.Packages = append(j.Packages, &models.PackageSpecInput{
					ID:        p.Local.ID,
					SourceURL: p.Remote.Repository.Path(),
				})
			}
		}
	}

	progress.SetTotal(len(j.Packages))

	for _, p := range j.Packages {
		if job.IsCancelled(ctx) {
			logger.Info("Cancelled updating packages")
			return nil
		}

		logger.Infof("Updating package %s", p.ID)
		taskDesc := fmt.Sprintf("Updating %s", p.ID)
		progress.ExecuteTask(taskDesc, func() {
			if err := j.installPackage(ctx, *p, progress); err != nil {
				logger.Errorf("Error updating package %s from %s: %v", p.ID, p.SourceURL, err)
			}
		})
	}

	if j.OnComplete != nil {
		j.OnComplete()
	}

	logger.Infof("Finished updating packages")
	return nil
}

type UninstallPackagesJob struct {
	PackagesJob
	Packages []*models.PackageSpecInput
}

func (j *UninstallPackagesJob) Execute(ctx context.Context, progress *job.Progress) error {
	progress.SetTotal(len(j.Packages))

	for _, p := range j.Packages {
		if job.IsCancelled(ctx) {
			logger.Info("Cancelled installing packages")
			return nil
		}

		logger.Infof("Uninstalling package %s", p.ID)
		taskDesc := fmt.Sprintf("Uninstalling %s", p.ID)
		progress.ExecuteTask(taskDesc, func() {
			if err := j.PackageManager.Uninstall(ctx, *p); err != nil {
				logger.Errorf("Error uninstalling package %s: %v", p.ID, err)
			}
		})
	}

	if j.OnComplete != nil {
		j.OnComplete()
	}

	logger.Infof("Finished uninstalling packages")
	return nil
}
