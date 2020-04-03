import React, { FunctionComponent, useState } from "react";
import Lightbox from "react-images";
import Gallery from "react-photo-gallery";
import * as GQL from "src/core/generated-graphql";

interface IProps {
  gallery: GQL.GalleryDataFragment;
}

export const GalleryViewer: FunctionComponent<IProps> = ({ gallery }) => {
  const [currentImage, setCurrentImage] = useState<number>(0);
  const [lightboxIsOpen, setLightboxIsOpen] = useState<boolean>(false);

  function openLightbox(
    _event: React.MouseEvent<Element>,
    obj: { index: number }
  ) {
    setCurrentImage(obj.index);
    setLightboxIsOpen(true);
  }
  function closeLightbox() {
    setCurrentImage(0);
    setLightboxIsOpen(false);
  }
  function gotoPrevious() {
    setCurrentImage(currentImage - 1);
  }
  function gotoNext() {
    setCurrentImage(currentImage + 1);
  }

  const photos = gallery.files.map((file) => ({
    src: file.path ?? "",
    caption: file.name ?? "",
  }));
  const thumbs = gallery.files.map((file) => ({
    src: `${file.path}?thumb=true` || "",
    width: 1,
    height: 1,
  }));

  return (
    <div>
      <Gallery photos={thumbs} columns={15} onClick={openLightbox} />
      <Lightbox
        images={photos}
        onClose={closeLightbox}
        onClickPrev={gotoPrevious}
        onClickNext={gotoNext}
        currentImage={currentImage}
        onClickImage={() =>
          window.open(photos[currentImage].src ?? "", "_blank")
        }
        isOpen={lightboxIsOpen}
        width={9999}
      />
    </div>
  );
};
