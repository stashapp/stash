export function cmToImperial(cm: number) {
  const cmInInches = 0.393700787;
  const inchesInFeet = 12;
  const inches = Math.floor(cm * cmInInches);
  const feet = Math.floor(inches / inchesInFeet);
  return [feet, inches % inchesInFeet];
}

export function kgToLbs(kg: number) {
  return Math.floor(kg * 2.20462262185);
}
