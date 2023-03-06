import {
  IInteractive,
  FunscriptPlayer
} from "./interactive";
import {
  ButtplugClient,
  ButtplugClientDevice,
  ButtplugBrowserWebsocketClientConnector
} from "buttplug";

export class ButtplugInteractive implements IInteractive {
  _connector: ButtplugBrowserWebsocketClientConnector;
  _client: ButtplugClient;
  _funscriptPlayer: FunscriptPlayer;

  constructor(scriptOffset: number = 0) {
    this._funscriptPlayer = new FunscriptPlayer(async (pos: number) => {
      await this.sendToDevice(pos);
    }, scriptOffset);
    this._connector = new ButtplugBrowserWebsocketClientConnector("ws://localhost:12345");
    this._client = new ButtplugClient(`Stash ${import.meta.env.VITE_APP_STASH_VERSION}`);
    this._client.addListener(
      "deviceadded",
      async (device: ButtplugClientDevice) => {
        console.log(`[buttplug] Device Connected: ${device.name}`);
        // TODO
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
    console.log('Buttplug construct');
  }

  enabled(): boolean {
    return true; // TODO?: this._connector.Connected;
  }

  async connect() {
    console.log('Buttplug.io connect');
    await this._client.connect(this._connector);
    await this._client.startScanning();
  }

  async disconnect() {
    // TODO
    return;
  }

  set scriptOffset(offset: number) {
    this._funscriptPlayer.offset = offset;
  }

  async uploadScript(funscriptPath: string) {
    if (!funscriptPath) {
      this._funscriptPlayer.funscript = undefined;
      return;
    }

    const json = await fetch(funscriptPath)
      .then((response) => response.json());

    console.log('[buttplug] Funscript:', json);
    this._funscriptPlayer.funscript = json;
  }

  async sendToDevice(pos: number) {
    console.log(`[buttplug] Action pos: ${pos}`);
    // TODO
  }

  async sync() {
    console.log(`[buttplug] Sync`);
    // TODO: Setting to 1 to ensure connect is called.
    //       Need to revisit this when focusing on the settings panel and context.
    return 1;
  }

  setServerTimeOffset(offset: number) {
    console.log(`[buttplug] ServerTimeOffset: ${offset}`);
    // TODO: I don't think anything is needed here (noop)
  }

  async play(position: number) {
    console.log(`[buttplug] Play position: ${position}`);
    this._funscriptPlayer.play(Math.trunc(position * 1000));
    return;
  }

  async pause() {
    console.log('[buttplug] Pause');
    this._funscriptPlayer.pause();
  }

  async ensurePlaying(position: number) {
    console.log(`[buttplug] Ensure play position: ${position}`);
    this._funscriptPlayer.playSync(Math.trunc(position * 1000));
  }

  async setLooping(looping: boolean) {
    console.log(`[buttplug] Looping: ${looping}`);
    // TODO: I don't think anything is needed here (noop)
  }
}
