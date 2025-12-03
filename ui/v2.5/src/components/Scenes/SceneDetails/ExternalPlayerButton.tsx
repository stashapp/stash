import { faExternalLinkAlt } from "@fortawesome/free-solid-svg-icons";
import React from "react";
import { Button } from "react-bootstrap";
import { useIntl } from "react-intl";
import { Icon } from "src/components/Shared/Icon";
import { objectTitle } from "src/core/files";
import { SceneDataFragment } from "src/core/generated-graphql";

export interface IExternalPlayerButtonProps {
  scene: SceneDataFragment;
}

export const ExternalPlayerButton: React.FC<IExternalPlayerButtonProps> = ({
  scene,
}) => {
  const isAndroid = /(android)/i.test(navigator.userAgent);
  const isAppleDevice = /(ipod|iphone|ipad)/i.test(navigator.userAgent);
  const intl = useIntl();

  const { paths } = scene;

  if (!paths || !paths.stream || (!isAndroid && !isAppleDevice))
    return <span />;

  const { stream } = paths;
  const title = objectTitle(scene);

  let url;
  const streamURL = new URL(stream);
  if (isAndroid) {
    const scheme = streamURL.protocol.slice(0, -1);
    streamURL.hash = `Intent;action=android.intent.action.VIEW;scheme=${scheme};type=video/mp4;S.title=${encodeURIComponent(
      title
    )};end`;

    // #4401 - not allowed to set the protocol from a "special" protocol to a non-special protocol
    url = streamURL
      .toString()
      .replace(new RegExp(`^${streamURL.protocol}`), "intent:");
  } else if (isAppleDevice) {
    streamURL.host = "x-callback-url";
    streamURL.port = "";
    streamURL.pathname = "stream";
    streamURL.search = `url=${encodeURIComponent(stream)}`;

    // #4401 - not allowed to set the protocol from a "special" protocol to a non-special protocol
    url = streamURL
      .toString()
      .replace(new RegExp(`^${streamURL.protocol}`), "vlc-x-callback:");
  }

  return (
    <Button
      className="minimal px-0 px-sm-2 pt-2"
      variant="secondary"
      title={intl.formatMessage({ id: "actions.open_in_external_player" })}
    >
      <a href={url}>
        <Icon icon={faExternalLinkAlt} color="white" />
      </a>
    </Button>
  );
};

export default ExternalPlayerButton;
