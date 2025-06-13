import { getPlayer } from "src/components/ScenePlayer/util";
import type { VideoJsPlayer } from "video.js";
import * as GQL from "src/core/generated-graphql";

export interface IDeviceSettings {
  connectionKey: string;
  scriptOffset: number;
  estimatedServerTimeOffset?: number;
  useStashHostedFunscript?: boolean;
  [key: string]: unknown;
}

export interface IInteractiveClientProviderOptions {
  handyKey: string;
  scriptOffset: number;
  defaultClientProvider?: IInteractiveClientProvider;
  stashConfig?: GQL.ConfigDataFragment;
}
export interface IInteractiveClientProvider {
  (options: IInteractiveClientProviderOptions): IInteractiveClient;
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
  getPlayer: () => VideoJsPlayer | undefined;
  interactiveClientProvider: IInteractiveClientProvider | undefined;
}
const InteractiveUtils = {
  // hook to allow to customize the interactive client
  interactiveClientProvider: undefined,
  // returns the active player
  getPlayer,
};

export default InteractiveUtils;
