import { useState } from "react";
import { useListContextOptional } from "src/components/List/ListProvider";

// Enum representing the possible sides for a drag operation.
export enum DragSide {
  BEFORE,
  AFTER,
}

/**
 * Hook to manage drag and move selection functionality.
 * Dragging while selecting will allow the user to select multiple items.
 * Dragging from the drag handle will allow the user to move the item or selected items.
 *
 * @param props - The properties for the hook.
 * @param props.selecting - Whether the one or more items have been selected.
 * @param props.selected - Whether this item is currently selected.
 * @param props.onSelectedChanged - Callback when the selected state changes.
 * @param props.objectId - The ID of this object.
 * @param props.onMove - Callback when a move operation occurs.
 *
 * @returns An object containing the drag event handlers and state.
 */
export function useDragMoveSelect(props: {
  selecting: boolean;
  selected: boolean;
  onSelectedChanged?: (selected: boolean, shiftKey: boolean) => void;
  objectId?: string;
  onMove?: (srcIds: string[], targetId: string, after: boolean) => void;
}) {
  const { selectedIds } = useListContextOptional();

  // true if the mouse is over the drag handle
  const [inHandle, setInHandle] = useState(false);

  // true if this is the source of a move operation
  const [moveSrc, setMoveSrc] = useState(false);
  // the target side for a move operation
  const [moveTarget, setMoveTarget] = useState<DragSide | undefined>();

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

  function doSetMoveTarget(event: React.DragEvent<HTMLElement>) {
    const isBefore =
      event.nativeEvent.offsetX < event.currentTarget.clientWidth / 2;
    if (isBefore && moveTarget !== DragSide.BEFORE) {
      setMoveTarget(DragSide.BEFORE);
    } else if (!isBefore && moveTarget !== DragSide.AFTER) {
      setMoveTarget(DragSide.AFTER);
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
      // don't allow move on self
      doSetMoveTarget(event);
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

    setMoveTarget(undefined);
  }

  function onDragOver(event: React.DragEvent<HTMLElement>) {
    // only set move target if move is allowed, or if this is not the source of the move
    if (event.dataTransfer.effectAllowed !== "move" || moveSrc) {
      return;
    }

    doSetMoveTarget(event);

    event.preventDefault();
  }

  function onDragEnd() {
    setMoveTarget(undefined);
    setMoveSrc(false);
  }

  function onDrop(event: React.DragEvent<HTMLElement>) {
    const ev = event;

    if (
      ev.dataTransfer.effectAllowed === "copy" ||
      !props.onMove ||
      !props.objectId
    ) {
      return;
    }

    const srcIds = ev.dataTransfer.getData("text/plain").split(",");
    const targetId = props.objectId;
    const after = moveTarget === DragSide.AFTER;

    props.onMove(srcIds, targetId, after);

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
