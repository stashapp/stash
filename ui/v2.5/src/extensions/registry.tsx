/**
 * Extension Registry
 * 
 * Manages registration and initialization of fork extensions.
 * This provides a clean way to add custom functionality without
 * modifying upstream components.
 */

import React, { createContext, useContext, ReactNode, useMemo } from "react";

/**
 * Extension definition
 */
export interface Extension {
  /** Unique identifier for the extension */
  id: string;
  /** Human-readable name */
  name: string;
  /** Extension version */
  version: string;
  /** Whether the extension is enabled */
  enabled: boolean;
  /** Optional provider component to wrap the app */
  Provider?: React.ComponentType<{ children: ReactNode }>;
  /** Optional initialization function */
  initialize?: () => void | Promise<void>;
}

/**
 * Extension registry state
 */
interface ExtensionRegistryState {
  extensions: Map<string, Extension>;
  isInitialized: boolean;
}

/**
 * Extension registry context value
 */
interface ExtensionRegistryContextValue {
  /** Get an extension by ID */
  getExtension: (id: string) => Extension | undefined;
  /** Check if an extension is enabled */
  isEnabled: (id: string) => boolean;
  /** Get all registered extensions */
  getAllExtensions: () => Extension[];
}

const ExtensionRegistryContext = createContext<ExtensionRegistryContextValue | null>(null);

// Global registry instance
const registry: ExtensionRegistryState = {
  extensions: new Map(),
  isInitialized: false,
};

/**
 * Register an extension
 */
export function registerExtension(extension: Extension): void {
  if (registry.isInitialized) {
    console.warn(
      `Extension "${extension.id}" registered after initialization. ` +
      `This extension may not be fully initialized.`
    );
  }
  registry.extensions.set(extension.id, extension);
}

/**
 * Get all registered extensions
 */
export function getRegisteredExtensions(): Extension[] {
  return Array.from(registry.extensions.values());
}

/**
 * Initialize all extensions
 */
async function initializeExtensions(): Promise<void> {
  if (registry.isInitialized) return;

  const extensions = Array.from(registry.extensions.values());
  
  for (const ext of extensions) {
    if (ext.enabled && ext.initialize) {
      try {
        await ext.initialize();
      } catch (error) {
        console.error(`Failed to initialize extension "${ext.id}":`, error);
      }
    }
  }

  registry.isInitialized = true;
}

/**
 * Compose all extension providers
 */
function composeProviders(
  providers: React.ComponentType<{ children: ReactNode }>[]
): React.ComponentType<{ children: ReactNode }> {
  return providers.reduce(
    (Composed, Provider) => {
      return function ComposedProvider({ children }: { children: ReactNode }) {
        return (
          <Composed>
            <Provider>{children}</Provider>
          </Composed>
        );
      };
    },
    ({ children }: { children: ReactNode }) => <>{children}</>
  );
}

interface ExtensionRegistryProviderProps {
  children: ReactNode;
}

/**
 * Extension Registry Provider
 * 
 * Wraps the application with all extension providers and provides
 * access to the extension registry.
 */
export function ExtensionRegistryProvider({ children }: ExtensionRegistryProviderProps) {
  // Initialize extensions on first render
  React.useEffect(() => {
    initializeExtensions();
  }, []);

  // Collect all extension providers
  const ExtensionProviders = useMemo(() => {
    const providers = Array.from(registry.extensions.values())
      .filter((ext) => ext.enabled && ext.Provider)
      .map((ext) => ext.Provider!);
    
    return composeProviders(providers);
  }, []);

  const contextValue = useMemo<ExtensionRegistryContextValue>(
    () => ({
      getExtension: (id) => registry.extensions.get(id),
      isEnabled: (id) => registry.extensions.get(id)?.enabled ?? false,
      getAllExtensions: () => Array.from(registry.extensions.values()),
    }),
    []
  );

  return (
    <ExtensionRegistryContext.Provider value={contextValue}>
      <ExtensionProviders>{children}</ExtensionProviders>
    </ExtensionRegistryContext.Provider>
  );
}

/**
 * Hook to access the extension registry
 */
export function useExtensionRegistry(): ExtensionRegistryContextValue {
  const context = useContext(ExtensionRegistryContext);
  if (!context) {
    throw new Error("useExtensionRegistry must be used within ExtensionRegistryProvider");
  }
  return context;
}

/**
 * Hook to check if a specific extension is available and enabled
 */
export function useExtension(id: string): Extension | undefined {
  const { getExtension, isEnabled } = useExtensionRegistry();
  const extension = getExtension(id);
  return extension && isEnabled(id) ? extension : undefined;
}

