# Scene Filename Parser

[This tool](/sceneFilenameParser) parses the scene filenames in your library and allows setting the metadata from those filenames.

## Parser Options

To use this tool, a filename pattern must be entered. The pattern accepts the following fields:

| Field | Remark |
|-------|--------|
| `title` | Text captured within is set as the title of the scene. |
|`ext`|Matches the end of the filename. It is not captured. Does not include the last `.` character.|
|`d`|Matches delimiter characters (`-_.`). Not captured.|
|`i`|Matches any ignored word entered in the `Ignored words` field. Ignored words are entered as space-delimited words. Not captured. Use this to match release artifacts like `DVDRip` or release groups.|
|`date`|Matches `yyyy-mm-dd` and sets the date of the scene.|
|`rating`|Matches a single digit and sets the rating of the scene.|
|`performer`| Sets the scene performer, based on the text captured.|
|`tag`| Sets the scene tag, based on the text captured.|
|`studio`| Sets the studio performer, based on the text captured.|
|`{}`|Matches any characters. Not captured.|

> **⚠️ Note:** `performer`, `tag` and `studio` fields will only match against Performers/Tags/Studios that already exist in the system.

The `performer`/`tag`/`studio` fields will remove any delimiter characters (`.-_`) before querying. Name matching is case-insensitive.

The following partial date fields are also supported. The date will only be set on the scene if a date string can be built using the partial date components:

| Field | Remark |
|-------|--------|
|`yyyy`|Four digit year|
|`yy`|Two digit year. Assumes the first two digits are `20`|
|`mm`|Two digit month|
|`mmm`|Three letter month, such as `Jan` (case-insensitive)|
|`dd`|Two digit date|

The following full date fields are supported, using the same partial date rules as above:

* `yyyymmdd`
* `yymmdd`
* `ddmmyyyy`
* `ddmmyy`
* `mmddyyyy`
* `mmddyy`

All of these fields are available from the `Add Field` button.

Title generation also has the following options:

| Option | Remark |
|--------|--------|
|Whitespace characters| These characters are replaced with whitespace (defaults to `._`, to handle filenames like `three.word.title.avi`|
|Capitalize title| capitalises the first letter of each word|

The fields to display can be customised with the `Display Fields` drop-down section. By default, any field with new/different values will be displayed.

## Applying the results

Once the options are correct, click on the `Find` button. The system will search for scenes that have filenames that match the given pattern.

The results are presented in a table showing the existing and generated values of the discovered fields, along with a checkbox to determine whether or not the field will be set on each scene. These fields can also be edited manually.

The `Apply` button updates the scenes based on the set fields.

> **⚠️ Note:** results are paged and the `Apply` button only applies to scenes on the current page.
