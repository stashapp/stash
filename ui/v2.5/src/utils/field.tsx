import React from "react";

interface ITextField {
  name: string;
  value?: string | null;
}

export const TextField: React.FC<ITextField> = ({ name, value }) => {
  if (!value) {
    return null;
  }

  return (
    <dl className="row mb-0">
      <dt className="col-3 col-xl-2">{name}:</dt>
      <dd className="col-9 col-xl-10">{value ?? undefined}</dd>
    </dl>
  );
};

interface IURLField {
  name: string;
  value?: string | null;
  url?: string | null;
}

export const URLField: React.FC<IURLField> = ({ name, value, url }) => {
  if (!value) {
    return null;
  }
  return (
    <dl className="row mb-0">
      <dt className="col-3 col-xl-2">{name}:</dt>
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
