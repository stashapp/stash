import { useCallback, useState } from "react";

export function useMultiSelect<T extends { id: number }>(items: T[]) {
  const [selectedIds, setSelectedIds] = useState<Set<number>>(new Set());

  const toggle = useCallback((id: number) => {
    setSelectedIds((prev) => {
      const next = new Set(prev);
      if (next.has(id)) {
        next.delete(id);
      } else {
        next.add(id);
      }
      return next;
    });
  }, []);

  const selectAll = useCallback(() => {
    setSelectedIds(new Set(items.map((i) => i.id)));
  }, [items]);

  const selectNone = useCallback(() => {
    setSelectedIds(new Set());
  }, []);

  return { selectedIds, toggle, selectAll, selectNone };
}
