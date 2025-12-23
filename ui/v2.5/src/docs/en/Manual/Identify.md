# Identify

The Identify task iterates through your Scenes and attempts to identify them using a selection of scraping sources. If a result is found in a source, the Scene is updated, and no further sources are checked for that scene.

This task is part of the advanced settings mode.

## Rules

- The task accepts one or more scraper sources, including stash-box instances and scene scrapers that support scraping via Scene Fragment. The order of the sources can be rearranged.
- The task iterates through the scraper sources in the provided order.
- If a result is found in a source, the Scene is updated, and further sources are not checked for that scene.

### Organized flag

Scenes that have the Organized flag added to them will not be modified by Identify. You can also use Organized flag status as a filter.

## Options

The following options can be configured:

| Option | Description |
|--------|-------------|
| Include male performers | If false, male performers will not be created or set on scenes. |
| Set cover images | If false, scene cover images will not be modified. |
| Set organized flag | If true, the organized flag is set to true when a scene is organized. |
| Skip matches that have more than one result | If this is not enabled and more than one result is returned, one will be randomly chosen to match |
| Tag skipped matches with | If the above option is set and a scene is skipped, this will add the tag so that you can filter for it in the Scene Tagger view and choose the correct match by hand |
| Skip single name performers with no disambiguation | If this is not enabled, performers that are often generic like Samantha or Olga will be matched |
| Tag skipped performers with | If the above option is set and a performer is skipped, this will add the tag so that you can filter for it in the Scene Tagger view and choose how you want to handle those performers |

### Field specific options

Each field may have a strategy. The behavior for each strategy is as follows:

| Strategy | Description |
|----------|-------------|
| Ignore | The field is not set. |
| Overwrite | Existing values are overwritten. |
| Merge (*default*) | For multi-value fields, adds to existing values. For single-value fields, only sets if not already set. |

For Studio, Performers, and Tags, an option is available to **Create Missing Objects**. This is enabled by default. When true, if a Studio/Performer/Tag is included during the identification process and does not exist in the system, it will be created.

## Running task

- **Identify...:** Run the Identify task on your entire library from the Tasks page.
- **Selective Identify:** Configure and run the Identify task on specific directories from Tasks > Identify... page. At the top of the page click folder icon to select directories.

## Logs

The result of the identification process for each scene is output to the log.
