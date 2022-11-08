import migration32 from "./32.md";
import migration38 from "./38.md";

type Module = typeof migration32;

export const migrationNotes: Record<number, Module> = {
  32: migration32,
  38: migration38,
};
