using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using Stash.Api.Services;
using Stash.Core.DTOs;
using Stash.Core.Entities;
using Stash.Core.Interfaces;
using Stash.Data;
using System.Net.Http;
using System.Text;
using System.Text.Json;

namespace Stash.Api.Controllers;

[ApiController]
[Route("api/metadata")]
public class MetadataController(
    IScanService scanService,
    IJobService jobService,
    IThumbnailService thumbnailService,
    IFingerprintService fingerprintService,
    IServiceScopeFactory scopeFactory,
    StashContext db,
    StashConfiguration config,
    ILogger<MetadataController> logger) : ControllerBase
{
    [HttpPost("scan")]
    public ActionResult<object> StartScan([FromBody] ScanOptionsDto? opts)
    {
        var jobId = scanService.StartScan(opts?.ScanGenerators ?? false);
        return Ok(new { jobId });
    }

    [HttpPost("generate")]
    public ActionResult<object> StartGenerate([FromBody] GenerateOptionsDto? opts)
    {
        var jobId = jobService.Enqueue("generate", "Generating content", async (progress, ct) =>
        {
            using var scope = scopeFactory.CreateScope();
            var dbCtx = scope.ServiceProvider.GetRequiredService<StashContext>();

            var scenes = opts?.SceneIds != null && opts.SceneIds.Count > 0
                ? await dbCtx.Scenes.Include(s => s.Files).ThenInclude(f => f.ParentFolder).Include(s => s.Files).ThenInclude(f => f.Fingerprints).Where(s => opts.SceneIds.Contains(s.Id)).ToListAsync(ct)
                : await dbCtx.Scenes.Include(s => s.Files).ThenInclude(f => f.ParentFolder).Include(s => s.Files).ThenInclude(f => f.Fingerprints).ToListAsync(ct);

            var total = scenes.Count;
            var processed = 0;
            var parallelism = Math.Max(1, Environment.ProcessorCount / 4 + 1);

            await Parallel.ForEachAsync(scenes, new ParallelOptions { MaxDegreeOfParallelism = parallelism, CancellationToken = ct }, async (scene, token) =>
            {
                var current = Interlocked.Increment(ref processed);
                progress.Report((double)current / total, $"Processing {scene.Title ?? "Untitled"} ({current}/{total})");

                var file = scene.Files.FirstOrDefault();
                if (file == null) return;
                var path = Path.Combine(file.ParentFolder?.Path ?? "", file.Basename);
                if (!System.IO.File.Exists(path)) return;

                if (opts?.Thumbnails != false)
                {
                    if (opts?.Overwrite == true)
                    {
                        // Delete existing thumbnail to force regeneration
                        var hash = Convert.ToHexStringLower(System.Security.Cryptography.SHA256.HashData(BitConverter.GetBytes(scene.Id)));
                        var thumbPath = Path.Combine(config.GeneratedPath, "screenshots", hash[..2], $"{scene.Id}.jpg");
                        if (System.IO.File.Exists(thumbPath)) System.IO.File.Delete(thumbPath);
                    }
                    await thumbnailService.GenerateSceneThumbnailAsync(scene.Id, null, token);
                }

                if (opts?.Previews == true)
                {
                    if (opts?.Overwrite == true)
                    {
                        var previewPath = thumbnailService.GetPreviewPath(scene.Id);
                        if (System.IO.File.Exists(previewPath)) System.IO.File.Delete(previewPath);
                    }
                    await thumbnailService.GenerateScenePreviewAsync(scene.Id, token);
                }

                if (opts?.Sprites == true)
                {
                    if (opts?.Overwrite == true)
                    {
                        var spritePath = thumbnailService.GetSpritePath(scene.Id);
                        var vttPath = thumbnailService.GetSpriteVttPath(scene.Id);
                        if (System.IO.File.Exists(spritePath)) System.IO.File.Delete(spritePath);
                        if (System.IO.File.Exists(vttPath)) System.IO.File.Delete(vttPath);
                    }
                    await thumbnailService.GenerateSceneSpriteAsync(scene.Id, token);
                }

                if (opts?.Phashes == true)
                {
                    var phash = await fingerprintService.ComputeVideoPhashAsync(path, file.Duration, token);
                    if (phash != null)
                    {
                        lock (file.Fingerprints)
                        {
                            var existing = file.Fingerprints.FirstOrDefault(f => f.Type == "phash");
                            if (existing != null) existing.Value = phash;
                            else file.Fingerprints.Add(new FileFingerprint { Type = "phash", Value = phash });
                        }
                    }
                }
            });

            await dbCtx.SaveChangesAsync(ct);
        });

        return Ok(new { jobId });
    }

    [HttpPost("auto-tag")]
    public ActionResult<object> StartAutoTag([FromBody] AutoTagOptionsDto? opts)
    {
        var jobId = jobService.Enqueue("auto-tag", "Auto-tagging content", async (progress, ct) =>
        {
            using var scope = scopeFactory.CreateScope();
            var dbCtx = scope.ServiceProvider.GetRequiredService<StashContext>();

            // Load all performers, studios, tags for matching
            var performers = await dbCtx.Performers.AsNoTracking().ToListAsync(ct);
            var studios = await dbCtx.Studios.AsNoTracking().ToListAsync(ct);
            var tags = await dbCtx.Tags.Include(t => t.Aliases).AsNoTracking().ToListAsync(ct);

            // Filter to configured subsets if provided
            if (opts?.Performers?.Count > 0)
                performers = performers.Where(p => opts.Performers.Any(n => p.Name.Contains(n, StringComparison.OrdinalIgnoreCase))).ToList();
            if (opts?.Studios?.Count > 0)
                studios = studios.Where(s => opts.Studios.Any(n => s.Name.Contains(n, StringComparison.OrdinalIgnoreCase))).ToList();
            if (opts?.Tags?.Count > 0)
                tags = tags.Where(t => opts.Tags.Any(n => t.Name.Contains(n, StringComparison.OrdinalIgnoreCase))).ToList();

            var scenes = await dbCtx.Scenes
                .Include(s => s.Files).ThenInclude(f => f.ParentFolder)
                .Include(s => s.SceneTags)
                .Include(s => s.ScenePerformers)
                .ToListAsync(ct);

            var total = scenes.Count;
            var tagged = 0;

            for (var i = 0; i < scenes.Count; i++)
            {
                ct.ThrowIfCancellationRequested();
                progress.Report((double)(i + 1) / total, $"Auto-tagging scenes ({i + 1}/{total})");

                var scene = scenes[i];
                var sceneFile = scene.Files.FirstOrDefault();
                if (sceneFile == null) continue;

                var path = Path.Combine(sceneFile.ParentFolder?.Path ?? "", sceneFile.Basename).ToLowerInvariant();
                var title = (scene.Title ?? "").ToLowerInvariant();
                var searchText = $"{path} {title}";

                // Match performers by name in file path/title
                var existingPerformerIds = scene.ScenePerformers.Select(sp => sp.PerformerId).ToHashSet();
                foreach (var performer in performers)
                {
                    if (existingPerformerIds.Contains(performer.Id)) continue;
                    if (searchText.Contains(performer.Name.ToLowerInvariant()))
                    {
                        scene.ScenePerformers.Add(new ScenePerformer { PerformerId = performer.Id, SceneId = scene.Id });
                        tagged++;
                    }
                }

                // Match studios by name in file path
                if (scene.StudioId == null)
                {
                    foreach (var studio in studios)
                    {
                        if (searchText.Contains(studio.Name.ToLowerInvariant()))
                        {
                            scene.StudioId = studio.Id;
                            tagged++;
                            break;
                        }
                    }
                }

                // Match tags by name or alias in file path/title
                var existingTagIds = scene.SceneTags.Select(st => st.TagId).ToHashSet();
                foreach (var tag in tags)
                {
                    if (existingTagIds.Contains(tag.Id)) continue;
                    var names = new List<string> { tag.Name.ToLowerInvariant() };
                    names.AddRange(tag.Aliases.Select(a => a.Alias.ToLowerInvariant()));

                    if (names.Any(n => searchText.Contains(n)))
                    {
                        scene.SceneTags.Add(new SceneTag { TagId = tag.Id, SceneId = scene.Id });
                        tagged++;
                    }
                }
            }

            await dbCtx.SaveChangesAsync(ct);
            logger.LogInformation("Auto-tag completed. {Tagged} associations created", tagged);
        });

        return Ok(new { jobId });
    }

    [HttpPost("clean")]
    public ActionResult<object> StartClean([FromBody] CleanOptionsDto? opts)
    {
        var jobId = jobService.Enqueue("clean", "Cleaning library", async (progress, ct) =>
        {
            using var scope = scopeFactory.CreateScope();
            var dbCtx = scope.ServiceProvider.GetRequiredService<StashContext>();

            var files = await dbCtx.Set<BaseFileEntity>()
                .Include(f => f.ParentFolder)
                .ToListAsync(ct);

            var total = files.Count;
            var cleaned = 0;

            for (var i = 0; i < files.Count; i++)
            {
                ct.ThrowIfCancellationRequested();
                progress.Report((double)(i + 1) / total, $"Checking files ({i + 1}/{total})");

                var file = files[i];
                var filePath = Path.Combine(file.ParentFolder?.Path ?? "", file.Basename);

                if (opts?.Paths?.Count > 0 && !opts.Paths.Any(p => filePath.StartsWith(p, StringComparison.OrdinalIgnoreCase)))
                    continue;

                if (!System.IO.File.Exists(filePath))
                {
                    if (opts?.DryRun != true)
                    {
                        dbCtx.Set<BaseFileEntity>().Remove(file);
                        cleaned++;
                    }
                    else
                    {
                        logger.LogInformation("[Dry Run] Would remove: {Path}", filePath);
                        cleaned++;
                    }
                }
            }

            if (opts?.DryRun != true)
                await dbCtx.SaveChangesAsync(ct);

            logger.LogInformation("Clean completed. {Cleaned} missing files {Action}", cleaned, opts?.DryRun == true ? "found" : "removed");
        });

        return Ok(new { jobId });
    }

    [HttpPost("export")]
    public ActionResult<object> StartExport([FromBody] ExportOptionsDto? opts)
    {
        var jobId = jobService.Enqueue("export", "Exporting metadata", async (progress, ct) =>
        {
            using var scope = scopeFactory.CreateScope();
            var dbCtx = scope.ServiceProvider.GetRequiredService<StashContext>();

            var exportPath = Path.Combine(config.GeneratedPath ?? Path.GetTempPath(), "export");
            Directory.CreateDirectory(exportPath);
            var timestamp = DateTime.UtcNow.ToString("yyyyMMdd_HHmmss");
            var exportFile = Path.Combine(exportPath, $"stash-export-{timestamp}.json");

            var exportData = new Dictionary<string, object>();

            if (opts?.IncludeScenes != false)
            {
                progress.Report(0.1, "Exporting scenes...");
                exportData["scenes"] = await dbCtx.Scenes
                    .Include(s => s.SceneTags).ThenInclude(st => st.Tag)
                    .Include(s => s.ScenePerformers).ThenInclude(sp => sp.Performer)
                    .Include(s => s.Studio)
                    .Include(s => s.Files).ThenInclude(f => f.Fingerprints)
                    .AsNoTracking()
                    .ToListAsync(ct);
            }

            if (opts?.IncludePerformers != false)
            {
                progress.Report(0.3, "Exporting performers...");
                exportData["performers"] = await dbCtx.Performers.AsNoTracking().ToListAsync(ct);
            }

            if (opts?.IncludeStudios != false)
            {
                progress.Report(0.5, "Exporting studios...");
                exportData["studios"] = await dbCtx.Studios.AsNoTracking().ToListAsync(ct);
            }

            if (opts?.IncludeTags != false)
            {
                progress.Report(0.6, "Exporting tags...");
                exportData["tags"] = await dbCtx.Tags.AsNoTracking().ToListAsync(ct);
            }

            if (opts?.IncludeGalleries != false)
            {
                progress.Report(0.7, "Exporting galleries...");
                exportData["galleries"] = await dbCtx.Galleries.AsNoTracking().ToListAsync(ct);
            }

            if (opts?.IncludeGroups != false)
            {
                progress.Report(0.8, "Exporting groups...");
                exportData["groups"] = await dbCtx.Groups.AsNoTracking().ToListAsync(ct);
            }

            progress.Report(0.9, "Writing export file...");
            var jsonOpts = new JsonSerializerOptions { WriteIndented = true };
            await System.IO.File.WriteAllTextAsync(exportFile, JsonSerializer.Serialize(exportData, jsonOpts), ct);

            logger.LogInformation("Export completed: {Path}", exportFile);
        });

        return Ok(new { jobId });
    }

    [HttpPost("clean-generated")]
    public ActionResult<object> CleanGenerated()
    {
        var jobId = jobService.Enqueue("clean-generated", "Cleaning generated files", async (progress, ct) =>
        {
            var generatedPath = config.GeneratedPath;
            if (string.IsNullOrEmpty(generatedPath) || !Directory.Exists(generatedPath))
            {
                logger.LogWarning("Generated path not configured or does not exist");
                return;
            }

            var dirs = new[] { "screenshots", "thumbnails", "previews", "sprites", "transcodes", "vtt" };
            var totalCleared = 0L;

            for (var i = 0; i < dirs.Length; i++)
            {
                ct.ThrowIfCancellationRequested();
                progress.Report((double)(i + 1) / dirs.Length, $"Cleaning {dirs[i]}...");

                var dir = Path.Combine(generatedPath, dirs[i]);
                if (!Directory.Exists(dir)) continue;

                foreach (var file in Directory.GetFiles(dir, "*", SearchOption.AllDirectories))
                {
                    var fi = new FileInfo(file);
                    totalCleared += fi.Length;
                    fi.Delete();
                }
            }

            logger.LogInformation("Cleaned generated files. Freed {Size} bytes", totalCleared);
        });

        return Ok(new { jobId });
    }

    [HttpPost("identify")]
    public ActionResult<object> StartIdentify([FromBody] IdentifyOptionsDto? opts)
    {
        var jobId = jobService.Enqueue("identify", "Identifying scenes", async (progress, ct) =>
        {
            using var scope = scopeFactory.CreateScope();
            var dbCtx = scope.ServiceProvider.GetRequiredService<StashContext>();
            var stashBoxSvc = scope.ServiceProvider.GetService<StashBoxService>();

            var scenes = opts?.SceneIds?.Count > 0
                ? await dbCtx.Scenes.Include(s => s.Files).ThenInclude(f => f.Fingerprints).Where(s => opts.SceneIds.Contains(s.Id)).ToListAsync(ct)
                : await dbCtx.Scenes.Include(s => s.Files).ThenInclude(f => f.Fingerprints).ToListAsync(ct);

            var total = scenes.Count;
            for (var i = 0; i < total; i++)
            {
                ct.ThrowIfCancellationRequested();
                progress.Report((double)(i + 1) / total, $"Identifying scene {i + 1}/{total}");

                var scene = scenes[i];
                var fingerprints = scene.Files.SelectMany(f => f.Fingerprints).ToList();
                if (fingerprints.Count == 0) continue;

                // Attempt StashBox identification
                if (stashBoxSvc != null)
                {
                    try
                    {
                        var matches = await stashBoxSvc.SearchScenesAsync(scene, null, null, ct);
                        if (matches.Count > 0)
                        {
                            var best = matches[0];
                            await stashBoxSvc.MergeSceneAsync(scene, best.Endpoint, best.Id, null, ct);
                        }
                    }
                    catch (Exception ex)
                    {
                        logger.LogWarning(ex, "StashBox identify failed for scene {SceneId}", scene.Id);
                    }
                }
            }

            await dbCtx.SaveChangesAsync(ct);
            logger.LogInformation("Identify completed for {Count} scenes", total);
        });

        return Ok(new { jobId });
    }

    [HttpPost("sync-fingerprints")]
    public ActionResult<object> SyncFingerprints([FromBody] SyncFingerprintsOptionsDto? opts)
    {
        var sourceUrl = opts?.SourceUrl ?? "http://localhost:3000/graphql";
        var apiKey = opts?.ApiKey;

        var jobId = jobService.Enqueue("sync-fingerprints", "Syncing fingerprints from source stash", async (progress, ct) =>
        {
            using var scope = scopeFactory.CreateScope();
            var dbCtx = scope.ServiceProvider.GetRequiredService<StashContext>();
            using var httpClient = new HttpClient();
            httpClient.Timeout = TimeSpan.FromSeconds(30);
            if (!string.IsNullOrEmpty(apiKey))
                httpClient.DefaultRequestHeaders.Add("ApiKey", apiKey);

            // Step 1: Fetch all fingerprints from the source stash, paging through results
            progress.Report(0, "Fetching fingerprints from source stash...");
            var oshashToPhash = new Dictionary<string, string>();
            var page = 1;
            var perPage = 100;
            var totalScenes = 0;
            var fetched = 0;

            do
            {
                ct.ThrowIfCancellationRequested();

                var graphqlQuery = new
                {
                    query = @"query FindScenes($filter: FindFilterType!) {
                        findScenes(filter: $filter) {
                            count
                            scenes {
                                files {
                                    fingerprints {
                                        type
                                        value
                                    }
                                }
                            }
                        }
                    }",
                    variables = new
                    {
                        filter = new { page, per_page = perPage, sort = "id", direction = "ASC" }
                    }
                };

                var jsonPayload = JsonSerializer.Serialize(graphqlQuery);
                var response = await httpClient.PostAsync(
                    sourceUrl,
                    new StringContent(jsonPayload, Encoding.UTF8, "application/json"),
                    ct);

                response.EnsureSuccessStatusCode();
                var responseJson = await response.Content.ReadAsStringAsync(ct);

                using var doc = JsonDocument.Parse(responseJson);
                var data = doc.RootElement.GetProperty("data").GetProperty("findScenes");
                totalScenes = data.GetProperty("count").GetInt32();

                foreach (var scene in data.GetProperty("scenes").EnumerateArray())
                {
                    foreach (var file in scene.GetProperty("files").EnumerateArray())
                    {
                        string? oshash = null;
                        string? phash = null;

                        foreach (var fp in file.GetProperty("fingerprints").EnumerateArray())
                        {
                            var type = fp.GetProperty("type").GetString();
                            var value = fp.GetProperty("value").GetString();
                            if (type == "oshash") oshash = value;
                            else if (type == "phash") phash = value;
                        }

                        if (oshash != null && phash != null)
                            oshashToPhash.TryAdd(oshash, phash);
                    }
                }

                fetched += perPage;
                page++;
                progress.Report(Math.Min(0.5, (double)fetched / Math.Max(totalScenes, 1)),
                    $"Fetched {Math.Min(fetched, totalScenes)}/{totalScenes} scenes from source...");
            }
            while (fetched < totalScenes);

            logger.LogInformation("Fetched {Count} oshash→phash mappings from source stash", oshashToPhash.Count);

            // Step 2: Load all files with fingerprints from our DB
            progress.Report(0.5, "Loading local scene fingerprints...");
            var localFiles = await dbCtx.Set<BaseFileEntity>()
                .Include(f => f.Fingerprints)
                .ToListAsync(ct);

            var updated = 0;
            var created = 0;
            var total = localFiles.Count;

            for (var i = 0; i < localFiles.Count; i++)
            {
                ct.ThrowIfCancellationRequested();
                var file = localFiles[i];
                var localOshash = file.Fingerprints.FirstOrDefault(f => f.Type == "oshash")?.Value;
                if (localOshash == null) continue;

                if (!oshashToPhash.TryGetValue(localOshash, out var sourcePhash)) continue;

                var existingPhash = file.Fingerprints.FirstOrDefault(f => f.Type == "phash");
                if (existingPhash != null)
                {
                    if (existingPhash.Value != sourcePhash)
                    {
                        existingPhash.Value = sourcePhash;
                        updated++;
                    }
                }
                else
                {
                    file.Fingerprints.Add(new FileFingerprint { FileId = file.Id, Type = "phash", Value = sourcePhash });
                    created++;
                }

                if ((i + 1) % 100 == 0)
                    progress.Report(0.5 + 0.5 * ((double)(i + 1) / total),
                        $"Processing files ({i + 1}/{total})...");
            }

            await dbCtx.SaveChangesAsync(ct);
            logger.LogInformation("Fingerprint sync completed. {Updated} updated, {Created} created from {Total} source mappings",
                updated, created, oshashToPhash.Count);
        });

        return Ok(new { jobId });
    }
}
