export const getStashIDs = (ids?: { stash_id: string; endpoint: string }[]) =>
  (ids ?? []).map(({ stash_id, endpoint }) => ({
    stash_id,
    endpoint,
  }));
