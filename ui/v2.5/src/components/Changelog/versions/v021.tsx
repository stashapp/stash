import React from "react";
import ReactMarkdown from "react-markdown";

const markup = `
### 🐛 Bug fixes
*  Fix max loop duration not working.
*  Fix URL sanitization on non-Chrome browsers.

`;

export default () => <ReactMarkdown source={markup} />;
