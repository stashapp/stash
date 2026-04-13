import { createContext, useContext, useEffect, useMemo, type ReactNode } from "react";
import { useQuery } from "@tanstack/react-query";
import { system } from "../api/client";
import type { StashConfig, SystemStatus } from "../api/types";

const defaultMenuItems = ["scenes", "images", "markers", "performers", "galleries", "studios", "tags", "groups"];

function normalizeRatingSystemType(value: string | undefined) {
  return value?.toLowerCase() === "decimal" ? "decimal" : "stars";
}

function normalizeRatingStarPrecision(value: string | undefined) {
  switch (value?.toLowerCase()) {
    case "half":
      return "half";
    case "quarter":
      return "quarter";
    case "tenth":
      return "tenth";
    default:
      return "full";
  }
}

function normalizeConfig(config: StashConfig): StashConfig {
  return {
    ...config,
    interface: {
      ...config.interface,
      menuItems: config.interface.menuItems.length > 0 ? config.interface.menuItems : defaultMenuItems,
    },
    ui: {
      ...config.ui,
      ratingSystemOptions: {
        type: normalizeRatingSystemType(config.ui.ratingSystemOptions.type),
        starPrecision: normalizeRatingStarPrecision(config.ui.ratingSystemOptions.starPrecision),
      },
    },
    scraping: {
      ...config.scraping,
      stashBoxes: config.scraping.stashBoxes ?? [],
    },
  };
}

interface AppConfigContextValue {
  config?: StashConfig;
  status?: SystemStatus;
  configLoading: boolean;
  statusLoading: boolean;
}

const AppConfigContext = createContext<AppConfigContextValue | null>(null);

export function AppConfigProvider({ children }: { children: ReactNode }) {
  const configQuery = useQuery({
    queryKey: ["system-config"],
    queryFn: system.getConfig,
  });

  const config = useMemo(() => {
    if (!configQuery.data) {
      return undefined;
    }

    return normalizeConfig(configQuery.data);
  }, [configQuery.data]);

  const statusQuery = useQuery({
    queryKey: ["system-status"],
    queryFn: system.status,
  });

  useEffect(() => {
    document.title = config?.ui.title?.trim() || "Stash";
  }, [config?.ui.title]);

  useEffect(() => {
    document.documentElement.lang = config?.interface.language || "en-US";
  }, [config?.interface.language]);

  return (
    <AppConfigContext.Provider
      value={{
        config,
        status: statusQuery.data,
        configLoading: configQuery.isLoading,
        statusLoading: statusQuery.isLoading,
      }}
    >
      {children}
    </AppConfigContext.Provider>
  );
}

export function useAppConfig() {
  const context = useContext(AppConfigContext);
  if (!context) {
    throw new Error("useAppConfig must be used within an AppConfigProvider");
  }

  return context;
}