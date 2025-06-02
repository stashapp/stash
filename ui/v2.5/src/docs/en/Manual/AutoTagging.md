# Auto Tagging

When media filepaths or filenames contain a Performer, Studio, Tag, or date, it is assigned those Performers, Studios, Tags, and dates. It will **only** tag based on Performer, Studio, and Tag names that exist in your database and recognized date formats.

When the Performer/Studio/Tag name has multiple words, the search will include paths/filenames where the Performer/Studio/Tag name is separated with `.`, `-`, `_`, and whitespace characters.

For example, auto tagging for performer `Jane Doe` will match the following filenames:

* `Jane.Doe.1.mp4`
* `Jane_Doe.2.mp4`
* `Jane-Doe.3.mp4`
* `Jane Doe.4.mp4`

Matching is case insensitive, and should only match exact wording within word boundaries. For example, the tag `Jane Doe` will not match `Maryjane-Doe` or `Jane-Doen`, but will match `Mary-Jane-Doe`, `Jane-Doe_n`, and `[OF]jane doe`. Dates are matched using common formats like YYYY-MM-DD, YYYYMMDD, and DD.MM.YYYY.

Auto tagging for specific Performers, Studios, and Tags can be performed from the individual Performer/Studio/Tag page.

> **Note:** Performer autotagging does not currently match on performer aliases.