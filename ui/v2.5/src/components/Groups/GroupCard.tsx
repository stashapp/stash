import React, { useEffect, useState } from "react";
import { Button, ButtonGroup } from "react-bootstrap";
import * as GQL from "src/core/generated-graphql";
import { GridCard, calculateCardWidth } from "../Shared/GridCard/GridCard";
import { HoverPopover } from "../Shared/HoverPopover";
import { Icon } from "../Shared/Icon";
import { SceneLink, TagLink } from "../Shared/TagLink";
import { TruncatedText } from "../Shared/TruncatedText";
import { FormattedMessage } from "react-intl";
import { RatingBanner } from "../Shared/RatingBanner";
import { faPlayCircle, faTag } from "@fortawesome/free-solid-svg-icons";
import ScreenUtils from "src/utils/screen";
import NavUtils from "src/utils/navigation";
import { PopoverCountButton } from "../Shared/PopoverCountButton";
import { Link } from "react-router-dom";

const ContainingGroups: React.FC<{
  group: GQL.GroupDataFragment;
}> = ({ group }) => {
  const containingGroups = group.containing_groups;

  if (containingGroups.length === 1) {
    const g = containingGroups[0].group;
    return (
      <div className="group-containing-groups">
        <FormattedMessage
          id="sub_group_of"
          values={{
            parent: <Link to={NavUtils.makeGroupUrl(g.id)}>{g.name}</Link>,
          }}
        />
      </div>
    );
  }

  if (containingGroups.length > 1) {
    return (
      <div className="group-containing-groups">
        <FormattedMessage
          id="sub_group_of"
          values={{
            parent: (
              <Link to={NavUtils.makeContainingGroupsUrl(group)}>
                {containingGroups.length}&nbsp;
                <FormattedMessage
                  id="countables.groups"
                  values={{ count: containingGroups.length }}
                />
              </Link>
            ),
          }}
        />
      </div>
    );
  }

  return null;
};

interface IProps {
  group: GQL.GroupDataFragment;
  containerWidth?: number;
  description?: number;
  selecting?: boolean;
  selected?: boolean;
  onSelectedChanged?: (selected: boolean, shiftKey: boolean) => void;
}

export const GroupCard: React.FC<IProps> = ({
  group,
  description,
  containerWidth,
  selecting,
  selected,
  onSelectedChanged,
}) => {
  const [cardWidth, setCardWidth] = useState<number>();

  useEffect(() => {
    if (!containerWidth || ScreenUtils.isMobile()) return;

    let preferredCardWidth = 250;
    let fittedCardWidth = calculateCardWidth(
      containerWidth,
      preferredCardWidth!
    );
    setCardWidth(fittedCardWidth);
  }, [containerWidth]);

  function maybeRenderSceneNumber() {
    if (!description) return;

    return (
      <>
        <hr />
        <span className="group-scene-number">
          <FormattedMessage id="scene" /> #{description}
        </span>
      </>
    );
  }

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

  function maybeRenderPopoverButtonGroup() {
    if (
      description ||
      group.scenes.length > 0 ||
      group.tags.length > 0 ||
      group.containing_groups.length > 0 ||
      group.sub_group_count > 0
    ) {
      return (
        <>
          {maybeRenderSceneNumber()}
          <hr />
          <ButtonGroup className="card-popovers">
            {maybeRenderScenesPopoverButton()}
            {maybeRenderTagPopoverButton()}
            {group.sub_group_count > 0 && (
              <PopoverCountButton
                count={group.sub_group_count}
                type="sub_group"
                url={NavUtils.makeSubGroupsUrl(group)}
              />
            )}
          </ButtonGroup>
        </>
      );
    }
  }

  return (
    <GridCard
      className="group-card"
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
          <ContainingGroups group={group} />
        </div>
      }
      selected={selected}
      selecting={selecting}
      onSelectedChanged={onSelectedChanged}
      popovers={maybeRenderPopoverButtonGroup()}
    />
  );
};
