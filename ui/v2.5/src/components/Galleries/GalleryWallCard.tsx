import React, { useState } from "react";
import { Form } from "react-bootstrap";
import { useIntl } from "react-intl";
import { Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import { TruncatedText } from "src/components/Shared/TruncatedText";
import TextUtils from "src/utils/text";
import { useGalleryLightbox } from "src/hooks/Lightbox/hooks";
import { galleryTitle } from "src/core/galleries";
import { RatingSystem } from "../Shared/Rating/RatingSystem";
import { GalleryPreviewScrubber } from "./GalleryPreviewScrubber";
import { useDragMoveSelect } from "../Shared/GridCard/dragMoveSelect";
import cx from "classnames";

const CLASSNAME = "GalleryWallCard";
const CLASSNAME_FOOTER = `${CLASSNAME}-footer`;
const CLASSNAME_IMG = `${CLASSNAME}-img`;
const CLASSNAME_TITLE = `${CLASSNAME}-title`;
const CLASSNAME_IMG_CONTAIN = `${CLASSNAME}-img-contain`;

interface IProps {
  gallery: GQL.SlimGalleryDataFragment;
  selected?: boolean;
  onSelectedChanged?: (selected: boolean, shiftKey: boolean) => void;
  selecting?: boolean;
}

type Orientation = "landscape" | "portrait";

function getOrientation(width: number, height: number): Orientation {
  return width > height ? "landscape" : "portrait";
}

const GalleryWallCard: React.FC<IProps> = ({
  gallery,
  selected,
  onSelectedChanged,
  selecting,
}) => {
  const intl = useIntl();
  const [coverOrientation, setCoverOrientation] =
    React.useState<Orientation>("landscape");
  const [imageOrientation, setImageOrientation] =
    React.useState<Orientation>("landscape");
  const showLightbox = useGalleryLightbox(gallery.id, gallery.chapters);

  const { dragProps } = useDragMoveSelect({
    selecting: selecting || false,
    selected: selected || false,
    onSelectedChanged: onSelectedChanged,
  });

  const cover = gallery?.paths.cover;

  function onCoverLoad(e: React.SyntheticEvent<HTMLImageElement, Event>) {
    const target = e.target as HTMLImageElement;
    setCoverOrientation(
      getOrientation(target.naturalWidth, target.naturalHeight)
    );
  }

  function onNonCoverLoad(e: React.SyntheticEvent<HTMLImageElement, Event>) {
    const target = e.target as HTMLImageElement;
    setImageOrientation(
      getOrientation(target.naturalWidth, target.naturalHeight)
    );
  }

  const [imgSrc, setImgSrc] = useState<string | undefined>(cover ?? undefined);
  const title = galleryTitle(gallery);
  const performerNames = gallery.performers.map((p) => p.name);
  const performers =
    performerNames.length >= 2
      ? [...performerNames.slice(0, -2), performerNames.slice(-2).join(" & ")]
      : performerNames;

  function handleCardClick(event: React.MouseEvent) {
    if (selecting && onSelectedChanged) {
      onSelectedChanged(!selected, event.shiftKey);
      return;
    }
    showLightboxStart();
  }

  async function showLightboxStart() {
    if (gallery.image_count === 0) {
      return;
    }

    showLightbox(0);
  }

  const imgClassname =
    imageOrientation !== coverOrientation ? CLASSNAME_IMG_CONTAIN : "";

  let shiftKey = false;

  return (
    <>
      <section
        className={`${CLASSNAME} ${CLASSNAME}-${coverOrientation} wall-item`}
        onClick={handleCardClick}
        onKeyPress={() => showLightboxStart()}
        role="button"
        tabIndex={0}
        {...dragProps}
      >
        {onSelectedChanged && (
          <Form.Control
            type="checkbox"
            className="wall-item-check mousetrap"
            checked={selected}
            onChange={() => onSelectedChanged(!selected, shiftKey)}
            onClick={(
              event: React.MouseEvent<HTMLInputElement, MouseEvent>
            ) => {
              shiftKey = event.shiftKey;
              event.stopPropagation();
            }}
          />
        )}
        <RatingSystem value={gallery.rating100} disabled withoutContext />
        <img
          loading="lazy"
          src={imgSrc}
          alt=""
          className={cx(CLASSNAME_IMG, imgClassname)}
          // set orientation based on cover only
          onLoad={imgSrc === cover ? onCoverLoad : onNonCoverLoad}
        />
        <div className="lineargradient">
          <footer className={CLASSNAME_FOOTER}>
            <Link
              to={`/galleries/${gallery.id}`}
              onClick={(e) => e.stopPropagation()}
            >
              {title && (
                <TruncatedText
                  text={title}
                  lineCount={1}
                  className={CLASSNAME_TITLE}
                />
              )}
              <TruncatedText text={performers.join(", ")} />
              <div>
                {gallery.date && TextUtils.formatFuzzyDate(intl, gallery.date)}
              </div>
            </Link>
          </footer>
          <GalleryPreviewScrubber
            previewPath={gallery.paths.preview}
            defaultPath={cover ?? ""}
            imageCount={gallery.image_count}
            onClick={(i) => {
              showLightbox(i);
            }}
            onPathChanged={setImgSrc}
          />
        </div>
      </section>
    </>
  );
};

export default GalleryWallCard;
