export interface IDeviceSettings {
  connectionKey: string;
  offset: number;
  estimatedServerTimeOffset?: number;
  useStashHostedFunscript?: boolean;
  [key: string]: unknown;
}
export interface IFunscript {
  actions: Array<IAction>;
  inverted: boolean;
  range: number;
}

export interface IAction {
  at: number;
  pos: number;
}

export interface IScriptData {
  type: string; // Script type (e.g., "funscript")
  url?: string; // URL to script if remote
  content?: IFunscript; // Script content if loaded directly
}
export interface IDevice {
  /**
   * Connection state
   */
  readonly isConnected: boolean;
  readonly isPlaying: boolean;

  /**
   * Connect to the device
   * @param config Optional configuration
   */
  connect(config?: Record<string, unknown>): Promise<boolean>;

  /**
   * Disconnect from the device
   */
  disconnect(): Promise<boolean>;

  /**
   * Get current device configuration
   */
  getConfig(): IDeviceSettings;

  /**
   * Update device configuration
   * @param config Partial configuration to update
   */
  updateConfig(config: Partial<IDeviceSettings>): Promise<boolean>;

  /**
   * Load a script for playback
   * @param scriptData Script data to load
   */
  loadScript(scriptData: IScriptData): Promise<boolean>;

  /**
   * Play the loaded script at the specified time
   * @param timeMs Current time in milliseconds
   * @param playbackRate Playback rate (1.0 = normal speed)
   * @param loop Whether to loop the script
   */
  play(timeMs: number, playbackRate?: number, loop?: boolean): Promise<boolean>;

  /**
   * Stop playback
   */
  stop(): Promise<boolean>;

  /**
   * Synchronize device time with provided time
   * @param timeMs Current time in milliseconds
   */
  syncTime(timeMs: number): Promise<number>;
}
