// Language data with code abbreviation (sorted alphabetically, Unknown at end)
export const languageData: { code: string; name: string }[] = [
  { code: "de", name: "Deutsche" },
  { code: "en", name: "English" },
  { code: "es", name: "Español" },
  { code: "fr", name: "Français" },
  { code: "nl", name: "Holandés" },
  { code: "it", name: "Italiano" },
  { code: "pt", name: "Português" },
  { code: "ja", name: "日本" },
  { code: "ko", name: "한국인" },
  { code: "ru", name: "Русский" },
  { code: "00", name: "Unknown" },
].sort((a, b) => {
  // Keep "Unknown" at the end
  if (a.code === "00") return 1;
  if (b.code === "00") return -1;
  return a.name.localeCompare(b.name);
});

// Legacy map for backward compatibility
export const languageMap = new Map<string, string>(
  languageData.map((lang) => [lang.code, lang.name])
);

// Get language code by name (for display as badge)
export const getLanguageCode = (name: string): string => {
  const lang = languageData.find((l) => l.name === name);
  return lang?.code.toUpperCase() ?? "";
};

export const valueToCode = (value?: string | null) => {
  if (!value) {
    return undefined;
  }

  return Array.from(languageMap.keys()).find((v) => {
    return languageMap.get(v) === value;
  });
};
