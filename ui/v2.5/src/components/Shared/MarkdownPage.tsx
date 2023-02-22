import React, { useEffect, useState } from "react";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";

interface IPageProps {
  // page is a markdown module
  page: string;
}

export const MarkdownPage: React.FC<IPageProps> = ({ page }) => {
  const [markdown, setMarkdown] = useState("");

  useEffect(() => {
    if (!markdown) {
      fetch(page)
        .then((res) => res.text())
        .then((text) => setMarkdown(text));
    }
  }, [page, markdown]);

  return (
    <ReactMarkdown className="markdown" remarkPlugins={[remarkGfm]}>
      {markdown}
    </ReactMarkdown>
  );
};
