import React, { useEffect, useMemo, useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { Remark } from "react-remark";
import remarkGfm from "remark-gfm";
import * as GQL from "src/core/generated-graphql";
import { TextField, URLField } from "src/utils/field";
import { getStoryMetadata } from "./storyMetadata";

interface IProps {
  gallery: GQL.GalleryDataFragment;
}

const markdownExtensions = new Set(["md", "markdown", "rst"]);

function fileExtension(path?: string) {
  if (!path) {
    return "";
  }

  const match = path.toLowerCase().match(/\.([^.]+)$/);
  return match?.[1] ?? "";
}

function looksLikeBBCode(text: string) {
  return /\[(?:b|i|s|u|quote|code|url|list|\*|img|email|color|size|font|center|left|right|indent|spoiler|h[1-6]|hr|br)[^\]]*\]/i.test(
    text
  );
}

function toBlockQuote(body: string) {
  return body
    .split(/\r?\n/)
    .map((line) => `> ${line}`)
    .join("\n");
}

function toList(body: string, ordered = false) {
  return body
    .split("[*]")
    .map((item) => item.trim())
    .filter(Boolean)
    .map((item, index) => `${ordered ? `${index + 1}.` : "-"} ${item}`)
    .join("\n");
}

function bbcodeToMarkdown(text: string) {
  return text
    .replace(/\r\n/g, "\n")
    .replace(/\[b\]([\s\S]*?)\[\/b\]/gi, "**$1**")
    .replace(/\[i\]([\s\S]*?)\[\/i\]/gi, "*$1*")
    .replace(/\[s\]([\s\S]*?)\[\/s\]/gi, "~~$1~~")
    .replace(/\[u\]([\s\S]*?)\[\/u\]/gi, "<ins>$1</ins>")
    .replace(/\[h([1-6])\]([\s\S]*?)\[\/h\1\]/gi, (_, level, body: string) => {
      return `${"#".repeat(Number(level))} ${body.trim()}`;
    })
    .replace(/\[br\s*\/?\]/gi, "  \n")
    .replace(/\[hr\s*\/?\]/gi, "\n---\n")
    .replace(/\[url=([^\]]+)\]([\s\S]*?)\[\/url\]/gi, "[$2]($1)")
    .replace(/\[url\]([\s\S]*?)\[\/url\]/gi, "<$1>")
    .replace(/\[email=([^\]]+)\]([\s\S]*?)\[\/email\]/gi, "[$2](mailto:$1)")
    .replace(/\[email\]([\s\S]*?)\[\/email\]/gi, "[$1](mailto:$1)")
    .replace(/\[img\]([\s\S]*?)\[\/img\]/gi, "![]($1)")
    .replace(/\[quote=([^\]]+)\]([\s\S]*?)\[\/quote\]/gi, (_, author, body: string) =>
      `> ${author} wrote:\n${toBlockQuote(body)}`
    )
    .replace(/\[quote\]([\s\S]*?)\[\/quote\]/gi, (_, body: string) =>
      toBlockQuote(body)
    )
    .replace(/\[code\]([\s\S]*?)\[\/code\]/gi, "\n```\n$1\n```\n")
    .replace(/\[spoiler\]([\s\S]*?)\[\/spoiler\]/gi, "> Spoiler:\n> $1")
    .replace(/\[(?:center|left|right)\]([\s\S]*?)\[\/(?:center|left|right)\]/gi, "$1")
    .replace(/\[indent\]([\s\S]*?)\[\/indent\]/gi, (_, body: string) =>
      toBlockQuote(body)
    )
    .replace(/\[(?:color|size|font)=[^\]]+\]([\s\S]*?)\[\/(?:color|size|font)\]/gi, "$1")
    .replace(/\[(?:color|size|font)\]([\s\S]*?)\[\/(?:color|size|font)\]/gi, "$1")
    .replace(/\[list=1\]([\s\S]*?)\[\/list\]/gi, (_, body: string) =>
      toList(body, true)
    )
    .replace(/\[list=a\]([\s\S]*?)\[\/list\]/gi, (_, body: string) =>
      toList(body, true)
    )
    .replace(/\[list\]([\s\S]*?)\[\/list\]/gi, (_, body: string) =>
      toList(body)
    );
}

export const GalleryTextPanel: React.FC<IProps> = ({ gallery }) => {
  const intl = useIntl();
  const [content, setContent] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const storyMetadata = getStoryMetadata(gallery);

  const primaryPath = gallery.files[0]?.path;
  const ext = fileExtension(primaryPath);
  const isMarkdown = markdownExtensions.has(ext);
  const renderAsMarkdown = isMarkdown || looksLikeBBCode(content);

  useEffect(() => {
    setContent("");
    setError("");
    setLoading(false);

    if (!gallery.paths.text) {
      return;
    }

    setLoading(true);
    fetch(gallery.paths.text)
      .then(async (res) => {
        if (!res.ok) {
          throw new Error(`Failed to load story text (${res.status})`);
        }

        return res.text();
      })
      .then((text) => setContent(text))
      .catch((e: Error) => setError(e.message))
      .finally(() => setLoading(false));
  }, [gallery.paths.text]);

  const markdown = useMemo(() => bbcodeToMarkdown(content), [content]);

  if (error) {
    return <div className="story-reader-error">{error}</div>;
  }

  if (loading) {
    return <div className="story-reader-empty" />;
  }

  if (!content) {
    return <div className="story-reader-empty" />;
  }

  return (
    <div className="story-reader">
      <h6>
        <FormattedMessage id="reader" />
      </h6>
      {storyMetadata.tagLine ? (
        <p className="story-reader-tag-line">{storyMetadata.tagLine}</p>
      ) : undefined}
      {storyMetadata.author ||
      storyMetadata.language ||
      storyMetadata.sourceUrl ||
      storyMetadata.sourceWebsite ||
      storyMetadata.audioUrl ? (
        <dl className="container details-list story-metadata-list">
          <TextField id="author" value={storyMetadata.author} />
          <TextField id="language" value={storyMetadata.language} />
          <URLField
            id="source_url"
            value={storyMetadata.sourceUrl}
            url={storyMetadata.sourceUrl}
            truncate
          />
          <TextField
            id="source_website"
            value={storyMetadata.sourceWebsite}
          />
          <URLField
            id="audio_link"
            value={storyMetadata.audioUrl}
            url={storyMetadata.audioUrl}
            truncate
          />
        </dl>
      ) : undefined}
      {storyMetadata.audioUrl ? (
        <div className="story-reader-audio-container">
          <audio
            className="story-reader-audio"
            controls
            preload="metadata"
            src={storyMetadata.audioUrl}
          />
        </div>
      ) : undefined}
      {renderAsMarkdown ? (
        <div className="markdown story-reader-content">
          <Remark remarkPlugins={[remarkGfm]}>{markdown}</Remark>
        </div>
      ) : (
        <pre className="story-reader-content">{content}</pre>
      )}
      {storyMetadata.backCoverUrl ? (
        <div className="story-back-cover">
          <h6>
            <FormattedMessage id="back_cover" />
          </h6>
          <img
            src={storyMetadata.backCoverUrl}
            alt={intl.formatMessage({ id: "back_cover" })}
          />
        </div>
      ) : undefined}
    </div>
  );
};
