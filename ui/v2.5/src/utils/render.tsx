function renderNonZero(count: number | undefined | null, element: JSX.Element) {
  if (!count) {
    return undefined;
  }

  return element;
}

export default renderNonZero;
