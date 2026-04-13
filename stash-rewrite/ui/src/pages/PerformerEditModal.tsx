import { useState, useEffect } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { performers, tags as tagsApi, entityImages } from "../api/client";
import { ImageInput } from "../components/ImageInput";
import type { Performer, PerformerUpdate } from "../api/types";
import { EditModal, Field, TextInput, TextArea, NumberInput, SaveButton } from "../components/EditModal";
import { RatingField } from "../components/Rating";

interface Props {
  performer: Performer;
  open: boolean;
  onClose: () => void;
}

const GENDER_OPTIONS = [
  { value: "Male", label: "Male" },
  { value: "Female", label: "Female" },
  { value: "TransMale", label: "Trans Male" },
  { value: "TransFemale", label: "Trans Female" },
  { value: "Intersex", label: "Intersex" },
  { value: "NonBinary", label: "Non-Binary" },
];

const CIRCUMCISED_OPTIONS = [
  { value: "Cut", label: "Cut" },
  { value: "Uncut", label: "Uncut" },
];

export function PerformerEditModal({ performer, open, onClose }: Props) {
  const queryClient = useQueryClient();

  const [name, setName] = useState(performer.name);
  const [disambiguation, setDisambiguation] = useState(performer.disambiguation || "");
  const [gender, setGender] = useState(performer.gender || "");
  const [birthdate, setBirthdate] = useState(performer.birthdate || "");
  const [ethnicity, setEthnicity] = useState(performer.ethnicity || "");
  const [country, setCountry] = useState(performer.country || "");
  const [eyeColor, setEyeColor] = useState(performer.eyeColor || "");
  const [hairColor, setHairColor] = useState(performer.hairColor || "");
  const [heightCm, setHeightCm] = useState<number | undefined>(performer.heightCm ?? undefined);
  const [weight, setWeight] = useState<number | undefined>(performer.weight ?? undefined);
  const [measurements, setMeasurements] = useState(performer.measurements || "");
  const [tattoos, setTattoos] = useState(performer.tattoos || "");
  const [piercings, setPiercings] = useState(performer.piercings || "");
  const [rating, setRating] = useState<number | undefined>(performer.rating ?? undefined);
  const [details, setDetails] = useState(performer.details || "");
  const [deathDate, setDeathDate] = useState(performer.deathDate || "");
  const [fakeTits, setFakeTits] = useState(performer.fakeTits || "");
  const [penisLength, setPenisLength] = useState<number | undefined>(performer.penisLength ?? undefined);
  const [circumcised, setCircumcised] = useState(performer.circumcised || "");
  const [careerStart, setCareerStart] = useState(performer.careerStart || "");
  const [careerEnd, setCareerEnd] = useState(performer.careerEnd || "");
  const [favorite, setFavorite] = useState(performer.favorite);
  const [ignoreAutoTag, setIgnoreAutoTag] = useState(performer.ignoreAutoTag ?? false);
  const [urls, setUrls] = useState(performer.urls.join("\n"));
  const [aliases, setAliases] = useState(performer.aliases.join(", "));
  const [selectedTagIds, setSelectedTagIds] = useState<number[]>(performer.tags.map((t) => t.id));
  const [tagSearch, setTagSearch] = useState("");

  const { data: allTags } = useQuery({
    queryKey: ["tags-all"],
    queryFn: () => tagsApi.find({ perPage: 500, sort: "name", direction: "asc" }),
  });

  useEffect(() => {
    setName(performer.name);
    setDisambiguation(performer.disambiguation || "");
    setGender(performer.gender || "");
    setBirthdate(performer.birthdate || "");
    setEthnicity(performer.ethnicity || "");
    setCountry(performer.country || "");
    setEyeColor(performer.eyeColor || "");
    setHairColor(performer.hairColor || "");
    setHeightCm(performer.heightCm ?? undefined);
    setWeight(performer.weight ?? undefined);
    setMeasurements(performer.measurements || "");
    setTattoos(performer.tattoos || "");
    setPiercings(performer.piercings || "");
    setRating(performer.rating ?? undefined);
    setDetails(performer.details || "");
    setDeathDate(performer.deathDate || "");
    setFakeTits(performer.fakeTits || "");
    setPenisLength(performer.penisLength ?? undefined);
    setCircumcised(performer.circumcised || "");
    setCareerStart(performer.careerStart || "");
    setCareerEnd(performer.careerEnd || "");
    setFavorite(performer.favorite);
    setIgnoreAutoTag(performer.ignoreAutoTag ?? false);
    setUrls(performer.urls.join("\n"));
    setAliases(performer.aliases.join(", "));
    setSelectedTagIds(performer.tags.map((t) => t.id));
  }, [performer]);

  const mutation = useMutation({
    mutationFn: (data: PerformerUpdate) => performers.update(performer.id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["performer", performer.id] });
      queryClient.invalidateQueries({ queryKey: ["performers"] });
      onClose();
    },
  });

  const handleSave = () => {
    const urlList = urls.split("\n").map((u) => u.trim()).filter(Boolean);
    const aliasList = aliases.split(",").map((a) => a.trim()).filter(Boolean);
    mutation.mutate({
      name,
      disambiguation: disambiguation || undefined,
      gender: gender || undefined,
      birthdate: birthdate || undefined,
      ethnicity: ethnicity || undefined,
      country: country || undefined,
      eyeColor: eyeColor || undefined,
      hairColor: hairColor || undefined,
      heightCm,
      weight,
      measurements: measurements || undefined,
      tattoos: tattoos || undefined,
      piercings: piercings || undefined,
      deathDate: deathDate || undefined,
      fakeTits: fakeTits || undefined,
      penisLength,
      circumcised: circumcised || undefined,
      careerStart: careerStart || undefined,
      careerEnd: careerEnd || undefined,
      rating,
      details: details || undefined,
      favorite,
      ignoreAutoTag,
      urls: urlList,
      aliases: aliasList,
      tagIds: selectedTagIds,
    });
  };

  const filteredTags = allTags?.items.filter(
    (t) => !selectedTagIds.includes(t.id) && t.name.toLowerCase().includes(tagSearch.toLowerCase())
  ) ?? [];
  const selectedTags = allTags?.items.filter((t) => selectedTagIds.includes(t.id)) ?? performer.tags;

  return (
    <EditModal title="Edit Performer" open={open} onClose={onClose}>
      <div className="grid grid-cols-[200px_1fr] gap-6">
        <ImageInput
          currentImageUrl={entityImages.performerImageUrl(performer.id)}
          onUpload={(file) => entityImages.uploadPerformerImage(performer.id, file)}
          onDelete={() => entityImages.deletePerformerImage(performer.id)}
          onSuccess={() => queryClient.invalidateQueries({ queryKey: ["performer", performer.id] })}
          label="Photo"
          aspectRatio="2/3"
        />
        <div className="space-y-4">
      <div className="grid grid-cols-2 gap-4">
        <Field label="Name *">
          <TextInput value={name} onChange={setName} placeholder="Performer name" />
        </Field>
        <Field label="Disambiguation">
          <TextInput value={disambiguation} onChange={setDisambiguation} placeholder="e.g. (2020s)" />
        </Field>
      </div>

      <div className="grid grid-cols-4 gap-4">
        <Field label="Gender">
          <select
            value={gender}
            onChange={(e) => setGender(e.target.value)}
            className="w-full bg-gray-800 border border-gray-700 rounded px-3 py-2 text-sm text-gray-200 focus:outline-none focus:border-blue-500"
          >
            <option value="">—</option>
            {GENDER_OPTIONS.map((o) => (
              <option key={o.value} value={o.value}>{o.label}</option>
            ))}
          </select>
        </Field>
        <Field label="Birthdate">
          <input
            type="date"
            value={birthdate}
            onChange={(e) => setBirthdate(e.target.value)}
            className="w-full bg-gray-800 border border-gray-700 rounded px-3 py-2 text-sm text-gray-200 focus:outline-none focus:border-blue-500"
          />
        </Field>
        <Field label="Death Date">
          <input
            type="date"
            value={deathDate}
            onChange={(e) => setDeathDate(e.target.value)}
            className="w-full bg-gray-800 border border-gray-700 rounded px-3 py-2 text-sm text-gray-200 focus:outline-none focus:border-blue-500"
          />
        </Field>
        <Field label="Country">
          <TextInput value={country} onChange={setCountry} placeholder="e.g. US" />
        </Field>
      </div>

      <div className="grid grid-cols-3 gap-4">
        <Field label="Ethnicity">
          <TextInput value={ethnicity} onChange={setEthnicity} />
        </Field>
        <Field label="Eye Color">
          <TextInput value={eyeColor} onChange={setEyeColor} />
        </Field>
        <Field label="Hair Color">
          <TextInput value={hairColor} onChange={setHairColor} />
        </Field>
      </div>

      <div className="grid grid-cols-3 gap-4">
        <Field label="Height (cm)">
          <NumberInput value={heightCm} onChange={setHeightCm} min={50} max={250} />
        </Field>
        <Field label="Weight (kg)">
          <NumberInput value={weight} onChange={setWeight} min={20} max={300} />
        </Field>
        <Field label="Measurements">
          <TextInput value={measurements} onChange={setMeasurements} placeholder="34D-24-34" />
        </Field>
      </div>

      <div className="grid grid-cols-2 gap-4">
        <Field label="Tattoos">
          <TextInput value={tattoos} onChange={setTattoos} />
        </Field>
        <Field label="Piercings">
          <TextInput value={piercings} onChange={setPiercings} />
        </Field>
      </div>

      <div className="grid grid-cols-3 gap-4">
        <Field label="Fake Tits">
          <TextInput value={fakeTits} onChange={setFakeTits} placeholder="e.g. Augmented" />
        </Field>
        <Field label="Penis Length (cm)">
          <NumberInput value={penisLength} onChange={setPenisLength} min={0} max={50} />
        </Field>
        <Field label="Circumcised">
          <select
            value={circumcised}
            onChange={(e) => setCircumcised(e.target.value)}
            className="w-full bg-gray-800 border border-gray-700 rounded px-3 py-2 text-sm text-gray-200 focus:outline-none focus:border-blue-500"
          >
            <option value="">—</option>
            {CIRCUMCISED_OPTIONS.map((o) => (
              <option key={o.value} value={o.value}>{o.label}</option>
            ))}
          </select>
        </Field>
      </div>

      <div className="grid grid-cols-2 gap-4">
        <Field label="Career Start">
          <input
            type="date"
            value={careerStart}
            onChange={(e) => setCareerStart(e.target.value)}
            className="w-full bg-gray-800 border border-gray-700 rounded px-3 py-2 text-sm text-gray-200 focus:outline-none focus:border-blue-500"
          />
        </Field>
        <Field label="Career End">
          <input
            type="date"
            value={careerEnd}
            onChange={(e) => setCareerEnd(e.target.value)}
            className="w-full bg-gray-800 border border-gray-700 rounded px-3 py-2 text-sm text-gray-200 focus:outline-none focus:border-blue-500"
          />
        </Field>
      </div>

      <Field label="Details">
        <TextArea value={details} onChange={setDetails} placeholder="Bio / notes" rows={2} />
      </Field>

      <div className="grid grid-cols-2 gap-4">
        <RatingField value={rating} onChange={setRating} />
        <Field label="Aliases (comma-separated)">
          <TextInput value={aliases} onChange={setAliases} placeholder="Alias 1, Alias 2" />
        </Field>
      </div>

      <Field label="URLs (one per line)">
        <TextArea value={urls} onChange={setUrls} placeholder="https://..." rows={2} />
      </Field>

      {/* Tags */}
      <Field label="Tags">
        <div className="flex flex-wrap gap-1.5 mb-2">
          {selectedTags.map((t) => (
            <span key={t.id} className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium bg-blue-900 text-blue-300">
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
              >{t.name}</button>
            ))}
          </div>
        )}
      </Field>

      <div className="flex gap-6 mb-4">
        <label className="flex items-center gap-2 text-sm text-gray-300 cursor-pointer">
          <input type="checkbox" checked={favorite} onChange={(e) => setFavorite(e.target.checked)} className="rounded border-gray-600 bg-gray-800" />
          Favorite
        </label>
        <label className="flex items-center gap-2 text-sm text-gray-300 cursor-pointer">
          <input type="checkbox" checked={ignoreAutoTag} onChange={(e) => setIgnoreAutoTag(e.target.checked)} className="rounded border-gray-600 bg-gray-800" />
          Ignore Auto Tag
        </label>
      </div>
        </div>{/* end right column */}
      </div>{/* end grid */}

      {mutation.error && (
        <div className="bg-red-900/50 border border-red-700 text-red-300 rounded p-2 mb-4 text-sm">
          {(mutation.error as Error).message}
        </div>
      )}

      <div className="flex justify-end gap-3">
        <button onClick={onClose} className="px-4 py-2 text-sm text-gray-400 hover:text-white">Cancel</button>
        <SaveButton loading={mutation.isPending} onClick={handleSave} />
      </div>
    </EditModal>
  );
}
