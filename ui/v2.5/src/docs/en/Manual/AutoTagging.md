# Auto Tag

Auto Tag automatically assigns Performers, Studios, and Tags to your media based on their names found in file paths or filenames. This task works for scenes, images, and galleries.

This task is part of the advanced settings mode.

## Rules

> **Important:** Auto Tag only works for names that already exist in your Stash database. It does not create new Performers, Studios, or Tags.

 - Multi-word names are matched when words appear in order and are separated by any of these characters: `.`, `-`, `_`, or whitespace. These separators are treated as word boundaries.
 - Matching is case-insensitive but requires complete words within word boundaries. Partial words or misspelled words will not match.
 - Auto Tag does not match performer aliases. Aliases will not be considered during matching.

### Examples (performer "Jane Doe")

**Matches:**

| Example | Explanation |
|---|---|
| `Jane.Doe.1.mp4` | Dot as separator. |
| `Jane_Doe.2.mp4` | Underscore as separator. |
| `Jane-Doe.3.mp4` | Hyphen as separator. |
| `Jane Doe.4.mp4` | Whitespace as separator. |
| `Mary-Jane-Doe` | Extra characters around word boundaries are allowed. |
| `Jane-Doe_n` | Extra characters around word boundaries are allowed. |
| `[OF]jane doe` | Extra characters around word boundaries are allowed. |

**Does not match:**

| Example | Explanation |
|---|---|
| `Maryjane-Doe` | Combined words without separator. |
| `Jane-Doen` | Spelling mismatch. |

### Organized flag

Scenes, images, and galleries that have the Organized flag added to them will not be modified by Auto Tag. You can also use Organized flag status as a filter.

### Ignore Auto Tag flag

Performers or Tags that have Ignore Auto Tag flag added to them will be skipped by the Auto Tag task.

## Running task

- **Auto Tag:** You can run the Auto Tag task on your entire library from the Tasks page.
- **Selective Auto Tag:** You can run the Auto Tag task on specific directories from the Tasks page.
- **Individual pages:** You can run Auto Tag tasks for specific Performers, Studios, and Tags from their respective pages.
