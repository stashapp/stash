# Locked fields

Individual scene fields can be locked once you know they are correct. A locked field cannot be changed by anything, including you, until it is unlocked.

## Locking a field

On the scene's Edit tab, click the lock icon to the right of a field to lock or unlock it, then save the scene. A locked field's input is disabled. Locked fields are listed on the scene's Details tab.

The following scene fields can be locked:

- Title
- Studio Code
- URLs
- Date
- Production Date
- Director
- Studio
- Details

An empty field can also be locked, to mark it as intentionally blank.

## What respects locked fields

| Feature | Behaviour |
|---------|-----------|
| Scene edit | Locked inputs are disabled. A field can be unlocked, changed and re-locked in a single save. |
| Scraping | Scraped values for locked fields are not shown or applied. |
| [Identify](/help/Identify.md) | Locked fields are skipped and logged. The other fields are still set. |
| [Scene Tagger](/help/Tagger.md) | Locked fields are excluded when saving. |
| [Auto Tag](/help/AutoTagging.md) | A locked studio is not set. |
| [Scene Filename Parser](/help/SceneFilenameParser.md) | Parsed values for a locked title, date or studio are ignored. |
| Bulk edit | Locked fields are skipped on the affected scenes, and a warning shows how many scenes are affected. |
| Merge | Locked fields on the destination scene keep their values. |
| GraphQL API | Updates that change a locked field fail with an error. |
| Import/export | Locked fields are included in the export and restored on import. |

## Locked fields and the Organized flag

The Organized flag stops automated tasks such as Identify and Auto Tag from modifying a whole scene, but still allows manual edits. Locked fields protect individual fields from all changes, including manual edits.
