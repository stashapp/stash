import React from "react";
import ReactMarkdown from "react-markdown";

const markup = `
### 🎨 Improvements
*  Show rating as stars in scene page.
*  Add reload scrapers button.

`;

export default () => <ReactMarkdown source={markup} />;
