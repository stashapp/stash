import { useState, useEffect, useCallback, useRef } from "react";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { jobs } from "../api/client";
import { HubConnectionBuilder, LogLevel } from "@microsoft/signalr";
import type { JobInfo } from "../api/types";
import { X, Loader2, CheckCircle, XCircle, Ban, Clock, Trash2 } from "lucide-react";

interface Props {
  open: boolean;
  onClose: () => void;
}

const statusIcon = (status: JobInfo["status"]) => {
  switch (status) {
    case "Running": return <Loader2 className="w-4 h-4 text-blue-400 animate-spin" />;
    case "Completed": return <CheckCircle className="w-4 h-4 text-green-400" />;
    case "Failed": return <XCircle className="w-4 h-4 text-red-400" />;
    case "Cancelled": return <Ban className="w-4 h-4 text-gray-500" />;
    default: return <Clock className="w-4 h-4 text-yellow-400" />;
  }
};

export function JobDrawer({ open, onClose }: Props) {
  const queryClient = useQueryClient();
  const [realtimeJobs, setRealtimeJobs] = useState<Map<string, JobInfo>>(new Map());
  const connectionRef = useRef<ReturnType<typeof HubConnectionBuilder.prototype.build> | null>(null);

  const { data: activeJobs } = useQuery({
    queryKey: ["jobs-active"],
    queryFn: jobs.list,
    refetchInterval: open ? 3000 : false,
  });

  const { data: jobHistory } = useQuery({
    queryKey: ["jobs-history"],
    queryFn: jobs.history,
    enabled: open,
  });

  // SignalR real-time updates
  useEffect(() => {
    const connection = new HubConnectionBuilder()
      .withUrl("/hubs/jobs")
      .withAutomaticReconnect()
      .configureLogging(LogLevel.Warning)
      .build();

    connection.on("JobUpdated", (job: JobInfo) => {
      setRealtimeJobs((prev) => {
        const next = new Map(prev);
        next.set(job.id, job);
        return next;
      });
      // Invalidate queries to stay in sync
      queryClient.invalidateQueries({ queryKey: ["jobs-active"] });
      queryClient.invalidateQueries({ queryKey: ["jobs-history"] });
      // When a job completes, invalidate content queries
      if (job.status === "Completed") {
        queryClient.invalidateQueries({ queryKey: ["scenes"] });
        queryClient.invalidateQueries({ queryKey: ["images"] });
        queryClient.invalidateQueries({ queryKey: ["galleries"] });
        queryClient.invalidateQueries({ queryKey: ["performers"] });
        queryClient.invalidateQueries({ queryKey: ["stats"] });
      }
    });

    connection.start().catch(() => {});
    connectionRef.current = connection;

    return () => {
      connection.stop();
    };
  }, [queryClient]);

  const handleCancel = useCallback(async (id: string) => {
    await jobs.cancel(id);
    queryClient.invalidateQueries({ queryKey: ["jobs-active"] });
    queryClient.invalidateQueries({ queryKey: ["jobs-history"] });
  }, [queryClient]);

  // Merge API jobs with real-time updates
  const mergedActive = activeJobs?.map((j) => realtimeJobs.get(j.id) ?? j) ?? [];
  // Also add any real-time jobs not in the API response
  for (const [id, job] of realtimeJobs) {
    if (
      (job.status === "Running" || job.status === "Pending") &&
      !mergedActive.find((j) => j.id === id)
    ) {
      mergedActive.push(job);
    }
  }

  const runningCount = mergedActive.filter((j) => j.status === "Running" || j.status === "Pending").length;

  if (!open) return null;

  return (
    <>
      {/* Backdrop */}
      <div className="fixed inset-0 bg-black/50 z-40" onClick={onClose} />

      {/* Drawer */}
      <div className="fixed right-0 top-0 h-full w-96 bg-gray-900 border-l border-gray-800 z-50 flex flex-col shadow-2xl">
        <div className="flex items-center justify-between px-4 py-3 border-b border-gray-800">
          <h2 className="font-semibold">
            Jobs {runningCount > 0 && <span className="text-blue-400 text-sm ml-1">({runningCount} active)</span>}
          </h2>
          <button onClick={onClose} className="text-gray-400 hover:text-white">
            <X className="w-5 h-5" />
          </button>
        </div>

        <div className="flex-1 overflow-y-auto">
          {/* Active jobs */}
          {mergedActive.length > 0 && (
            <div className="p-4">
              <h3 className="text-xs font-semibold text-gray-500 uppercase mb-2">Active</h3>
              <div className="space-y-2">
                {mergedActive.map((job) => (
                  <JobCard key={job.id} job={job} onCancel={handleCancel} />
                ))}
              </div>
            </div>
          )}

          {/* History */}
          {jobHistory && jobHistory.length > 0 && (
            <div className="p-4 border-t border-gray-800">
              <h3 className="text-xs font-semibold text-gray-500 uppercase mb-2">History</h3>
              <div className="space-y-2">
                {jobHistory.map((job) => (
                  <JobCard key={job.id} job={job} />
                ))}
              </div>
            </div>
          )}

          {mergedActive.length === 0 && (!jobHistory || jobHistory.length === 0) && (
            <div className="p-8 text-center text-gray-600 text-sm">
              No jobs running or in history
            </div>
          )}
        </div>
      </div>
    </>
  );
}

function JobCard({ job, onCancel }: { job: JobInfo; onCancel?: (id: string) => void }) {
  const progressPct = Math.round(job.progress * 100);

  return (
    <div className="bg-gray-800 rounded-lg p-3">
      <div className="flex items-start justify-between gap-2">
        <div className="flex items-center gap-2 flex-1 min-w-0">
          {statusIcon(job.status)}
          <div className="min-w-0 flex-1">
            <p className="text-sm font-medium truncate">{job.description}</p>
            {job.subTask && (
              <p className="text-xs text-gray-500 truncate mt-0.5">{job.subTask}</p>
            )}
          </div>
        </div>
        {(job.status === "Running" || job.status === "Pending") && onCancel && (
          <button onClick={() => onCancel(job.id)} className="text-gray-500 hover:text-red-400 flex-shrink-0">
            <Trash2 className="w-3.5 h-3.5" />
          </button>
        )}
      </div>

      {job.status === "Running" && (
        <div className="mt-2">
          <div className="h-1.5 bg-gray-700 rounded-full overflow-hidden">
            <div
              className="h-full bg-blue-500 rounded-full transition-all duration-300"
              style={{ width: `${progressPct}%` }}
            />
          </div>
          <p className="text-xs text-gray-500 mt-1">{progressPct}%</p>
        </div>
      )}

      {job.error && (
        <p className="text-xs text-red-400 mt-1 truncate">{job.error}</p>
      )}

      {job.completedAt && (
        <p className="text-xs text-gray-600 mt-1">
          {new Date(job.completedAt).toLocaleTimeString()}
        </p>
      )}
    </div>
  );
}

// Export a hook for the navbar badge
export function useJobCount() {
  const [count, setCount] = useState(0);

  useEffect(() => {
    const connection = new HubConnectionBuilder()
      .withUrl("/hubs/jobs")
      .withAutomaticReconnect()
      .configureLogging(LogLevel.None)
      .build();

    let activeIds = new Set<string>();

    connection.on("JobUpdated", (job: JobInfo) => {
      if (job.status === "Running" || job.status === "Pending") {
        activeIds.add(job.id);
      } else {
        activeIds.delete(job.id);
      }
      setCount(activeIds.size);
    });

    // Also poll once on mount
    jobs.list().then((list) => {
      activeIds = new Set(list.filter((j) => j.status === "Running" || j.status === "Pending").map((j) => j.id));
      setCount(activeIds.size);
    }).catch(() => {});

    connection.start().catch(() => {});

    return () => { connection.stop(); };
  }, []);

  return count;
}
