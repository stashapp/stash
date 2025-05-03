import { getPlayerPosition } from "../../components/ScenePlayer/util";
import { IDevice, IDeviceSettings, IScriptData } from "./device";
import { HandyDevice } from "./handy-device";
import Handy from "thehandy";

// Interactive currently uses the Handy API, but could be expanded to use buttplug.io
// via buttplugio/buttplug-rs-ffi's WASM module.
export class Interactive {
  _device: IDevice;
  _handy: Handy;

  constructor(handyKey: string, scriptOffset: number) {
    this._device = new HandyDevice(handyKey, scriptOffset);
    this._handy = (this._device as HandyDevice).handy;
  }

  set device(device: IDevice) {
    this._device = device;
  }

  get device() {
    return this._device;
  }

  async connect() {
    await this._device.connect();
  }

  set handyKey(key: string) {
    this._device
      .updateConfig({
        configKey: key,
      })
      .catch(console.error);
  }

  get handyKey(): string {
    return this._device.getConfig().connectionKey;
  }

  get useStashHostedFunscript(): boolean {
    return this._device.getConfig().useStashHostedFunscript ?? false;
  }

  async uploadScript(funscriptPath: string, apiKey?: string) {
    const config = this._device.getConfig();
    if (!(config.connectionKey && funscriptPath)) {
      return;
    }
    const scriptData: IScriptData = {
      type: "funscript",
    };
    if (config.useStashHostedFunscript) {
      let funscriptUrl = funscriptPath.replace(
        "/funscript",
        "/interactive_csv"
      );
      if (typeof apiKey !== "undefined" && apiKey !== "") {
        var url = new URL(funscriptUrl);
        url.searchParams.append("apikey", apiKey);
        funscriptUrl = url.toString();
      }
      scriptData.url = funscriptUrl;
    } else {
      scriptData.content = await fetch(funscriptPath).then((response) =>
        response.json()
      );
    }

    await this._device.loadScript(scriptData);
  }

  async sync() {
    return this._device.syncTime(getPlayerPosition() ?? 0);
  }

  async settings(config: Partial<IDeviceSettings>) {
    await this._device.updateConfig(config);
  }

  async play(position: number) {
    if (!this._device.isConnected) {
      return;
    }

    await this._device.play(position * 1000);
  }

  async pause() {
    await this._device.stop();
  }

  async ensurePlaying(position: number) {
    if (this._device) {
      return;
    }
    await this.play(position);
  }

  async setLooping(looping: boolean) {
    if (!this._device.isConnected) {
      return;
    }
    await this._device.updateConfig({
      looping,
    });
  }
}
