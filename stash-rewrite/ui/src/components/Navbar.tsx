import { Film, Users, Building2, Tags, Image, ImageIcon, Layers, Settings, BarChart3, Activity, Clapperboard, Bookmark, HelpCircle } from "lucide-react";
import { useState } from "react";
import { JobDrawer, useJobCount } from "./JobDrawer";
import { GlobalSearch } from "./GlobalSearch";
import { useRouteRegistry } from "../router/RouteRegistry";
import { useAppConfig } from "../state/AppConfigContext";
import { KeyboardShortcutsDialog } from "./KeyboardShortcutsDialog";

interface NavbarProps {
  currentPage: string;
  navigate: (r: any) => void;
}

const navItems = [
  { page: "scenes", label: "Scenes", icon: Film },
  { page: "images", label: "Images", icon: ImageIcon },
  { page: "markers", label: "Markers", icon: Bookmark },
  { page: "galleries", label: "Galleries", icon: Image },
  { page: "performers", label: "Performers", icon: Users },
  { page: "studios", label: "Studios", icon: Building2 },
  { page: "tags", label: "Tags", icon: Tags },
  { page: "groups", label: "Groups", icon: Layers },
];

const DETAIL_PARENT_PAGE: Record<string, string> = {
  scene: "scenes",
  image: "images",
  performer: "performers",
  gallery: "galleries",
  studio: "studios",
  tag: "tags",
  group: "groups",
};

export function Navbar({ currentPage, navigate }: NavbarProps) {
  const [jobDrawerOpen, setJobDrawerOpen] = useState(false);
  const [helpOpen, setHelpOpen] = useState(false);
  const jobCount = useJobCount();
  const { routes } = useRouteRegistry();
  const { config } = useAppConfig();
  const activePage = DETAIL_PARENT_PAGE[currentPage] ?? currentPage;

  const enabledMenuItems = config?.interface.menuItems.length
    ? new Set(config.interface.menuItems)
    : null;

  const extensionNavItems = routes
    .filter((r) => r.navItem)
    .map((r) => ({ page: r.navItem!.page, label: r.navItem!.label, icon: r.navItem!.icon, order: r.navItem!.order ?? 99 }))
    .sort((a, b) => a.order - b.order);

  const allNavItems = [...navItems.filter((item) => !enabledMenuItems || enabledMenuItems.has(item.page)), ...extensionNavItems];

  return (
    <nav className="stash-navbar bg-plex-nav sticky top-0 z-50 shadow-lg shadow-black/30">
      <div className="max-w-[1800px] mx-auto px-4">
        <div className="flex items-center h-12">
          {/* Logo */}
          <button
            onClick={() => navigate({ page: "home" })}
            className="flex items-center gap-2 mr-6 shrink-0"
          >
            <Clapperboard className="w-5 h-5 text-plex-accent" />
            <span className="text-base font-bold uppercase tracking-wider text-plex-text hover:text-plex-accent transition-colors">
              Stash
            </span>
          </button>

          {/* Nav items */}
          <div className="flex items-center gap-0.5 overflow-x-auto">
            {allNavItems.map(({ page, label, icon: Icon }) => (
              <button
                key={page}
                onClick={() => navigate({ page })}
                className={`flex items-center gap-1.5 px-3 py-1.5 rounded text-sm font-medium whitespace-nowrap transition-colors ${
                  activePage === page
                    ? "text-plex-accent"
                    : "text-plex-text-secondary hover:text-plex-text"
                }`}
              >
                {Icon && <Icon className="w-4 h-4" />}
                {label}
              </button>
            ))}
          </div>

          {/* Spacer */}
          <div className="flex-1" />

          <GlobalSearch navigate={navigate} />

          {/* Actions */}
          <div className="flex items-center gap-1 shrink-0">
            {jobCount > 0 && (
              <button
                onClick={() => setJobDrawerOpen(true)}
                className="relative p-2 rounded text-plex-text-secondary hover:text-plex-text"
                title="Jobs"
              >
                <Activity className="w-[18px] h-[18px]" />
                <span className="absolute -top-0.5 -right-0.5 bg-plex-accent text-white text-[10px] rounded-full w-4 h-4 flex items-center justify-center font-bold">
                  {jobCount}
                </span>
              </button>
            )}
            <button
              onClick={() => navigate({ page: "stats" })}
              className={`p-2 rounded ${
                currentPage === "stats" ? "text-plex-accent" : "text-plex-text-secondary hover:text-plex-text"
              }`}
              title="Statistics"
            >
              <BarChart3 className="w-[18px] h-[18px]" />
            </button>
            <button
              onClick={() => setHelpOpen(true)}
              className="p-2 rounded text-plex-text-secondary hover:text-plex-text"
              title="Keyboard Shortcuts (?)"
            >
              <HelpCircle className="w-[18px] h-[18px]" />
            </button>
            <button
              onClick={() => navigate({ page: "settings" })}
              className={`p-2 rounded ${
                currentPage === "settings" ? "text-plex-accent" : "text-plex-text-secondary hover:text-plex-text"
              }`}
              title="Settings"
            >
              <Settings className="w-[18px] h-[18px]" />
            </button>
          </div>
        </div>
      </div>
      <JobDrawer open={jobDrawerOpen} onClose={() => setJobDrawerOpen(false)} />
      <KeyboardShortcutsDialog open={helpOpen} onClose={() => setHelpOpen(false)} />
    </nav>
  );
}
