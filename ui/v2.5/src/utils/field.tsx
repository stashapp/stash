import React from "react";
import { FormattedMessage } from "react-intl";
import { TruncatedText } from "../components/Shared";

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
}

export const URLField: React.FC<IURLField> = ({
  id,
  name,
  value,
  url,
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
        {url ? (
          <a href={url} target="_blank" rel="noopener noreferrer">
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
