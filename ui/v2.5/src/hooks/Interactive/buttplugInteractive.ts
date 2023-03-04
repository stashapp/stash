import {
  IInteractive
} from "./interactive";
import {
  ButtplugClient,
  ButtplugClientDevice,
  ButtplugBrowserWebsocketClientConnector
} from "buttplug";

export class ButtplugInteractive implements IInteractive {
  _scriptOffset: number;
  _connector: ButtplugBrowserWebsocketClientConnector;
  _client: ButtplugClient;

  constructor(scriptOffset: number = 0) {
    this._scriptOffset = scriptOffset;
    this._connector = new ButtplugBrowserWebsocketClientConnector("ws://localhost:12345");
    this._client = new ButtplugClient("Stash - An organizer for your porn");
    this._client.addListener(
      "deviceadded",
      async (device: ButtplugClientDevice) => {
        console.log(`[buttplug] Device Connected: ${device.name}`);
        //devices.current.push(device);
        // setDeviceDatas((deviceDatas) => [
        //   ...deviceDatas,
        //   {
        //     intensities: { vibration: 0, rotation: 0 },
        //   },
        // ]);
      }
    );
    this._client.addListener("deviceremoved", (device) =>
      console.log(`[buttplug] Device Removed: ${device.name}`)
    );
  }

  enabled(): boolean {
    return true; //this._connector.Connected;
  }

  async connect() {
    await this._client.connect(this._connector);
    await this._client.startScanning();
  }

  set scriptOffset(offset: number) {
    this._scriptOffset = offset;
  }

  async uploadScript(funscriptPath: string) {
    if (!funscriptPath) {
      return;
    }

    const json = await fetch(funscriptPath)
      .then((response) => response.json());

    // TODO
    console.log('[buttplug] Funscript:', json);
    return;
  }

  async sync() {
    return 0;
  }

  setServerTimeOffset(offset: number) {
    console.log(`[buttplug] ServerTimeOffset: ${offset}`);
    return;
  }

  async play(position: number) {
    console.log(`[buttplug] Play position: ${position}`);
    return;
  }

  async pause() {
    console.log('[buttplug] Pause');
    return;
  }

  async ensurePlaying(position: number) {
    console.log(`[buttplug] Ensure play position: ${position}`);
    return;
  }

  async setLooping(looping: boolean) {
    console.log(`[buttplug] Looping: ${looping}`);
    return;
  }
}
