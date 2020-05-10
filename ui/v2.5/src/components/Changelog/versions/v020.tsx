import React from "react";
import ReactMarkdown from "react-markdown";

const markup = `
💥 **Note: After upgrading performance will be degraded until a full [scan](/settings?tab=tasks) has been completed.**
&nbsp;

#### Major Changes:
* ✨ [Movies](/movies) are now supported.
* 💄 Responsive layout for mobile phones.
* ⚡️ Performance improvements and improved video support.
* 📝 Support for localized text, dates and numbers.

#### Full list of changes:

* ✨ Add support for image scraping.
* ✨ Allow user to regenerate scene cover based on timestamp.
* ♻️ Replace Blueprint with react-bootstrap.
* 🐛 Update performer image in UI when it's replaced.
* 🐛 Fix performer height filter.
* 🐛 Fix error when viewing scenes related to objects with illegal characters in name.
* ✨ Autoassociate galleries to scenes when scanning.
* ✨ Configurable scraper user agent string.
* ✨ Backup database if a migration is needed.
* ✨ Add modes for performer/tag for bulk scene editing.
* ✨ Add gender support for performer.
* 🐛 Fix to allow scene to be removed when attached to a movie.
* 🐛 Make ethnicity freetext and fix freeones ethnicity panic.
* 💄 Add image count to gallery list.
* ✨ Add SVG studio image support, and studio image caching.
* 🐛 Fix to filter on movies from performer filter to movie filter.
* 🎨 Update prettier to v2.0.1 and enable for SCSS.
* 💄 Add library size to main stats page.
* ✨ Enable sorting for galleries.
* ✨ Add scene rating to scene filename parser.
* ✨ Replace basic auth with cookie authentication.
* 🐛 Added various missing filters to performer page.
* ✨ Add detection of container/video_codec/audio_codec compatibility for live file streaming or transcoding.
* 🐛 Delete marker preview on marker change or delete.
* 🐛 Prefer modified performer image over scraped one.
* 🐛 Don't redirect login to migrate page.
* 🐛 Performer and Movie UI fixes and improvements.
* 🐛 Include gender in performer scraper results.
* ⚡️ Add slim endpoints for entities to speed up filters.
* 🐛 Include scene o-counter in import/export.
* ⚡️ Export performance optimization.
* ✨ Move image with cover.jpg in name to first place in Galleries.
* ✨ Add "reshuffle button" when sortby is random.
* ✨ Implement clean for missing galleries.
* 💄 Add random male performer image.
* ✨ Add parser support for 3-letter month.
* 💄 Add index/total count to end of pagination buttons.
* 🐛 Make image extension check in zip files case insensitive.
* ♻️ Refactor build.
* 💄 Add flags for performer countries.
* 🐛 Freeones scraper tweaks.
* ✅ Querybuilder integration tests.
* ✨ Add is-missing tags filter.
* 💄 Overhaul look and feel of folder select.
* 📝 Add changelog to start page.
`;

export default () => <ReactMarkdown source={markup} />;
