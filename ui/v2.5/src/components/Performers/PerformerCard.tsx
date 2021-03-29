import React from "react";
import { Link } from "react-router-dom";
import { FormattedMessage } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import { NavUtils, TextUtils } from "src/utils";
import {
  BasicCard,
  CountryFlag,
  HoverPopover,
  Icon,
  TagLink,
  TruncatedText,
} from "src/components/Shared";
import { Button, ButtonGroup } from "react-bootstrap";

interface IPerformerCardProps {
  performer: GQL.PerformerDataFragment;
  ageFromDate?: string;
  selecting?: boolean;
  selected?: boolean;
  onSelectedChanged?: (selected: boolean, shiftKey: boolean) => void;
}

export const PerformerCard: React.FC<IPerformerCardProps> = ({
  performer,
  ageFromDate,
  selecting,
  selected,
  onSelectedChanged,
}) => {
  const age = TextUtils.age(performer.birthdate, ageFromDate);
  const ageString = `${age} years old${ageFromDate ? " in this scene." : "."}`;

  function maybeRenderFavoriteBanner() {
    if (performer.favorite === false) {
      return;
    }
    return (
      <div className="rating-banner rating-5">
        <FormattedMessage id="favourite" defaultMessage="Favourite" />
      </div>
    );
  }

  function maybeRenderScenesPopoverButton() {
    if (!performer.scene_count) return;

    return (
      <Link to={NavUtils.makePerformerScenesUrl(performer)}>
        <Button className="minimal">
          <Icon icon="play-circle" />
          <span>{performer.scene_count}</span>
        </Button>
      </Link>
    );
  }

  function maybeRenderTagPopoverButton() {
    if (performer.tags.length <= 0) return;

    const popoverContent = performer.tags.map((tag) => (
      <TagLink key={tag.id} tagType="performer" tag={tag} />
    ));

    return (
      <HoverPopover placement="bottom" content={popoverContent}>
        <Button className="minimal">
          <Icon icon="tag" />
          <span>{performer.tags.length}</span>
        </Button>
      </HoverPopover>
    );
  }

  function maybeRenderPopoverButtonGroup() {
    if (performer.scene_count || performer.tags.length > 0) {
      return (
        <>
          <hr />
          <ButtonGroup className="card-popovers">
            {maybeRenderScenesPopoverButton()}
            {maybeRenderTagPopoverButton()}
          </ButtonGroup>
        </>
      );
    }
  }

  return (
    <BasicCard
      className="performer-card"
      url={`/performers/${performer.id}`}
      image={
        <>
          <img
            className="performer-card-image"
            alt={performer.name ?? ""}
            src={performer.image_path ?? ""}
          />
          {maybeRenderFavoriteBanner()}
        </>
      }
      details={
        <>
          <h5>
            <TruncatedText text={performer.name} />
          </h5>
          {age !== 0 ? <div className="text-muted">{ageString}</div> : ""}
          <Link to={NavUtils.makePerformersCountryUrl(performer)}>
            <CountryFlag country={performer.country} />
          </Link>
          {maybeRenderPopoverButtonGroup()}
        </>
      }
      selected={selected}
      selecting={selecting}
      onSelectedChanged={onSelectedChanged}
    />
  );
};
