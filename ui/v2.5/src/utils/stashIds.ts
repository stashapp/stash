 export const getStashIDs = (
  ids?: { stash_id: string; endpoint: string; updated_at: string }[]
) =>
  (ids ?? []).map(({ stash_id, endpoint, updated_at }) => ({
    stash_id,
    endpoint,
    updated_at,
  }));

// UUID regex pattern to detect StashIDs (supports v4 and v7)
const UUID_PATTERN =
  /^[0-9a-f]{8}-[0-9a-f]{4}-[47][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i;

/**
 * Separates a list of inputs into names and StashIDs based on UUID pattern matching
 * @param inputs - Array of strings that could be either names or StashIDs
 * @returns Object containing separate arrays for names and stashIds
 */
export const separateNamesAndStashIds = (
  inputs: string[]
): { names: string[]; stashIds: string[] } => {
  const names: string[] = [];
  const stashIds: string[] = [];

  inputs.forEach((input) => {
    if (UUID_PATTERN.test(input)) {
      stashIds.push(input);
    } else {
      names.push(input);
    }
  });

  return { names, stashIds };
};
