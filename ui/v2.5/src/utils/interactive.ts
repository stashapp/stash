import { Handy, HandyUtils } from "thehandy";
import { HandyMode, HsspSetupResult } from "thehandy/lib/types";

interface IFunscript {
  actions: Array<IAction>;
  inverted: boolean;
  range: number;
}

interface IAction {
  at: number;
  pos: number;
}

// Utility function to convert one range of values to another
function convertRange(
  value: number,
  fromLow: number,
  fromHigh: number,
  toLow: number,
  toHigh: number
) {
  return ((value - fromLow) * (toHigh - toLow)) / (fromHigh - fromLow) + toLow;
}

// Converting to CSV first instead of uploading Funscripts is required
// Reference for Funscript format:
// https://pkg.go.dev/github.com/funjack/launchcontrol/protocol/funscript
function convertFunscriptToCSV(funscript: IFunscript) {
  const lineTerminator = "\r\n";
  if (funscript?.actions?.length > 0) {
    return funscript.actions.reduce((prev: string, curr: IAction) => {
      var { pos } = curr;
      // If it's inverted in the Funscript, we flip it because
      // the Handy doesn't have inverted support
      if (funscript.inverted === true) {
        pos = convertRange(curr.pos, 0, 100, 100, 0);
      }
      // in APIv2; the Handy maintains it's own slide range
      // (ref: https://staging.handyfeeling.com/api/handy/v2/docs/#/SLIDE )
      // so if a range is specified in the Funscript, we convert it to the
      // full range and let the Handy's settings take precedence
      if (funscript.range) {
        pos = convertRange(curr.pos, 0, funscript.range, 0, 100);
      }
      return `${prev}${curr.at},${pos}${lineTerminator}`;
    }, `#Created by stash.app ${new Date().toUTCString()}\n`);
  }
  throw new Error("Not a valid funscript");
}

// Interactive currently uses the Handy API, but could be expanded to use buttplug.io
// via buttplugio/buttplug-rs-ffi's WASM module.
export class Interactive {
  private _connected: boolean;
  private _playing: boolean;
  private _scriptOffset: number;
  private _handy: Handy;

  constructor(handyKey: string, scriptOffset: number) {
    this._handy = new Handy();
    this._handy.connectionKey = handyKey;
    this._scriptOffset = scriptOffset;
    this._connected = false;
    this._playing = false;
  }

  get handyKey(): string {
    return this._handy.connectionKey;
  }

  async uploadScript(funscriptPath: string) {
    if (!(this._handy.connectionKey && funscriptPath)) {
      return;
    }
    // Calibrates the latency between the browser client and the Handy server's
    // This is done before a script upload to ensure a synchronized experience
    await this._handy.getServerTimeOffset();

    const csv = await fetch(funscriptPath)
      .then((response) => response.json())
      .then((json) => convertFunscriptToCSV(json));
    const fileName = `${Math.round(Math.random() * 100000000)}.csv`;
    const csvFile = new File([csv], fileName);
    const tempURL = await HandyUtils.uploadCsv(csvFile).then(
      (response) => response.url
    );

    await this._handy.setMode(HandyMode.hssp);

    this._connected = await this._handy
      .setHsspSetup(tempURL)
      .then((result) => result === HsspSetupResult.downloaded);
  }

  async play(position: number) {
    if (!this._connected) {
      return;
    }

    this._playing = await this._handy
      .setHsspPlay(
        Math.round(position * 1000 + this._scriptOffset),
        this._handy.estimatedServerTimeOffset + Date.now() // our guess of the Handy server's UNIX epoch time
      )
      .then(() => true);
  }

  async pause() {
    if (!this._connected) {
      return;
    }
    this._playing = await this._handy.setHsspStop().then(() => false);
  }

  async ensurePlaying(position: number) {
    if (this._playing) {
      return;
    }
    await this.play(position);
  }
}
