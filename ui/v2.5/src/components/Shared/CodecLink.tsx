import React, { useMemo } from "react";
import { Link } from "react-router-dom";
import NavUtils from "src/utils/navigation";

interface ICommonLinkProps {
  link: string;
  className?: string;
}

const CommonLinkComponent: React.FC<ICommonLinkProps> = ({
  link,
  children,
}) => {
  return <Link to={link}>{children}</Link>;
};

interface ICodecLinkProps {
  codec: string;
  codecType: "audio_codec" | "video_codec";
  linkType?: "scene";
  className?: string;
}

export const CodecLink: React.FC<ICodecLinkProps> = ({
  codec,
  codecType,
  className,
}) => {
  const link = useMemo(() => {
    switch (codecType) {
      case "audio_codec":
        return NavUtils.makeCodecScenesUrl(codec, "audio_codec");
      case "video_codec":
        return NavUtils.makeCodecScenesUrl(codec, "video_codec");
    }
  }, [codec, codecType]);

  const title = codec || "";

  return (
    <CommonLinkComponent link={link} className={className}>
      <span>{title}</span>
    </CommonLinkComponent>
  );
};
