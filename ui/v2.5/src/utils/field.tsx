import React from "react";
import { FormattedMessage } from "react-intl";

interface ITextField {
  id?: string | undefined;
  name: string;
  value?: string | null;
}

export const TextField: React.FC<ITextField> = ({ id, name, value }) => {
  if (!value) {
    return null;
  }

  return (
    <dl className="row mb-0">
      <dt className="col-3 col-xl-2"><FormattedMessage id={id} defaultMessage={name} />:</dt>
      <dd className="col-9 col-xl-10">{value ?? undefined}</dd>
    </dl>
  );
};

interface IURLField {
  id?: string | undefined;
  name: string;
  value?: string | null;
  url?: string | null;
}

export const URLField: React.FC<IURLField> = ({ id, name, value, url }) => {
  if (!value) {
    return null;
  }
  return (
    <dl className="row mb-0">
      <dt className="col-3 col-xl-2"><FormattedMessage id={id} defaultMessage={name} />:</dt>
      <dd className="col-9 col-xl-10">
        {url ? (
          <a href={url} target="_blank" rel="noopener noreferrer">
            {value}
          </a>
        ) : undefined}
      </dd>
    </dl>
  );
};
