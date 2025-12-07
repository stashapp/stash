/**
 * Facets Extension
 * 
 * Provides dynamic filter counts for sidebar filters across all list pages.
 * 
 * The actual facet counting logic lives in:
 * - src/hooks/useFacetCounts.ts (the hooks)
 * - src/extensions/lists/*.tsx (the list components that use them)
 * 
 * This module provides:
 * - Extension registration
 * - Re-exports of enhanced list components
 */

import { registerExtension } from "../registry";

// Extension ID
export const FACETS_EXTENSION_ID = "facets";

// Register the facets extension
registerExtension({
  id: FACETS_EXTENSION_ID,
  name: "Facets System",
  version: "1.0.0",
  enabled: true,
  initialize: () => {
    console.log("[Facets Extension] Initialized");
  },
});

// Re-export enhanced list components
export * from "./enhanced";
