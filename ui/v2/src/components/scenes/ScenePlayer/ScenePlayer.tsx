import { Hotkey, Hotkeys, HotkeysTarget } from "@blueprintjs/core";
import React, { Component, FunctionComponent } from "react";
import ReactJWPlayer from "react-jw-player";
import * as GQL from "../../../core/generated-graphql";
import { SceneHelpers } from "../helpers";
import { ScenePlayerScrubber } from "./ScenePlayerScrubber";
import videojs from "video.js";
import "video.js/dist/video-js.css";
import { StashService } from "../../../core/StashService";

interface IScenePlayerProps {
  scene: GQL.SceneDataFragment;
  timestamp: number;
  onReady?: any;
  onSeeked?: any;
  onTime?: any;
  config?: GQL.ConfigInterfaceDataFragment;
}
interface IScenePlayerState {
  scrubberPosition: number;
}

interface IVideoJSPlayerProps extends IScenePlayerProps {
  videoJSOptions: videojs.PlayerOptions
}

export class VideoJSPlayer extends React.Component<IVideoJSPlayerProps> {
  private player: any;
  private videoNode: any;

  constructor(props: IVideoJSPlayerProps) {
    super(props);
  }

  componentDidMount() {
    this.player = videojs(this.videoNode, this.props.videoJSOptions);

    SceneHelpers.registerJSPlayer(this.player);

    this.player.src(this.props.scene.paths.stream);

    // hack duration
    this.player.duration = () => { return this.props.scene.file.duration; };
    this.player.start = 0;
    this.player.oldCurrentTime = this.player.currentTime;
    this.player.currentTime = (time: any) => { 
      if( time == undefined )
      {
        return this.player.oldCurrentTime() + this.player.start;
      }
      this.player.start = time;
      this.player.oldCurrentTime(0);
      this.player.src(this.props.scene.paths.stream + "?start=" + time);
      this.player.play();

      return this;
    };

    // dirty hack - make this player look like JWPlayer
    this.player.seek = this.player.currentTime;
    this.player.getPosition = this.player.currentTime;

    this.player.ready(() => {
      this.player.on("timeupdate", () => {
        this.props.onTime();
      });

      this.player.on("seeked", () => {
        this.props.onSeeked();
      });

      this.props.onReady();
    });
  }

  componentWillUnmount() {
    if (this.player) {
      this.player.dispose();
      SceneHelpers.deregisterJSPlayer();
    }
  }

  render() {
    return (
      <div>
        <div className="vjs-player" data-vjs-player>
          <video 
              ref={ node => this.videoNode = node } 
              className="video-js vjs-default-skin vjs-big-play-centered" 
              poster={this.props.scene.paths.screenshot}
              controls 
              preload="auto">
          </video>
        </div>
      </div>
    );
  }
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

    this.state = {scrubberPosition: 0};
  }

  public componentDidUpdate(prevProps: IScenePlayerProps) {
    if (prevProps.timestamp !== this.props.timestamp) {
      this.player.seek(this.props.timestamp);
    }
  }

  renderPlayer() {
    if (this.props.scene.is_streamable) {
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
    } else {
      // don't render videoJS until config is loaded
      if (this.props.config) {
        const config = this.makeVideoJSConfig(this.props.scene);
        return (
          <VideoJSPlayer 
              videoJSOptions={config}
              scene={this.props.scene}
              timestamp={this.props.timestamp}
              onReady={this.onReady}
              onSeeked={this.onSeeked}
              onTime={this.onTime}>
          </VideoJSPlayer>
        )
      }
    }
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
      aspectratio: "16:9",
      width: "100%",
      floating: {
        dismissible: true,
      },
      cast: {},
      primary: "html5",
      autostart: this.props.config ? this.props.config.autostartVideo : false,
      repeat: repeat,
      playbackRateControls: true,
      playbackRates: [0.75, 1, 1.5, 2, 3, 4],
    };
  }

  private makeVideoJSConfig(scene: GQL.SceneDataFragment) {
    if (!scene.paths.stream) { return {}; }

    let repeat = this.shouldRepeat(scene);

    return {
      autoplay: this.props.config ? this.props.config.autostartVideo : false,
      loop: repeat,
    };
  } 

  private onReady() {
    this.player = SceneHelpers.getPlayer();
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

export const ScenePlayer: FunctionComponent<IScenePlayerProps> = (props: IScenePlayerProps) => {
    const config = StashService.useConfiguration();

    return <ScenePlayerImpl {...props} config={config.data && config.data.configuration ? config.data.configuration.interface : undefined}/>
}
