import migration32 from "./32.md";
import migration39 from "./39.md";
import migration48 from "./48.md";
import migration58 from "./58.md";
import migration60 from "./60.md";

export const migrationNotes: Record<number, string> = {
  32: migration32,
  39: migration39,
  48: migration48,
  58: migration58,
  60: migration60,
};
