import React, { useMemo } from "react";
import { Button, ButtonGroup } from "react-bootstrap";
import * as GQL from "src/core/generated-graphql";
import { PatchComponent } from "src/patch";
import { GridCard } from "../Shared/GridCard/GridCard";
import { HoverPopover } from "../Shared/HoverPopover";
import { Icon } from "../Shared/Icon";
import { SceneLink, TagLink } from "../Shared/TagLink";
import { TruncatedText } from "../Shared/TruncatedText";
import { FormattedMessage } from "react-intl";
import { RatingBanner } from "../Shared/RatingBanner";
import { faPlayCircle, faTag } from "@fortawesome/free-solid-svg-icons";
import { RelatedGroupPopoverButton } from "./RelatedGroupPopover";
import { OCounterButton } from "../Shared/CountButton";

const Description: React.FC<{
  sceneNumber?: number;
  description?: string;
}> = ({ sceneNumber, description }) => {
  if (!sceneNumber && !description) return null;

  return (
    <>
      <hr />
      {sceneNumber !== undefined && (
        <span className="group-scene-number">
          <FormattedMessage id="scene" /> #{sceneNumber}
        </span>
      )}
      {description !== undefined && (
        <span className="group-containing-group-description">
          {description}
        </span>
      )}
    </>
  );
};

interface IProps {
  group: GQL.ListGroupDataFragment;
  cardWidth?: number;
  sceneNumber?: number;
  selecting?: boolean;
  selected?: boolean;
  zoomIndex?: number;
  onSelectedChanged?: (selected: boolean, shiftKey: boolean) => void;
  fromGroupId?: string;
  onMove?: (srcIds: string[], targetId: string, after: boolean) => void;
}

export const GroupCard: React.FC<IProps> = PatchComponent(
  "GroupCard",
  ({
    group,
    sceneNumber,
    cardWidth,
    selecting,
    selected,
    zoomIndex,
    onSelectedChanged,
    fromGroupId,
    onMove,
  }) => {
    const groupDescription = useMemo(() => {
      if (!fromGroupId) {
        return undefined;
      }

      const containingGroup = group.containing_groups.find(
        (cg) => cg.group.id === fromGroupId
      );

      return containingGroup?.description ?? undefined;
    }, [fromGroupId, group.containing_groups]);

    function maybeRenderScenesPopoverButton() {
      if (group.scenes.length === 0) return;

      const popoverContent = group.scenes.map((scene) => (
        <SceneLink key={scene.id} scene={scene} />
      ));

      return (
        <HoverPopover
          className="scene-count"
          placement="bottom"
          content={popoverContent}
        >
          <Button className="minimal">
            <Icon icon={faPlayCircle} />
            <span>{group.scenes.length}</span>
          </Button>
        </HoverPopover>
      );
    }

    function maybeRenderTagPopoverButton() {
      if (group.tags.length <= 0) return;

      const popoverContent = group.tags.map((tag) => (
        <TagLink key={tag.id} linkType="group" tag={tag} />
      ));

      return (
        <HoverPopover placement="bottom" content={popoverContent}>
          <Button className="minimal tag-count">
            <Icon icon={faTag} />
            <span>{group.tags.length}</span>
          </Button>
        </HoverPopover>
      );
    }

    function maybeRenderOCounter() {
      if (!group.o_counter) return;

      return <OCounterButton value={group.o_counter} />;
    }

    function maybeRenderPopoverButtonGroup() {
      if (
        sceneNumber ||
        groupDescription ||
        group.scenes.length > 0 ||
        group.tags.length > 0 ||
        group.containing_groups.length > 0 ||
        group.sub_group_count > 0
      ) {
        return (
          <>
            <Description
              sceneNumber={sceneNumber}
              description={groupDescription}
            />
            <hr />
            <ButtonGroup className="card-popovers">
              {maybeRenderScenesPopoverButton()}
              {maybeRenderTagPopoverButton()}
              {(group.sub_group_count > 0 ||
                group.containing_groups.length > 0) && (
                <RelatedGroupPopoverButton group={group} />
              )}
              {maybeRenderOCounter()}
            </ButtonGroup>
          </>
        );
      }
    }

    return (
      <GridCard
        className={`group-card zoom-${zoomIndex}`}
        objectId={group.id}
        onMove={onMove}
        url={`/groups/${group.id}`}
        width={cardWidth}
        title={group.name}
        linkClassName="group-card-header"
        image={
          <>
            <img
              loading="lazy"
              className="group-card-image"
              alt={group.name ?? ""}
              src={group.front_image_path ?? ""}
            />
            <RatingBanner rating={group.rating100} />
          </>
        }
        details={
          <div className="group-card__details">
            <span className="group-card__date">{group.date}</span>
            <TruncatedText
              className="group-card__description"
              text={group.synopsis}
              lineCount={3}
            />
          </div>
        }
        selected={selected}
        selecting={selecting}
        onSelectedChanged={onSelectedChanged}
        popovers={maybeRenderPopoverButtonGroup()}
      />
    );
  }
);
