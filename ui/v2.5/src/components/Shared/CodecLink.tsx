import { Badge } from "react-bootstrap";
import React, { useMemo } from "react";
import { Link } from "react-router-dom";
import cx from "classnames";
import NavUtils from "src/utils/navigation";
import { faVolumeHigh, faVideo } from "@fortawesome/free-solid-svg-icons";
import { Icon } from "../Shared/Icon";

interface ICommonLinkProps {
  link: string;
  className?: string;
}

const CommonLinkComponent: React.FC<ICommonLinkProps> = ({
  link,
  className,
  children,
}) => {
  return (
    <Badge className={cx("tag-item", className)} variant="secondary">
      <Link to={link}>{children}</Link>
    </Badge>
  );
};

interface ICodecLinkProps {
  codec: string;
  codecType: "audio_codec" | "video_codec";
  fileIndex: number;
  filesLength: number;
  linkType?: "scene";
  className?: string;
}

export const CodecLink: React.FC<ICodecLinkProps> = ({
  codec,
  codecType,
  fileIndex,
  filesLength,
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

  fileIndex += 1;
  const index = JSON.stringify(fileIndex);

  return (
    <CommonLinkComponent link={link} className={className}>
      {filesLength > 1 ? `${index} / ${filesLength} | ` : ""}
      <span>{title} | </span>
      {codecType === "audio_codec" ? (
        <Icon icon={faVolumeHigh} className="tag-icon" />
      ) : codecType === "video_codec" ? (
        <Icon icon={faVideo} className="tag-icon" />
      ) : (
        ""
      )}
      <Link to={link}></Link>
    </CommonLinkComponent>
  );
};
