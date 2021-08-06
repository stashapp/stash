import React from "react";
import { FormattedMessage } from "react-intl";

interface ITextField {
  id?: string;
  name?: string;
  value?: string | null;
}

export const TextField: React.FC<ITextField> = ({ id, name, value }) => {
  if (!value) {
    return null;
  }
  return (
    <>
      <dt>{id ? <FormattedMessage id={id} defaultMessage={name} /> : name}:</dt>
      <dd>{value ?? undefined}</dd>
    </>
  );
};

interface IURLField {
  id?: string;
  name?: string;
  value?: string | null;
  url?: string | null;
}

export const URLField: React.FC<IURLField> = ({ id, name, value, url }) => {
  if (!value) {
    return null;
  }
  return (
    <>
      <dt>{id ? <FormattedMessage id={id} defaultMessage={name} /> : name}:</dt>
      <dd>
        {url ? (
          <a href={url} target="_blank" rel="noopener noreferrer">
            {value}
          </a>
        ) : undefined}
      </dd>
    </>
  );
};
