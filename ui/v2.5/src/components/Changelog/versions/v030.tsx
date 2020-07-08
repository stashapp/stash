import React from "react";
import ReactMarkdown from "react-markdown";

const markup = `
### ✨ New Features
*  Add tag thumbnails, tags grid view and tag page.
*  Add post-scrape dialog.
*  Add various keyboard shortcuts (see manual).
*  Support deleting multiple scenes.
*  Add in-app help manual.
*  Add support for custom served folders.
*  Add support for parent/child studios.

### 🎨 Improvements
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

`;

export default () => <ReactMarkdown source={markup} />;
