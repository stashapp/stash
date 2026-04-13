using Stash.Core.Events;

namespace Stash.Core.Interfaces;

public enum JobStatus
{
    Pending,
    Running,
    Completed,
    Failed,
    Cancelled
}

public record JobInfo(
    string Id,
    string Type,
    string Description,
    JobStatus Status,
    double Progress,
    string? SubTask,
    DateTime StartedAt,
    DateTime? CompletedAt,
    string? Error);

public interface IJobService
{
    string Enqueue(string type, string description, Func<IJobProgress, CancellationToken, Task> work);
    bool Cancel(string jobId);
    JobInfo? GetJob(string jobId);
    IReadOnlyList<JobInfo> GetAllJobs();
    IReadOnlyList<JobInfo> GetJobHistory();
}

public interface IJobProgress
{
    void Report(double progress, string? subTask = null);
}

public interface IScanService
{
    string StartScan(bool scanGenerators = false);
}

public interface IAutoTagService
{
    string StartAutoTag(IEnumerable<string>? performerIds = null, IEnumerable<string>? studioIds = null, IEnumerable<string>? tagIds = null);
}

public interface ICleanService
{
    string StartClean(bool dryRun = false);
}

public interface IBackupService
{
    string StartBackup();
    Task<string?> GetLatestBackupPathAsync(CancellationToken ct = default);
}

public interface IStreamService
{
    Task<(Stream stream, string contentType, long? fileSize)?> GetSceneStream(int sceneId, CancellationToken ct = default);
    Task<(Stream stream, string contentType)?> GetSceneScreenshot(int sceneId, double? seconds, CancellationToken ct = default);
}
