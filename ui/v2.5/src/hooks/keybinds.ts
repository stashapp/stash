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
  const firstDigit = useRef<string | undefined>(undefined);
  const sequenceTimeout = useRef<ReturnType<typeof setTimeout>>();

  useEffect(() => {
    if (!isVisible) return;

    function endSequence() {
      mode.current = "idle";
      firstDigit.current = undefined;
      if (sequenceTimeout.current) {
        clearTimeout(sequenceTimeout.current);
        sequenceTimeout.current = undefined;
      }
    }

    // "r", the digits and "`" are bound for isVisible's lifetime and gate
    // their behaviour on `mode`, instead of being bound and unbound per
    // sequence.
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
      firstDigit.current = undefined;

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
          if (firstDigit.current !== undefined) {
            let combined = parseInt(firstDigit.current + i.toString(), 10);
            if (combined === 0) {
              combined = 100;
            }

            setRatingRef.current(combined);
            endSequence();
          } else {
            firstDigit.current = i.toString();
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
