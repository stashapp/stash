import Mousetrap from "mousetrap";
import { useEffect, useRef } from "react";
import { RatingSystemType } from "src/utils/rating";

export function useRatingKeybinds(
  isVisible: boolean,
  ratingSystem: RatingSystemType | undefined,
  setRating: (v: number) => void,
  mousetrap: Pick<Mousetrap.MousetrapInstance, "bind" | "unbind"> = Mousetrap
) {
  const firstChar = useRef<string | undefined>(undefined);
  const ratingTimeout = useRef<ReturnType<typeof setTimeout>>();

  // (Re)start the 1s window after each "r"-initiated sequence, cancelling any
  // pending unbind. Without this, pressing "r" again before the window elapses
  // leaves the earlier timeout scheduled, which then unbinds the digit keys
  // mid-sequence and drops the next keypress (e.g. a quick "r 3" then "r 4").
  function restartRatingTimeout(unbind: () => void) {
    if (ratingTimeout.current) clearTimeout(ratingTimeout.current);
    ratingTimeout.current = setTimeout(() => {
      ratingTimeout.current = undefined;
      unbind();
    }, 1000);
  }

  const starRatingShortcuts: { [char: string]: number } = {
    "0": NaN,
    "1": 20,
    "2": 40,
    "3": 60,
    "4": 80,
    "5": 100,
  };

  function handleStarRatingKeybinds() {
    for (const key in starRatingShortcuts) {
      mousetrap.bind(key, () => setRating(starRatingShortcuts[key]));
    }

    restartRatingTimeout(() => {
      for (const key in starRatingShortcuts) {
        mousetrap.unbind(key);
      }
    });
  }

  function handleDecimalKeybinds() {
    // start each sequence fresh so a new "r" doesn't combine with a digit left
    // buffered from a previous, abandoned sequence
    firstChar.current = undefined;

    mousetrap.bind("`", () => {
      setRating(NaN);
    });

    for (let i = 0; i <= 9; ++i) {
      mousetrap.bind(i.toString(), () => {
        if (firstChar.current !== undefined) {
          let combined = parseInt(firstChar.current + i.toString(), 10);
          if (combined === 0) {
            combined = 100;
          }

          setRating(combined);
          firstChar.current = undefined;
        } else {
          firstChar.current = i.toString();
        }
      });
    }

    restartRatingTimeout(() => {
      firstChar.current = undefined;

      mousetrap.unbind("`");
      for (let i = 0; i <= 9; ++i) {
        mousetrap.unbind(i.toString());
      }
    });
  }

  useEffect(() => {
    if (!isVisible) return;

    mousetrap.bind("r", () => {
      // numeric keypresses get caught by jwplayer, so blur the element
      // if the rating sequence is started
      if (document.activeElement instanceof HTMLElement) {
        document.activeElement.blur();
      }

      if (!ratingSystem || ratingSystem === RatingSystemType.Stars) {
        return handleStarRatingKeybinds();
      } else {
        return handleDecimalKeybinds();
      }
    });

    return () => {
      mousetrap.unbind("r");
    };
  });
}
