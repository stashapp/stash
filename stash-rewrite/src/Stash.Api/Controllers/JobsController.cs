using Microsoft.AspNetCore.Mvc;
using Stash.Api.Services;
using Stash.Core.Interfaces;

namespace Stash.Api.Controllers;

[ApiController]
[Route("api/[controller]")]
public class JobsController(IJobService jobService, IScanService scanService, IThumbnailService thumbnailService, IFingerprintService fingerprintService, IAutoTagService autoTagService, ICleanService cleanService, IBackupService backupService) : ControllerBase
{
    [HttpGet]
    public ActionResult<IReadOnlyList<JobInfo>> GetJobs()
        => Ok(jobService.GetAllJobs());

    [HttpGet("history")]
    public ActionResult<IReadOnlyList<JobInfo>> GetHistory()
        => Ok(jobService.GetJobHistory());

    [HttpGet("{jobId}")]
    public ActionResult<JobInfo> GetJob(string jobId)
    {
        var job = jobService.GetJob(jobId);
        return job != null ? Ok(job) : NotFound();
    }

    [HttpDelete("{jobId}")]
    public IActionResult CancelJob(string jobId)
    {
        return jobService.Cancel(jobId) ? Ok() : NotFound();
    }

    [HttpPost("scan")]
    public ActionResult<object> StartScan([FromQuery] bool generatePreviews = false)
    {
        var jobId = scanService.StartScan(generatePreviews);
        return Accepted(new { jobId });
    }

    [HttpPost("generate-thumbnails")]
    public ActionResult<object> GenerateThumbnails()
    {
        var jobId = thumbnailService.StartGenerateAllThumbnails();
        return Accepted(new { jobId });
    }

    [HttpPost("generate-scene-phashes")]
    public ActionResult<object> GenerateScenePhashes()
    {
        var jobId = fingerprintService.StartGenerateScenePhashes();
        return Accepted(new { jobId });
    }

    [HttpPost("generate-image-phashes")]
    public ActionResult<object> GenerateImagePhashes()
    {
        var jobId = fingerprintService.StartGenerateImagePhashes();
        return Accepted(new { jobId });
    }

    [HttpPost("auto-tag")]
    public ActionResult<object> StartAutoTag([FromBody] AutoTagRequest? request = null)
    {
        var jobId = autoTagService.StartAutoTag(request?.PerformerIds, request?.StudioIds, request?.TagIds);
        return Accepted(new { jobId });
    }

    [HttpPost("clean")]
    public ActionResult<object> StartClean([FromQuery] bool dryRun = false)
    {
        var jobId = cleanService.StartClean(dryRun);
        return Accepted(new { jobId });
    }

    [HttpPost("backup")]
    public ActionResult<object> StartBackup()
    {
        var jobId = backupService.StartBackup();
        return Accepted(new { jobId });
    }

    [HttpGet("backup/latest")]
    public async Task<ActionResult<object>> GetLatestBackup()
    {
        var path = await backupService.GetLatestBackupPathAsync();
        return path != null ? Ok(new { path }) : NotFound();
    }
}

public class AutoTagRequest
{
    public IEnumerable<string>? PerformerIds { get; set; }
    public IEnumerable<string>? StudioIds { get; set; }
    public IEnumerable<string>? TagIds { get; set; }
}
