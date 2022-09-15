// This is a modified version of ImageCard which has additional tagging utilities on the bottom.
import React, { MouseEvent, useMemo } from "react";
import { Button, ButtonGroup } from "react-bootstrap";
import cx from "classnames";
import * as GQL from "src/core/generated-graphql";
import { Icon, TagLink, HoverPopover, SweatDrops } from "src/components/Shared";
import { PerformerPopoverButton } from "../Shared/PerformerPopoverButton";
import { GridCard } from "../Shared/GridCard";
import { RatingStars } from "../Scenes/SceneDetails/RatingStars";
import { useImageUpdate } from "../../core/StashService";
import { useToast } from "src/hooks";
import { faBox, faSearch, faTag } from "@fortawesome/free-solid-svg-icons";
import { objectTitle } from "src/core/files";

interface IImageRaterProps {
  image: GQL.SlimImageDataFragment;
  selecting?: boolean;
  selected?: boolean | undefined;
  zoomIndex: number;
  onSelectedChanged?: (selected: boolean, shiftKey: boolean) => void;
  onPreview?: (ev: MouseEvent) => void;
}

export const ImageRater: React.FC<IImageRaterProps> = (
  props: IImageRaterProps
) => {
  const file = useMemo(
    () => (props.image.files.length > 0 ? props.image.files[0] : undefined),
    [props.image]
  );
  const [updateImage] = useImageUpdate();
  const Toast = useToast();
  const [isUpdating, setIsUpdating] = React.useState(false);

  function maybeRenderTagPopoverButton() {
    if (props.image.tags.length <= 0) return;

    const popoverContent = props.image.tags.map((tag) => (
      <TagLink key={tag.id} tag={tag} tagType="image" />
    ));

    return (
      <HoverPopover
        className="tag-count"
        placement="bottom"
        content={popoverContent}
      >
        <Button className="minimal">
          <Icon icon={faTag} />
          <span>{props.image.tags.length}</span>
        </Button>
      </HoverPopover>
    );
  }

  function maybeRenderPerformerPopoverButton() {
    if (props.image.performers.length <= 0) return;

    return <PerformerPopoverButton performers={props.image.performers} />;
  }

  function maybeRenderOCounter() {
    if (props.image.o_counter) {
      return (
        <div className="o-count">
          <Button className="minimal">
            <span className="fa-icon">
              <SweatDrops />
            </span>
            <span>{props.image.o_counter}</span>
          </Button>
        </div>
      );
    }
  }

  function maybeRenderOrganized() {
    if (props.image.organized) {
      return (
        <div className="organized">
          <Button className="minimal">
            <Icon icon={faBox} />
          </Button>
        </div>
      );
    }
  }

  function maybeRenderPopoverButtonGroup() {
    if (
      props.image.tags.length > 0 ||
      props.image.performers.length > 0 ||
      props.image.o_counter ||
      props.image.organized
    ) {
      return (
        <>
          <hr />
          <ButtonGroup className="card-popovers">
            {maybeRenderTagPopoverButton()}
            {maybeRenderPerformerPopoverButton()}
            {maybeRenderOCounter()}
            {maybeRenderOrganized()}
          </ButtonGroup>
        </>
      );
    }
  }

  function isPortrait() {
    const width = file?.width ? file.width : 0;
    const height = file?.height ? file.height : 0;
    return height > width;
  }

  async function onSave(input: GQL.ImageUpdateInput) {
    setIsUpdating(true);
    try {
      await updateImage({
        variables: {
          input,
        },
      });
    } catch (e) {
      Toast.error(e);
    }
    setIsUpdating(false);
  }

  function getImageInput(rating: number | undefined): GQL.ImageUpdateInput {
    return {
      id: props.image.id,
      rating: rating,
    };
  }

  // Removed Rating Banner           <RatingBanner rating={props.image.rating} />
  return (
    <GridCard
      className={`image-card zoom-${props.zoomIndex}`}
      url={`/images/${props.image.id}`}
      title={objectTitle(props.image)}
      linkClassName="image-card-link"
      image={
        <>
          <div className={cx("image-card-preview", { portrait: isPortrait() })}>
            <img
              className="image-card-preview-image"
              alt={props.image.title ?? ""}
              src={props.image.paths.thumbnail ?? ""}
            />
            {props.onPreview ? (
              <div className="preview-button">
                <Button onClick={props.onPreview}>
                  <Icon icon={faSearch} />
                </Button>
              </div>
            ) : undefined}
          </div>
        </>
      }
      popovers={maybeRenderPopoverButtonGroup()}
      details={
        <div className={cx("image-card-rater-stars")}>
          <RatingStars
            value={props.image.rating ?? undefined}
            onSetRating={(value) => {
              onSave(getImageInput(value));
            }}
            disabled={isUpdating}
          />
        </div>
      }
    />
  );
};
