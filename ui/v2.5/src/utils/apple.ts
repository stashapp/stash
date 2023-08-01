export function isPlatfornUniquelyRenderByApple() {
  return (
    /(ipad)/i.test(navigator.userAgent) ||
    /(macintosh.*safari)/i.test(navigator.userAgent)
  );
}
