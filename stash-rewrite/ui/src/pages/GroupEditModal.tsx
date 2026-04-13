import { useState, useEffect } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { groups, studios, tags as tagsApi, entityImages } from "../api/client";
import type { Group, GroupUpdate, Tag } from "../api/types";
import { EditModal, Field, TextInput, TextArea, NumberInput, SaveButton } from "../components/EditModal";
import { ImageInput } from "../components/ImageInput";
import { RatingField } from "../components/Rating";

interface Props {
  group: Group;
  open: boolean;
  onClose: () => void;
}

export function GroupEditModal({ group, open, onClose }: Props) {
  const queryClient = useQueryClient();

  const [name, setName] = useState(group.name);
  const [aliases, setAliases] = useState(group.aliases ?? "");
  const [director, setDirector] = useState(group.director ?? "");
  const [date, setDate] = useState(group.date ?? "");
  const [duration, setDuration] = useState<number | undefined>(group.duration ?? undefined);
  const [rating, setRating] = useState<number | undefined>(group.rating ?? undefined);
  const [studioId, setStudioId] = useState<number | undefined>(group.studioId ?? undefined);
  const [synopsis, setSynopsis] = useState(group.synopsis ?? "");
  const [urls, setUrls] = useState(group.urls.join("\n"));
  const [selectedTagIds, setSelectedTagIds] = useState<number[]>(group.tags.map((t) => t.id));

  // Tag search
  const [tagSearch, setTagSearch] = useState("");
  const { data: allTags } = useQuery({
    queryKey: ["tags-all"],
    queryFn: () => tagsApi.find({ perPage: 500, sort: "name", direction: "asc" }),
  });

  // Studios
  const { data: allStudios } = useQuery({
    queryKey: ["studios-all"],
    queryFn: () => studios.find({ perPage: 500, sort: "name", direction: "asc" }),
  });

  useEffect(() => {
    setName(group.name);
    setAliases(group.aliases ?? "");
    setDirector(group.director ?? "");
    setDate(group.date ?? "");
    setDuration(group.duration ?? undefined);
    setRating(group.rating ?? undefined);
    setStudioId(group.studioId ?? undefined);
    setSynopsis(group.synopsis ?? "");
    setUrls(group.urls.join("\n"));
    setSelectedTagIds(group.tags.map((t) => t.id));
  }, [group]);

  const mutation = useMutation({
    mutationFn: (data: GroupUpdate) => groups.update(group.id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["group", group.id] });
      queryClient.invalidateQueries({ queryKey: ["groups"] });
      onClose();
    },
  });

  const handleSave = () => {
    const urlList = urls.split("\n").map((u) => u.trim()).filter(Boolean);
    mutation.mutate({
      name,
      aliases: aliases || undefined,
      director: director || undefined,
      date: date || undefined,
      duration,
      rating,
      studioId,
      synopsis: synopsis || undefined,
      urls: urlList,
      tagIds: selectedTagIds,
    });
  };

  const filteredTags = allTags?.items.filter(
    (t) => !selectedTagIds.includes(t.id) && t.name.toLowerCase().includes(tagSearch.toLowerCase())
  ) ?? [];
  const selectedTags = allTags?.items.filter((t) => selectedTagIds.includes(t.id)) ?? group.tags;

  return (
    <EditModal title={`Edit Group: ${group.name}`} open={open} onClose={onClose}>
      <div className="grid grid-cols-2 gap-4 mb-4">
        <ImageInput
          currentImageUrl={entityImages.groupFrontImageUrl(group.id)}
          onUpload={(file) => entityImages.uploadGroupFrontImage(group.id, file)}
          onDelete={() => entityImages.deleteGroupFrontImage(group.id)}
          onSuccess={() => queryClient.invalidateQueries({ queryKey: ["group", group.id] })}
          label="Front Cover"
          aspectRatio="2/3"
        />
        <ImageInput
          currentImageUrl={entityImages.groupBackImageUrl(group.id)}
          onUpload={(file) => entityImages.uploadGroupBackImage(group.id, file)}
          onDelete={() => entityImages.deleteGroupBackImage(group.id)}
          onSuccess={() => queryClient.invalidateQueries({ queryKey: ["group", group.id] })}
          label="Back Cover"
          aspectRatio="2/3"
        />
      </div>
      <div className="grid grid-cols-2 gap-4">
        <Field label="Name *">
          <TextInput value={name} onChange={setName} placeholder="Group name" />
        </Field>
        <Field label="Aliases">
          <TextInput value={aliases} onChange={setAliases} placeholder="Alternative names" />
        </Field>
      </div>

      <div className="grid grid-cols-2 gap-4">
        <Field label="Director">
          <TextInput value={director} onChange={setDirector} placeholder="Director name" />
        </Field>
        <Field label="Date">
          <input
            type="date"
            value={date}
            onChange={(e) => setDate(e.target.value)}
            className="w-full bg-gray-800 border border-gray-700 rounded px-3 py-2 text-sm text-gray-200 focus:outline-none focus:border-blue-500"
          />
        </Field>
      </div>

      <div className="grid grid-cols-3 gap-4">
        <Field label="Duration (seconds)">
          <NumberInput value={duration} onChange={setDuration} min={0} />
        </Field>
        <RatingField value={rating} onChange={setRating} />
        <Field label="Studio">
          <select
            value={studioId ?? ""}
            onChange={(e) => setStudioId(e.target.value ? Number(e.target.value) : undefined)}
            className="w-full bg-gray-800 border border-gray-700 rounded px-3 py-2 text-sm text-gray-200 focus:outline-none focus:border-blue-500"
          >
            <option value="">None</option>
            {allStudios?.items.map((s) => (
              <option key={s.id} value={s.id}>{s.name}</option>
            ))}
          </select>
        </Field>
      </div>

      <Field label="Synopsis">
        <TextArea value={synopsis} onChange={setSynopsis} placeholder="Group synopsis / description" rows={4} />
      </Field>

      <Field label="URLs (one per line)">
        <TextArea value={urls} onChange={setUrls} placeholder="https://..." rows={2} />
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

      <div className="flex justify-end gap-3 mt-4">
        <button onClick={onClose} className="px-4 py-2 text-sm text-gray-400 hover:text-white">Cancel</button>
        <SaveButton loading={mutation.isPending} onClick={handleSave} />
      </div>
    </EditModal>
  );
}
