import { PropsWithChildren } from "react";
import { useLightbox } from "./hooks";
import { ILightboxImage } from "./types";
import { Button } from "react-bootstrap";
import { PatchComponent } from "src/patch";

export const LightboxLink: React.FC<
  PropsWithChildren<{ images?: ILightboxImage[] | undefined; index?: number }>
> = PatchComponent("LightboxLink", ({ images, index, children }) => {
  const showLightbox = useLightbox();

  if (!images || images.length === 0) {
    return <>{children}</>;
  }

  return (
    <Button
      variant="link"
      onClick={() => showLightbox({ images, initialIndex: index })}
    >
      {children}
    </Button>
  );
});
