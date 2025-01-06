export const getStashIDs = (
  ids?: { stash_id: string; endpoint: string; updated_at: string }[]
) =>
  (ids ?? []).map(({ stash_id, endpoint, updated_at }) => ({
    stash_id,
    endpoint,
    updated_at,
  }));
