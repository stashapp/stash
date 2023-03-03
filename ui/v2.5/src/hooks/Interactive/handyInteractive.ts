import {
  IInteractive,
  convertFunscriptToCSV
} from "./interactive";
import Handy from "thehandy";
import {
  HandyMode,
  HsspSetupResult,
  CsvUploadResponse,
  HandyFirmwareStatus,
} from "thehandy/lib/types";

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

// Interactive currently uses the Handy API, but could be expanded to use buttplug.io
// via buttplugio/buttplug-rs-ffi's WASM module.
export class HandyInteractive implements IInteractive {
  _connected: boolean;
  _playing: boolean;
  _scriptOffset: number;
  _handy: Handy;

  constructor(scriptOffset: number = 0) {
    this._handy = new Handy();
    this._scriptOffset = scriptOffset;
    this._connected = false;
    this._playing = false;
  }

  enabled(): boolean {
    return (this._handy.connectionKey !== "");
  }

  async connect() {
    const connected = await this._handy.getConnected();
    if (!connected) {
      throw new Error("Handy not connected");
    }

    // check the firmware and make sure it's compatible
    const info = await this._handy.getInfo();
    if (info.fwStatus === HandyFirmwareStatus.updateRequired) {
      throw new Error("Handy firmware update required");
    }
  }

  set handyKey(key: string) {
    this._handy.connectionKey = key;
  }

  get handyKey(): string {
    return this._handy.connectionKey;
  }

  set scriptOffset(offset: number) {
    this._scriptOffset = offset;
  }

  async uploadScript(funscriptPath: string) {
    if (!(this._handy.connectionKey && funscriptPath)) {
      return;
    }

    const csv = await fetch(funscriptPath)
      .then((response) => response.json())
      .then((json) => convertFunscriptToCSV(json));
    const fileName = `${Math.round(Math.random() * 100000000)}.csv`;
    const csvFile = new File([csv], fileName);

    const tempURL = await uploadCsv(csvFile).then((response) => response.url);

    await this._handy.setMode(HandyMode.hssp);

    this._connected = await this._handy
      .setHsspSetup(tempURL)
      .then((result) => result === HsspSetupResult.downloaded);
  }

  async sync() {
    return this._handy.getServerTimeOffset();
  }

  setServerTimeOffset(offset: number) {
    this._handy.estimatedServerTimeOffset = offset;
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

  async setLooping(looping: boolean) {
    if (!this._connected) {
      return;
    }
    this._handy.setHsspLoop(looping);
  }
}
