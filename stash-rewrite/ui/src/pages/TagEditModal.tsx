import { useState, useEffect } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { tags, entityImages } from "../api/client";
import type { TagDetail, TagUpdate, Tag } from "../api/types";
import { EditModal, Field, TextInput, TextArea, SaveButton } from "../components/EditModal";
import { ImageInput } from "../components/ImageInput";
import { CustomFieldsEditor } from "../components/shared";

interface Props {
  tag: TagDetail;
  open: boolean;
  onClose: () => void;
}

export function TagEditModal({ tag, open, onClose }: Props) {
  const queryClient = useQueryClient();

  const [name, setName] = useState(tag.name);
  const [sortName, setSortName] = useState(tag.sortName ?? "");
  const [description, setDescription] = useState(tag.description ?? "");
  const [favorite, setFavorite] = useState(tag.favorite);
  const [ignoreAutoTag, setIgnoreAutoTag] = useState(tag.ignoreAutoTag);
  const [aliases, setAliases] = useState(tag.aliases.join("\n"));
  const [selectedParentIds, setSelectedParentIds] = useState<number[]>(tag.parents.map((t) => t.id));
  const [selectedChildIds, setSelectedChildIds] = useState<number[]>(tag.children.map((t) => t.id));

  // Tag search for parents
  const [parentSearch, setParentSearch] = useState("");
  // Tag search for children
  const [childSearch, setChildSearch] = useState("");
  const [customFields, setCustomFields] = useState<Record<string, string>>(
    Object.fromEntries(Object.entries(tag.customFields ?? {}).map(([k, v]) => [k, String(v ?? "")]))
  );

  const { data: allTags } = useQuery({
    queryKey: ["tags-all"],
    queryFn: () => tags.find({ perPage: 500, sort: "name", direction: "asc" }),
  });

  useEffect(() => {
    setName(tag.name);
    setSortName(tag.sortName ?? "");
    setDescription(tag.description ?? "");
    setFavorite(tag.favorite);
    setIgnoreAutoTag(tag.ignoreAutoTag);
    setAliases(tag.aliases.join("\n"));
    setSelectedParentIds(tag.parents.map((t) => t.id));
    setSelectedChildIds(tag.children.map((t) => t.id));
    setCustomFields(Object.fromEntries(Object.entries(tag.customFields ?? {}).map(([k, v]) => [k, String(v ?? "")])));
  }, [tag]);

  const mutation = useMutation({
    mutationFn: (data: TagUpdate) => tags.update(tag.id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["tag", tag.id] });
      queryClient.invalidateQueries({ queryKey: ["tags"] });
      onClose();
    },
  });

  const handleSave = () => {
    const aliasList = aliases.split("\n").map((a) => a.trim()).filter(Boolean);
    mutation.mutate({
      name,
      sortName: sortName || undefined,
      description: description || undefined,
      favorite,
      ignoreAutoTag,
      aliases: aliasList,
      parentIds: selectedParentIds,
      childIds: selectedChildIds,
      customFields: Object.keys(customFields).length > 0 ? customFields : undefined,
    });
  };

  // Exclude self from parent/child options
  const excludedIds = new Set([tag.id, ...selectedParentIds, ...selectedChildIds]);

  const filteredParents = allTags?.items.filter(
    (t) => !excludedIds.has(t.id) && t.name.toLowerCase().includes(parentSearch.toLowerCase())
  ) ?? [];

  const filteredChildren = allTags?.items.filter(
    (t) => !excludedIds.has(t.id) && t.name.toLowerCase().includes(childSearch.toLowerCase())
  ) ?? [];

  const selectedParents = allTags?.items.filter((t) => selectedParentIds.includes(t.id)) ?? tag.parents;
  const selectedChildren = allTags?.items.filter((t) => selectedChildIds.includes(t.id)) ?? tag.children;

  return (
    <EditModal title={`Edit Tag: ${tag.name}`} open={open} onClose={onClose}>
      <ImageInput
        currentImageUrl={entityImages.tagImageUrl(tag.id)}
        onUpload={(file) => entityImages.uploadTagImage(tag.id, file)}
        onDelete={() => entityImages.deleteTagImage(tag.id)}
        onSuccess={() => queryClient.invalidateQueries({ queryKey: ["tag", tag.id] })}
        label="Image"
        aspectRatio="16/9"
      />
      <Field label="Name *">
        <TextInput value={name} onChange={setName} placeholder="Tag name" />
      </Field>

      <Field label="Sort Name">
        <TextInput value={sortName} onChange={setSortName} placeholder="Custom sort name (optional)" />
      </Field>

      <Field label="Description">
        <TextArea value={description} onChange={setDescription} placeholder="Tag description" rows={3} />
      </Field>

      <Field label="Aliases (one per line)">
        <TextArea value={aliases} onChange={setAliases} placeholder="Alternative names, one per line" rows={2} />
      </Field>

      <div className="flex items-center gap-4 mb-4">
        <label className="flex items-center gap-2 text-sm">
          <input type="checkbox" checked={favorite} onChange={(e) => setFavorite(e.target.checked)} className="rounded bg-gray-800 border-gray-700" />
          Favorite
        </label>
        <label className="flex items-center gap-2 text-sm">
          <input type="checkbox" checked={ignoreAutoTag} onChange={(e) => setIgnoreAutoTag(e.target.checked)} className="rounded bg-gray-800 border-gray-700" />
          Ignore Auto Tag
        </label>
      </div>

      <Field label="Custom Fields">
        <CustomFieldsEditor value={customFields} onChange={setCustomFields} />
      </Field>

      {/* Parent Tags */}
      <Field label="Parent Tags">
        <div className="flex flex-wrap gap-1.5 mb-2">
          {selectedParents.map((t) => (
            <span
              key={t.id}
              className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium bg-green-900 text-green-300"
            >
              {t.name}
              <button onClick={() => setSelectedParentIds(selectedParentIds.filter((id) => id !== t.id))} className="hover:text-white">×</button>
            </span>
          ))}
        </div>
        <input
          type="text"
          value={parentSearch}
          onChange={(e) => setParentSearch(e.target.value)}
          placeholder="Search parent tags..."
          className="w-full bg-gray-800 border border-gray-700 rounded px-3 py-1.5 text-sm text-gray-200 focus:outline-none focus:border-blue-500 mb-1"
        />
        {parentSearch && filteredParents.length > 0 && (
          <div className="max-h-32 overflow-y-auto bg-gray-800 rounded border border-gray-700">
            {filteredParents.slice(0, 10).map((t) => (
              <button
                key={t.id}
                onClick={() => { setSelectedParentIds([...selectedParentIds, t.id]); setParentSearch(""); }}
                className="block w-full text-left px-3 py-1.5 text-sm text-gray-300 hover:bg-gray-700"
              >
                {t.name}
              </button>
            ))}
          </div>
        )}
      </Field>

      {/* Child Tags */}
      <Field label="Child Tags">
        <div className="flex flex-wrap gap-1.5 mb-2">
          {selectedChildren.map((t) => (
            <span
              key={t.id}
              className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium bg-purple-900 text-purple-300"
            >
              {t.name}
              <button onClick={() => setSelectedChildIds(selectedChildIds.filter((id) => id !== t.id))} className="hover:text-white">×</button>
            </span>
          ))}
        </div>
        <input
          type="text"
          value={childSearch}
          onChange={(e) => setChildSearch(e.target.value)}
          placeholder="Search child tags..."
          className="w-full bg-gray-800 border border-gray-700 rounded px-3 py-1.5 text-sm text-gray-200 focus:outline-none focus:border-blue-500 mb-1"
        />
        {childSearch && filteredChildren.length > 0 && (
          <div className="max-h-32 overflow-y-auto bg-gray-800 rounded border border-gray-700">
            {filteredChildren.slice(0, 10).map((t) => (
              <button
                key={t.id}
                onClick={() => { setSelectedChildIds([...selectedChildIds, t.id]); setChildSearch(""); }}
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
