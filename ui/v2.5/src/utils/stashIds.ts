import * as GQL from "src/core/generated-graphql";

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

/**
 * Utility to add or update a StashID in an array.
 * If a StashID with the same endpoint exists, it will be replaced.
 * Otherwise, the new StashID will be appended.
 */
export const addUpdateStashID = (
  existingStashIDs: GQL.StashIdInput[],
  newItem: GQL.StashIdInput,
  allowMultiple: boolean = false
): GQL.StashIdInput[] => {
  const existingIndex = existingStashIDs.findIndex(
    (s) => s.endpoint === newItem.endpoint
  );

  if (!allowMultiple && existingIndex >= 0) {
    const newStashIDs = [...existingStashIDs];
    newStashIDs[existingIndex] = newItem;
    return newStashIDs;
  }

  // ensure we don't add duplicates if allowMultiple is true
  if (
    allowMultiple &&
    existingStashIDs.some(
      (s) => s.endpoint === newItem.endpoint && s.stash_id === newItem.stash_id
    )
  ) {
    return existingStashIDs;
  }

  return [...existingStashIDs, newItem];
};
