import React from "react";
import { useIntl } from "react-intl";
import { Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import { RatingStars, TruncatedText } from "src/components/Shared";
import { TextUtils } from "src/utils";
import { useGalleryLightbox } from "src/hooks";
import { galleryTitle } from "src/core/galleries";

const CLASSNAME = "GalleryWallCard";
const CLASSNAME_FOOTER = `${CLASSNAME}-footer`;
const CLASSNAME_IMG = `${CLASSNAME}-img`;
const CLASSNAME_TITLE = `${CLASSNAME}-title`;

interface IProps {
  gallery: GQL.SlimGalleryDataFragment;
}

const GalleryWallCard: React.FC<IProps> = ({ gallery }) => {
  const intl = useIntl();
  const showLightbox = useGalleryLightbox(gallery.id);

  const coverFile = gallery?.cover?.files.length
    ? gallery.cover.files[0]
    : undefined;

  const orientation =
    (coverFile?.width ?? 0) > (coverFile?.height ?? 0)
      ? "landscape"
      : "portrait";
  const cover = gallery?.cover?.paths.thumbnail ?? "";
  const title = galleryTitle(gallery);
  const performerNames = gallery.performers.map((p) => p.name);
  const performers =
    performerNames.length >= 2
      ? [...performerNames.slice(0, -2), performerNames.slice(-2).join(" & ")]
      : performerNames;

  return (
    <>
      <section
        className={`${CLASSNAME} ${CLASSNAME}-${orientation}`}
        onClick={showLightbox}
        onKeyPress={showLightbox}
        role="button"
        tabIndex={0}
      >
        <RatingStars rating={gallery.rating} />
        <img src={cover} alt="" className={CLASSNAME_IMG} />
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
              {gallery.date && TextUtils.formatDate(intl, gallery.date)}
            </div>
          </Link>
        </footer>
      </section>
    </>
  );
};

export default GalleryWallCard;
