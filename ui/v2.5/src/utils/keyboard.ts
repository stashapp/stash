export function keyboardClickHandler(onClick: () => void) {
  function onKeyDown(e: React.KeyboardEvent<HTMLAnchorElement>) {
    if (e.key === "Enter" || e.key === " ") {
      onClick();
    }
  }

  return onKeyDown;
}
