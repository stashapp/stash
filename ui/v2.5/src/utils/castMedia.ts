// Choosing a stream for a cast device, and making its URL reachable.
//
// The scene player cannot forward what it is playing: Stash labels every
// direct stream video/mp4 even when the file is Matroska, which the browser
// tolerates and a Chromecast does not.
//
// Kept free of Video.js, the Cast SDK and the DOM, so the choices are testable.

export interface ICastStream {
  url: string;
  label?: string | null;
}

export interface ICastFile {
  // Container as verified at scan time, not as claimed by the extension.
  format: string;
  video_codec: string;
  audio_codec: string;
  duration: number;
}

export interface ICastSource {
  url: string;
  contentType: string;
}

const MIME_MP4 = "video/mp4";
const MIME_HLS = "application/x-mpegURL";

const ORIGINAL = "ORIGINAL";

// ffprobe codec names, as stored on the file.
const CAST_VIDEO_CODECS = new Set(["h264", "avc1"]);
const CAST_AUDIO_CODECS = new Set(["aac", "mp3", "mp4a"]);
// Container name as resolved by MatchContainer, which folds m4v and mov in.
const CAST_CONTAINER = "mp4";

export function isLocalHost(hostname: string): boolean {
  return (
    hostname === "localhost" ||
    hostname === "127.0.0.1" ||
    hostname === "[::1]" ||
    hostname === "::1"
  );
}

function isPrivateIPv4(ip: string): boolean {
  const parts = ip.split(".").map(Number);
  if (parts.length !== 4 || parts.some((n) => Number.isNaN(n))) return false;
  if (parts[0] === 10) return true;
  if (parts[0] === 192 && parts[1] === 168) return true;
  return parts[0] === 172 && parts[1] >= 16 && parts[1] <= 31;
}

/** A private LAN address is preferred: the device is on the same network. */
export function pickLanIPv4(ips: readonly string[] | null | undefined): string {
  if (!ips?.length) return "";
  return ips.find(isPrivateIPv4) ?? ips[0];
}

interface IParsedStream {
  url: string;
  pathname: string;
  // "ORIGINAL", "FULL_HD", ... Absent on the direct stream.
  resolution: string | null;
}

function parseStreams(
  streams: readonly ICastStream[],
  base: string
): IParsedStream[] {
  const parsed: IParsedStream[] = [];

  for (const stream of streams) {
    try {
      const url = new URL(stream.url, base);
      parsed.push({
        url: stream.url,
        pathname: url.pathname,
        resolution: url.searchParams.get("resolution"),
      });
    } catch {
      // a stream we cannot parse is a stream we cannot cast
    }
  }

  return parsed;
}

/** Whether the device can play the file as-is, whatever the endpoint claims. */
function isDirectFriendly(file?: ICastFile): boolean {
  if (!file) return false;

  return (
    CAST_VIDEO_CODECS.has(file.video_codec.toLowerCase()) &&
    CAST_AUDIO_CODECS.has(file.audio_codec.toLowerCase()) &&
    file.format.toLowerCase() === CAST_CONTAINER
  );
}

interface ICandidate {
  suffix: string;
  contentType: string;
  // Skips the candidate entirely when it returns false.
  applies?: (file?: ICastFile) => boolean;
}

// Tried in order; within each, the original resolution wins over a downscale.
// The resolution lives in the query string, not the label - the
// original-resolution endpoints are labelled plainly ("MP4", "HLS").
const CANDIDATES: ICandidate[] = [
  { suffix: "/stream", contentType: MIME_MP4, applies: isDirectFriendly },
  { suffix: ".mp4", contentType: MIME_MP4 },
  { suffix: ".m3u8", contentType: MIME_HLS },
];

/** The stream to hand to the device, or null when none of them will play. */
export function pickCastSource(
  streams: readonly ICastStream[],
  file?: ICastFile,
  base: string = window.location.href
): ICastSource | null {
  const parsed = parseStreams(streams, base);

  for (const { suffix, contentType, applies } of CANDIDATES) {
    if (applies && !applies(file)) continue;

    const matching = parsed.filter((s) => s.pathname.endsWith(suffix));
    const found =
      matching.find((s) => s.resolution === ORIGINAL) ?? matching[0];

    if (found) {
      return { url: found.url, contentType };
    }
  }

  return null;
}

/** Points a localhost URL somewhere the device can reach. */
export function rewriteCastUrl(
  rawUrl: string,
  lanIp: string,
  base: string = window.location.href
): string {
  const url = new URL(rawUrl, base);

  if (lanIp && isLocalHost(url.hostname)) {
    url.hostname = lanIp;
  }

  return url.toString();
}
