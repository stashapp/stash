export interface ICastStream {
  url: string;
}

export interface ICastFile {
  // The container detected at scan time, which can differ from the extension.
  format: string;
  video_codec: string;
  audio_codec: string;
}

export interface ICastSource {
  url: string;
  contentType: string;
}

const CAST_VIDEO_CODECS = new Set(["h264", "avc1", "hevc", "h265"]);
const CAST_AUDIO_CODECS = new Set(["aac", "mp3", "mp4a"]);

function isLocalHost(hostname: string): boolean {
  return ["localhost", "127.0.0.1", "[::1]", "::1"].includes(hostname);
}

function isPrivateIPv4(ip: string): boolean {
  const parts = ip.split(".").map(Number);
  if (parts.length !== 4 || parts.some((n) => Number.isNaN(n))) return false;
  if (parts[0] === 10) return true;
  if (parts[0] === 192 && parts[1] === 168) return true;
  return parts[0] === 172 && parts[1] >= 16 && parts[1] <= 31;
}

export function pickLanIPv4(ips: readonly string[] | null | undefined): string {
  if (!ips?.length) return "";
  return ips.find(isPrivateIPv4) ?? ips[0];
}

function canPlayDirect(file?: ICastFile): boolean {
  return (
    !!file &&
    file.format.toLowerCase() === "mp4" &&
    CAST_VIDEO_CODECS.has(file.video_codec.toLowerCase()) &&
    CAST_AUDIO_CODECS.has(file.audio_codec.toLowerCase())
  );
}

function findStream(
  streams: readonly ICastStream[],
  suffix: string,
  base: string
): string | undefined {
  let fallback: string | undefined;

  for (const stream of streams) {
    let url: URL;
    try {
      url = new URL(stream.url, base);
    } catch {
      continue;
    }

    if (!url.pathname.endsWith(suffix)) continue;
    if (url.searchParams.get("resolution") === "ORIGINAL") return stream.url;
    fallback ??= stream.url;
  }

  return fallback;
}

// Every stream endpoint claims video/mp4, so the direct stream is only sent when the file's codecs say the device can play it.
export function pickCastSource(
  streams: readonly ICastStream[],
  file?: ICastFile,
  base: string = window.location.href
): ICastSource | null {
  const direct = canPlayDirect(file)
    ? findStream(streams, "/stream", base)
    : undefined;
  if (direct) return { url: direct, contentType: "video/mp4" };

  const hls = findStream(streams, ".m3u8", base);
  return hls ? { url: hls, contentType: "application/x-mpegURL" } : null;
}

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
