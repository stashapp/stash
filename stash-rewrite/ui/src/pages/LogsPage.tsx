import { useState, useEffect, useRef } from "react";
import { useQuery } from "@tanstack/react-query";
import { logs, type LogEntry } from "../api/client";
import { HubConnectionBuilder, LogLevel } from "@microsoft/signalr";
import { Trash2, Pause, Play, ArrowDown } from "lucide-react";

const LEVEL_COLORS: Record<string, string> = {
  Verbose: "text-gray-600",
  Debug: "text-gray-500",
  Information: "text-blue-400",
  Warning: "text-yellow-400",
  Error: "text-red-400",
  Fatal: "text-red-600 font-bold",
};

const LEVEL_BADGES: Record<string, string> = {
  Verbose: "bg-gray-800 text-gray-500",
  Debug: "bg-gray-800 text-gray-400",
  Information: "bg-blue-900/50 text-blue-400",
  Warning: "bg-yellow-900/50 text-yellow-400",
  Error: "bg-red-900/50 text-red-400",
  Fatal: "bg-red-900 text-red-300",
};

export function LogsPage() {
  const [entries, setEntries] = useState<LogEntry[]>([]);
  const [paused, setPaused] = useState(false);
  const [levelFilter, setLevelFilter] = useState<string>("all");
  const logEndRef = useRef<HTMLDivElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const [autoScroll, setAutoScroll] = useState(true);

  // Load initial logs
  const { data: initialLogs } = useQuery({
    queryKey: ["logs-initial"],
    queryFn: () => logs.recent(undefined, 500),
    staleTime: Infinity,
  });

  useEffect(() => {
    if (initialLogs) setEntries(initialLogs);
  }, [initialLogs]);

  // Real-time via SignalR
  useEffect(() => {
    const connection = new HubConnectionBuilder()
      .withUrl("/hubs/logs")
      .withAutomaticReconnect()
      .configureLogging(LogLevel.Warning)
      .build();

    connection.on("LogReceived", (entry: LogEntry) => {
      if (!paused) {
        setEntries((prev) => {
          const next = [...prev, entry];
          return next.length > 1000 ? next.slice(-500) : next;
        });
      }
    });

    connection.start().catch(() => {});
    return () => { connection.stop(); };
  }, [paused]);

  // Auto scroll
  useEffect(() => {
    if (autoScroll) logEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [entries, autoScroll]);

  const handleScroll = () => {
    if (containerRef.current) {
      const { scrollTop, scrollHeight, clientHeight } = containerRef.current;
      setAutoScroll(scrollHeight - scrollTop - clientHeight < 50);
    }
  };

  const filtered = levelFilter === "all"
    ? entries
    : entries.filter((e) => e.level === levelFilter);

  return (
    <div className="max-w-6xl mx-auto flex flex-col" style={{ height: "calc(100vh - 5rem)" }}>
      <div className="flex items-center justify-between mb-4">
        <h1 className="text-2xl font-bold">Logs</h1>
        <div className="flex items-center gap-2">
          <select
            value={levelFilter}
            onChange={(e) => setLevelFilter(e.target.value)}
            className="bg-gray-800 border border-gray-700 rounded px-2 py-1.5 text-sm text-gray-200"
          >
            <option value="all">All Levels</option>
            <option value="Debug">Debug</option>
            <option value="Information">Info</option>
            <option value="Warning">Warning</option>
            <option value="Error">Error</option>
          </select>
          <button
            onClick={() => setPaused(!paused)}
            className={`flex items-center gap-1.5 px-3 py-1.5 rounded-md text-sm transition-colors ${
              paused ? "bg-green-700 hover:bg-green-600" : "bg-gray-800 hover:bg-gray-700"
            }`}
          >
            {paused ? <Play className="w-3.5 h-3.5" /> : <Pause className="w-3.5 h-3.5" />}
            {paused ? "Resume" : "Pause"}
          </button>
          <button
            onClick={() => setEntries([])}
            className="flex items-center gap-1.5 px-3 py-1.5 rounded-md text-sm bg-gray-800 hover:bg-gray-700"
          >
            <Trash2 className="w-3.5 h-3.5" /> Clear
          </button>
          {!autoScroll && (
            <button
              onClick={() => { setAutoScroll(true); logEndRef.current?.scrollIntoView(); }}
              className="flex items-center gap-1.5 px-3 py-1.5 rounded-md text-sm bg-blue-600 hover:bg-blue-500"
            >
              <ArrowDown className="w-3.5 h-3.5" /> Scroll to bottom
            </button>
          )}
        </div>
      </div>

      <div
        ref={containerRef}
        onScroll={handleScroll}
        className="flex-1 bg-gray-950 border border-gray-800 rounded-lg overflow-y-auto font-mono text-xs"
      >
        {filtered.length === 0 ? (
          <div className="flex items-center justify-center h-full text-gray-600">
            No log entries
          </div>
        ) : (
          <div className="p-2 space-y-px">
            {filtered.map((entry, i) => (
              <div key={i} className="flex gap-2 py-0.5 hover:bg-gray-900/50 px-2 rounded">
                <span className="text-gray-600 flex-shrink-0 w-20">
                  {new Date(entry.timestamp).toLocaleTimeString()}
                </span>
                <span className={`flex-shrink-0 w-14 text-center rounded px-1 ${LEVEL_BADGES[entry.level] || "text-gray-400"}`}>
                  {entry.level.substring(0, 4).toUpperCase()}
                </span>
                <span className={LEVEL_COLORS[entry.level] || "text-gray-300"}>
                  {entry.message}
                  {entry.exception && (
                    <span className="text-red-400 ml-2">[{entry.exception}]</span>
                  )}
                </span>
              </div>
            ))}
            <div ref={logEndRef} />
          </div>
        )}
      </div>
    </div>
  );
}
