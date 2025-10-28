import React, { PropsWithChildren } from "react";

export const DetailTitle: React.FC<
  PropsWithChildren<{
    name: string;
    disambiguation?: string;
    classNamePrefix: string;
  }>
> = ({ name, disambiguation, classNamePrefix, children }) => {
  return (
    <h2>
      <span className={`${classNamePrefix}-name`}>{name}</span>
      {disambiguation && (
        <span className={`${classNamePrefix}-disambiguation`}>
          {` (${disambiguation})`}
        </span>
      )}
      {children}
    </h2>
  );
};
