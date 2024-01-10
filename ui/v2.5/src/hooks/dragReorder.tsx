import { useEffect, useState } from "react";

export function useDragReorder<T>(list: T[], setList: (list: T[]) => void) {
  const [stageList, setStageList] = useState(list);
  const [dragIndex, setDragIndex] = useState<number | undefined>();

  useEffect(() => {
    setStageList(list);
  }, [list]);

  function onDragStart(event: React.DragEvent<HTMLElement>, index: number) {
    event.dataTransfer.effectAllowed = "move";
    setDragIndex(index);
  }

  function onDragOver(event: React.DragEvent<HTMLElement>, index?: number) {
    if (dragIndex !== undefined && index !== undefined && index !== dragIndex) {
      const newList = [...stageList];
      const moved = newList.splice(dragIndex, 1);
      newList.splice(index, 0, moved[0]);
      setStageList(newList);
      setDragIndex(index);
    }

    event.dataTransfer.dropEffect = "move";
    event.preventDefault();
  }

  function onDragOverDefault(event: React.DragEvent<HTMLDivElement>) {
    event.dataTransfer.dropEffect = "move";
    event.preventDefault();
  }

  function onDrop() {
    // assume we've already set the temp source list
    // feed it up
    setList(stageList);
    setDragIndex(undefined);
  }

  return {
    stageList,
    onDragStart,
    onDragOver,
    onDragOverDefault,
    onDrop,
  };
}
