import migration32 from "./32.md";

type Module = typeof migration32;

export const migrationNotes: Record<number, Module> = {
  32: migration32,
};
