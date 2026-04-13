import { useDeferredValue, useEffect, useMemo, useRef, useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { Building2, Film, FolderOpen, ImageIcon, Layers, Loader2, Search, Tag, Users } from "lucide-react";
import { galleries, groups, images, performers, scenes, studios, tags } from "../api/client";

interface Props {
  navigate: (r: any) => void;
}

type SearchGroup = {
  key: string;
  label: string;
  icon: React.ComponentType<{ className?: string }>;
  items: { id: number; title: string; subtitle?: string; route: any }[];
};

export function GlobalSearch({ navigate }: Props) {
  const [term, setTerm] = useState("");
  const [open, setOpen] = useState(false);
  const deferredTerm = useDeferredValue(term.trim());
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const onPointerDown = (event: MouseEvent) => {
      if (containerRef.current && !containerRef.current.contains(event.target as Node)) {
        setOpen(false);
      }
    };
    document.addEventListener("mousedown", onPointerDown);
    return () => document.removeEventListener("mousedown", onPointerDown);
  }, []);

  const { data, isFetching } = useQuery({
    queryKey: ["global-search", deferredTerm],
    enabled: deferredTerm.length >= 2,
    queryFn: async () => {
      const query = { q: deferredTerm, perPage: 5, direction: "desc" as const };
      const [sceneRes, performerRes, studioRes, tagRes, galleryRes, imageRes, groupRes] = await Promise.all([
        scenes.find(query),
        performers.find({ ...query, sort: "name", direction: "asc" }),
        studios.find({ ...query, sort: "name", direction: "asc" }),
        tags.find({ ...query, sort: "name", direction: "asc" }),
        galleries.find({ ...query, sort: "title", direction: "asc" }),
        images.find({ ...query, sort: "title", direction: "asc" }),
        groups.find({ ...query, sort: "name", direction: "asc" }),
      ]);

      const groupsData: SearchGroup[] = [
        {
          key: "scenes",
          label: "Scenes",
          icon: Film,
          items: sceneRes.items.map((item) => ({
            id: item.id,
            title: item.title || item.files[0]?.basename || `Scene ${item.id}`,
            subtitle: item.studioName || item.date || undefined,
            route: { page: "scene", id: item.id },
          })),
        },
        {
          key: "performers",
          label: "Performers",
          icon: Users,
          items: performerRes.items.map((item) => ({
            id: item.id,
            title: item.name,
            subtitle: item.disambiguation || undefined,
            route: { page: "performer", id: item.id },
          })),
        },
        {
          key: "studios",
          label: "Studios",
          icon: Building2,
          items: studioRes.items.map((item) => ({
            id: item.id,
            title: item.name,
            subtitle: item.parentName || undefined,
            route: { page: "studio", id: item.id },
          })),
        },
        {
          key: "tags",
          label: "Tags",
          icon: Tag,
          items: tagRes.items.map((item) => ({
            id: item.id,
            title: item.name,
            subtitle: item.description || undefined,
            route: { page: "tag", id: item.id },
          })),
        },
        {
          key: "galleries",
          label: "Galleries",
          icon: FolderOpen,
          items: galleryRes.items.map((item) => ({
            id: item.id,
            title: item.title || `Gallery ${item.id}`,
            subtitle: item.studioName || item.date || undefined,
            route: { page: "gallery", id: item.id },
          })),
        },
        {
          key: "images",
          label: "Images",
          icon: ImageIcon,
          items: imageRes.items.map((item) => ({
            id: item.id,
            title: item.title || `Image ${item.id}`,
            subtitle: item.studioName || item.date || undefined,
            route: { page: "image", id: item.id },
          })),
        },
        {
          key: "groups",
          label: "Groups",
          icon: Layers,
          items: groupRes.items.map((item) => ({
            id: item.id,
            title: item.name,
            subtitle: item.studioName || item.date || undefined,
            route: { page: "group", id: item.id },
          })),
        },
      ];

      return groupsData.filter((group) => group.items.length > 0);
    },
  });

  const flatResults = useMemo(() => (data ?? []).flatMap((group) => group.items), [data]);

  const handleSelect = (route: any) => {
    navigate(route);
    setOpen(false);
    setTerm("");
  };

  return (
    <div ref={containerRef} className="relative hidden xl:block">
      <div className="relative">
        <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-plex-text-muted" />
        <input
          value={term}
          onChange={(event) => {
            setTerm(event.target.value);
            setOpen(true);
          }}
          onFocus={() => setOpen(true)}
          onKeyDown={(event) => {
            if (event.key === "Escape") {
              setOpen(false);
              return;
            }
            if (event.key === "Enter" && flatResults.length > 0) {
              event.preventDefault();
              handleSelect(flatResults[0].route);
            }
          }}
          placeholder="Search all..."
          className="w-72 rounded-lg border border-plex-border bg-plex-input py-1.5 pl-9 pr-3 text-sm text-plex-text placeholder:text-plex-text-muted focus:border-plex-accent focus:outline-none"
        />
      </div>

      {open && (
        <div className="absolute right-0 top-full z-50 mt-2 w-[30rem] overflow-hidden rounded-xl border border-plex-border bg-plex-card shadow-2xl shadow-black/40">
          <div className="border-b border-plex-border px-3 py-2 text-[11px] uppercase tracking-wider text-plex-text-muted">
            Global Search
          </div>

          {deferredTerm.length < 2 ? (
            <div className="px-4 py-6 text-sm text-plex-text-secondary">Type at least 2 characters to search scenes, performers, studios, tags, galleries, images, and groups.</div>
          ) : isFetching ? (
            <div className="flex items-center gap-2 px-4 py-6 text-sm text-plex-text-secondary">
              <Loader2 className="h-4 w-4 animate-spin" /> Searching...
            </div>
          ) : !data || data.length === 0 ? (
            <div className="px-4 py-6 text-sm text-plex-text-secondary">No results found for "{deferredTerm}".</div>
          ) : (
            <div className="max-h-[28rem] overflow-y-auto">
              {data.map((group) => {
                const Icon = group.icon;
                return (
                  <div key={group.key} className="border-b border-plex-border last:border-b-0">
                    <div className="flex items-center gap-2 px-3 py-2 text-[11px] font-semibold uppercase tracking-wider text-plex-text-muted">
                      <Icon className="h-3.5 w-3.5" />
                      {group.label}
                    </div>
                    <div className="pb-2">
                      {group.items.map((item) => (
                        <button
                          key={`${group.key}-${item.id}`}
                          onClick={() => handleSelect(item.route)}
                          className="flex w-full items-start gap-3 px-3 py-2 text-left hover:bg-plex-surface"
                        >
                          <Icon className="mt-0.5 h-4 w-4 shrink-0 text-plex-accent" />
                          <span className="min-w-0 flex-1">
                            <span className="block truncate text-sm text-plex-text">{item.title}</span>
                            {item.subtitle && <span className="block truncate text-xs text-plex-text-secondary">{item.subtitle}</span>}
                          </span>
                        </button>
                      ))}
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </div>
      )}
    </div>
  );
}