import React, { PropsWithChildren } from "react";

interface IProps {
  className?: string;
  header: string;
  link: JSX.Element;
}

export const RecommendationRow: React.FC<PropsWithChildren<IProps>> = ({
  className,
  header,
  link,
  children,
}) => (
  <div className={`recommendation-row ${className}`}>
    <div className="recommendation-row-head">
      <div>
        <h2>{header}</h2>
      </div>
      {link}
    </div>
    {children}
  </div>
);
