import migration32 from "./32.md";
import migration39 from "./39.md";

type Module = typeof migration32;

export const migrationNotes: Record<number, Module> = {
  32: migration32,
  39: migration39,
};
