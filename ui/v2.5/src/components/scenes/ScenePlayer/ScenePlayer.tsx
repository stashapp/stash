import React from "react";
import ReactJWPlayer from "react-jw-player";
import { HotKeys } from "react-hotkeys";
import * as GQL from "src/core/generated-graphql";
import { StashService } from "src/core/StashService";
import { SceneHelpers } from "../helpers";
import { ScenePlayerScrubber } from "./ScenePlayerScrubber";

interface IScenePlayerProps {
  scene: GQL.SceneDataFragment;
  timestamp: number;
  autoplay?: boolean;
  onReady?: any;
  onSeeked?: any;
  onTime?: any;
  config?: GQL.ConfigInterfaceDataFragment;
}
interface IScenePlayerState {
  scrubberPosition: number;
}

const KeyMap = {
  NUM0: "0",
  NUM1: "1",
  NUM2: "2",
  SPACE: " "
};

export class ScenePlayerImpl extends React.Component<
  IScenePlayerProps,
  IScenePlayerState
> {
  private player: any;
  private lastTime = 0;

  private KeyHandlers = {
    NUM0: () => {
      this.onReset();
    },
    NUM1: () => {
      this.onDecrease();
    },
    NUM2: () => {
      this.onIncrease();
    },
    SPACE: () => {
      this.onPause();
    }
  };

  constructor(props: IScenePlayerProps) {
    super(props);
    this.onReady = this.onReady.bind(this);
    this.onSeeked = this.onSeeked.bind(this);
    this.onTime = this.onTime.bind(this);

    this.onScrubberSeek = this.onScrubberSeek.bind(this);
    this.onScrubberScrolled = this.onScrubberScrolled.bind(this);

    this.state = { scrubberPosition: 0 };
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
    if (this.player.getState().paused) this.player.play();
    else this.player.pause();
  }

  private onReady() {
    this.player = SceneHelpers.getPlayer();
    if (this.props.timestamp > 0) {
      this.player.seek(this.props.timestamp);
    }
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
    }
  }

  private onScrubberSeek(seconds: number) {
    this.player.seek(seconds);
  }

  private onScrubberScrolled() {
    this.player.pause();
  }

  private shouldRepeat(scene: GQL.SceneDataFragment) {
    const maxLoopDuration = this.props.config
      ? this.props.config.maximumLoopDuration
      : 0;
    return (
      !!scene.file.duration &&
      !!maxLoopDuration &&
      scene.file.duration < maxLoopDuration
    );
  }

  private makeJWPlayerConfig(scene: GQL.SceneDataFragment) {
    if (!scene.paths.stream) {
      return {};
    }

    const repeat = this.shouldRepeat(scene);
    let getDurationHook: (() => GQL.Maybe<number>) | undefined;
    let seekHook:
      | ((seekToPosition: number, _videoTag: any) => void)
      | undefined;
    let getCurrentTimeHook: ((_videoTag: any) => number) | undefined;

    if (!this.props.scene.is_streamable) {
      getDurationHook = () => {
        return this.props.scene.file.duration;
      };

      seekHook = (seekToPosition: number, _videoTag: any) => {
        // eslint-disable-next-line no-param-reassign
        _videoTag.start = seekToPosition;
        // eslint-disable-next-line no-param-reassign
        _videoTag.src = `${this.props.scene.paths.stream}?start=${seekToPosition}`;
        _videoTag.play();
      };

      getCurrentTimeHook = (_videoTag: any) => {
        const start = _videoTag.start || 0;
        return _videoTag.currentTime + start;
      };
    }

    const ret = {
      file: scene.paths.stream,
      image: scene.paths.screenshot,
      tracks: [
        {
          file: scene.paths.vtt,
          kind: "thumbnails"
        },
        {
          file: scene.paths.chapters_vtt,
          kind: "chapters"
        }
      ],
      aspectratio: "16:9",
      width: "100%",
      floating: {
        dismissible: true
      },
      cast: {},
      primary: "html5",
      autostart:
        this.props.autoplay ||
        (this.props.config ? this.props.config.autostartVideo : false),
      repeat,
      playbackRateControls: true,
      playbackRates: [0.75, 1, 1.5, 2, 3, 4],
      getDurationHook,
      seekHook,
      getCurrentTimeHook
    };

    return ret;
  }

  renderPlayer() {
    const config = this.makeJWPlayerConfig(this.props.scene);
    return (
      <ReactJWPlayer
        playerId={SceneHelpers.getJWPlayerId()}
        playerScript="/jwplayer/jwplayer.js"
        customProps={config}
        onReady={this.onReady}
        onSeeked={this.onSeeked}
        onTime={this.onTime}
      />
    );
  }

  public render() {
    return (
      <HotKeys keyMap={KeyMap} handlers={this.KeyHandlers}>
        <div id="jwplayer-container">
          {this.renderPlayer()}
          <ScenePlayerScrubber
            scene={this.props.scene}
            position={this.state.scrubberPosition}
            onSeek={this.onScrubberSeek}
            onScrolled={this.onScrubberScrolled}
          />
        </div>
      </HotKeys>
    );
  }
}

export const ScenePlayer: React.FC<IScenePlayerProps> = (
  props: IScenePlayerProps
) => {
  const config = StashService.useConfiguration();

  return (
    <ScenePlayerImpl
      {...props}
      config={
        config.data && config.data.configuration
          ? config.data.configuration.interface
          : undefined
      }
    />
  );
};
