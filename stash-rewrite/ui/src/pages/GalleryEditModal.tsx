import { useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { galleries, studios, tags, performers } from "../api/client";
import type { Gallery, GalleryUpdate, Studio, Tag, Performer } from "../api/types";
import { EditModal, Field, TextInput, TextArea, SaveButton } from "../components/EditModal";
import { Search, X } from "lucide-react";
import { RatingField } from "../components/Rating";

interface Props {
  gallery: Gallery;
  open: boolean;
  onClose: () => void;
}

export function GalleryEditModal({ gallery, open, onClose }: Props) {
  const qc = useQueryClient();
  const [form, setForm] = useState({
    title: gallery.title ?? "",
    code: gallery.code ?? "",
    date: gallery.date ?? "",
    details: gallery.details ?? "",
    photographer: gallery.photographer ?? "",
    rating: gallery.rating,
    organized: gallery.organized,
    studioId: gallery.studioId,
    urls: gallery.urls.join("\n"),
    tagIds: gallery.tags.map((t) => t.id),
    performerIds: gallery.performers.map((p) => p.id),
  });

  const [tagSearch, setTagSearch] = useState("");
  const [performerSearch, setPerformerSearch] = useState("");

  const { data: studioList } = useQuery({ queryKey: ["studios-all"], queryFn: () => studios.find({ perPage: 200 }) });
  const { data: tagResults } = useQuery({ queryKey: ["tags-search", tagSearch], queryFn: () => tags.find({ q: tagSearch, perPage: 20 }), enabled: tagSearch.length > 0 });
  const { data: performerResults } = useQuery({ queryKey: ["performers-search", performerSearch], queryFn: () => performers.find({ q: performerSearch, perPage: 20 }), enabled: performerSearch.length > 0 });

  const selectedTags = gallery.tags.filter((t) => form.tagIds.includes(t.id));
  const selectedPerformers = gallery.performers.filter((p) => form.performerIds.includes(p.id));

  const mutation = useMutation({
    mutationFn: (data: GalleryUpdate) => galleries.update(gallery.id, data),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["gallery", gallery.id] });
      qc.invalidateQueries({ queryKey: ["galleries"] });
      onClose();
    },
  });

  const save = () => {
    mutation.mutate({
      title: form.title || undefined,
      code: form.code || undefined,
      date: form.date || undefined,
      details: form.details || undefined,
      photographer: form.photographer || undefined,
      rating: form.rating,
      organized: form.organized,
      studioId: form.studioId,
      urls: form.urls ? form.urls.split("\n").map((u) => u.trim()).filter(Boolean) : [],
      tagIds: form.tagIds,
      performerIds: form.performerIds,
    });
  };

  return (
    <EditModal title={`Edit Gallery: ${gallery.title || "Untitled"}`} open={open} onClose={onClose}>
      <div className="grid grid-cols-2 gap-4">
        <div className="col-span-2">
          <Field label="Title">
            <TextInput value={form.title} onChange={(v) => setForm({ ...form, title: v })} />
          </Field>
        </div>
        <Field label="Code">
          <TextInput value={form.code} onChange={(v) => setForm({ ...form, code: v })} />
        </Field>
        <Field label="Date">
          <TextInput value={form.date} onChange={(v) => setForm({ ...form, date: v })} placeholder="YYYY-MM-DD" />
        </Field>
        <Field label="Photographer">
          <TextInput value={form.photographer} onChange={(v) => setForm({ ...form, photographer: v })} />
        </Field>
        <RatingField value={form.rating} onChange={(v) => setForm({ ...form, rating: v })} />
        <Field label="Studio">
          <select
            value={form.studioId ?? ""}
            onChange={(e) => setForm({ ...form, studioId: e.target.value ? Number(e.target.value) : undefined })}
            className="w-full bg-gray-800 border border-gray-700 rounded px-3 py-2 text-sm text-gray-200 focus:outline-none focus:border-blue-500"
          >
            <option value="">None</option>
            {studioList?.items.map((s) => <option key={s.id} value={s.id}>{s.name}</option>)}
          </select>
        </Field>
        <div className="flex items-end pb-4">
          <label className="flex items-center gap-2 text-sm">
            <input type="checkbox" checked={form.organized} onChange={(e) => setForm({ ...form, organized: e.target.checked })} className="rounded bg-gray-800 border-gray-700" />
            Organized
          </label>
        </div>
      </div>
      <Field label="Details">
        <TextArea value={form.details} onChange={(v) => setForm({ ...form, details: v })} rows={3} />
      </Field>
      <Field label="URLs (one per line)">
        <TextArea value={form.urls} onChange={(v) => setForm({ ...form, urls: v })} rows={2} />
      </Field>

      {/* Tags picker */}
      <Field label="Tags">
        <div className="flex flex-wrap gap-1 mb-2">
          {selectedTags.map((t) => (
            <span key={t.id} className="bg-blue-600/30 text-blue-300 text-xs px-2 py-0.5 rounded-full flex items-center gap-1">
              {t.name}
              <X className="w-3 h-3 cursor-pointer" onClick={() => setForm({ ...form, tagIds: form.tagIds.filter((id) => id !== t.id) })} />
            </span>
          ))}
        </div>
        <div className="relative">
          <Search className="w-3.5 h-3.5 absolute left-2 top-2.5 text-gray-500" />
          <input
            type="text" value={tagSearch} onChange={(e) => setTagSearch(e.target.value)} placeholder="Search tags..."
            className="w-full bg-gray-800 border border-gray-700 rounded pl-8 pr-3 py-2 text-sm text-gray-200 focus:outline-none focus:border-blue-500"
          />
          {tagSearch && tagResults && (
            <div className="absolute z-10 w-full mt-1 bg-gray-800 border border-gray-700 rounded max-h-32 overflow-y-auto">
              {tagResults.items.filter((t) => !form.tagIds.includes(t.id)).map((t) => (
                <div key={t.id} onClick={() => { setForm({ ...form, tagIds: [...form.tagIds, t.id] }); setTagSearch(""); }}
                  className="px-3 py-1.5 text-sm hover:bg-gray-700 cursor-pointer">{t.name}</div>
              ))}
            </div>
          )}
        </div>
      </Field>

      {/* Performers picker */}
      <Field label="Performers">
        <div className="flex flex-wrap gap-1 mb-2">
          {selectedPerformers.map((p) => (
            <span key={p.id} className="bg-purple-600/30 text-purple-300 text-xs px-2 py-0.5 rounded-full flex items-center gap-1">
              {p.name}
              <X className="w-3 h-3 cursor-pointer" onClick={() => setForm({ ...form, performerIds: form.performerIds.filter((id) => id !== p.id) })} />
            </span>
          ))}
        </div>
        <div className="relative">
          <Search className="w-3.5 h-3.5 absolute left-2 top-2.5 text-gray-500" />
          <input
            type="text" value={performerSearch} onChange={(e) => setPerformerSearch(e.target.value)} placeholder="Search performers..."
            className="w-full bg-gray-800 border border-gray-700 rounded pl-8 pr-3 py-2 text-sm text-gray-200 focus:outline-none focus:border-blue-500"
          />
          {performerSearch && performerResults && (
            <div className="absolute z-10 w-full mt-1 bg-gray-800 border border-gray-700 rounded max-h-32 overflow-y-auto">
              {performerResults.items.filter((p) => !form.performerIds.includes(p.id)).map((p) => (
                <div key={p.id} onClick={() => { setForm({ ...form, performerIds: [...form.performerIds, p.id] }); setPerformerSearch(""); }}
                  className="px-3 py-1.5 text-sm hover:bg-gray-700 cursor-pointer">{p.name}</div>
              ))}
            </div>
          )}
        </div>
      </Field>

      <div className="flex justify-end mt-4">
        <SaveButton loading={mutation.isPending} onClick={save} />
      </div>
    </EditModal>
  );
}
