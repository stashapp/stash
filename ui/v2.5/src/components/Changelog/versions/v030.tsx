import React from "react";
import ReactMarkdown from "react-markdown";

const markup = `
### ✨ New Features
*  Add support for parent/child studios.

### 🎨 Improvements
*  Show rating as stars in scene page.
*  Add reload scrapers button.

`;

export default () => <ReactMarkdown source={markup} />;
