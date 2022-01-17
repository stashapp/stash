import React from "react";
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
  const scenePath = scene.path;
  const sceneTitle = scene.title ?? TextUtils.fileNameFromPath(scenePath);
  const sceneStream = scene.paths.stream ?? "";
  const sceneStream_org = scene.paths.stream_org ?? "";

  let url, prompt;
  const streamURL = new URL(sceneStream);
  if (isAndroid) {
    const scheme = streamURL.protocol.slice(0, -1);
    streamURL.hash = `Intent;action=android.intent.action.VIEW;scheme=${scheme};type=video/mp4;S.title=${encodeURI(
      sceneTitle
    )};end`;
    streamURL.protocol = "intent";
    url = streamURL.toString();
    prompt = "Click here to open the video in Android's default media player";
  } else if (isAppleDevice) {
    streamURL.host = "x-callback-url";
    streamURL.port = "";
    streamURL.pathname = "stream";
    streamURL.search = `url=${encodeURIComponent(sceneStream)}`;
    streamURL.protocol = "vlc-x-callback";
    url = streamURL.toString();
    prompt = "Click here to open the video in iOS's default media player";
  } else {
    url = sceneStream_org;
    prompt =
      "Click here if you're using Android/iOS/DeoVR/HereSphere.\n" +
      "In Windows or MacOS, drag me to an external media player, e.g. VLC, MPV.";
  }

  return (
    <div className="btn-group btn" title={prompt}>
      <a href={url} target="_self">
        <span>
          <Icon icon="external-link-alt" color="white" />
        </span>
        <span className="ml-2">Launch / Drag Me</span>
      </a>
    </div>
  );
};
