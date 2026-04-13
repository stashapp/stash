import { useState, useCallback, useRef, useEffect } from "react";
import { createPortal } from "react-dom";
import {
  X,
  ChevronLeft,
  ChevronRight,
  ChevronUp,
  ChevronDown,
  Trash2,
  Shuffle,
  Repeat,
  List,
  SkipBack,
  SkipForward,
} from "lucide-react";
import { scenes } from "../api/client";

export interface SceneQueueProps {
  scenes: { id: number; title?: string; duration?: number; screenshotUrl?: string }[];
  initialIndex?: number;
  onClose: () => void;
  onNavigate: (route: any) => void;
}

function formatQueueDuration(seconds?: number): string {
  if (!seconds || seconds <= 0) return "--:--";
  const h = Math.floor(seconds / 3600);
  const m = Math.floor((seconds % 3600) / 60);
  const s = Math.floor(seconds % 60);
  if (h > 0) return `${h}:${m.toString().padStart(2, "0")}:${s.toString().padStart(2, "0")}`;
  return `${m}:${s.toString().padStart(2, "0")}`;
}

export function SceneQueue({
  scenes: initialScenes,
  initialIndex = 0,
  onClose,
  onNavigate,
}: SceneQueueProps) {
  const [queue, setQueue] = useState(initialScenes);
  const [currentIndex, setCurrentIndex] = useState(
    Math.min(initialIndex, initialScenes.length - 1),
  );
  const [queueOpen, setQueueOpen] = useState(true);
  const [repeat, setRepeat] = useState(false);
  const videoRef = useRef<HTMLVideoElement>(null);
  const queueItemRefs = useRef<Map<number, HTMLDivElement>>(new Map());

  const current = queue[currentIndex];

  // Scroll current item into view
  useEffect(() => {
    if (current) {
      const el = queueItemRefs.current.get(current.id);
      el?.scrollIntoView({ block: "nearest", behavior: "smooth" });
    }
  }, [currentIndex, current]);

  const goTo = useCallback(
    (idx: number) => {
      if (idx < 0 || idx >= queue.length) {
        if (repeat && queue.length > 0) {
          setCurrentIndex(0);
        }
        return;
      }
      setCurrentIndex(idx);
    },
    [queue.length, repeat],
  );

  const goPrev = useCallback(() => goTo(currentIndex - 1), [goTo, currentIndex]);
  const goNext = useCallback(() => goTo(currentIndex + 1), [goTo, currentIndex]);

  const handleEnded = useCallback(() => {
    if (currentIndex < queue.length - 1) {
      goTo(currentIndex + 1);
    } else if (repeat) {
      goTo(0);
    }
  }, [currentIndex, queue.length, repeat, goTo]);

  const removeFromQueue = useCallback(
    (idx: number) => {
      setQueue((prev) => {
        const next = prev.filter((_, i) => i !== idx);
        if (next.length === 0) {
          onClose();
          return prev;
        }
        return next;
      });
      if (idx < currentIndex) {
        setCurrentIndex((i) => i - 1);
      } else if (idx === currentIndex) {
        setCurrentIndex((i) => Math.min(i, queue.length - 2));
      }
    },
    [currentIndex, queue.length, onClose],
  );

  const moveUp = useCallback(
    (idx: number) => {
      if (idx <= 0) return;
      setQueue((prev) => {
        const next = [...prev];
        [next[idx - 1], next[idx]] = [next[idx], next[idx - 1]];
        return next;
      });
      if (idx === currentIndex) setCurrentIndex(idx - 1);
      else if (idx - 1 === currentIndex) setCurrentIndex(idx);
    },
    [currentIndex],
  );

  const moveDown = useCallback(
    (idx: number) => {
      if (idx >= queue.length - 1) return;
      setQueue((prev) => {
        const next = [...prev];
        [next[idx], next[idx + 1]] = [next[idx + 1], next[idx]];
        return next;
      });
      if (idx === currentIndex) setCurrentIndex(idx + 1);
      else if (idx + 1 === currentIndex) setCurrentIndex(idx);
    },
    [currentIndex, queue.length],
  );

  const shuffleQueue = useCallback(() => {
    setQueue((prev) => {
      const currentScene = prev[currentIndex];
      const rest = prev.filter((_, i) => i !== currentIndex);
      for (let i = rest.length - 1; i > 0; i--) {
        const j = Math.floor(Math.random() * (i + 1));
        [rest[i], rest[j]] = [rest[j], rest[i]];
      }
      return [currentScene, ...rest];
    });
    setCurrentIndex(0);
  }, [currentIndex]);

  // Keyboard shortcuts
  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      if (e.key === "Escape") {
        e.preventDefault();
        onClose();
      }
    };
    window.addEventListener("keydown", handler);
    return () => window.removeEventListener("keydown", handler);
  }, [onClose]);

  if (!current) return <></>;

  const streamUrl = scenes.streamUrl(current.id);

  return createPortal(
    <div className="fixed inset-0 z-50 flex bg-black">
      {/* Main player area */}
      <div className={`flex flex-col flex-1 min-w-0 ${queueOpen ? "mr-80" : ""} transition-all duration-200`}>
        {/* Top bar */}
        <div className="flex items-center justify-between px-4 py-2 bg-black/80">
          <div className="flex items-center gap-3 min-w-0">
            <button onClick={onClose} className="p-1.5 text-white/70 hover:text-white rounded hover:bg-white/10">
              <X size={18} />
            </button>
            <span className="text-white text-sm font-medium truncate">
              {current.title || "Untitled"}{" "}
              <span className="text-white/50">
                ({currentIndex + 1}/{queue.length})
              </span>
            </span>
          </div>
          <div className="flex items-center gap-1">
            <button
              onClick={shuffleQueue}
              className="p-1.5 text-white/70 hover:text-white rounded hover:bg-white/10"
              title="Shuffle queue"
            >
              <Shuffle size={16} />
            </button>
            <button
              onClick={() => setRepeat((r) => !r)}
              className={`p-1.5 rounded hover:bg-white/10 ${repeat ? "text-plex-accent" : "text-white/70 hover:text-white"}`}
              title={repeat ? "Repeat on" : "Repeat off"}
            >
              <Repeat size={16} />
            </button>
            <button
              onClick={() => setQueueOpen((o) => !o)}
              className="p-1.5 text-white/70 hover:text-white rounded hover:bg-white/10"
              title={queueOpen ? "Hide queue" : "Show queue"}
            >
              <List size={16} />
            </button>
          </div>
        </div>

        {/* Video player */}
        <div className="flex-1 flex items-center justify-center bg-black">
          <video
            ref={videoRef}
            key={streamUrl}
            controls
            autoPlay
            className="max-h-full max-w-full"
            onEnded={handleEnded}
          >
            <source src={streamUrl} type="video/mp4" />
          </video>
        </div>

        {/* Transport bar */}
        <div className="flex items-center justify-center gap-4 px-4 py-2 bg-black/80">
          <button
            onClick={goPrev}
            disabled={currentIndex === 0 && !repeat}
            className="p-2 text-white/70 hover:text-white disabled:text-white/20 rounded hover:bg-white/10 disabled:hover:bg-transparent"
            title="Previous scene"
          >
            <SkipBack size={20} />
          </button>
          <button
            onClick={goNext}
            disabled={currentIndex >= queue.length - 1 && !repeat}
            className="p-2 text-white/70 hover:text-white disabled:text-white/20 rounded hover:bg-white/10 disabled:hover:bg-transparent"
            title="Next scene"
          >
            <SkipForward size={20} />
          </button>
        </div>
      </div>

      {/* Queue sidebar */}
      {queueOpen && (
        <div className="fixed right-0 top-0 bottom-0 w-80 flex flex-col bg-plex-bg border-l border-plex-border">
          <div className="px-4 py-3 border-b border-plex-border flex items-center justify-between">
            <h2 className="text-sm font-semibold text-plex-text">Queue ({queue.length})</h2>
            <button
              onClick={() => setQueueOpen(false)}
              className="p-1 text-plex-text-muted hover:text-plex-text rounded hover:bg-plex-border/50"
            >
              <X size={14} />
            </button>
          </div>
          <div className="flex-1 overflow-y-auto">
            {queue.map((scene, idx) => (
              <div
                key={`${scene.id}-${idx}`}
                ref={(el) => {
                  if (el) queueItemRefs.current.set(scene.id, el);
                }}
                className={`flex items-center gap-2 px-3 py-2 cursor-pointer border-l-2 transition-colors ${
                  idx === currentIndex
                    ? "border-l-plex-accent bg-plex-accent/10"
                    : "border-l-transparent hover:bg-plex-surface"
                }`}
                onClick={() => goTo(idx)}
              >
                {/* Thumbnail */}
                <div className="w-16 h-9 flex-shrink-0 rounded overflow-hidden bg-plex-surface">
                  {scene.screenshotUrl ? (
                    <img
                      src={scene.screenshotUrl}
                      alt=""
                      className="w-full h-full object-cover"
                      loading="lazy"
                    />
                  ) : (
                    <div className="w-full h-full flex items-center justify-center text-plex-text-muted text-[10px]">
                      No thumb
                    </div>
                  )}
                </div>

                {/* Info */}
                <div className="flex-1 min-w-0">
                  <p className="text-xs text-plex-text truncate">
                    {scene.title || "Untitled"}
                  </p>
                  <p className="text-[10px] text-plex-text-muted">
                    {formatQueueDuration(scene.duration)}
                  </p>
                </div>

                {/* Actions */}
                <div className="flex flex-col gap-0.5 flex-shrink-0" onClick={(e) => e.stopPropagation()}>
                  <button
                    onClick={() => moveUp(idx)}
                    disabled={idx === 0}
                    className="p-0.5 text-plex-text-muted hover:text-plex-text disabled:opacity-25"
                    title="Move up"
                  >
                    <ChevronUp size={12} />
                  </button>
                  <button
                    onClick={() => moveDown(idx)}
                    disabled={idx === queue.length - 1}
                    className="p-0.5 text-plex-text-muted hover:text-plex-text disabled:opacity-25"
                    title="Move down"
                  >
                    <ChevronDown size={12} />
                  </button>
                </div>
                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    removeFromQueue(idx);
                  }}
                  className="p-1 text-plex-text-muted hover:text-red-400 flex-shrink-0"
                  title="Remove from queue"
                >
                  <Trash2 size={12} />
                </button>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>,
    document.body,
  );
}
