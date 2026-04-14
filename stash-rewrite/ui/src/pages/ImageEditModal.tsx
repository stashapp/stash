import { useState, useEffect } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { images, tags as tagsApi, performers as performersApi, studios as studiosApi, galleries as galleriesApi } from "../api/client";
import type { Image, ImageCreate } from "../api/types";
import { EditModal, Field, TextInput, TextArea, SaveButton } from "../components/EditModal";
import { RatingField } from "../components/Rating";
import { CustomFieldsEditor } from "../components/shared";

interface ImageEditProps {
  image: Image;
  open: boolean;
  onClose: () => void;
}

interface ImageCreateProps {
  open: boolean;
  onClose: () => void;
  onCreated: (id: number) => void;
}

interface ImageFormState {
  title: string;
  code: string;
  details: string;
  photographer: string;
  date: string;
  rating: number | undefined;
  organized: boolean;
  urls: string;
  studioId: number | undefined;
  selectedTagIds: number[];
  selectedPerformerIds: number[];
  selectedGalleryIds: number[];
  customFields: Record<string, string>;
}

interface ImageMetadataModalProps {
  title: string;
  open: boolean;
  onClose: () => void;
  initialState: ImageFormState;
  onSubmit: (data: ImageCreate) => void;
  isPending: boolean;
  error: Error | null;
  image?: Image;
}

const EMPTY_FORM_STATE: ImageFormState = {
  title: "",
  code: "",
  details: "",
  photographer: "",
  date: "",
  rating: undefined,
  organized: false,
  urls: "",
  studioId: undefined,
  selectedTagIds: [],
  selectedPerformerIds: [],
  selectedGalleryIds: [],
  customFields: {},
};

function toFormState(image?: Image): ImageFormState {
  if (!image) {
    return {
      ...EMPTY_FORM_STATE,
      selectedTagIds: [],
      selectedPerformerIds: [],
    };
  }

  return {
    title: image.title || "",
    code: image.code || "",
    details: image.details || "",
    photographer: image.photographer || "",
    date: image.date || "",
    rating: image.rating ?? undefined,
    organized: image.organized,
    urls: image.urls.join("\n"),
    studioId: image.studioId ?? undefined,
    selectedTagIds: image.tags.map((tag) => tag.id),
    selectedPerformerIds: image.performers.map((performer) => performer.id),
    selectedGalleryIds: image.galleryIds ?? [],
    customFields: Object.fromEntries(Object.entries(image.customFields ?? {}).map(([k, v]) => [k, String(v ?? "")])),
  };
}

function cloneFormState(state: ImageFormState): ImageFormState {
  return {
    ...state,
    selectedTagIds: [...state.selectedTagIds],
    selectedPerformerIds: [...state.selectedPerformerIds],
    selectedGalleryIds: [...state.selectedGalleryIds],
    customFields: { ...state.customFields },
  };
}

function ImageMetadataModal({ title, open, onClose, initialState, onSubmit, isPending, error, image }: ImageMetadataModalProps) {
  const [form, setForm] = useState<ImageFormState>(() => cloneFormState(initialState));
  const [tagSearch, setTagSearch] = useState("");
  const [perfSearch, setPerfSearch] = useState("");
  const [gallerySearch, setGallerySearch] = useState("");

  const { data: allTags } = useQuery({
    queryKey: ["tags-all"],
    queryFn: () => tagsApi.find({ perPage: 500, sort: "name", direction: "asc" }),
    enabled: open,
  });

  const { data: allPerformers } = useQuery({
    queryKey: ["performers-all"],
    queryFn: () => performersApi.find({ perPage: 500, sort: "name", direction: "asc" }),
    enabled: open,
  });

  const { data: allStudios } = useQuery({
    queryKey: ["studios-all"],
    queryFn: () => studiosApi.find({ perPage: 500, sort: "name", direction: "asc" }),
    enabled: open,
  });

  const { data: allGalleries } = useQuery({
    queryKey: ["galleries-all"],
    queryFn: () => galleriesApi.find({ perPage: 500, sort: "title", direction: "asc" }),
    enabled: open,
  });

  useEffect(() => {
    if (!open) return;
    setForm(cloneFormState(initialState));
    setTagSearch("");
    setPerfSearch("");
    setGallerySearch("");
  }, [initialState, open]);

  const handleSave = () => {
    const urlList = form.urls.split("\n").map((url) => url.trim()).filter(Boolean);
    onSubmit({
      title: form.title.trim() || undefined,
      code: form.code.trim() || undefined,
      details: form.details.trim() || undefined,
      photographer: form.photographer.trim() || undefined,
      date: form.date || undefined,
      rating: form.rating,
      organized: form.organized,
      studioId: form.studioId,
      urls: urlList,
      tagIds: form.selectedTagIds,
      performerIds: form.selectedPerformerIds,
      galleryIds: form.selectedGalleryIds,
      customFields: Object.keys(form.customFields).length > 0 ? form.customFields : undefined,
    });
  };

  const filteredTags = allTags?.items.filter(
    (tag) => !form.selectedTagIds.includes(tag.id) && tag.name.toLowerCase().includes(tagSearch.toLowerCase())
  ) ?? [];

  const filteredPerformers = allPerformers?.items.filter(
    (performer) => !form.selectedPerformerIds.includes(performer.id) && performer.name.toLowerCase().includes(perfSearch.toLowerCase())
  ) ?? [];

  const selectedTags = allTags?.items.filter((tag) => form.selectedTagIds.includes(tag.id)) ?? image?.tags ?? [];
  const selectedPerformers = allPerformers?.items.filter((performer) => form.selectedPerformerIds.includes(performer.id)) ?? image?.performers ?? [];

  return (
    <EditModal title={title} open={open} onClose={onClose}>
      <div className="grid grid-cols-2 gap-4">
        <Field label="Title">
          <TextInput value={form.title} onChange={(value) => setForm({ ...form, title: value })} placeholder="Image title" />
        </Field>
        <Field label="Date">
          <input
            type="date"
            value={form.date}
            onChange={(e) => setForm({ ...form, date: e.target.value })}
            className="w-full bg-gray-800 border border-gray-700 rounded px-3 py-2 text-sm text-gray-200 focus:outline-none focus:border-blue-500"
          />
        </Field>
      </div>

      <div className="grid grid-cols-2 gap-4">
        <Field label="Code">
          <TextInput value={form.code} onChange={(value) => setForm({ ...form, code: value })} placeholder="Image code" />
        </Field>
        <Field label="Photographer">
          <TextInput value={form.photographer} onChange={(value) => setForm({ ...form, photographer: value })} placeholder="Photographer name" />
        </Field>
      </div>

      <Field label="Details">
        <TextArea value={form.details} onChange={(value) => setForm({ ...form, details: value })} placeholder="Image description" />
      </Field>

      <div className="grid grid-cols-2 gap-4">
        <RatingField value={form.rating} onChange={(value) => setForm({ ...form, rating: value })} />
        <Field label="Studio">
          <select
            value={form.studioId ?? ""}
            onChange={(e) => setForm({ ...form, studioId: e.target.value ? Number(e.target.value) : undefined })}
            className="w-full bg-gray-800 border border-gray-700 rounded px-3 py-2 text-sm text-gray-200 focus:outline-none focus:border-blue-500"
          >
            <option value="">None</option>
            {allStudios?.items.map((s) => (
              <option key={s.id} value={s.id}>{s.name}</option>
            ))}
          </select>
        </Field>
      </div>

      <Field label="URLs (one per line)">
        <TextArea value={form.urls} onChange={(value) => setForm({ ...form, urls: value })} placeholder="https://..." rows={2} />
      </Field>

      <Field label="Tags">
        <div className="flex flex-wrap gap-1.5 mb-2">
          {selectedTags.map((tag) => (
            <span key={tag.id} className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium bg-blue-900 text-blue-300">
              {tag.name}
              <button onClick={() => setForm({ ...form, selectedTagIds: form.selectedTagIds.filter((id) => id !== tag.id) })} className="hover:text-white">×</button>
            </span>
          ))}
        </div>
        <input
          type="text"
          value={tagSearch}
          onChange={(e) => setTagSearch(e.target.value)}
          placeholder="Search tags..."
          className="w-full bg-gray-800 border border-gray-700 rounded px-3 py-1.5 text-sm text-gray-200 focus:outline-none focus:border-blue-500 mb-1"
        />
        {tagSearch && filteredTags.length > 0 && (
          <div className="max-h-32 overflow-y-auto bg-gray-800 rounded border border-gray-700">
            {filteredTags.slice(0, 10).map((tag) => (
              <button
                key={tag.id}
                onClick={() => { setForm({ ...form, selectedTagIds: [...form.selectedTagIds, tag.id] }); setTagSearch(""); }}
                className="block w-full text-left px-3 py-1.5 text-sm text-gray-300 hover:bg-gray-700"
              >{tag.name}</button>
            ))}
          </div>
        )}
      </Field>

      <Field label="Performers">
        <div className="flex flex-wrap gap-1.5 mb-2">
          {selectedPerformers.map((performer) => (
            <span key={performer.id} className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium bg-purple-900 text-purple-300">
              {performer.name}
              <button onClick={() => setForm({ ...form, selectedPerformerIds: form.selectedPerformerIds.filter((id) => id !== performer.id) })} className="hover:text-white">×</button>
            </span>
          ))}
        </div>
        <input
          type="text"
          value={perfSearch}
          onChange={(e) => setPerfSearch(e.target.value)}
          placeholder="Search performers..."
          className="w-full bg-gray-800 border border-gray-700 rounded px-3 py-1.5 text-sm text-gray-200 focus:outline-none focus:border-blue-500 mb-1"
        />
        {perfSearch && filteredPerformers.length > 0 && (
          <div className="max-h-32 overflow-y-auto bg-gray-800 rounded border border-gray-700">
            {filteredPerformers.slice(0, 10).map((performer) => (
              <button
                key={performer.id}
                onClick={() => { setForm({ ...form, selectedPerformerIds: [...form.selectedPerformerIds, performer.id] }); setPerfSearch(""); }}
                className="block w-full text-left px-3 py-1.5 text-sm text-gray-300 hover:bg-gray-700"
              >
                {performer.name}{performer.disambiguation ? ` (${performer.disambiguation})` : ""}
              </button>
            ))}
          </div>
        )}
      </Field>

      {/* Galleries */}
      <Field label="Galleries">
        <div className="flex flex-wrap gap-1.5 mb-2">
          {form.selectedGalleryIds.map((gid) => {
            const g = allGalleries?.items.find((gal) => gal.id === gid);
            return (
              <span key={gid} className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium bg-emerald-900 text-emerald-300">
                {g?.title || `Gallery #${gid}`}
                <button onClick={() => setForm({ ...form, selectedGalleryIds: form.selectedGalleryIds.filter((id) => id !== gid) })} className="hover:text-white">×</button>
              </span>
            );
          })}
        </div>
        <input
          type="text"
          value={gallerySearch}
          onChange={(e) => setGallerySearch(e.target.value)}
          placeholder="Search galleries..."
          className="w-full bg-gray-800 border border-gray-700 rounded px-3 py-1.5 text-sm text-gray-200 focus:outline-none focus:border-blue-500 mb-1"
        />
        {gallerySearch && (allGalleries?.items.filter(
          (g) => !form.selectedGalleryIds.includes(g.id) && (g.title || "").toLowerCase().includes(gallerySearch.toLowerCase())
        ) ?? []).length > 0 && (
          <div className="max-h-32 overflow-y-auto bg-gray-800 rounded border border-gray-700">
            {(allGalleries?.items.filter(
              (g) => !form.selectedGalleryIds.includes(g.id) && (g.title || "").toLowerCase().includes(gallerySearch.toLowerCase())
            ) ?? []).slice(0, 10).map((g) => (
              <button
                key={g.id}
                onClick={() => { setForm({ ...form, selectedGalleryIds: [...form.selectedGalleryIds, g.id] }); setGallerySearch(""); }}
                className="block w-full text-left px-3 py-1.5 text-sm text-gray-300 hover:bg-gray-700"
              >
                {g.title || "Untitled"}
              </button>
            ))}
          </div>
        )}
      </Field>

      <label className="flex items-center gap-2 text-sm text-gray-300 mb-4 cursor-pointer">
        <input type="checkbox" checked={form.organized} onChange={(e) => setForm({ ...form, organized: e.target.checked })} className="rounded border-gray-600 bg-gray-800" />
        Organized
      </label>

      <Field label="Custom Fields">
        <CustomFieldsEditor value={form.customFields} onChange={(v) => setForm({ ...form, customFields: v })} />
      </Field>

      {error && (
        <div className="bg-red-900/50 border border-red-700 text-red-300 rounded p-2 mb-4 text-sm">
          {error.message}
        </div>
      )}

      <div className="flex justify-end gap-3">
        <button onClick={onClose} className="px-4 py-2 text-sm text-gray-400 hover:text-white">Cancel</button>
        <SaveButton loading={isPending} onClick={handleSave} />
      </div>
    </EditModal>
  );
}

export function ImageEditModal({ image, open, onClose }: ImageEditProps) {
  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: (data: ImageCreate) => images.update(image.id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["image", image.id] });
      queryClient.invalidateQueries({ queryKey: ["images"] });
      onClose();
    },
  });

  return (
    <ImageMetadataModal
      title="Edit Image"
      open={open}
      onClose={onClose}
      initialState={toFormState(image)}
      onSubmit={(data) => mutation.mutate(data)}
      isPending={mutation.isPending}
      error={mutation.error as Error | null}
      image={image}
    />
  );
}

export function ImageCreateModal({ open, onClose, onCreated }: ImageCreateProps) {
  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: (data: ImageCreate) => images.create(data),
    onSuccess: (created) => {
      queryClient.invalidateQueries({ queryKey: ["images"] });
      onClose();
      if (created?.id) onCreated(created.id);
    },
  });

  return (
    <ImageMetadataModal
      title="Create Image"
      open={open}
      onClose={onClose}
      initialState={EMPTY_FORM_STATE}
      onSubmit={(data) => mutation.mutate(data)}
      isPending={mutation.isPending}
      error={mutation.error as Error | null}
    />
  );
}
