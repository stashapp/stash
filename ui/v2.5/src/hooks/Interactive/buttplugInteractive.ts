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

  constructor(scriptOffset: number = 0) {
    this._scriptOffset = scriptOffset;
  }

  enabled(): boolean {
    return true;
  }

  async connect() {
    const connector = new ButtplugBrowserWebsocketClientConnector("ws://localhost:12345");
    const client = new ButtplugClient("Device Control Example");
    client.addListener(
      "deviceadded",
      async (device: ButtplugClientDevice) => {
        console.log(`Device Connected: ${device.name}`);
        //devices.current.push(device);
        // setDeviceDatas((deviceDatas) => [
        //   ...deviceDatas,
        //   {
        //     intensities: { vibration: 0, rotation: 0 },
        //   },
        // ]);
      }
    );
    client.addListener("deviceremoved", (device) =>
      console.log(`Device Removed: ${device.name}`)
    );
    await client.connect(connector);
    await client.startScanning();
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
    console.log(json);
    return;
  }

  async sync() {
    return 0;
  }

  setServerTimeOffset(offset: number) {
    return;
  }

  async play(position: number) {
    return;
  }

  async pause() {
    return;
  }

  async ensurePlaying(position: number) {
    return;
  }

  async setLooping(looping: boolean) {
    return;
  }
}
