import { useState } from "react";
import { useListContextOptional } from "src/components/List/ListProvider";

export function useDragMoveSelect(props: {
  selecting: boolean;
  selected: boolean;
  onSelectedChanged?: (selected: boolean, shiftKey: boolean) => void;
  objectId?: string;
  onMove?: (srcIds: string[], targetId: string) => void;
}) {
  const { selectedIds } = useListContextOptional();

  const [inHandle, setInHandle] = useState(false);
  const [moveSrc, setMoveSrc] = useState(false);
  const [moveTarget, setMoveTarget] = useState(false);

  const canSelect = props.onSelectedChanged && props.selecting;
  const canMove = !!props.objectId && props.onMove && inHandle;
  const draggable = canSelect || canMove;

  function onDragStart(event: React.DragEvent<HTMLElement>) {
    if (!draggable) {
      event.preventDefault();
      return;
    }

    if (!inHandle && props.selecting) {
      event.dataTransfer.setData("text/plain", "");
      // event.dataTransfer.setDragImage(new Image(), 0, 0);
      event.dataTransfer.effectAllowed = "copy";
      event.stopPropagation();
    } else if (inHandle && props.objectId) {
      if (selectedIds.size > 1 && selectedIds.has(props.objectId)) {
        // moving all selected
        const movingIds = Array.from(selectedIds.values()).join(",");
        event.dataTransfer.setData("text/plain", movingIds);
      } else {
        // moving single
        setMoveSrc(true);
        event.dataTransfer.setData("text/plain", props.objectId);
      }
      event.dataTransfer.effectAllowed = "move";
      event.stopPropagation();
    }
  }

  function onDragEnter(event: React.DragEvent<HTMLElement>) {
    const ev = event;
    const shiftKey = false;

    if (ev.dataTransfer.effectAllowed === "copy") {
      if (!props.onSelectedChanged) {
        return;
      }

      if (props.selecting && !props.selected) {
        props.onSelectedChanged(true, shiftKey);
      }

      ev.dataTransfer.dropEffect = "copy";
      ev.preventDefault();
    } else if (ev.dataTransfer.effectAllowed === "move" && !moveSrc) {
      setMoveTarget(true);
      ev.dataTransfer.dropEffect = "move";
      ev.preventDefault();
    } else {
      ev.dataTransfer.dropEffect = "none";
    }
  }

  function onDragLeave(event: React.DragEvent<HTMLElement>) {
    if (event.currentTarget.contains(event.relatedTarget as Node)) {
      return;
    }

    setMoveTarget(false);
  }

  function onDragOver(event: React.DragEvent<HTMLElement>) {
    if (event.dataTransfer.effectAllowed === "move" && moveSrc) {
      return;
    }
    event.preventDefault();
  }

  function onDragEnd() {
    setMoveTarget(false);
    setMoveSrc(false);
  }

  function onDrop(event: React.DragEvent<HTMLElement>) {
    const ev = event;

    if (ev.dataTransfer.effectAllowed === "copy") {
      return;
    }

    // TODO: move logic

    onDragEnd();
  }

  return {
    inHandle,
    setInHandle,
    moveTarget,
    dragProps: {
      draggable: draggable || undefined,
      onDragStart,
      onDragEnter,
      onDragLeave,
      onDragOver,
      onDragEnd,
      onDrop,
    },
  };
}
