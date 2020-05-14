import React from "react";
import ReactMarkdown from "react-markdown";

const markup = `
#### 💥 **Note: After upgrading performance will be degraded until a full [scan](/settings?tab=tasks) has been completed.**

&nbsp;
### ✨ New Features
*  Movies are now supported.
*  Responsive layout for mobile phones.
*  Add support for image scraping.
*  Allow user to regenerate scene cover based on timestamp.
*  Autoassociate galleries to scenes when scanning.
*  Configurable scraper user agent string.
*  Backup database if a migration is needed.
*  Add modes for performer/tag for bulk scene editing.
*  Add gender support for performer.
*  Add SVG studio image support, and studio image caching.
*  Enable sorting for galleries.
*  Add scene rating to scene filename parser.
*  Replace basic auth with cookie authentication.
*  Add detection of container/video_codec/audio_codec compatibility for live file streaming or transcoding.
*  Move image with cover.jpg in name to first place in Galleries.
*  Add "reshuffle button" when sortby is random.
*  Implement clean for missing galleries.
*  Add parser support for 3-letter month.
*  Add is-missing tags filter.

### 🎨 Improvements
*  Performance improvements and improved video support.
*  Support for localized text, dates and numbers.
*  Replace Blueprint with react-bootstrap.
*  Add image count to gallery list.
*  Update prettier to v2.0.1 and enable for SCSS.
*  Add library size to main stats page.
*  Add slim endpoints for entities to speed up filters.
*  Export performance optimization.
*  Add random male performer image.
*  Added various missing filters to performer page.
*  Add index/total count to end of pagination buttons.
*  Refactor build.
*  Add flags for performer countries.
*  Querybuilder integration tests.
*  Overhaul look and feel of folder select.
*  Add changelog to start page.

### 🐛 Bug fixes
*  Update performer image in UI when it's replaced.
*  Fix performer height filter.
*  Fix error when viewing scenes related to objects with illegal characters in name.
*  Fix to allow scene to be removed when attached to a movie.
*  Make ethnicity freetext and fix freeones ethnicity panic.
*  Fix to filter on movies from performer filter to movie filter.
*  Delete marker preview on marker change or delete.
*  Prefer modified performer image over scraped one.
*  Don't redirect login to migrate page.
*  Performer and Movie UI fixes and improvements.
*  Include gender in performer scraper results.
*  Include scene o-counter in import/export.
*  Make image extension check in zip files case insensitive.
*  Freeones scraper tweaks.

`;

export default () => <ReactMarkdown source={markup} />;
