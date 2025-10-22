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
    
   
    // Clear the flag immediately after reading
    sessionStorage.removeItem(NAVIGATION_KEY);
    
    // Only scroll to top if:
    // 1. It's the first mount of this component instance
    // 2. AND it's a forward navigation (not back/forward button)
    if (isFirstMount.current && isForwardNavigation) {
      window.scrollTo(0, 0);
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
        return;
      }
      
      // CRITICAL: Only save if we're still on the same page
      // This prevents saving scroll position when navigating away
      if (window.location.pathname !== currentPathname || window.location.search !== currentSearch) {
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
        
          // If loading state changed back to true, abort restoration
          if (isStillLoading) {
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
          
            // Wait a bit before enabling scroll listener
            setTimeout(() => {
              isRestoringRef.current = false;
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
      }

      hasRestoredRef.current = true;
    }

    // Save scroll position on scroll (debounced)
    // This is a backup mechanism in case the user scrolls without navigating
    window.addEventListener("scroll", handleScroll);

    return () => {
      window.removeEventListener("scroll", handleScroll);
      clearTimeout(scrollTimeout);
    };
  }, [location.pathname, location.search, scrollKey, loading]);
}

// Helper function to mark forward navigation
// Call this before navigating to a detail page
export function markForwardNavigation() {
  sessionStorage.setItem(NAVIGATION_KEY, "true");
}
