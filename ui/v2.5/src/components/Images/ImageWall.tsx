import React, { MouseEvent, useMemo, useRef } from "react";
import { Button, ButtonGroup } from "react-bootstrap";
import { PerformerPopoverButton } from "../Shared/PerformerPopoverButton";
import {
  faBox,
  faImages,
  faPlus,
  faTag,
} from "@fortawesome/free-solid-svg-icons";
import { Icon, TagLink, HoverPopover, SweatDrops } from "src/components/Shared";
import * as GQL from "src/core/generated-graphql";
import { objectTitle } from "src/core/files";
import { Link } from "react-router-dom";

interface IImageWallProps {
  images: GQL.SlimImageDataFragment[];
  onChangePage: (page: number) => void;
  currentPage: number;
  pageCount: number;
  handleImageOpen: (index: number) => void;
}

export const ImageWall: React.FC<IImageWallProps> = ({
  images,
  handleImageOpen,
}) => {
  const thumbs = images.map((image, index) => {
    function maybeRenderTagPopoverButton() {
      if (image.tags.length <= 0) return;

      const popoverContent = image.tags.map((tag) => (
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
            <span>{image.tags.length}</span>
          </Button>
        </HoverPopover>
      );
    }

    function maybeRenderPerformerPopoverButton() {
      if (image.performers.length <= 0) return;

      return <PerformerPopoverButton performers={image.performers} />;
    }

    function maybeRenderOCounter() {
      if (image.o_counter) {
        return (
          <div className="o-count">
            <Button className="minimal">
              <span className="fa-icon">
                <SweatDrops />
              </span>
              <span>{image.o_counter}</span>
            </Button>
          </div>
        );
      }
    }

    function maybeRenderGallery() {
      if (image.galleries.length <= 0) return;

      const popoverContent = image.galleries.map((gallery) => (
        <TagLink key={gallery.id} gallery={gallery} />
      ));

      return (
        <HoverPopover
          className="gallery-count"
          placement="bottom"
          content={popoverContent}
        >
          <Button className="minimal">
            <Icon icon={faImages} />
            <span>{image.galleries.length}</span>
          </Button>
        </HoverPopover>
      );
    }

    function maybeRenderOrganized() {
      if (image.organized) {
        return (
          <div className="organized">
            <Button className="minimal">
              <Icon icon={faBox} />
            </Button>
          </div>
        );
      }
    }

    function renderEdit() {
      return (
        <div className="edit">
          <Button className="minimal">
            <Icon icon={faPlus} />
          </Button>
        </div>
      );
    }
    const ref = useRef(null);
    const handleTouchEnd = (event) => {
        if(ref.current === document.activeElement){
            alert("If i would to display this, it would be nice!")
        }
        else{
            ref.current.focus();
            console.log(ref.current)
        }
        event.stopPropagation()
        event.preventDefault();
      }; 
    function renderPopoverButtonGroup() {
      if (
        image.tags.length > 0 ||
        image.performers.length > 0 ||
        image.o_counter ||
        image.galleries.length > 0 ||
        image.organized
      ) {
        return (
          <>
            <Link
              to={`/images/${image.id}`}
              onClick={(e) => e.stopPropagation()}
              onTouchEnd={handleTouchEnd}
            >
              <ButtonGroup className="wall-popovers" ref={ref}>
                {maybeRenderTagPopoverButton()}
                {maybeRenderPerformerPopoverButton()}
                {maybeRenderOCounter()}
                {maybeRenderGallery()}
                {maybeRenderOrganized()}
              </ButtonGroup>
            </Link>
          </>
        );
      } else {
        return (
          <>
            <Link
              to={`/images/${image.id}`}
              onClick={(e) => e.stopPropagation()}
              onTouchEnd={handleTouchEnd}
            >
              <ButtonGroup className="wall-popovers">
                {renderEdit()}
              </ButtonGroup>
            </Link>
          </>
        );
      }
    }
    return (
      <>
        <div
          role="link"
          tabIndex={index}
          key={image.id}
          onClick={() => handleImageOpen(index)}
          onKeyPress={() => handleImageOpen(index)}
          className="wall-image"
        >
          <img
            src={image.paths.thumbnail ?? ""}
            loading="lazy"
            className="gallery-image"
            alt={objectTitle(image)}
          />
          {renderPopoverButtonGroup()}
        </div>
      </>
    );
  });

  return (
    <div className="gallery">
      <div className="flexbin">{thumbs}</div>
    </div>
  );
};
