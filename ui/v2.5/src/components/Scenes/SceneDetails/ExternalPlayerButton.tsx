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

  const { stream, stream_org } = paths;
  const sceneTitle = title ?? TextUtils.fileNameFromPath(path);

  let url, prompt;
  const streamURL = new URL(stream);
  if (isAndroid) {
    const scheme = streamURL.protocol.slice(0, -1);
    streamURL.hash = `Intent;action=android.intent.action.VIEW;scheme=${scheme};type=video/mp4;S.title=${encodeURI(
      sceneTitle
    )};end`;
    streamURL.protocol = "intent";
    url = streamURL.toString();
    prompt = "Open in Android's media player"
  } else if (isAppleDevice) {
    streamURL.host = "x-callback-url";
    streamURL.port = "";
    streamURL.pathname = "stream";
    streamURL.search = `url=${encodeURIComponent(stream)}`;
    streamURL.protocol = "vlc-x-callback";
    url = streamURL.toString();
    prompt = "Open in iOS media player"
  } else {
    url = stream_org;
    prompt = "Drag this link to an external media player, e.g. VLC, MPV.\nIf you are using DeoVR, click the button directly." 
  }

  return (
        <div className="btn-group btn" title="Click here if you're using Android/iOS/DeoVR/HereSphere.&#10;In Windows or MacOS, drag me to an external media player, e.g. VLC, MPV.">
          <a href={url} target="_self">
              <span><Icon icon="external-link-alt" color="white" /></span>
              <span className="ml-2">Launch &#47; Drag Me</span>
          </a>
        </div>
  );
};
