using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Logging;
using Stash.Core.Entities;
using Stash.Core.Interfaces;
using Stash.Data;

namespace Stash.Api.Services;

public class AutoTagService(
    IJobService jobService,
    IServiceScopeFactory scopeFactory,
    ILogger<AutoTagService> logger) : IAutoTagService
{
    public string StartAutoTag(IEnumerable<string>? performerIds = null, IEnumerable<string>? studioIds = null, IEnumerable<string>? tagIds = null)
    {
        return jobService.Enqueue("auto-tag", "Auto-tagging scenes by filename", async (progress, ct) =>
        {
            using var scope = scopeFactory.CreateScope();
            var db = scope.ServiceProvider.GetRequiredService<StashContext>();

            // Load all performers with aliases
            var performers = await db.Performers
                .Include(p => p.Aliases)
                .Select(p => new { p.Id, p.Name, Aliases = p.Aliases.Select(a => a.Alias).ToList() })
                .ToListAsync(ct);

            // Load all studios
            var studios = await db.Studios
                .Select(s => new { s.Id, s.Name })
                .ToListAsync(ct);

            // Load all tags with aliases
            var tags = await db.Tags
                .Include(t => t.Aliases)
                .Select(t => new { t.Id, t.Name, Aliases = t.Aliases.Select(a => a.Alias).ToList() })
                .ToListAsync(ct);

            // Get all scenes with file paths (need ParentFolder for Path computation)
            var scenes = await db.Scenes
                .Include(s => s.Files).ThenInclude(f => f.ParentFolder)
                .Include(s => s.ScenePerformers)
                .Include(s => s.SceneTags)
                .ToListAsync(ct);

            int matched = 0;
            int total = scenes.Count;

            for (int i = 0; i < total; i++)
            {
                ct.ThrowIfCancellationRequested();
                var scene = scenes[i];
                var file = scene.Files.FirstOrDefault();
                if (file == null) continue;

                var filePath = file.Path.ToLowerInvariant();
                var basename = System.IO.Path.GetFileNameWithoutExtension(filePath).ToLowerInvariant();
                bool changed = false;

                // Match performers by name or alias
                foreach (var p in performers)
                {
                    if (scene.ScenePerformers.Any(sp => sp.PerformerId == p.Id)) continue;
                    var names = new List<string> { p.Name.ToLowerInvariant() };
                    names.AddRange(p.Aliases.Select(a => a.ToLowerInvariant()));

                    if (names.Any(n => n.Length > 2 && (filePath.Contains(n, StringComparison.Ordinal) || basename.Contains(n, StringComparison.Ordinal))))
                    {
                        scene.ScenePerformers.Add(new ScenePerformer { SceneId = scene.Id, PerformerId = p.Id });
                        changed = true;
                    }
                }

                // Match studios by name
                foreach (var s in studios)
                {
                    if (scene.StudioId.HasValue) continue;
                    var name = s.Name.ToLowerInvariant();
                    if (name.Length > 2 && (filePath.Contains(name, StringComparison.Ordinal) || basename.Contains(name, StringComparison.Ordinal)))
                    {
                        scene.StudioId = s.Id;
                        changed = true;
                        break;
                    }
                }

                // Match tags by name or alias
                foreach (var t in tags)
                {
                    if (scene.SceneTags.Any(st => st.TagId == t.Id)) continue;
                    var names = new List<string> { t.Name.ToLowerInvariant() };
                    names.AddRange(t.Aliases.Select(a => a.ToLowerInvariant()));

                    if (names.Any(n => n.Length > 2 && (filePath.Contains(n, StringComparison.Ordinal) || basename.Contains(n, StringComparison.Ordinal))))
                    {
                        scene.SceneTags.Add(new SceneTag { SceneId = scene.Id, TagId = t.Id });
                        changed = true;
                    }
                }

                if (changed) matched++;
                progress.Report((double)(i + 1) / total, $"({i + 1}/{total}) {file.Basename}");
            }

            await db.SaveChangesAsync(ct);
            logger.LogInformation("Auto-tag complete: {Matched}/{Total} scenes updated", matched, total);
        });
    }
}
