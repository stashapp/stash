import Mousetrap from "mousetrap";
import { useEffect, useRef } from "react";
import { RatingSystemType } from "src/utils/rating";

const starRatingShortcuts: { [char: string]: number } = {
  "0": NaN,
  "1": 20,
  "2": 40,
  "3": 60,
  "4": 80,
  "5": 100,
};

type RatingSequenceMode = "idle" | "star" | "decimal";

export function useRatingKeybinds(
  isVisible: boolean,
  ratingSystem: RatingSystemType | undefined,
  setRating: (v: number) => void,
  mousetrap: Pick<Mousetrap.MousetrapInstance, "bind" | "unbind"> = Mousetrap
) {
  // setRating/ratingSystem are recreated every render by every caller (they
  // close over the currently displayed entity). Reading them through refs,
  // updated unconditionally on each render, lets the bind effect below key
  // only off isVisible/mousetrap while still always acting on the latest
  // values -- rebinding "r" itself isn't needed to pick up a fresh setRating.
  const setRatingRef = useRef(setRating);
  setRatingRef.current = setRating;
  const ratingSystemRef = useRef(ratingSystem);
  ratingSystemRef.current = ratingSystem;

  const mode = useRef<RatingSequenceMode>("idle");
  const firstChar = useRef<string | undefined>(undefined);
  const sequenceTimeout = useRef<ReturnType<typeof setTimeout>>();

  useEffect(() => {
    if (!isVisible) return;

    function endSequence() {
      mode.current = "idle";
      firstChar.current = undefined;
      if (sequenceTimeout.current) {
        clearTimeout(sequenceTimeout.current);
        sequenceTimeout.current = undefined;
      }
    }

    // "r", the digits and "`" are bound unconditionally for isVisible's
    // lifetime, and gate their behaviour on `mode` instead of being bound
    // and unbound per sequence. Callers pass a setRating closure they don't
    // memoize, so this effect only depends on isVisible/mousetrap -- if it
    // depended on setRating too, an unrelated re-render could tear down and
    // rebind this effect mid-sequence, unbinding the digit keys before the
    // 1s window elapses.
    mousetrap.bind("r", () => {
      // numeric keypresses get caught by jwplayer, so blur the element
      // if the rating sequence is started
      if (document.activeElement instanceof HTMLElement) {
        document.activeElement.blur();
      }

      mode.current =
        !ratingSystemRef.current ||
        ratingSystemRef.current === RatingSystemType.Stars
          ? "star"
          : "decimal";
      firstChar.current = undefined;

      if (sequenceTimeout.current) clearTimeout(sequenceTimeout.current);
      sequenceTimeout.current = setTimeout(endSequence, 1000);
    });

    mousetrap.bind("`", () => {
      if (mode.current !== "decimal") return;
      setRatingRef.current(NaN);
      endSequence();
    });

    for (let i = 0; i <= 9; ++i) {
      mousetrap.bind(i.toString(), () => {
        if (mode.current === "star") {
          const value = starRatingShortcuts[i.toString()];
          if (value === undefined) return;
          setRatingRef.current(value);
          endSequence();
        } else if (mode.current === "decimal") {
          if (firstChar.current !== undefined) {
            let combined = parseInt(firstChar.current + i.toString(), 10);
            if (combined === 0) {
              combined = 100;
            }

            setRatingRef.current(combined);
            endSequence();
          } else {
            firstChar.current = i.toString();
          }
        }
      });
    }

    return () => {
      mousetrap.unbind("r");
      mousetrap.unbind("`");
      for (let i = 0; i <= 9; ++i) {
        mousetrap.unbind(i.toString());
      }
      endSequence();
    };
  }, [isVisible, mousetrap]);
}
