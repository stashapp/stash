import React from "react";
import ReactMarkdown from "react-markdown";

const markup = `
### 🐛 Bug fixes
Fix version checking.
`;

export default () => <ReactMarkdown source={markup} />;
