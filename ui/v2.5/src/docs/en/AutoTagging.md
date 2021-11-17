# Auto Tagging

This task matches your Performers, Studios, and Tags against your media, based on names only. It finds Scenes, Images, and Galleries where the path or filename contains the Performer/Studio/Tag. 

For each scene it finds that matches, it sets the applicable field. It will **only** tag based on performers, studios, and tags that already exist in your database. In order to completely identify and gather information about the scenes in your collection, you will need to use the Tagger view and/or Scraping tools.

When the Performer/Studio/Tag name has multiple words, the search will include paths/filenames where the Performer/Studio/Tag name is separated with `.`, `-` or `_` characters, as well as whitespace.

For example, auto tagging for performer `Jane Doe` will match the following filenames:
* `Jane.Doe.1.mp4`
* `Jane_Doe.2.mp4`
* `Jane-Doe.3.mp4`
* `Jane Doe.4.mp4`

Matching is case insensitive, and should only match exact wording within word boundaries. For example, `Jane Doe` will not match `Maryjane-Doe`, but will match `Mary-Jane-Doe`.

Auto tagging for only specific Performers, Studios and Tags can be performed from the individual Performer/Studio/Tag page.
