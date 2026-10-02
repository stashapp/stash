import { UAParser } from "ua-parser-js";

export interface ICastMedia {
  url: string;
  contentType: string;
  title: string;
  duration?: number;
  startTime?: number;
  autoplay?: boolean;
}

const CONNECTED = "CONNECTED";

// __onGCastApiAvailable fires once per page, and setOptions re-triggers
// receiver discovery, so the configured context is resolved once and shared.
let contextPromise: Promise<CastContext | null> | null = null;

/**
 * Whether this browser can run the Cast sender SDK. Only Chromium-based ones
 * can; on iOS every browser is WebKit, so none can, whatever the name says.
 */
export function isCastSenderSupported(): boolean {
  const { browser, os } = UAParser();

  if (os.name?.includes("iOS")) return false;

  return /Chrome|Chromium|Edge|Opera/.test(browser.name ?? "");
}

function getCastContext(): CastContext | null {
  return window.cast?.framework.CastContext.getInstance() ?? null;
}

function configure(): CastContext | null {
  const context = getCastContext();
  const chromeCast = window.chrome?.cast;
  if (!context || !chromeCast) return null;

  context.setOptions({
    receiverApplicationId: chromeCast.media.DEFAULT_MEDIA_RECEIVER_APP_ID,
    // Only join sessions this origin started, so Stash does not take over a
    // cast someone else on the network is running.
    autoJoinPolicy: chromeCast.AutoJoinPolicy.ORIGIN_SCOPED,
  });

  return context;
}

/**
 * Resolves a configured CastContext, or null when the SDK is unavailable -
 * unsupported browser, blocked script, or an origin the SDK refuses. Never
 * rejects: a missing context and an unsupported browser both mean "cannot
 * cast".
 */
function loadCastSdk(): Promise<CastContext | null> {
  contextPromise ??= new Promise<CastContext | null>((resolve) => {
    if (window.cast?.framework) {
      resolve(configure());
      return;
    }

    window.__onGCastApiAvailable = (available) => {
      resolve(available ? configure() : null);
    };
  });

  return contextPromise;
}

/**
 * This page's connection to a cast device. The SDK keeps one context per page,
 * so the listeners registered here are removed again on destroy - a disposed
 * scene player must not keep reacting to the device.
 */
export class CastSession {
  private watching: CastContext | null = null;

  public onConnectionChange: () => void = () => undefined;

  public isConnected(): boolean {
    return getCastContext()?.getCastState() === CONNECTED;
  }

  /** Safe to call repeatedly: the SDK is only entered on the first call. */
  public async watch(): Promise<void> {
    if (this.watching) return;

    const context = await loadCastSdk();
    // Another call may have finished while this one awaited.
    if (!context || this.watching) return;
    this.watching = context;

    const events = window.cast?.framework.CastContextEventType;
    if (events) {
      context.addEventListener(events.CAST_STATE_CHANGED, this.syncState);
      context.addEventListener(events.SESSION_STATE_CHANGED, this.syncState);
    }

    this.syncState();
  }

  /** Opens the device picker. Dismissing it arrives here as a rejection. */
  public async connect(): Promise<void> {
    const context = await loadCastSdk();
    if (!context) {
      throw new Error("Cast SDK did not load");
    }

    if (!context.getCurrentSession()) {
      await context.requestSession();
    }

    this.syncState();
  }

  public disconnect(): void {
    getCastContext()?.endCurrentSession(true);
    this.syncState();
  }

  public async load(media: ICastMedia): Promise<void> {
    const chromeCast = window.chrome?.cast;
    const session = getCastContext()?.getCurrentSession();

    if (!chromeCast || !session) {
      throw new Error("No Cast session");
    }

    const info = new chromeCast.media.MediaInfo(media.url, media.contentType);
    info.streamType = chromeCast.media.StreamType.BUFFERED;
    if (media.duration) {
      info.duration = media.duration;
    }

    const metadata = new chromeCast.media.GenericMediaMetadata();
    metadata.metadataType = chromeCast.media.MetadataType.GENERIC;
    metadata.title = media.title;
    info.metadata = metadata;

    const request = new chromeCast.media.LoadRequest(info);
    request.autoplay = media.autoplay ?? true;
    request.currentTime = media.startTime ?? 0;

    await session.loadMedia(request);
  }

  /** Stops listening. The device keeps playing. */
  public destroy(): void {
    const events = window.cast?.framework.CastContextEventType;
    if (this.watching && events) {
      this.watching.removeEventListener(
        events.CAST_STATE_CHANGED,
        this.syncState
      );
      this.watching.removeEventListener(
        events.SESSION_STATE_CHANGED,
        this.syncState
      );
    }

    this.watching = null;
    this.onConnectionChange = () => undefined;
  }

  private syncState = () => this.onConnectionChange();
}
