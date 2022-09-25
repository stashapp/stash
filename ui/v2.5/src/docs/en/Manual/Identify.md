# Identify

This task iterates through your Scenes and attempts to identify the scene using a selection of scraping sources.

This task accepts one or more scraper sources. Valid scraper sources for the Identify task are stash-box instances, and scene scrapers which support scraping via Scene Fragment. The order of the sources may be rearranged.

For each Scene, the Identify task iterates through the scraper sources, in the order provided, and tries to identify the scene using each source. If a result is found in a source, then the Scene is updated, and no further sources are checked for that scene.

## Options

The following options can be set:

| Option | Description |
|--------|-------------|
| Include male performers | If false, then male performers will not be created or set on scenes. |
| Set cover images | If false, then scene cover images will not be modified. |
| Set organised flag | If true, the organised flag is set to true when a scene is organised. |

Field specific options may be set as well. Each field may have a Strategy. The behaviour for each strategy value is as follows:

| Strategy | Description |
|----------|-------------|
| Ignore | Not set. |
| Overwrite | Overwrite existing value. |
| Merge (*default*) | For multi-value fields, adds to existing values. For single-value fields, only sets if not already set. |

For Studio, Performers and Tags, an option is also available to Create Missing objects. This is false by default. When true, if a Studio/Performer/Tag is included during the identification process and does not exist in the system, then it will be created.

Default Options are applied to all sources unless overridden in specific source options. 

The result of the identification process for each scene is output to the log.
