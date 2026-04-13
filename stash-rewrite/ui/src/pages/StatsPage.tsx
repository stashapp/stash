import { useQuery } from "@tanstack/react-query";
import { system, scenes, performers, images as imagesApi, studios, tags as tagsApi } from "../api/client";
import type { Scene, Performer, Image, Studio, Tag } from "../api/types";
import {
  Film, Users, Tag as TagIcon, Building2, Images, FolderOpen, Layers, HardDrive,
  Clock, Play, ChevronLeft, ChevronRight, Star, Calendar, TrendingUp,
} from "lucide-react";
import { useRef, useState, useCallback } from "react";

interface StatsPageProps {
  onNavigate?: (route: any) => void;
}

export function StatsPage({ onNavigate }: StatsPageProps) {
  const { data: stats, isLoading } = useQuery({
    queryKey: ["stats"],
    queryFn: system.stats,
  });

  const { data: recentScenes } = useQuery({
    queryKey: ["recent-scenes"],
    queryFn: () => scenes.find({ perPage: 16, sort: "created_at", direction: "desc" }),
  });

  const { data: recentPerformers } = useQuery({
    queryKey: ["recent-performers"],
    queryFn: () => performers.find({ perPage: 16, sort: "created_at", direction: "desc" }),
  });

  const { data: recentImages } = useQuery({
    queryKey: ["recent-images"],
    queryFn: () => imagesApi.find({ perPage: 16, sort: "updated_at", direction: "desc" }),
  });

  // Additional recommendation rows
  const { data: recentlyPlayed } = useQuery({
    queryKey: ["recently-played-scenes"],
    queryFn: () => scenes.find({ perPage: 16, sort: "last_played_at", direction: "desc" }),
  });

  const { data: mostPlayed } = useQuery({
    queryKey: ["most-played-scenes"],
    queryFn: () => scenes.find({ perPage: 16, sort: "play_count", direction: "desc" }),
  });

  const { data: recentlyReleased } = useQuery({
    queryKey: ["recently-released-scenes"],
    queryFn: () => scenes.find({ perPage: 16, sort: "date", direction: "desc" }),
  });

  const { data: topRated } = useQuery({
    queryKey: ["top-rated-scenes"],
    queryFn: () => scenes.find({ perPage: 16, sort: "rating", direction: "desc" }),
  });

  const { data: recentStudios } = useQuery({
    queryKey: ["recent-studios"],
    queryFn: () => studios.find({ perPage: 16, sort: "created_at", direction: "desc" }),
  });

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500" />
      </div>
    );
  }

  if (!stats) return null;

  const formatBytes = (bytes: number) => {
    if (bytes === 0) return "0 B";
    const k = 1024;
    const sizes = ["B", "KB", "MB", "GB", "TB"];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return `${(bytes / Math.pow(k, i)).toFixed(1)} ${sizes[i]}`;
  };

  const formatDuration = (seconds: number) => {
    const h = Math.floor(seconds / 3600);
    const m = Math.floor((seconds % 3600) / 60);
    if (h > 0) return `${h.toLocaleString()}h ${m}m`;
    return `${m}m`;
  };

  const statCards = [
    { label: "Scenes", value: stats.sceneCount, icon: Film, color: "text-blue-400", bg: "bg-blue-400/10", page: "scenes" },
    { label: "Images", value: stats.imageCount, icon: Images, color: "text-green-400", bg: "bg-green-400/10", page: "images" },
    { label: "Galleries", value: stats.galleryCount, icon: FolderOpen, color: "text-yellow-400", bg: "bg-yellow-400/10", page: "galleries" },
    { label: "Performers", value: stats.performerCount, icon: Users, color: "text-pink-400", bg: "bg-pink-400/10", page: "performers" },
    { label: "Studios", value: stats.studioCount, icon: Building2, color: "text-purple-400", bg: "bg-purple-400/10", page: "studios" },
    { label: "Tags", value: stats.tagCount, icon: TagIcon, color: "text-orange-400", bg: "bg-orange-400/10", page: "tags" },
    { label: "Groups", value: stats.groupCount, icon: Layers, color: "text-cyan-400", bg: "bg-cyan-400/10", page: "groups" },
    { label: "Total Size", value: formatBytes(stats.totalFileSize), icon: HardDrive, color: "text-gray-400", bg: "bg-gray-400/10" },
  ];

  return (
    <div className="space-y-8">
      {/* Stats Grid */}
      <div>
        <h1 className="text-2xl font-bold mb-4">Dashboard</h1>
        <div className="grid grid-cols-2 sm:grid-cols-4 gap-3">
          {statCards.map((card) => (
            <button
              key={card.label}
              onClick={() => card.page && onNavigate?.({ page: card.page })}
              className={`bg-gray-900 rounded-lg p-4 text-left transition-colors ${card.page ? "hover:bg-gray-800 cursor-pointer" : "cursor-default"}`}
            >
              <div className="flex items-center gap-3">
                <div className={`p-2 rounded-lg ${card.bg}`}>
                  <card.icon className={`w-5 h-5 ${card.color}`} />
                </div>
                <div>
                  <div className="text-xl font-bold">{typeof card.value === "number" ? card.value.toLocaleString() : card.value}</div>
                  <div className="text-xs text-gray-500">{card.label}</div>
                </div>
              </div>
            </button>
          ))}
        </div>
      </div>

      {/* Playback stats */}
      {stats.totalPlayDuration > 0 && (
        <div className="grid grid-cols-2 gap-3">
          <div className="bg-gray-900 rounded-lg p-4 flex items-center gap-3">
            <div className="p-2 rounded-lg bg-indigo-400/10">
              <Clock className="w-5 h-5 text-indigo-400" />
            </div>
            <div>
              <div className="text-xl font-bold">{formatDuration(stats.totalPlayDuration)}</div>
              <div className="text-xs text-gray-500">Total Play Duration</div>
            </div>
          </div>
          <div className="bg-gray-900 rounded-lg p-4 flex items-center gap-3">
            <div className="p-2 rounded-lg bg-emerald-400/10">
              <Play className="w-5 h-5 text-emerald-400" />
            </div>
            <div>
              <div className="text-xl font-bold">{stats.sceneCount > 0 ? formatDuration(stats.totalPlayDuration / stats.sceneCount) : "0m"}</div>
              <div className="text-xs text-gray-500">Avg Play Time / Scene</div>
            </div>
          </div>
        </div>
      )}

      {/* Recently Added Scenes */}
      {recentScenes && recentScenes.items.length > 0 && (
        <Carousel
          title="Recently Added Scenes"
          onSeeAll={() => onNavigate?.({ page: "scenes" })}
        >
          {recentScenes.items.map((scene) => (
            <SceneCard key={scene.id} scene={scene} onNavigate={onNavigate} />
          ))}
        </Carousel>
      )}

      {/* Recently Added Performers */}
      {recentPerformers && recentPerformers.items.length > 0 && (
        <Carousel
          title="Recently Added Performers"
          onSeeAll={() => onNavigate?.({ page: "performers" })}
        >
          {recentPerformers.items.map((p) => (
            <PerformerCard key={p.id} performer={p} onNavigate={onNavigate} />
          ))}
        </Carousel>
      )}

      {/* Recently Added Images */}
      {recentImages && recentImages.items.length > 0 && (
        <Carousel
          title="Recently Added Images"
          onSeeAll={() => onNavigate?.({ page: "images" })}
        >
          {recentImages.items.map((img) => (
            <ImageCard key={img.id} image={img} onNavigate={onNavigate} />
          ))}
        </Carousel>
      )}

      {/* Recently Released Scenes (by scene date) */}
      {recentlyReleased && recentlyReleased.items.filter(s => s.date).length > 0 && (
        <Carousel
          title="Recently Released"
          onSeeAll={() => onNavigate?.({ page: "scenes" })}
        >
          {recentlyReleased.items.filter(s => s.date).map((scene) => (
            <SceneCard key={scene.id} scene={scene} onNavigate={onNavigate} />
          ))}
        </Carousel>
      )}

      {/* Recently Played */}
      {recentlyPlayed && recentlyPlayed.items.filter(s => s.playCount > 0).length > 0 && (
        <Carousel
          title="Recently Played"
          onSeeAll={() => onNavigate?.({ page: "scenes" })}
        >
          {recentlyPlayed.items.filter(s => s.playCount > 0).map((scene) => (
            <SceneCard key={scene.id} scene={scene} onNavigate={onNavigate} />
          ))}
        </Carousel>
      )}

      {/* Most Played */}
      {mostPlayed && mostPlayed.items.filter(s => s.playCount > 0).length > 0 && (
        <Carousel
          title="Most Played"
          onSeeAll={() => onNavigate?.({ page: "scenes" })}
        >
          {mostPlayed.items.filter(s => s.playCount > 0).map((scene) => (
            <SceneCard key={scene.id} scene={scene} onNavigate={onNavigate} />
          ))}
        </Carousel>
      )}

      {/* Top Rated */}
      {topRated && topRated.items.filter(s => s.rating != null && s.rating > 0).length > 0 && (
        <Carousel
          title="Top Rated"
          onSeeAll={() => onNavigate?.({ page: "scenes" })}
        >
          {topRated.items.filter(s => s.rating != null && s.rating > 0).map((scene) => (
            <SceneCard key={scene.id} scene={scene} onNavigate={onNavigate} />
          ))}
        </Carousel>
      )}

      {/* Recently Added Studios */}
      {recentStudios && recentStudios.items.length > 0 && (
        <Carousel
          title="Recently Added Studios"
          onSeeAll={() => onNavigate?.({ page: "studios" })}
        >
          {recentStudios.items.map((studio) => (
            <StudioCard key={studio.id} studio={studio} onNavigate={onNavigate} />
          ))}
        </Carousel>
      )}
    </div>
  );
}

// ===== Carousel Component =====
function Carousel({ title, onSeeAll, children }: { title: string; onSeeAll?: () => void; children: React.ReactNode }) {
  const scrollRef = useRef<HTMLDivElement>(null);
  const [canScrollLeft, setCanScrollLeft] = useState(false);
  const [canScrollRight, setCanScrollRight] = useState(true);

  const checkScroll = useCallback(() => {
    const el = scrollRef.current;
    if (!el) return;
    setCanScrollLeft(el.scrollLeft > 0);
    setCanScrollRight(el.scrollLeft < el.scrollWidth - el.clientWidth - 1);
  }, []);

  const scroll = (dir: "left" | "right") => {
    const el = scrollRef.current;
    if (!el) return;
    const amount = el.clientWidth * 0.8;
    el.scrollBy({ left: dir === "left" ? -amount : amount, behavior: "smooth" });
    setTimeout(checkScroll, 350);
  };

  return (
    <div>
      <div className="flex items-center justify-between mb-3">
        <h2 className="text-lg font-semibold">{title}</h2>
        <div className="flex items-center gap-2">
          <button
            onClick={() => scroll("left")}
            disabled={!canScrollLeft}
            className="p-1 rounded hover:bg-gray-800 disabled:opacity-30 disabled:cursor-default transition-colors"
          >
            <ChevronLeft className="w-5 h-5" />
          </button>
          <button
            onClick={() => scroll("right")}
            disabled={!canScrollRight}
            className="p-1 rounded hover:bg-gray-800 disabled:opacity-30 disabled:cursor-default transition-colors"
          >
            <ChevronRight className="w-5 h-5" />
          </button>
          {onSeeAll && (
            <button onClick={onSeeAll} className="text-xs text-blue-400 hover:text-blue-300 ml-2">
              See All
            </button>
          )}
        </div>
      </div>
      <div
        ref={scrollRef}
        onScroll={checkScroll}
        className="flex gap-3 overflow-x-auto scrollbar-hide snap-x snap-mandatory pb-2"
        style={{ scrollbarWidth: "none", msOverflowStyle: "none" }}
      >
        {children}
      </div>
    </div>
  );
}

// ===== Scene Card =====
function SceneCard({ scene, onNavigate }: { scene: Scene; onNavigate?: (r: { page: string; id?: number }) => void }) {
  const formatDur = (s?: number) => {
    if (!s) return "";
    const files = scene.files;
    const dur = files?.[0]?.duration ?? 0;
    if (dur <= 0) return "";
    const h = Math.floor(dur / 3600);
    const m = Math.floor((dur % 3600) / 60);
    const sec = Math.floor(dur % 60);
    return h > 0 ? `${h}:${m.toString().padStart(2, "0")}:${sec.toString().padStart(2, "0")}` : `${m}:${sec.toString().padStart(2, "0")}`;
  };

  const duration = scene.files?.[0]?.duration;

  return (
    <button
      onClick={() => onNavigate?.({ page: "scene", id: scene.id })}
      className="flex-shrink-0 w-56 snap-start group text-left"
    >
      <div className="relative aspect-video rounded-lg overflow-hidden bg-gray-800 mb-2">
        <img
          src={scenes.screenshotUrl(scene.id)}
          alt={scene.title || ""}
          className="w-full h-full object-cover transition-transform duration-300 group-hover:scale-105"
          loading="lazy"
        />
        {duration && duration > 0 && (
          <span className="absolute bottom-1 right-1 bg-black/80 text-white text-[10px] font-mono px-1.5 py-0.5 rounded">
            {formatDur(duration)}
          </span>
        )}
      </div>
      <div className="text-sm font-medium truncate group-hover:text-blue-400 transition-colors">
        {scene.title || scene.files?.[0]?.basename || `Scene ${scene.id}`}
      </div>
      <div className="text-xs text-gray-500 truncate">
        {scene.date || new Date(scene.createdAt).toLocaleDateString()}
        {scene.performers.length > 0 && ` · ${scene.performers.map((p) => p.name).join(", ")}`}
      </div>
    </button>
  );
}

// ===== Performer Card =====
function PerformerCard({ performer, onNavigate }: { performer: Performer; onNavigate?: (r: { page: string; id?: number }) => void }) {
  return (
    <button
      onClick={() => onNavigate?.({ page: "performer", id: performer.id })}
      className="flex-shrink-0 w-36 snap-start group text-left"
    >
      <div className="relative aspect-[2/3] rounded-lg overflow-hidden bg-gray-800 mb-2">
        {performer.imagePath ? (
          <img
            src={performer.imagePath}
            alt={performer.name}
            className="w-full h-full object-cover transition-transform duration-300 group-hover:scale-105"
            loading="lazy"
          />
        ) : (
          <div className="flex items-center justify-center w-full h-full">
            <Users className="w-10 h-10 text-gray-600" />
          </div>
        )}
        {performer.sceneCount > 0 && (
          <span className="absolute bottom-1 right-1 bg-black/80 text-white text-[10px] px-1.5 py-0.5 rounded">
            {performer.sceneCount} scenes
          </span>
        )}
      </div>
      <div className="text-sm font-medium truncate group-hover:text-blue-400 transition-colors">
        {performer.name}
      </div>
      {performer.disambiguation && (
        <div className="text-xs text-gray-500 truncate">{performer.disambiguation}</div>
      )}
    </button>
  );
}

// ===== Image Card =====
function ImageCard({ image, onNavigate }: { image: Image; onNavigate?: (r: { page: string; id?: number }) => void }) {
  return (
    <button
      onClick={() => onNavigate?.({ page: "image", id: image.id })}
      className="flex-shrink-0 w-44 snap-start group text-left"
    >
      <div className="relative aspect-square rounded-lg overflow-hidden bg-gray-800 mb-2">
        <img
          src={imagesApi.thumbnailUrl(image.id)}
          alt={image.title || ""}
          className="w-full h-full object-cover transition-transform duration-300 group-hover:scale-105"
          loading="lazy"
        />
      </div>
      <div className="text-sm font-medium truncate group-hover:text-blue-400 transition-colors">
        {image.title || `Image ${image.id}`}
      </div>
    </button>
  );
}

// ===== Studio Card =====
function StudioCard({ studio, onNavigate }: { studio: Studio; onNavigate?: (r: { page: string; id?: number }) => void }) {
  return (
    <button
      onClick={() => onNavigate?.({ page: "studio", id: studio.id })}
      className="flex-shrink-0 w-40 snap-start group text-left"
    >
      <div className="relative aspect-video rounded-lg overflow-hidden bg-gray-800 mb-2 flex items-center justify-center">
        {studio.imagePath ? (
          <img
            src={studio.imagePath}
            alt={studio.name}
            className="w-full h-full object-contain p-2 transition-transform duration-300 group-hover:scale-105"
            loading="lazy"
          />
        ) : (
          <Building2 className="w-10 h-10 text-gray-600" />
        )}
      </div>
      <div className="text-sm font-medium truncate group-hover:text-blue-400 transition-colors">
        {studio.name}
      </div>
      {studio.sceneCount > 0 && (
        <div className="text-xs text-gray-500">{studio.sceneCount} scenes</div>
      )}
    </button>
  );
}
