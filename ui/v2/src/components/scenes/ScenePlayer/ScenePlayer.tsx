import { Hotkey, Hotkeys, HotkeysTarget } from "@blueprintjs/core";
import React from "react";
import ReactJWPlayer from "react-jw-player";
import * as GQL from "../../../core/generated-graphql";
import { SceneHelpers } from "../helpers";
import { ScenePlayerScrubber } from "./ScenePlayerScrubber";

interface IScenePlayerProps {
  scene: GQL.SceneDataFragment;
  timestamp: number;
}
interface IScenePlayerState {
  scrubberPosition: number;
}

@HotkeysTarget
export class ScenePlayer extends React.Component<IScenePlayerProps, IScenePlayerState> {
  private player: any;
  private lastTime = 0;

  constructor(props: IScenePlayerProps) {
    super(props);
    this.onReady = this.onReady.bind(this);
    this.onSeeked = this.onSeeked.bind(this);
    this.onTime = this.onTime.bind(this);

    this.onScrubberSeek = this.onScrubberSeek.bind(this);
    this.onScrubberScrolled = this.onScrubberScrolled.bind(this);

    this.state = {scrubberPosition: 0};
  }

  public componentDidUpdate(prevProps: IScenePlayerProps) {
    if (prevProps.timestamp !== this.props.timestamp) {
      this.player.seek(this.props.timestamp);
    }
  }

  public render() {
    const config = this.makeConfig(this.props.scene);
    return (
      <>
        <div id="jwplayer-container">
          <ReactJWPlayer
            playerId={SceneHelpers.getJWPlayerId()}
            playerScript="https://content.jwplatform.com/libraries/QRX6Y71b.js"
            customProps={config}
            onReady={this.onReady}
            onSeeked={this.onSeeked}
            onTime={this.onTime}
          />
          <ScenePlayerScrubber
            scene={this.props.scene}
            position={this.state.scrubberPosition}
            onSeek={this.onScrubberSeek}
            onScrolled={this.onScrubberScrolled}
          />
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

  private makeConfig(scene: GQL.SceneDataFragment) {
    if (!scene.paths.stream) { return {}; }
    return {
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
      primary: "html5",
      autostart: false,
      playbackRateControls: true,
      playbackRates: [0.75, 1, 1.5, 2, 3, 4],
    };
  }

  private onReady() {
    this.player = SceneHelpers.getJWPlayer();
    if (this.props.timestamp > 0) {
      this.player.seek(this.props.timestamp);
    }
  }

  private onSeeked() {
    const position = this.player.getPosition();
    this.setState({scrubberPosition: position});
    this.player.play();
  }

  private onTime(data: any) {
    const position = this.player.getPosition();
    const difference = Math.abs(position - this.lastTime);
    if (difference > 1) {
      this.lastTime = position;
      this.setState({scrubberPosition: position});
    }
  }

  private onScrubberSeek(seconds: number) {
    this.player.seek(seconds);
  }

  private onScrubberScrolled() {
    this.player.pause();
  }
}
