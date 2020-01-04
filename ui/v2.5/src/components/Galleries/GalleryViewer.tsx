import React, { FunctionComponent, useState } from "react";
import Lightbox from "react-images";
import Gallery from "react-photo-gallery";
import * as GQL from "../../core/generated-graphql";

interface IProps {
  gallery: GQL.GalleryDataFragment;
}

export const GalleryViewer: FunctionComponent<IProps> = (props: IProps) => {
  const [currentImage, setCurrentImage] = useState<number>(0);
  const [lightboxIsOpen, setLightboxIsOpen] = useState<boolean>(false);

  function openLightbox(event: any, obj: any) {
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

  const photos = props.gallery.files.map((file) => ({src: file.path || "", caption: file.name}));
  const thumbs = props.gallery.files.map((file) => ({src: `${file.path}?thumb=true` || "", width: 1, height: 1}));
  return (
    <div>
      <Gallery photos={thumbs} columns={15} onClick={openLightbox} />
      <Lightbox
        images={photos}
        onClose={closeLightbox}
        onClickPrev={gotoPrevious}
        onClickNext={gotoNext}
        currentImage={currentImage}
        isOpen={lightboxIsOpen}
        onClickImage={() => window.open(photos[currentImage].src, "_blank")}
        width={9999}
      />
    </div>
  );
};
