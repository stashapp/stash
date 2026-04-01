import { DeviceManager, HandyDevice, loadScript } from "ive-connect";
import type { IDeviceSettings } from "./utils";

// Interactive currently uses the Handy API, but could be expanded to use buttplug.io
// via buttplugio/buttplug-rs-ffi's WASM module.
export class Interactive {
  _playing: boolean;
  _scriptOffset: number;
  _manager = new DeviceManager();
  _handyDevice: HandyDevice;
  _connectionKey: string;

  constructor(handyKey: string, scriptOffset: number) {
    this._scriptOffset = scriptOffset;
    this._connectionKey = handyKey;
    this._handyDevice = new HandyDevice({
      connectionKey: this._connectionKey,
    });
    this._handyDevice.updateConfig({
      offset: this._scriptOffset,
    });
    this._playing = false;
  }

  async updateConfig() {
    this._handyDevice.updateConfig({
      connectionKey: this._connectionKey,
      offset: this._scriptOffset,
    });
  }

  get connected() {
    return this._handyDevice.isConnected;
  }
  get playing() {
    return this._handyDevice.isPlaying;
  }

  async connect() {
    const connected = await this._handyDevice.connect({
      offset: this._scriptOffset,
    });
    if (!connected) {
      throw new Error("Handy not connected");
    }
  }

  set handyKey(key: string) {
    this._connectionKey = key;
    this.updateConfig();
  }

  set scriptOffset(offset: number) {
    this._scriptOffset = offset;
    this.updateConfig();
  }

  async uploadScript(funscriptPath: string, apiKey?: string) {
    // append apikey if necessary
    var funscriptURL = funscriptPath;
    if (typeof apiKey !== "undefined" && apiKey !== "") {
      const url = new URL(funscriptPath);
      url.searchParams.append("apikey", apiKey);
      funscriptPath = url.toString();
    }

    const result = await loadScript({
      type: "funscript",
      url: funscriptURL,
    });
    if (result.error) {
      throw new Error(result.error);
    }
  }

  sync() {
    // only function that handles offset is updateConfig
    return this._handyDevice.api.getServerTimeOffset()
  }

  async configure(config: Partial<IDeviceSettings>) {
    this._scriptOffset = config.scriptOffset ?? this._scriptOffset;
    this.handyKey = config.connectionKey ?? this.handyKey;
  }

  async play(position: number, loop: boolean = false) {
    if (!this.connected) {
      return;
    }

    this._playing = await this._handyDevice.play(
      Math.round(position * 1000 + this._scriptOffset),
      1.0, // playback rate
      loop
    );
  }

  async pause() {
    if (!this.connected) {
      return;
    }
    // returns boolean about success, not playing state
    this._playing = await this._handyDevice.stop().then((res) => !!res);
  }

  async ensurePlaying(position: number) {
    if (this._playing) {
      return;
    }
    await this.play(position);
  }

  async setLooping(looping: boolean) {
    if (!this.connected) {
      return;
    }
    this._handyDevice.hspSetLoop(looping);
  }
}
