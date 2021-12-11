import React from "react";
import { Button } from "react-bootstrap";
import { Icon } from "src/components/Shared";
import { SceneDataFragment } from "src/core/generated-graphql";
import { TextUtils } from "src/utils";

export interface IExternalPlayerButtonProps {
  scene: SceneDataFragment;
}

export const ExternalPlayerButton: React.FC<IExternalPlayerButtonProps> = ({
  scene,
}) => {
  const isAndroid = /(android)/i.test(navigator.userAgent);
  const isAppleDevice = /(ipod|iphone|ipad)/i.test(navigator.userAgent);

  const { paths, path, title } = scene;

  if (!paths || !paths.stream || (!isAndroid && !isAppleDevice))
    return <span />;

  const { stream } = paths;
  const sceneTitle = title ?? TextUtils.fileNameFromPath(path);

  let url;
  const streamURL = new URL(stream);
  if (isAndroid) {
    const scheme = streamURL.protocol.slice(0, -1);
    streamURL.hash = `Intent;action=android.intent.action.VIEW;scheme=${scheme};type=video/mp4;S.title=${encodeURI(
      sceneTitle
    )};end`;
    streamURL.protocol = "intent";
    url = streamURL.toString();
  } else if (isAppleDevice) {
    streamURL.host = "x-callback-url";
    streamURL.port = "";
    streamURL.pathname = "stream";
    streamURL.search = `url=${encodeURIComponent(stream)}`;
    streamURL.protocol = "vlc-x-callback";
    url = streamURL.toString();
  }

  return (
    <Button
      className="minimal"
      variant="secondary"
      title="Open in external player"
    >
      <a href={url}>
        <Icon icon="external-link-alt" color="white" />
      </a>
    </Button>
  );
};
