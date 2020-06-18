import React from "react";
import ReactMarkdown from "react-markdown";

const markup = `
### ✨ New Features
*  Add support for parent/child studios.

### 🎨 Improvements
*  Improved the layout of the scene page.
*  Show rating as stars in scene page.
*  Add reload scrapers button.

`;

export default () => <ReactMarkdown source={markup} />;
