export function isPlatformUniquelyRenderedByApple() {
  return (
    /(ipad)/i.test(navigator.userAgent) ||
    /(macintosh.*safari)/i.test(navigator.userAgent)
  );
}
