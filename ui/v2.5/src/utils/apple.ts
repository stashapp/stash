import { UAParser } from "ua-parser-js";

export function isPlatformUniquelyRenderedByApple() {
  // OS name on iPads show up as iOS or Mac OS depending on the browser.
  // ua-parser-js v2 reports "macOS", v1 reported "Mac OS" - match both.
  const osName = UAParser().os.name?.toLowerCase();
  const isiOS = osName?.includes("ios");
  const isMacOS = osName?.includes("mac os") || osName?.includes("macos");
  const isSafari = UAParser().browser.name?.includes("Safari");
  return isiOS || (isMacOS && isSafari);
}
