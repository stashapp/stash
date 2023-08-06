import UAParser from "ua-parser-js";

export function isPlatformUniquelyRenderedByApple() {
  // OS name on iPads show up as iOS or Max OS depending on the browser.
  const isiOS = UAParser().os.name?.includes("iOS");
  const isMacOS = UAParser().os.name?.includes("Mac OS");
  const isSafari = UAParser().browser.name?.includes("Safari");
  return isiOS || (isMacOS && isSafari);
}
