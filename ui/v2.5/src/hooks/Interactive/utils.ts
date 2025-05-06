import { getPlayer } from "../../components/ScenePlayer/util";

export interface IDeviceSettings {
  connectionKey: string;
  scriptOffset: number;
  estimatedServerTimeOffset?: number;
  useStashHostedFunscript?: boolean;
  [key: string]: unknown;
}

export interface IInteractiveClientProvider {
  (
    handyKey: string,
    scriptOffset: number,
    defaultClientProvider?: IInteractiveClientProvider
  ): IInteractiveClient;
}

/**
 * Interface that is used for InteractiveProvider
 */
export interface IInteractiveClient {
  connect(): Promise<void>;
  handyKey: string;
  uploadScript: (funscriptPath: string, apiKey?: string) => Promise<void>;
  sync(): Promise<number>;
  configure(config: Partial<IDeviceSettings>): Promise<void>;
  play(position: number): Promise<void>;
  pause(): Promise<void>;
  ensurePlaying(position: number): Promise<void>;
  setLooping(looping: boolean): Promise<void>;
  readonly connected: boolean;
  readonly playing: boolean;
}

export interface IInteractiveUtils {
  getPlayer: typeof getPlayer;
  interactiveClientProvider: IInteractiveClientProvider;
}
const InteractiveUtils = {
  // hook to allow to customize the interactive client
  interactiveClientProvider: undefined,
  // returns the active player
  getPlayer,
};

export default InteractiveUtils;
