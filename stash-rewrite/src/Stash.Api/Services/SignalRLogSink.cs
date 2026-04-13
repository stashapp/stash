using System.Collections.Concurrent;
using Microsoft.AspNetCore.SignalR;
using Serilog.Core;
using Serilog.Events;
using Stash.Api.Hubs;

namespace Stash.Api.Services;

public class SignalRLogSink : ILogEventSink
{
    private static IHubContext<LogHub>? _hubContext;
    private static readonly ConcurrentQueue<LogEntry> _recentLogs = new();
    private const int MaxLogs = 500;

    public static void SetHubContext(IHubContext<LogHub> hubContext)
    {
        _hubContext = hubContext;
    }

    public static IReadOnlyList<LogEntry> GetRecentLogs()
    {
        return [.. _recentLogs];
    }

    public void Emit(LogEvent logEvent)
    {
        var entry = new LogEntry
        {
            Timestamp = logEvent.Timestamp.UtcDateTime.ToString("o"),
            Level = logEvent.Level.ToString(),
            Message = logEvent.RenderMessage(),
            Exception = logEvent.Exception?.Message
        };

        _recentLogs.Enqueue(entry);
        while (_recentLogs.Count > MaxLogs)
            _recentLogs.TryDequeue(out _);

        _hubContext?.Clients.All.SendAsync("LogReceived", entry);
    }
}

public record LogEntry
{
    public string Timestamp { get; init; } = "";
    public string Level { get; init; } = "";
    public string Message { get; init; } = "";
    public string? Exception { get; init; }
}
