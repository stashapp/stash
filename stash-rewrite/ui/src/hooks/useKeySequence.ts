import { useEffect, useRef, useCallback } from "react";

type KeyBinding = {
  keys: string; // e.g. "g s", "d d", "r 5", "e", "Space"
  action: () => void;
  /** If true, this binding works even when an input/textarea is focused */
  global?: boolean;
};

/**
 * Multi-key sequence keyboard shortcut hook (like vim motions).
 * Supports single keys ("e"), two-key sequences ("g s"), and modifier combos ("Ctrl+Home").
 * 
 * Buffer resets after 800ms of no input.
 */
export function useKeySequence(bindings: KeyBinding[], enabled = true) {
  const bufferRef = useRef<string[]>([]);
  const timerRef = useRef<ReturnType<typeof setTimeout> | undefined>(undefined);

  const clearBuffer = useCallback(() => {
    bufferRef.current = [];
    if (timerRef.current) clearTimeout(timerRef.current);
  }, []);

  useEffect(() => {
    if (!enabled) return;

    const handler = (e: KeyboardEvent) => {
      const tag = (e.target as HTMLElement)?.tagName;
      const inInput = tag === "INPUT" || tag === "TEXTAREA" || tag === "SELECT";

      // Normalize key name
      let key = e.key;
      if (key === " ") key = "Space";

      // Build modifier prefix
      const mods: string[] = [];
      if (e.ctrlKey || e.metaKey) mods.push("Ctrl");
      if (e.altKey) mods.push("Alt");
      if (e.shiftKey && key.length > 1) mods.push("Shift"); // Only for special keys, not letters
      const fullKey = [...mods, key].join("+");

      // Push to buffer
      bufferRef.current.push(fullKey);
      if (timerRef.current) clearTimeout(timerRef.current);

      const bufferStr = bufferRef.current.join(" ");

      // Check for exact match
      const match = bindings.find((b) => {
        if (inInput && !b.global) return false;
        return b.keys === bufferStr;
      });

      if (match) {
        e.preventDefault();
        e.stopPropagation();
        match.action();
        clearBuffer();
        return;
      }

      // Check if buffer could still be a prefix of any binding
      const couldMatch = bindings.some((b) => {
        if (inInput && !b.global) return false;
        return b.keys.startsWith(bufferStr);
      });

      if (couldMatch) {
        timerRef.current = setTimeout(clearBuffer, 800);
      } else {
        clearBuffer();
      }
    };

    window.addEventListener("keydown", handler);
    return () => {
      window.removeEventListener("keydown", handler);
      if (timerRef.current) clearTimeout(timerRef.current);
    };
  }, [bindings, enabled, clearBuffer]);
}
