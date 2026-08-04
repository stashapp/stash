import { useEffect, useLayoutEffect, useRef } from "react";
import { useHistory, useLocation } from "react-router-dom";

// the initial history entry has no key
const initialKey = "@@initial";

// scroll positions for each history entry, keyed by location key. Kept in
// session storage so that the positions survive a page reload - the history
// entries do.
const storageKey = "stash-scroll-positions";

function loadPositions() {
  try {
    const stored = sessionStorage.getItem(storageKey);
    if (stored) return new Map<string, number>(JSON.parse(stored));
  } catch {
    // ignore - unparseable or unavailable storage
  }
  return new Map<string, number>();
}

const positions = loadPositions();

// the number of history entries to remember positions for
const maxPositions = 50;

function setPosition(key: string, position: number) {
  positions.set(key, position);

  while (positions.size > maxPositions) {
    positions.delete(positions.keys().next().value!);
  }
}

// positions are updated as the user scrolls, so they are only written out when
// navigating or leaving the page
function savePositions() {
  try {
    sessionStorage.setItem(storageKey, JSON.stringify([...positions]));
  } catch {
    // ignore - storage may be unavailable or full
  }
}

// how long to keep trying to restore the scroll position while the page
// content is still loading in
const restoreTimeout = 3000;

// how long to keep the restored position applied - the list components scroll
// the page around as their contents settle
const settleTime = 750;

let restoreCount = 0;

// true while a scroll position is being restored. Automatic scrolling (such as
// scrolling to the top of a list when its page changes) should be suppressed
// while this is the case.
export function isRestoringScroll() {
  return restoreCount > 0;
}

// Scrolls to the given position, keeping it applied while the page settles.
//
// The position can't be applied while the document is still too short - which
// is the case while the page content is being fetched and rendered - and once
// it can be applied, the content rendering in may scroll it away again. Keep
// it applied until the page has settled, the user scrolls themselves, or we
// give up.
//
// Returns a function to abort the restoration.
export function scrollToWhenReady(target: number) {
  const start = Date.now();
  let reached: number | undefined;
  let frame = 0;
  let stopped = false;

  restoreCount++;

  function stop() {
    if (stopped) return;
    stopped = true;
    restoreCount--;
    cancelAnimationFrame(frame);
    window.removeEventListener("wheel", stop);
    window.removeEventListener("touchstart", stop);
    window.removeEventListener("keydown", stop);
  }

  function attempt() {
    const now = Date.now();

    if (Math.abs(window.scrollY - target) > 1) {
      window.scrollTo(0, target);
    }

    if (Math.abs(window.scrollY - target) <= 1) {
      // applied - keep it that way until the page has settled
      if (reached === undefined) reached = now;
      if (now - reached > settleTime) {
        stop();
        return;
      }
    } else if (now - start > restoreTimeout) {
      // the page never got tall enough - give up
      stop();
      return;
    }

    frame = requestAnimationFrame(attempt);
  }

  // abandon the restoration if the user scrolls in the meantime
  window.addEventListener("wheel", stop, { passive: true });
  window.addEventListener("touchstart", stop, { passive: true });
  window.addEventListener("keydown", stop);

  attempt();

  return stop;
}

// Restores the scroll position when navigating back or forward.
//
// The browser's native restoration doesn't work for the list views: at the
// point the history entry is popped, the page content hasn't been fetched and
// rendered yet, so the document is too short and the position is clamped to
// the top. Instead we record the position for each history entry ourselves and
// re-apply it once the page is able to scroll there.
export function useScrollRestoration() {
  const location = useLocation();
  const history = useHistory();

  const key = location.key ?? initialKey;
  const keyRef = useRef(key);
  const abortRestore = useRef<() => void>();

  // take over scroll restoration from the browser
  useEffect(() => {
    if (!("scrollRestoration" in window.history)) return;

    const original = window.history.scrollRestoration;
    window.history.scrollRestoration = "manual";
    return () => {
      window.history.scrollRestoration = original;
    };
  }, []);

  // record the scroll position of the current history entry
  useEffect(() => {
    function onScroll() {
      // ignore the scroll events generated while restoring - they would
      // overwrite the position being restored to
      if (isRestoringScroll()) return;

      setPosition(keyRef.current, window.scrollY);
    }

    window.addEventListener("scroll", onScroll, { passive: true });
    window.addEventListener("pagehide", savePositions);
    return () => {
      window.removeEventListener("scroll", onScroll);
      window.removeEventListener("pagehide", savePositions);
    };
  }, []);

  useEffect(() => {
    return () => abortRestore.current?.();
  }, []);

  // biome-ignore lint/correctness/useExhaustiveDependencies: intentionally only running when the history entry changes
  useLayoutEffect(() => {
    const prevKey = keyRef.current;
    keyRef.current = key;

    if (history.action === "REPLACE") {
      // the same history entry under a new key - the lists replace the URL to
      // keep it in sync with their filter, and do so while the position is
      // still being restored. Carry the position over and leave the
      // restoration running.
      const position = positions.get(prevKey);
      if (position !== undefined) {
        setPosition(key, position);
      }
      savePositions();
      return;
    }

    abortRestore.current?.();
    savePositions();

    if (history.action !== "POP") {
      // new entry - it starts wherever the browser leaves us
      setPosition(key, window.scrollY);
      return;
    }

    const target = positions.get(key) ?? 0;
    if (target === 0 && window.scrollY === 0) return;

    abortRestore.current = scrollToWhenReady(target);
  }, [key]);
}
