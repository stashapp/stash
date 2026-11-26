# Locked fields

Individual scene fields can be locked once you know they are correct. A locked field cannot be changed by anything, including you, until it is unlocked.

## Locking a field

On the scene's Edit tab, click the lock icon inside a field, at its right edge, to lock or unlock it, then save the scene. A locked field's input is disabled and greyed out. Locked fields are listed on the scene's Details tab.

The following scene fields can be locked:

- Title
- Studio Code
- URLs
- Date
- Production Date
- Director
- Studio
- Performers
- Tags
- Groups
- Galleries
- Details

An empty field can also be locked, to mark it as intentionally blank.

## Rules

- A locked field blocks manual edits, including edits made from the other side of the link. For example, a scene with locked galleries cannot be added to or removed from a gallery by editing that gallery's scenes.
- All automation (Identify, Auto Tag, the Scene Tagger, scraping and the Scene Filename Parser) skips locked fields, and does not create missing performers, tags or studios for a locked field.
- Deleting or merging a linked performer, tag, gallery, group or studio still updates the scene, even when that field is locked.

## What respects locked fields

| Feature | Behaviour |
|---------|-----------|
| Scene edit | Locked inputs are disabled. A field can be unlocked, changed and re-locked in a single save. |
| Editing the other side | Scenes with a locked field cannot be added to or removed from that object, eg. from a gallery's edit page. |
| Scraping | Scraped values for locked fields are not shown or applied. |
| [Identify](/help/Identify.md) | Locked fields are skipped and logged. The other fields are still set. |
| [Scene Tagger](/help/Tagger.md) | Locked fields are excluded when saving. Missing performers, tags and studios are not offered for creation or linking when that field is locked. |
| [Auto Tag](/help/AutoTagging.md) | Locked performers, tags and studio are not changed. |
| [Scene Filename Parser](/help/SceneFilenameParser.md) | Parsed values for a locked title, date, studio, performers or tags are ignored. |
| Bulk edit | Locked fields are skipped on the affected scenes. A warning shows how many scenes are affected, and after saving a notice shows how many scenes were skipped for each field. |
| Merge | Locked fields on the destination scene keep their values. For locked performers, tags, groups and galleries, the destination keeps its own links and does not gain those of the merged scenes. |
| GraphQL API | Updates that change a locked field fail with an error. |
| Import/export | Locked fields are included in the export and restored on import. |

## Locked fields and the Organized flag

The Organized flag stops automated tasks such as Identify and Auto Tag from modifying a whole scene, but still allows manual edits. Locked fields protect individual fields from all changes, including manual edits.
