const hexToBinary = (hex: string) =>
  hex
    .split("")
    .map((i) => parseInt(i, 16).toString(2).padStart(4, "0"))
    .join("");

export const distance = (a: string, b: string): number | undefined => {
  if (a.length !== b.length) return;

  const aBinary = hexToBinary(a);
  const bBinary = hexToBinary(b);

  let counter = 0;
  for (let i = 0; i < aBinary.length; i++) {
    if (aBinary[i] !== bBinary[i]) counter++;
  }

  return counter;
};
