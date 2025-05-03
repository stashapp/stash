import Handy from "thehandy";

import {
  CsvUploadResponse,
  HandyFirmwareStatus,
  HandyMode,
  HsspSetupResult,
} from "thehandy/lib/types";
import {
  IAction,
  IDevice,
  IDeviceSettings,
  IFunscript,
  IScriptData,
} from "./device";

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
      if (funscript.inverted) {
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

// copied from https://github.com/defucilis/thehandy/blob/main/src/HandyUtils.ts
// since HandyUtils is not exported.
// License is listed as MIT. No copyright notice is provided in original.
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
async function uploadCsv(
  csv: File,
  filename?: string
): Promise<CsvUploadResponse> {
  const url = "https://www.handyfeeling.com/api/sync/upload?local=true";
  if (!filename) filename = "script_" + new Date().valueOf() + ".csv";
  const formData = new FormData();
  formData.append("syncFile", csv, filename);
  const response = await fetch(url, {
    method: "post",
    body: formData,
  });
  const newUrl = await response.json();
  return newUrl;
}

export class HandyDevice implements IDevice {
  _connected: boolean;
  _handy: Handy;
  _playing: boolean;
  _settings: IDeviceSettings;

  constructor(handyKey: string, scriptOffset: number) {
    this._handy = new Handy();

    this._settings = {
      offset: scriptOffset,
      connectionKey: handyKey,
      useStashHostedFunscript: false,
    };
    this._connected = false;
    this._playing = false;
    this._handy.connectionKey = handyKey;
  }
  get handy() {
    return this._handy;
  }

  get isConnected() {
    return this._connected;
  }

  get isPlaying() {
    return this._playing;
  }

  disconnect(): Promise<boolean> {
    throw new Error("Method not implemented.");
  }

  getConfig(): IDeviceSettings {
    return { ...this._settings };
  }

  async updateConfig(config: Partial<IDeviceSettings>): Promise<boolean> {
    if (config.looping) {
      await this._handy.setHsspLoop(config.looping as boolean);
    }
    if (config.estimatedServerTimeOffset) {
      this._handy.estimatedServerTimeOffset = config.estimatedServerTimeOffset;
    }
    this._settings = {
      ...this._settings,
      ...config,
    };
    this._handy.connectionKey = this._settings.connectionKey;
    return true;
  }

  async loadScript(scriptData: IScriptData): Promise<boolean> {
    let { url } = scriptData;
    if (scriptData.content) {
      const csv = convertFunscriptToCSV(scriptData.content as IFunscript);
      const fileName = `${Math.round(Math.random() * 100000000)}.csv`;
      const csvFile = new File([csv], fileName);
      url = await uploadCsv(csvFile).then((response) => response.url);
    }
    if (!url) {
      return false;
    }
    await this._handy.setMode(HandyMode.hssp);

    this._connected = await this._handy
      .setHsspSetup(url)
      .then((result) => result === HsspSetupResult.downloaded);
    return this._connected;
  }

  async stop(): Promise<boolean> {
    await this.pause();
    return true;
  }

  syncTime(): Promise<number> {
    return this._handy.getServerTimeOffset();
  }

  /**
   * Connect to the device
   * @param config Optional configuration
   */
  async connect(): Promise<boolean> {
    const connected = await this._handy.getConnected();
    if (!connected) {
      throw new Error("Handy not connected");
    }

    // check the firmware and make sure it's compatible
    const info = await this._handy.getInfo();
    if (info.fwStatus === HandyFirmwareStatus.updateRequired) {
      throw new Error("Handy firmware update required");
    }

    const slideInfo = await this._handy.getSlideSettings();
    this._settings = {
      ...this._settings,
      stroke: slideInfo,
    };
    return this.isConnected;
  }

  async play(timeMs: number): Promise<boolean> {
    this._playing = await this._handy
      .setHsspPlay(
        Math.round(timeMs + this._settings.offset),
        this._handy.estimatedServerTimeOffset + Date.now() // our guess of the Handy server's UNIX epoch time
      )
      .then(() => true);
    return this.isPlaying;
  }

  async pause() {
    if (!this._connected) {
      return;
    }
    this._playing = await this._handy.setHsspStop().then(() => false);
    return this.isPlaying;
  }
}
