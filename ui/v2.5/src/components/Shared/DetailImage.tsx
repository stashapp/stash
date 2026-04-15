import { useLayoutEffect, useRef } from "react";
import { PatchComponent } from "src/patch";
import { remToPx } from "src/utils/units";

const DEFAULT_WIDTH = Math.round(remToPx(30));

// Props used by the <img> element
type IDetailImageProps = JSX.IntrinsicElements["img"];

export const DetailImage = PatchComponent(
  "DetailImage",
  (props: IDetailImageProps) => {
    const imgRef = useRef<HTMLImageElement>(null);

    function fixWidth() {
      const img = imgRef.current;
      if (!img) return;

      // prevent SVG's w/o intrinsic size from rendering as 0x0
      if (img.naturalWidth === 0) {
        // If the naturalWidth is zero, it means the image either hasn't loaded yet
        // or we're on Firefox and it is an SVG w/o an intrinsic size.
        // So set the width to our fallback width.
        img.setAttribute("width", String(DEFAULT_WIDTH));
      } else {
        // If we have a `naturalWidth`, this could either be the actual intrinsic width
        // of the image, or the image is an SVG w/o an intrinsic size and we're on Chrome or Safari,
        // which seem to return a size calculated in some browser-specific way.
        // Worse yet, once rendered, Safari will then return the value of `img.width` as `img.naturalWidth`,
        // so we need to clone the image to disconnect it from the DOM, and then get the `naturalWidth` of the clone,
        // in order to always return the same `naturalWidth` for a given src.
        const i = img.cloneNode() as HTMLImageElement;
        img.setAttribute("width", String(i.naturalWidth || DEFAULT_WIDTH));
      }
    }

    useLayoutEffect(() => {
      fixWidth();
    }, [props.src]);

    return <img ref={imgRef} onLoad={() => fixWidth()} {...props} />;
  }
);
