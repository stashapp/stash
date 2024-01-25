import React from "react";
import { FormattedMessage } from "react-intl";
import { ExternalLink } from "src/components/Shared/ExternalLink";
import { TruncatedText } from "src/components/Shared/TruncatedText";

interface ITextField {
  id?: string;
  name?: string;
  abbr?: string | null;
  value?: string | null;
  truncate?: boolean | null;
}

export const TextField: React.FC<ITextField> = ({
  id,
  name,
  value,
  abbr,
  truncate,
  children,
}) => {
  if (!value && !children) {
    return null;
  }

  const message = (
    <>{id ? <FormattedMessage id={id} defaultMessage={name} /> : name}:</>
  );

  return (
    <>
      <dt>{abbr ? <abbr title={abbr}>{message}</abbr> : message}</dt>
      <dd>
        {value ? truncate ? <TruncatedText text={value} /> : value : children}
      </dd>
    </>
  );
};

interface IURLField {
  id?: string;
  name?: string;
  abbr?: string | null;
  value?: string | null;
  url?: string | null;
  truncate?: boolean | null;
  target?: string;
  // use for internal links
  trusted?: boolean;
}

export const URLField: React.FC<IURLField> = ({
  id,
  name,
  value,
  url,
  abbr,
  truncate,
  target = "_blank",
  trusted,
}) => {
  if (!value) {
    return null;
  }

  const message = (
    <>{id ? <FormattedMessage id={id} defaultMessage={name} /> : name}:</>
  );

  function maybeRenderUrl() {
    if (!url) return;

    const children = truncate ? <TruncatedText text={value} /> : value;

    if (trusted) {
      return (
        <a href={url} target={target}>
          {children}
        </a>
      );
    } else {
      return (
        <ExternalLink href={url} target={target}>
          {children}
        </ExternalLink>
      );
    }
  }

  return (
    <>
      <dt>{abbr ? <abbr title={abbr}>{message}</abbr> : message}</dt>
      <dd>{maybeRenderUrl()}</dd>
    </>
  );
};

interface IURLsField {
  id?: string;
  name?: string;
  abbr?: string | null;
  urls?: string[] | null;
  truncate?: boolean | null;
  target?: string;
  // use for internal links
  trusted?: boolean;
}

export const URLsField: React.FC<IURLsField> = ({
  id,
  name,
  urls,
  abbr,
  truncate,
  target = "_blank",
  trusted,
}) => {
  if (!urls || !urls.length) {
    return null;
  }

  const message = (
    <>{id ? <FormattedMessage id={id} defaultMessage={name} /> : name}:</>
  );

  const renderUrls = () => {
    return urls.map((url, i) => {
      if (!url) return;

      const children = truncate ? <TruncatedText text={url} /> : url;

      if (trusted) {
        return (
          <a key={i} href={url} target={target}>
            {children}
          </a>
        );
      } else {
        return (
          <ExternalLink key={i} href={url} target={target}>
            {children}
          </ExternalLink>
        );
      }
    });
  };

  return (
    <>
      <dt>{abbr ? <abbr title={abbr}>{message}</abbr> : message}</dt>
      <dd>
        <dl>{renderUrls()}</dl>
      </dd>
    </>
  );
};
