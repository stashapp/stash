import React from "react";
import { Button } from "react-bootstrap";
import { Icon } from "src/components/Shared";
import { SceneDataFragment } from "src/core/generated-graphql";
import { TextUtils } from "src/utils";

export interface IExternalPlayerButtonProps {
  scene?: SceneDataFragment;
}

export const ExternalPlayerButton: React.FC<IExternalPlayerButtonProps> = ({ scene }) => {
  if (!scene) return <span />;

  const icon = <Icon icon="external-link-alt" color="white" />;
  const m3u = Buffer.from(`#EXTM3U\n#EXTINF:${scene.file.duration},${scene.title ?? TextUtils.fileNameFromPath(scene.path)}\n${scene.paths.stream}`).toString('base64');
  let link;
  if(/(android)/i.test(navigator.userAgent)) {
    const scheme = scene.paths?.stream?.match(/https?/)?.[0] ?? '';
    link =  <a href={`intent${scene.paths?.stream?.slice(scheme.length)}#Intent;action=android.intent.action.VIEW;scheme=${scheme};type=video/mp4;end`}>{icon}</a>;
  }
  else if(/(ipod|iphone|ipad)/i.test(navigator.userAgent))
    link =  <a href={`data:application/vnd.apple.mpegurl;base64,${m3u}`}>{icon}</a>;
  else
    link =  <a href={`data:audio/x-mpegurl;base64,${m3u}`}>{icon}</a>;

  return (
    <Button
      className="minimal"
      variant="secondary"
      title="Open in external player"
    >
      {link}
    </Button>
  );
};
