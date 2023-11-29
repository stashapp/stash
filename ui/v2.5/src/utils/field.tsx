import React from "react";
import { FormattedMessage } from "react-intl";
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
  children,
  target,
  trusted,
}) => {
  if (!value && !children) {
    return null;
  }

  const message = (
    <>{id ? <FormattedMessage id={id} defaultMessage={name} /> : name}:</>
  );

  const rel = !trusted ? "noopener noreferrer" : undefined;

  return (
    <>
      <dt>{abbr ? <abbr title={abbr}>{message}</abbr> : message}</dt>
      <dd>
        {url ? (
          <a href={url} target={target || "_blank"} rel={rel}>
            {value ? (
              truncate ? (
                <TruncatedText text={value} />
              ) : (
                value
              )
            ) : (
              children
            )}
          </a>
        ) : undefined}
      </dd>
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
  target,
  trusted,
}) => {
  const values = urls ?? [];
  if (!values.length) {
    return null;
  }

  const message = (
    <>{id ? <FormattedMessage id={id} defaultMessage={name} /> : name}:</>
  );

  const rel = !trusted ? "noopener noreferrer" : undefined;

  return (
    <>
      <dt>{abbr ? <abbr title={abbr}>{message}</abbr> : message}</dt>
      <dd>
        <dl>
          {values.map((url, i) => (
            <a key={i} href={url} target={target || "_blank"} rel={rel}>
              {truncate ? <TruncatedText text={url} /> : url}
            </a>
          ))}
        </dl>
      </dd>
    </>
  );
};
