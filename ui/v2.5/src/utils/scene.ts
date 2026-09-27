// The date a performer's age is calculated from for a scene. The production
// date is when the scene was actually shot, so it is preferred over the
// release date where one is set. Empty strings are treated as unset, since
// older databases can contain them.
export function sceneAgeFromDate(
  productionDate?: string | null,
  date?: string | null
) {
  return productionDate || date || undefined;
}
