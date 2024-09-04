import { PropsWithChildren } from "react";
import { useLightbox } from "./hooks";
import { ILightboxImage } from "./types";
import { Button } from "react-bootstrap";

export const LightboxLink: React.FC<
  PropsWithChildren<{ images?: ILightboxImage[] | undefined; index?: number }>
> = ({ images, index, children }) => {
  const showLightbox = useLightbox({
    images,
  });

  if (!images || images.length === 0) {
    return <>{children}</>;
  }

  return (
    <Button variant="link" onClick={() => showLightbox(index)}>
      {children}
    </Button>
  );
};
