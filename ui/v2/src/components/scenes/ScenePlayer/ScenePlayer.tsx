import { Hotkey, Hotkeys, HotkeysTarget } from "@blueprintjs/core";
import React, { FunctionComponent } from "react";
import ReactJWPlayer from "react-jw-player";
import * as GQL from "../../../core/generated-graphql";
import { SceneHelpers } from "../helpers";
import { ScenePlayerScrubber } from "./ScenePlayerScrubber";
import { StashService } from "../../../core/StashService";
import { IOCounterButtonProps, OCounterButton } from "../OCounterButton";

interface IScenePlayerProps {
  scene: GQL.SceneDataFragment;
  timestamp: number;
  autoplay?: boolean;
  onReady?: any;
  onSeeked?: any;
  onTime?: any;
  config?: GQL.ConfigInterfaceDataFragment;
  oCounter: IOCounterButtonProps;
}
interface IScenePlayerState {
  scrubberPosition: number;
  oCounterMenuOpen: boolean;
}

@HotkeysTarget
export class ScenePlayerImpl extends React.Component<IScenePlayerProps, IScenePlayerState> {
  private player: any;
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
      oCounterMenuOpen: false
    };
  }

  public componentDidUpdate(prevProps: IScenePlayerProps) {
    if (prevProps.timestamp !== this.props.timestamp) {
      this.player.seek(this.props.timestamp);
    }
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
      <>
        <div id="jwplayer-container">
          {this.renderPlayer()}
          <ScenePlayerScrubber
            scene={this.props.scene}
            position={this.state.scrubberPosition}
            onSeek={this.onScrubberSeek}
            onScrolled={this.onScrubberScrolled}
          />
          <div id="o-counter-container" className={this.state.oCounterMenuOpen ? "menu-open" : undefined}>
            <OCounterButton 
              {...this.props.oCounter} 
              onMenuOpened={() => this.setOCounterMenuOpen(true)}
              onMenuClosed={() => this.setOCounterMenuOpen(false)}
            />
          </div>
        </div>
      </>
    );
  }

  public renderHotkeys() {
    const onIncrease = () => {
      const currentPlaybackRate = !!this.player ? this.player.getPlaybackRate() : 1;
      this.player.setPlaybackRate(currentPlaybackRate + 0.5);
    };
    const onDecrease = () => {
      const currentPlaybackRate = !!this.player ? this.player.getPlaybackRate() : 1;
      this.player.setPlaybackRate(currentPlaybackRate - 0.5);
    };
    const onReset = () => { this.player.setPlaybackRate(1); };

    return (
      <Hotkeys>
        <Hotkey
          global={true}
          combo="num2"
          label="Increase playback speed"
          preventDefault={true}
          onKeyDown={onIncrease}
        />
        <Hotkey
          global={true}
          combo="num1"
          label="Decrease playback speed"
          preventDefault={true}
          onKeyDown={onDecrease}
        />
        <Hotkey
          global={true}
          combo="num0"
          label="Reset playback speed"
          preventDefault={true}
          onKeyDown={onReset}
        />
      </Hotkeys>
    );
  }

  private shouldRepeat(scene: GQL.SceneDataFragment) {
    let maxLoopDuration = this.props.config ? this.props.config.maximumLoopDuration : 0;
    return !!scene.file.duration && !!maxLoopDuration && scene.file.duration < maxLoopDuration;
  }

  private makeJWPlayerConfig(scene: GQL.SceneDataFragment) {
    if (!scene.paths.stream) { return {}; }

    let repeat = this.shouldRepeat(scene);
    let getDurationHook: (() => GQL.Maybe<number>) | undefined = undefined;
    let seekHook: ((seekToPosition: number, _videoTag: any) => void) | undefined = undefined;
    let getCurrentTimeHook: ((_videoTag: any) => number) | undefined = undefined;

    if (!this.props.scene.is_streamable) {
      getDurationHook = () => { 
        return this.props.scene.file.duration; 
      };

      seekHook = (seekToPosition: number, _videoTag: any) => {
        _videoTag.start = seekToPosition;
        _videoTag.src = (this.props.scene.paths.stream + "?start=" + seekToPosition);
        _videoTag.play();
      };

      getCurrentTimeHook = (_videoTag: any) => {
        let start = _videoTag.start || 0;
        return _videoTag.currentTime + start;
      }
    }

    let ret = {
      file: scene.paths.stream,
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
      aspectratio: "16:9",
      width: "100%",
      floating: {
        dismissible: true,
      },
      cast: {},
      primary: "html5",
      autostart: this.props.autoplay || (this.props.config ? this.props.config.autostartVideo : false),
      repeat: repeat,
      playbackRateControls: true,
      playbackRates: [0.75, 1, 1.5, 2, 3, 4],
      getDurationHook: getDurationHook,
      seekHook: seekHook,
      getCurrentTimeHook: getCurrentTimeHook
    };

    return ret;
  }

  private onReady() {
    this.player = SceneHelpers.getPlayer();
    if (this.props.timestamp > 0) {
      this.player.seek(this.props.timestamp);
    }
  }

  private setScrubberPosition(position: any) {
    this.setState({scrubberPosition: position, oCounterMenuOpen: this.state.oCounterMenuOpen});
  }

  private setOCounterMenuOpen(value : boolean) {
    console.log("Setting menu open to : " + value);
    this.setState({scrubberPosition: this.state.scrubberPosition, oCounterMenuOpen: value});
  }

  private onSeeked() {
    const position = this.player.getPosition();
    this.setScrubberPosition(position);
    this.player.play();
  }

  private onTime(data: any) {
    const position = this.player.getPosition();
    const difference = Math.abs(position - this.lastTime);
    if (difference > 1) {
      this.lastTime = position;
      this.setScrubberPosition(position);
    }
  }

  private onScrubberSeek(seconds: number) {
    this.player.seek(seconds);
  }

  private onScrubberScrolled() {
    this.player.pause();
  }
}

export const ScenePlayer: FunctionComponent<IScenePlayerProps> = (props: IScenePlayerProps) => {
    const config = StashService.useConfiguration();

    return (
      <>
      <ScenePlayerImpl {...props} config={config.data && config.data.configuration ? config.data.configuration.interface : undefined}/>
      </>
    );
}
