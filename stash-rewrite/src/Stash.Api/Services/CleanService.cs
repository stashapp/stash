using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Logging;
using Stash.Core.Interfaces;
using Stash.Data;

namespace Stash.Api.Services;

public class CleanService(
    IJobService jobService,
    IServiceScopeFactory scopeFactory,
    ILogger<CleanService> logger) : ICleanService
{
    public string StartClean(bool dryRun = false)
    {
        return jobService.Enqueue("clean", dryRun ? "Cleaning (dry run)" : "Cleaning library", async (progress, ct) =>
        {
            using var scope = scopeFactory.CreateScope();
            var db = scope.ServiceProvider.GetRequiredService<StashContext>();

            // Find scenes whose files no longer exist on disk
            var scenes = await db.Scenes
                .Include(s => s.Files).ThenInclude(f => f.ParentFolder)
                .ToListAsync(ct);

            var orphanSceneIds = new List<int>();
            int total = scenes.Count;

            for (int i = 0; i < total; i++)
            {
                ct.ThrowIfCancellationRequested();
                var scene = scenes[i];
                var file = scene.Files.FirstOrDefault();

                if (file == null || !File.Exists(file.Path))
                {
                    orphanSceneIds.Add(scene.Id);
                }

                progress.Report((double)(i + 1) / total, $"Checking ({i + 1}/{total})");
            }

            // Find images whose files no longer exist
            var images = await db.Images
                .Include(i => i.Files).ThenInclude(f => f.ParentFolder)
                .ToListAsync(ct);

            var orphanImageIds = new List<int>();
            for (int i = 0; i < images.Count; i++)
            {
                ct.ThrowIfCancellationRequested();
                var img = images[i];
                var file = img.Files.FirstOrDefault();
                if (file == null || !File.Exists(file.Path))
                {
                    orphanImageIds.Add(img.Id);
                }
            }

            // Find galleries whose folders no longer exist
            var galleries = await db.Galleries
                .Include(g => g.Folder)
                .Where(g => g.FolderId != null)
                .ToListAsync(ct);

            var orphanGalleryIds = new List<int>();
            foreach (var gallery in galleries)
            {
                ct.ThrowIfCancellationRequested();
                if (gallery.Folder != null && !Directory.Exists(gallery.Folder.Path))
                {
                    orphanGalleryIds.Add(gallery.Id);
                }
            }

            logger.LogInformation("Clean found {Scenes} orphaned scenes, {Images} orphaned images, {Galleries} orphaned galleries",
                orphanSceneIds.Count, orphanImageIds.Count, orphanGalleryIds.Count);

            if (dryRun)
            {
                logger.LogInformation("Dry run - no changes made");
                return;
            }

            // Remove orphaned records
            if (orphanSceneIds.Count > 0)
            {
                await db.Scenes.Where(s => orphanSceneIds.Contains(s.Id)).ExecuteDeleteAsync(ct);
                logger.LogInformation("Removed {Count} orphaned scenes", orphanSceneIds.Count);
            }

            if (orphanImageIds.Count > 0)
            {
                await db.Images.Where(im => orphanImageIds.Contains(im.Id)).ExecuteDeleteAsync(ct);
                logger.LogInformation("Removed {Count} orphaned images", orphanImageIds.Count);
            }

            if (orphanGalleryIds.Count > 0)
            {
                await db.Galleries.Where(g => orphanGalleryIds.Contains(g.Id)).ExecuteDeleteAsync(ct);
                logger.LogInformation("Removed {Count} orphaned galleries", orphanGalleryIds.Count);
            }
        });
    }
}
