import { createContext, useContext, useState, useCallback, type ReactNode } from "react";

interface SceneQueueState {
  sceneIds: number[];
  currentIndex: number;
  autoplay: boolean;
}

interface SceneQueueContextValue {
  queue: SceneQueueState | null;
  setQueue: (ids: number[], currentId: number) => void;
  clearQueue: () => void;
  currentId: number | null;
  prevId: number | null;
  nextId: number | null;
  hasPrev: boolean;
  hasNext: boolean;
  goToIndex: (index: number) => number | null;
  toggleAutoplay: () => void;
  autoplay: boolean;
  queueLength: number;
  currentPosition: number;
}

const SceneQueueContext = createContext<SceneQueueContextValue | null>(null);

export function SceneQueueProvider({ children }: { children: ReactNode }) {
  const [queue, setQueueState] = useState<SceneQueueState | null>(null);

  const setQueue = useCallback((ids: number[], currentId: number) => {
    const idx = ids.indexOf(currentId);
    setQueueState({ sceneIds: ids, currentIndex: idx >= 0 ? idx : 0, autoplay: false });
  }, []);

  const clearQueue = useCallback(() => setQueueState(null), []);

  const currentId = queue ? queue.sceneIds[queue.currentIndex] ?? null : null;
  const prevId = queue && queue.currentIndex > 0 ? queue.sceneIds[queue.currentIndex - 1] : null;
  const nextId = queue && queue.currentIndex < queue.sceneIds.length - 1 ? queue.sceneIds[queue.currentIndex + 1] : null;

  const goToIndex = useCallback((index: number) => {
    if (!queue || index < 0 || index >= queue.sceneIds.length) return null;
    const id = queue.sceneIds[index];
    setQueueState({ ...queue, currentIndex: index });
    return id;
  }, [queue]);

  const toggleAutoplay = useCallback(() => {
    setQueueState((prev) => prev ? { ...prev, autoplay: !prev.autoplay } : null);
  }, []);

  return (
    <SceneQueueContext.Provider
      value={{
        queue,
        setQueue,
        clearQueue,
        currentId,
        prevId,
        nextId,
        hasPrev: prevId !== null,
        hasNext: nextId !== null,
        goToIndex,
        toggleAutoplay,
        autoplay: queue?.autoplay ?? false,
        queueLength: queue?.sceneIds.length ?? 0,
        currentPosition: queue ? queue.currentIndex + 1 : 0,
      }}
    >
      {children}
    </SceneQueueContext.Provider>
  );
}

export function useSceneQueue() {
  const ctx = useContext(SceneQueueContext);
  if (!ctx) throw new Error("useSceneQueue must be used within SceneQueueProvider");
  return ctx;
}
