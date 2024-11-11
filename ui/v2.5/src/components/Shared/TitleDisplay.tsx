import React from "react";

interface ITitleDisplayProps {
  text: string;
  className?: string;
}

export const TitleDisplay: React.FC<ITitleDisplayProps> = ({
  text,
  className,
}) => {
  if (!text) return null;

  return (
    <div className={className}>
      {text.split("\n").map((line, i) => (
        <React.Fragment key={i}>
          {line}
          {i < text.split("\n").length - 1 && <br />}
        </React.Fragment>
      ))}
    </div>
  );
};
