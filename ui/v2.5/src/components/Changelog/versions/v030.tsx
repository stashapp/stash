import React from "react";
import ReactMarkdown from "react-markdown";

const markup = `
### ✨ New Features
*  Add support for parent/child studios.

### 🎨 Improvements
*  Show pagination at top as well as bottom of the page.
*  Add split xpath post-processing action.
*  Improved the layout of the scene page.
*  Show rating as stars in scene page.
*  Add reload scrapers button.

`;

export default () => <ReactMarkdown source={markup} />;
