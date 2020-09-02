import React, { useEffect, useState } from "react";
import ReactMarkdown from "react-markdown";

interface IPageProps {
  // page is a markdown module
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  page: any;
}

export const Page: React.FC<IPageProps> = ({ page }) => {
  const [markdown, setMarkdown] = useState("");

  useEffect(() => {
    if (!markdown) {
      fetch(page)
        .then((res) => res.text())
        .then((text) => setMarkdown(text));
    }
  }, [page, markdown]);

  return <ReactMarkdown className="markdown" source={markdown} />;
};
