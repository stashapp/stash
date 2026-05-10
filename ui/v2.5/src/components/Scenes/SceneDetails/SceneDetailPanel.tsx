import React, { useEffect, useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import TextUtils from "src/utils/text";
import { TagLink } from "src/components/Shared/TagLink";
import { PerformerCard } from "src/components/Performers/PerformerCard";
import { sortPerformers } from "src/core/performers";
import { DirectorLink } from "src/components/Shared/Link";
import { CustomFields } from "src/components/Shared/CustomFields";
import { RatingSystem } from "src/components/Shared/Rating/RatingSystem";
import { useToast } from "src/hooks/Toast";

interface ISceneDetailProps {
  scene: GQL.SceneDataFragment;
}

export const SceneDetailPanel: React.FC<ISceneDetailProps> = (props) => {
  const intl = useIntl();
  const Toast = useToast();

  // Local rating state keyed by tag id so the UI updates immediately on click.
  // Seeded from props on mount and whenever the scene's tag_ratings change.
  const [tagRatings, setTagRatings] = useState<Map<string, number | null>>(
    () => new Map(props.scene.tag_ratings.map((r) => [r.tag.id, r.rating100]))
  );
  useEffect(() => {
    setTagRatings(
      new Map(props.scene.tag_ratings.map((r) => [r.tag.id, r.rating100]))
    );
  }, [props.scene.tag_ratings]);

  const [setSceneTagRating] = GQL.useSceneSetTagRatingMutation();

  async function updateTagRating(tagId: string, value: number | null) {
    const previous = tagRatings.get(tagId) ?? null;
    setTagRatings((prev) => {
      const next = new Map(prev);
      if (value === null) next.delete(tagId);
      else next.set(tagId, value);
      return next;
    });
    try {
      await setSceneTagRating({
        variables: {
          scene_id: props.scene.id,
          tag_id: tagId,
          rating100: value,
        },
      });
    } catch (e) {
      // Roll back the optimistic update on failure.
      setTagRatings((prev) => {
        const next = new Map(prev);
        if (previous === null) next.delete(tagId);
        else next.set(tagId, previous);
        return next;
      });
      Toast.error(e);
    }
  }

  function renderDetails() {
    if (!props.scene.details || props.scene.details === "") return;
    return (
      <>
        <h6>
          <FormattedMessage id="details" />:{" "}
        </h6>
        <p className="pre">{props.scene.details}</p>
      </>
    );
  }

  function renderTags() {
    if (props.scene.tags.length === 0) return;
    const tags = props.scene.tags.map((tag) => {
      if (!tag.supports_numeric_rating) {
        return <TagLink key={tag.id} tag={tag} />;
      }
      return (
        <div key={tag.id} className="scene-tag-rating-row">
          <TagLink tag={tag} />
          <RatingSystem
            value={tagRatings.get(tag.id) ?? null}
            onSetRating={(v) => updateTagRating(tag.id, v)}
          />
        </div>
      );
    });
    return (
      <>
        <h6>
          <FormattedMessage
            id="countables.tags"
            values={{ count: props.scene.tags.length }}
          />
        </h6>
        {tags}
      </>
    );
  }

  function renderPerformers() {
    if (props.scene.performers.length === 0) return;
    const performers = sortPerformers(props.scene.performers);
    const cards = performers.map((performer) => (
      <PerformerCard
        key={performer.id}
        performer={performer}
        ageFromDate={props.scene.date ?? undefined}
      />
    ));

    return (
      <>
        <h6>
          <FormattedMessage
            id="countables.performers"
            values={{ count: props.scene.performers.length }}
          />
        </h6>
        <div className="row justify-content-center scene-performers">
          {cards}
        </div>
      </>
    );
  }

  // filename should use entire row if there is no studio
  const sceneDetailsWidth = props.scene.studio ? "col-9" : "col-12";

  return (
    <>
      <div className="row">
        <div className={`${sceneDetailsWidth} col-12 scene-details`}>
          <h6>
            <FormattedMessage id="created_at" />:{" "}
            {TextUtils.formatDateTime(intl, props.scene.created_at)}{" "}
          </h6>
          <h6>
            <FormattedMessage id="updated_at" />:{" "}
            {TextUtils.formatDateTime(intl, props.scene.updated_at)}{" "}
          </h6>
          {props.scene.code && (
            <h6>
              <FormattedMessage id="scene_code" />: {props.scene.code}{" "}
            </h6>
          )}
          {props.scene.director && (
            <h6>
              <FormattedMessage id="director" />:{" "}
              <DirectorLink director={props.scene.director} linkType="scene" />
            </h6>
          )}
        </div>
      </div>
      <div className="row">
        <div className="col-12">
          {renderDetails()}
          {renderTags()}
          {renderPerformers()}
          <CustomFields values={props.scene.custom_fields} fullWidth />
        </div>
      </div>
    </>
  );
};

export default SceneDetailPanel;
