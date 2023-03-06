import { IInteractive } from "./interactive";
import { FunscriptPlayer } from "./funscriptPlayer";
import {
  ButtplugClient,
  ButtplugClientDevice,
  ButtplugBrowserWebsocketClientConnector
} from "buttplug";

export class ButtplugInteractive implements IInteractive {
  _connector: ButtplugBrowserWebsocketClientConnector;
  _client: ButtplugClient;
  _funscriptPlayer: FunscriptPlayer;

  constructor(wsUri: string = "ws://localhost:12345", scriptOffset: number = 0) {
    this._funscriptPlayer = new FunscriptPlayer(async (pos: number) => {
      await this.sendToDevice(pos);
    }, scriptOffset);
    this._connector = new ButtplugBrowserWebsocketClientConnector(wsUri);
    this._client = new ButtplugClient(`Stash ${import.meta.env.VITE_APP_STASH_VERSION}`);
    this._client.addListener("deviceadded", (device: ButtplugClientDevice) => {
      console.log(`[buttplug] Device Connected: ${device.name}`, device);
    });
    this._client.addListener("deviceremoved", (device: ButtplugClientDevice) => {
      console.log(`[buttplug] Device Removed: ${device.name}`);
    });
    console.log('Buttplug construct');
  }

  enabled(): boolean {
    return true; // TODO?: this._connector.Connected;
  }

  async connect() {
    console.log('Buttplug.io connect');
    await this._client.connect(this._connector);
    await this._client.startScanning();
    await this._client.stopScanning();
  }

  async disconnect() {
    await this._client.disconnect();
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
    for (const device of this._client.devices) {
      await device.linear(pos);
      /**
       * Getting following error from Intiface Central v2.3.0 and buttplug-js v3.1.1:
       *
       * [E] Global Loggy: Got invalid messages from remote Buttplug connection
       * - Message: Text("[{\"LinearCmd\":{\"Id\":312,\"DeviceIndex\":0,\"Vectors\":[{\"Index\":0,\"Position\":28}]}}]")
       * - Error: JsonSerializerError("Error during JSON Schema Validation
       * - Message: [{\"LinearCmd\":{\"DeviceIndex\":0,\"Id\":312,\"Vectors\":[{\"Index\":0,\"Position\":28}]}}]
       * - Error: [ValidationError { instance: Array [Object {\"LinearCmd\": Object {\"DeviceIndex\": Number(0), \"Id\": Number(312), \"Vectors\": Array [Object {\"Index\": Number(0), \"Position\": Number(28)}]}}], kind: AnyOf, instance_path: JSONPointer([]), schema_path: JSONPointer([Keyword(\"anyOf\")]) }]")
       *
       * Also, `device.linearAttributes` does not return any info, even with having messageAttributes.LinearCmd
       *    https://github.com/buttplugio/buttplug-js/blob/3.1.1/src/client/ButtplugClientDevice.ts#L212
       *
       * TODO: Submit issue
       */
    }
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
