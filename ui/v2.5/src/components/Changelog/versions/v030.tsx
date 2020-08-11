import React from "react";
import ReactMarkdown from "react-markdown";

const markup = `
#### 💥 **Note: After upgrading, the next scan will populate all scenes with oshash hashes. MD5 calculation can be disabled after populating the oshash for all scenes. See \`Hashing Algorithms\` in the \`Configuration\` section of the manual for details. **

### ✨ New Features
*  Add support for scraping movie details.
*  Add support for JSON scrapers.
*  Add support for plugin tasks.
*  Add oshash algorithm for hashing scene video files. Enabled by default on new systems.
*  Support (re-)generation of generated content for specific scenes.
*  Add tag thumbnails, tags grid view and tag page.
*  Add post-scrape dialog.
*  Add various keyboard shortcuts (see manual).
*  Support deleting multiple scenes.
*  Add in-app help manual.
*  Add support for custom served folders.
*  Add support for parent/child studios.

### 🎨 Improvements
*  Allow free-editing of scene movie number.
*  Allow adding performers and studios from selectors.
*  Add support for chrome dp in xpath scrapers.
*  Allow customisation of preview video generation.
*  Add support for live transcoding in Safari.
*  Add mapped and fixed post-processing scraping options.
*  Add random sorting for performers.
*  Search for files which have low or upper case supported filename extensions.
*  Add dialog when pasting movie images.
*  Allow click and click-drag selection after selecting scene.
*  Added multi-scene edit dialog.
*  Moved images to separate tables, increasing performance.
*  Add gallery grid view.
*  Add is-missing scene filter for gallery query.
*  Don't import galleries with no images, and delete galleries with no images during clean.
*  Show pagination at top as well as bottom of the page.
*  Add split xpath post-processing action.
*  Improved the layout of the scene page.
*  Show rating as stars in scene page.
*  Add reload scrapers button.

### 🐛 Bug fixes
*  Fix formatted dates using incorrect timezone.

`;

export default () => <ReactMarkdown source={markup} />;
