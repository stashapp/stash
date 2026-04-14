import { useState, useEffect } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { studios, tags as tagsApi, entityImages } from "../api/client";
import type { Studio, StudioUpdate, Tag } from "../api/types";
import { EditModal, Field, TextInput, TextArea, SaveButton } from "../components/EditModal";
import { ImageInput } from "../components/ImageInput";
import { RatingField } from "../components/Rating";
import { CustomFieldsEditor } from "../components/shared";

interface Props {
  studio: Studio;
  open: boolean;
  onClose: () => void;
}

export function StudioEditModal({ studio, open, onClose }: Props) {
  const queryClient = useQueryClient();

  const [name, setName] = useState(studio.name);
  const [details, setDetails] = useState(studio.details ?? "");
  const [rating, setRating] = useState<number | undefined>(studio.rating ?? undefined);
  const [favorite, setFavorite] = useState(studio.favorite);
  const [ignoreAutoTag, setIgnoreAutoTag] = useState(studio.ignoreAutoTag);
  const [organized, setOrganized] = useState(studio.organized);
  const [urls, setUrls] = useState(studio.urls.join("\n"));
  const [aliases, setAliases] = useState(studio.aliases.join("\n"));
  const [parentId, setParentId] = useState<number | undefined>(studio.parentId ?? undefined);
  const [selectedTagIds, setSelectedTagIds] = useState<number[]>(studio.tags.map((t) => t.id));

  // Tag search
  const [tagSearch, setTagSearch] = useState("");
  const [customFields, setCustomFields] = useState<Record<string, string>>(
    Object.fromEntries(Object.entries(studio.customFields ?? {}).map(([k, v]) => [k, String(v ?? "")]))
  );
  const { data: allTags } = useQuery({
    queryKey: ["tags-all"],
    queryFn: () => tagsApi.find({ perPage: 500, sort: "name", direction: "asc" }),
  });

  // Parent studio list
  const { data: allStudios } = useQuery({
    queryKey: ["studios-all"],
    queryFn: () => studios.find({ perPage: 500, sort: "name", direction: "asc" }),
  });

  useEffect(() => {
    setName(studio.name);
    setDetails(studio.details ?? "");
    setRating(studio.rating ?? undefined);
    setFavorite(studio.favorite);
    setIgnoreAutoTag(studio.ignoreAutoTag);
    setOrganized(studio.organized);
    setUrls(studio.urls.join("\n"));
    setAliases(studio.aliases.join("\n"));
    setParentId(studio.parentId ?? undefined);
    setSelectedTagIds(studio.tags.map((t) => t.id));
    setCustomFields(Object.fromEntries(Object.entries(studio.customFields ?? {}).map(([k, v]) => [k, String(v ?? "")])));
  }, [studio]);

  const mutation = useMutation({
    mutationFn: (data: StudioUpdate) => studios.update(studio.id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["studio", studio.id] });
      queryClient.invalidateQueries({ queryKey: ["studios"] });
      onClose();
    },
  });

  const handleSave = () => {
    const urlList = urls.split("\n").map((u) => u.trim()).filter(Boolean);
    const aliasList = aliases.split("\n").map((a) => a.trim()).filter(Boolean);
    mutation.mutate({
      name,
      details: details || undefined,
      rating,
      favorite,
      ignoreAutoTag,
      organized,
      parentId,
      urls: urlList,
      aliases: aliasList,
      tagIds: selectedTagIds,
      customFields: Object.keys(customFields).length > 0 ? customFields : undefined,
    });
  };

  const filteredTags = allTags?.items.filter(
    (t) => !selectedTagIds.includes(t.id) && t.name.toLowerCase().includes(tagSearch.toLowerCase())
  ) ?? [];
  const selectedTags = allTags?.items.filter((t) => selectedTagIds.includes(t.id)) ?? studio.tags;

  // Exclude self from parent studio options
  const parentStudioOptions = allStudios?.items.filter((s) => s.id !== studio.id) ?? [];

  return (
    <EditModal title={`Edit Studio: ${studio.name}`} open={open} onClose={onClose}>
      <div className="flex gap-6 mb-4">
        <ImageInput
          currentImageUrl={entityImages.studioImageUrl(studio.id)}
          onUpload={(file) => entityImages.uploadStudioImage(studio.id, file)}
          onDelete={() => entityImages.deleteStudioImage(studio.id)}
          onSuccess={() => queryClient.invalidateQueries({ queryKey: ["studio", studio.id] })}
          label="Logo"
          aspectRatio="1/1"
          className="w-40"
        />
        <div className="flex-1 space-y-4">
      <div className="grid grid-cols-2 gap-4">
        <Field label="Name *">
          <TextInput value={name} onChange={setName} placeholder="Studio name" />
        </Field>
        <Field label="Parent Studio">
          <select
            value={parentId ?? ""}
            onChange={(e) => setParentId(e.target.value ? Number(e.target.value) : undefined)}
            className="w-full bg-gray-800 border border-gray-700 rounded px-3 py-2 text-sm text-gray-200 focus:outline-none focus:border-blue-500"
          >
            <option value="">None</option>
            {parentStudioOptions.map((s) => (
              <option key={s.id} value={s.id}>{s.name}</option>
            ))}
          </select>
        </Field>
      </div>

      <Field label="Details">
        <TextArea value={details} onChange={setDetails} placeholder="Studio description" rows={3} />
      </Field>

      <div className="grid grid-cols-2 gap-4">
        <RatingField value={rating} onChange={setRating} />
        <div className="flex items-end gap-4 pb-4">
          <label className="flex items-center gap-2 text-sm">
            <input type="checkbox" checked={favorite} onChange={(e) => setFavorite(e.target.checked)} className="rounded bg-gray-800 border-gray-700" />
            Favorite
          </label>
          <label className="flex items-center gap-2 text-sm">
            <input type="checkbox" checked={ignoreAutoTag} onChange={(e) => setIgnoreAutoTag(e.target.checked)} className="rounded bg-gray-800 border-gray-700" />
            Ignore Auto Tag
          </label>
          <label className="flex items-center gap-2 text-sm">
            <input type="checkbox" checked={organized} onChange={(e) => setOrganized(e.target.checked)} className="rounded bg-gray-800 border-gray-700" />
            Organized
          </label>
        </div>
      </div>

      <Field label="URLs (one per line)">
        <TextArea value={urls} onChange={setUrls} placeholder="https://..." rows={2} />
      </Field>

      <Field label="Aliases (one per line)">
        <TextArea value={aliases} onChange={setAliases} placeholder="Alias names, one per line" rows={2} />
      </Field>

      {/* Tags */}
      <Field label="Tags">
        <div className="flex flex-wrap gap-1.5 mb-2">
          {selectedTags.map((t) => (
            <span
              key={t.id}
              className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium bg-blue-900 text-blue-300"
            >
              {t.name}
              <button onClick={() => setSelectedTagIds(selectedTagIds.filter((id) => id !== t.id))} className="hover:text-white">×</button>
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
            {filteredTags.slice(0, 10).map((t) => (
              <button
                key={t.id}
                onClick={() => { setSelectedTagIds([...selectedTagIds, t.id]); setTagSearch(""); }}
                className="block w-full text-left px-3 py-1.5 text-sm text-gray-300 hover:bg-gray-700"
              >
                {t.name}
              </button>
            ))}
          </div>
        )}
      </Field>

      <Field label="Custom Fields">
        <CustomFieldsEditor value={customFields} onChange={setCustomFields} />
      </Field>

      </div></div>
      <div className="flex justify-end gap-3 mt-4">
        <button onClick={onClose} className="px-4 py-2 text-sm text-gray-400 hover:text-white">Cancel</button>
        <SaveButton loading={mutation.isPending} onClick={handleSave} />
      </div>
    </EditModal>
  );
}
