import { useEffect, useRef } from "react";
import { useLocation } from "react-router-dom";

const NAVIGATION_KEY = "stash_forward_navigation";
const SCROLL_POSITION_KEY_PREFIX = "stash_scroll_position_";

export function useScrollToTopOnMount() {
  const isFirstMount = useRef(true);
  const location = useLocation();

  useEffect(() => {
    // Check if this is a forward navigation (user clicked a link)
    // vs a back navigation (user clicked browser back button)
    const isForwardNavigation = sessionStorage.getItem(NAVIGATION_KEY) === "true";
    
    // Get navigation type from Performance API as additional info
    const navigationEntry = performance.getEntriesByType("navigation")[0] as PerformanceNavigationTiming;
    const navigationType = navigationEntry?.type;
    
    console.debug("[useScrollToTopOnMount] Debug Info:", {
      isFirstMount: isFirstMount.current,
      isForwardNavigation,
      navigationType,
      pathname: location.pathname,
      currentScrollY: window.scrollY,
      willScroll: isFirstMount.current && isForwardNavigation
    });
    
    // Clear the flag immediately after reading
    sessionStorage.removeItem(NAVIGATION_KEY);
    
    // Only scroll to top if:
    // 1. It's the first mount of this component instance
    // 2. AND it's a forward navigation (not back/forward button)
    if (isFirstMount.current && isForwardNavigation) {
      console.debug("[useScrollToTopOnMount] ✅ Scrolling to top (forward navigation)");
      window.scrollTo(0, 0);
    } else {
      console.debug("[useScrollToTopOnMount] ⏭️  Skipping scroll (back navigation or no flag)");
    }
    
    isFirstMount.current = false;
  }, [location.pathname]);
}

// Hook for list pages to preserve scroll position on back navigation
export function useScrollRestoration(loading = false) {
  const location = useLocation();
  const scrollKey = `${SCROLL_POSITION_KEY_PREFIX}${location.pathname}${location.search}`;
  const hasRestoredRef = useRef(false);
  const isRestoringRef = useRef(false);
  const loadingRef = useRef(loading);
  
  // Track loading state changes
  useEffect(() => {
    loadingRef.current = loading;
  }, [loading]);

  useEffect(() => {
    // Don't attempt restoration while page is loading
    if (loading) {
      console.debug("[useScrollRestoration] ⏸️  Waiting for loading to complete");
      return;
    }
    
    // Check if this is a back navigation
    const isForwardNavigation = sessionStorage.getItem(NAVIGATION_KEY) === "true";
    
    // Capture the current pathname at the time of effect setup
    const currentPathname = location.pathname;
    const currentSearch = location.search;
    
    let scrollTimeout: NodeJS.Timeout;
    const handleScroll = () => {
      // Don't save scroll position while we're restoring it
      if (isRestoringRef.current) {
        console.debug("[useScrollRestoration] ⏸️  Ignoring scroll event during restoration");
        return;
      }
      
      // CRITICAL: Only save if we're still on the same page
      // This prevents saving scroll position when navigating away
      if (window.location.pathname !== currentPathname || window.location.search !== currentSearch) {
        console.debug("[useScrollRestoration] ⏭️  Skipping save - page has changed:", {
          expected: currentPathname + currentSearch,
          actual: window.location.pathname + window.location.search
        });
        return;
      }
      
      clearTimeout(scrollTimeout);
      scrollTimeout = setTimeout(() => {
        const currentPosition = window.scrollY;
        
        // Extract page number from URL
        const urlParams = new URLSearchParams(location.search);
        const currentPage = parseInt(urlParams.get("p") || "1", 10);
        
        // Save both scroll position and page number (format: "scrollY|pageNumber")
        const dataToSave = `${currentPosition}|${currentPage}`;
        sessionStorage.setItem(scrollKey, dataToSave);
        console.debug("[useScrollRestoration] 💾 Saved scroll position (on scroll):", {
          pathname: location.pathname,
          position: currentPosition,
          page: currentPage
        });
      }, 100);
    };
    
    if (!isForwardNavigation && !hasRestoredRef.current) {
      // This is a back navigation, try to restore scroll position
      const savedData = sessionStorage.getItem(scrollKey);
      if (savedData) {
        // Parse saved data (format: "scrollY|pageNumber" or just "scrollY" for backward compatibility)
        const parts = savedData.split("|");
        const position = parseInt(parts[0], 10);
        const savedPage = parts[1] ? parseInt(parts[1], 10) : null;
        
        console.debug("[useScrollRestoration] 🔄 Restoring scroll position:", {
          pathname: location.pathname,
          position,
          savedPage,
          currentDocHeight: document.documentElement.scrollHeight,
          currentScrollY: window.scrollY
        });
        
        // Mark that we're restoring to prevent saving during restoration
        isRestoringRef.current = true;
        
        // Function to attempt scroll restoration
        const attemptRestore = (attempts = 0, lastDocHeight = 0, stableCount = 0) => {
          const maxAttempts = 50; // Try for up to 5 seconds (50 * 100ms)
          const docHeight = document.documentElement.scrollHeight;
          const windowHeight = window.innerHeight;
          const maxScroll = docHeight - windowHeight;
          const docHeightStable = docHeight === lastDocHeight;
          
          // Count how many times the height has been stable
          const newStableCount = docHeightStable ? stableCount + 1 : 0;
          
          // Check if loading state changed during restoration attempts
          const isStillLoading = loadingRef.current;
          
          console.debug("[useScrollRestoration] 🔍 Attempt", attempts + 1, {
            targetPosition: position,
            currentDocHeight: docHeight,
            lastDocHeight,
            maxScroll,
            canScroll: maxScroll >= position,
            docHeightStable,
            stableCount: newStableCount,
            isStillLoading
          });
          
          // If loading state changed back to true, abort restoration
          if (isStillLoading) {
            console.debug("[useScrollRestoration] ⏸️  Loading state changed, aborting restoration");
            isRestoringRef.current = false;
            return;
          }
          
          // We should restore if:
          // 1. Document height is stable AND can scroll to target position
          // 2. OR document height has been stable for 5+ attempts (500ms) AND page height is reasonable
          // 3. OR we've reached max attempts
          const canScrollToTarget = maxScroll >= position;
          // Page height should be at least 50% of target position to be considered reasonable
          const pageHeightReasonable = docHeight >= position * 0.5;
          const pageStableEnough = newStableCount >= 5 && pageHeightReasonable; // Height stable for 500ms
          const shouldRestore = (docHeightStable && canScrollToTarget) || pageStableEnough || attempts >= maxAttempts;
          
          if (shouldRestore) {
            // Scroll to target position or max scroll, whichever is smaller
            const targetScroll = Math.min(position, maxScroll);
            window.scrollTo(0, targetScroll);
            
            let reason = "unknown";
            if (docHeightStable && canScrollToTarget) {
              reason = "can scroll to target";
            } else if (pageStableEnough) {
              reason = `page stable (height unchanged for 500ms, height: ${docHeight}px >= ${(position * 0.5).toFixed(0)}px)`;
            } else if (attempts >= maxAttempts) {
              reason = "max attempts reached";
            }
            
            console.debug("[useScrollRestoration] ✅ Scrolled to position:", {
              targetPosition: position,
              actualPosition: window.scrollY,
              maxScroll,
              docHeight,
              attempts: attempts + 1,
              stableCount: newStableCount,
              pageHeightReasonable,
              reason
            });
            
            // Wait a bit before enabling scroll listener
            setTimeout(() => {
              isRestoringRef.current = false;
              console.debug("[useScrollRestoration] ✅ Restoration complete, scroll listener enabled");
            }, 150);
          } else {
            // Page height is still changing, try again
            setTimeout(() => attemptRestore(attempts + 1, docHeight, newStableCount), 100);
          }
        };
        
        // Start attempting to restore after a short delay to let React render
        requestAnimationFrame(() => {
          setTimeout(() => attemptRestore(0, 0, 0), 50);
        });
      } else {
        console.debug("[useScrollRestoration] ℹ️  No saved position found for:", location.pathname);
      }
      hasRestoredRef.current = true;
    } else if (isForwardNavigation) {
      console.debug("[useScrollRestoration] ➡️  Forward navigation, not restoring scroll");
      // DON'T remove NAVIGATION_KEY here - let useScrollToTopOnMount handle it
      // This prevents race condition where list page clears the flag before detail page reads it
    }

    // Save scroll position on scroll (debounced)
    // This is a backup mechanism in case the user scrolls without navigating
    window.addEventListener("scroll", handleScroll);

    return () => {
      window.removeEventListener("scroll", handleScroll);
      clearTimeout(scrollTimeout);
      // DON'T save on unmount - the scroll position is already saved by App.tsx
      // when PUSH action happens, and saving here would overwrite it with 0
      console.debug("[useScrollRestoration] 🧹 Cleanup (not saving on unmount)");
    };
  }, [location.pathname, location.search, scrollKey, loading]);
}

// Helper function to mark forward navigation
// Call this before navigating to a detail page
export function markForwardNavigation() {
  sessionStorage.setItem(NAVIGATION_KEY, "true");
  console.debug("[markForwardNavigation] 🔖 Marked forward navigation");
}
