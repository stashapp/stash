/* eslint-disable @typescript-eslint/no-explicit-any */
import React from "react";
import ReactJWPlayer from "react-jw-player";
import * as GQL from "src/core/generated-graphql";
import { JWUtils, ScreenUtils } from "src/utils";
import { ConfigurationContext } from "src/hooks/Config";
import { ScenePlayerScrubber } from "./ScenePlayerScrubber";
import { Interactive } from "../../utils/interactive";

/*
fast-forward svg derived from https://github.com/jwplayer/jwplayer/blob/master/src/assets/SVG/rewind-10.svg
Flipped horizontally, then flipped '10' numerals horizontally.

Creative Common License: https://github.com/jwplayer/jwplayer/blob/master/LICENSE
*/
const ffSVG = `
<svg xmlns="http://www.w3.org/2000/svg" class="jw-svg-icon jw-svg-icon-rewind" viewBox="0 0 240 240" focusable="false">
  <path d="M185,135.6c-3.7-6.3-10.4-10.3-17.7-10.6c-7.3,0.3-14,4.3-17.7,10.6c-8.6,14.2-8.6,32.1,0,46.3c3.7,6.3,10.4,10.3,17.7,10.6
  c7.3-0.3,14-4.3,17.7-10.6C193.6,167.6,193.6,149.8,185,135.6z M167.3,182.8c-7.8,0-14.4-11-14.4-24.1s6.6-24.1,14.4-24.1
  s14.4,11,14.4,24.1S175.2,182.8,167.3,182.8z M123.9,192.5v-51l-4.8,4.8l-6.8-6.8l13-13c1.9-1.9,4.9-1.9,6.8,0
  c0.9,0.9,1.4,2.1,1.4,3.4v62.7L123.9,192.5z M22.7,57.4h130.1V38.1c0-5.3,3.6-7.2,8-4.3l41.8,27.9c1.2,0.6,2.1,1.5,2.7,2.7
  c1.4,3,0.2,6.5-2.7,8l-41.8,27.9c-4.4,2.9-8,1-8-4.3V76.7H37.1v96.4h48.2v19.3H22.6c-2.6,0-4.8-2.2-4.8-4.8V62.3
  C17.8,59.6,20,57.4,22.7,57.4z">
  </path>
</svg>
`;

interface IScenePlayerProps {
  className?: string;
  scene: GQL.SceneDataFragment;
  sceneStreams: GQL.SceneStreamEndpoint[];
  timestamp: number;
  autoplay?: boolean;
  onReady?: () => void;
  onSeeked?: () => void;
  onTime?: () => void;
  onComplete?: () => void;
  config?: GQL.ConfigInterfaceDataFragment;
}
interface IScenePlayerState {
  scrubberPosition: number;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  config: Record<string, any>;
  interactiveClient: Interactive;
}
export class ScenePlayerImpl extends React.Component<
  IScenePlayerProps,
  IScenePlayerState
> {
  private static isDirectStream(src?: string) {
    if (!src) {
      return false;
    }

    const url = new URL(src);
    return url.pathname.endsWith("/stream");
  }

  // Typings for jwplayer are, unfortunately, very lacking
  private player: any;
  private playlist: any;
  private lastTime = 0;

  constructor(props: IScenePlayerProps) {
    super(props);
    this.onReady = this.onReady.bind(this);
    this.onSeeked = this.onSeeked.bind(this);
    this.onTime = this.onTime.bind(this);

    this.onScrubberSeek = this.onScrubberSeek.bind(this);
    this.onScrubberScrolled = this.onScrubberScrolled.bind(this);
    this.state = {
      scrubberPosition: 0,
      config: this.makeJWPlayerConfig(props.scene),
      interactiveClient: new Interactive(
        this.props.config?.handyKey || "",
        this.props.config?.funscriptOffset || 0
      ),
    };

    // Default back to Direct Streaming
    localStorage.removeItem("jwplayer.qualityLabel");
  }
  public UNSAFE_componentWillReceiveProps(props: IScenePlayerProps) {
    if (props.scene !== this.props.scene) {
      this.setState((state) => ({
        ...state,
        config: this.makeJWPlayerConfig(this.props.scene),
      }));
    }
  }

  public componentDidUpdate(prevProps: IScenePlayerProps) {
    if (prevProps.timestamp !== this.props.timestamp) {
      this.player.seek(this.props.timestamp);
    }
  }

  onIncrease() {
    const currentPlaybackRate = this.player ? this.player.getPlaybackRate() : 1;
    this.player.setPlaybackRate(currentPlaybackRate + 0.5);
  }
  onDecrease() {
    const currentPlaybackRate = this.player ? this.player.getPlaybackRate() : 1;
    this.player.setPlaybackRate(currentPlaybackRate - 0.5);
  }

  onReset() {
    this.player.setPlaybackRate(1);
  }
  onPause() {
    if (this.player.getState().paused) {
      this.player.play();
    } else {
      this.player.pause();
    }
  }

  private addForwardButton() {
    // add forward button: https://github.com/jwplayer/jwplayer/issues/3894
    const playerContainer = document.querySelector(
      `#${JWUtils.playerID}`
    ) as HTMLElement;

    // display icon
    const rewindContainer = playerContainer.querySelector(
      ".jw-display-icon-rewind"
    ) as HTMLElement;
    const forwardContainer = rewindContainer.cloneNode(true) as HTMLElement;
    const forwardDisplayButton = forwardContainer.querySelector(
      ".jw-icon-rewind"
    ) as HTMLElement;
    forwardDisplayButton.innerHTML = ffSVG;
    forwardDisplayButton.ariaLabel = "Forward 10 Seconds";
    const nextContainer = playerContainer.querySelector(
      ".jw-display-icon-next"
    ) as HTMLElement;
    (nextContainer.parentNode as HTMLElement).insertBefore(
      forwardContainer,
      nextContainer
    );

    // control bar icon
    const buttonContainer = playerContainer.querySelector(
      ".jw-button-container"
    ) as HTMLElement;
    const rewindControlBarButton = buttonContainer.querySelector(
      ".jw-icon-rewind"
    ) as HTMLElement;
    const forwardControlBarButton = rewindControlBarButton.cloneNode(
      true
    ) as HTMLElement;
    forwardControlBarButton.innerHTML = ffSVG;
    forwardControlBarButton.ariaLabel = "Forward 10 Seconds";
    (rewindControlBarButton.parentNode as HTMLElement).insertBefore(
      forwardControlBarButton,
      rewindControlBarButton.nextElementSibling
    );

    // add onclick handlers
    [forwardDisplayButton, forwardControlBarButton].forEach((button) => {
      button.onclick = () => {
        this.player.seek(this.player.getPosition() + 10);
      };
    });
  }

  private onReady() {
    this.player = JWUtils.getPlayer();
    this.addForwardButton();

    this.player.on("error", (err: any) => {
      if (err && err.code === 224003) {
        // When jwplayer has been requested to play but the browser doesn't support the video format.
        this.handleError(true);
      }
    });

    //
    this.player.on("meta", (metadata: any) => {
      if (
        metadata.metadataType === "media" &&
        !metadata.width &&
        !metadata.height
      ) {
        // Occurs during preload when videos with supported audio/unsupported video are preloaded.
        // Treat this as a decoding error and try the next source without playing.
        // However on Safari we get an media event when m3u8 is loaded which needs to be ignored.
        const currentFile = this.player.getPlaylistItem().file;
        if (currentFile != null && !currentFile.includes("m3u8")) {
          const state = this.player.getState();
          const play = state === "buffering" || state === "playing";
          this.handleError(play);
        }
      }
    });

    this.player.on("firstFrame", () => {
      if (this.props.timestamp > 0) {
        this.player.seek(this.props.timestamp);
      }
    });

    this.player.on("play", () => {
      if (this.props.scene.interactive) {
        this.state.interactiveClient.play(this.player.getPosition());
      }
    });

    this.player.on("pause", () => {
      if (this.props.scene.interactive) {
        this.state.interactiveClient.pause();
      }
    });

    if (this.props.scene.interactive) {
      this.state.interactiveClient.uploadScript(
        this.props.scene.paths.funscript || ""
      );
    }

    this.player.getContainer().focus();
  }

  private onSeeked() {
    const position = this.player.getPosition();
    this.setState({ scrubberPosition: position });
    this.player.play();
  }

  private onTime() {
    const position = this.player.getPosition();
    const difference = Math.abs(position - this.lastTime);
    if (difference > 1) {
      this.lastTime = position;
      this.setState({ scrubberPosition: position });
      if (this.props.scene.interactive) {
        this.state.interactiveClient.ensurePlaying(position);
      }
    }
  }

  private onComplete() {
    if (this.props?.onComplete) {
      this.props.onComplete();
    }
  }

  private onScrubberSeek(seconds: number) {
    this.player.seek(seconds);
  }

  private onScrubberScrolled() {
    this.player.pause();
  }

  private handleError(play: boolean) {
    const currentFile = this.player.getPlaylistItem();
    if (currentFile) {
      // eslint-disable-next-line no-console
      console.log(`Source failed: ${currentFile.file}`);
    }

    if (this.tryNextStream()) {
      // eslint-disable-next-line no-console
      console.log(
        `Trying next source in playlist: ${this.playlist.sources[0].file}`
      );
      this.player.load(this.playlist);
      if (play) {
        this.player.play();
      }
    }
  }

  private shouldRepeat(scene: GQL.SceneDataFragment) {
    const maxLoopDuration = this.props?.config?.maximumLoopDuration ?? 0;
    return (
      !!scene.file.duration &&
      !!maxLoopDuration &&
      scene.file.duration < maxLoopDuration
    );
  }

  private tryNextStream() {
    if (this.playlist.sources.length > 1) {
      this.playlist.sources.shift();
      return true;
    }

    return false;
  }

  private makePlaylist() {
    const { scene } = this.props;

    return {
      image: scene.paths.screenshot,
      tracks: [
        {
          file: scene.paths.vtt,
          kind: "thumbnails",
        },
        {
          file: scene.paths.chapters_vtt,
          kind: "chapters",
        },
      ],
      sources: this.props.sceneStreams.map((s) => {
        return {
          file: s.url,
          type: s.mime_type,
          label: s.label,
        };
      }),
    };
  }

  private makeJWPlayerConfig(scene: GQL.SceneDataFragment) {
    if (!scene.paths.stream) {
      return {};
    }

    const repeat = this.shouldRepeat(scene);
    const getDurationHook = () => {
      return this.props.scene.file.duration ?? null;
    };

    const seekHook = (seekToPosition: number, _videoTag: HTMLVideoElement) => {
      if (!_videoTag.src || _videoTag.src.endsWith(".m3u8")) {
        return false;
      }

      if (ScenePlayerImpl.isDirectStream(_videoTag.src)) {
        if (_videoTag.dataset.start) {
          /* eslint-disable-next-line no-param-reassign */
          _videoTag.dataset.start = "0";
        }

        // direct stream - fall back to default
        return false;
      }

      // remove the start parameter
      const srcUrl = new URL(_videoTag.src);
      srcUrl.searchParams.delete("start");

      /* eslint-disable no-param-reassign */
      _videoTag.dataset.start = seekToPosition.toString();
      srcUrl.searchParams.append("start", seekToPosition.toString());
      _videoTag.src = srcUrl.toString();
      /* eslint-enable no-param-reassign */

      _videoTag.play();

      // return true to indicate not to fall through to default
      return true;
    };

    const getCurrentTimeHook = (_videoTag: HTMLVideoElement) => {
      const start = Number.parseFloat(_videoTag.dataset?.start ?? "0");
      return _videoTag.currentTime + start;
    };

    this.playlist = this.makePlaylist();

    // TODO: leverage the floating.mode option after upgrading JWPlayer
    const extras: any = {};

    if (!ScreenUtils.isMobile()) {
      extras.floating = {
        dismissible: true,
      };
    }

    const ret = {
      playlist: this.playlist,
      image: scene.paths.screenshot,
      width: "100%",
      height: "100%",
      cast: {},
      primary: "html5",
      preload: "none",
      autostart:
        this.props.autoplay ||
        (this.props.config ? this.props.config.autostartVideo : false) ||
        this.props.timestamp > 0,
      repeat,
      playbackRateControls: true,
      playbackRates: [0.75, 1, 1.5, 2, 3, 4],
      getDurationHook,
      seekHook,
      getCurrentTimeHook,
      ...extras,
    };

    return ret;
  }

  public render() {
    let className =
      this.props.className ?? "w-100 col-sm-9 m-sm-auto no-gutter";
    const sceneFile = this.props.scene.file;

    if (
      sceneFile.height &&
      sceneFile.width &&
      sceneFile.height > sceneFile.width
    ) {
      className += " portrait";
    }

    return (
      <div id="jwplayer-container" className={className}>
        <ReactJWPlayer
          playerId={JWUtils.playerID}
          playerScript="jwplayer/jwplayer.js"
          customProps={this.state.config}
          onReady={this.onReady}
          onSeeked={this.onSeeked}
          onTime={this.onTime}
          onOneHundredPercent={() => this.onComplete()}
          className="video-wrapper"
        />
        <ScenePlayerScrubber
          scene={this.props.scene}
          position={this.state.scrubberPosition}
          onSeek={this.onScrubberSeek}
          onScrolled={this.onScrubberScrolled}
        />
      </div>
    );
  }
}

export const ScenePlayer: React.FC<IScenePlayerProps> = (
  props: IScenePlayerProps
) => {
  const { configuration } = React.useContext(ConfigurationContext);

  return (
    <ScenePlayerImpl
      {...props}
      config={configuration ? configuration.interface : undefined}
    />
  );
};
