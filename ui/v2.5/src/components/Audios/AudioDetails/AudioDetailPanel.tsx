import React from "react";
import { FormattedMessage, useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import TextUtils from "src/utils/text";
import { TagLink } from "src/components/Shared/TagLink";
import { PerformerCard } from "src/components/Performers/PerformerCard";
import { sortPerformers } from "src/core/performers";
import { CustomFields } from "src/components/Shared/CustomFields";

interface IAudioDetailProps {
  audio: GQL.AudioDataFragment;
}

export const AudioDetailPanel: React.FC<IAudioDetailProps> = (props) => {
  const intl = useIntl();

  function renderDetails() {
    if (!props.audio.details || props.audio.details === "") return;
    return (
      <>
        <h6>
          <FormattedMessage id="details" />:{" "}
        </h6>
        <p className="pre">{props.audio.details}</p>
      </>
    );
  }

  function renderTags() {
    if (props.audio.tags.length === 0) return;
    const tags = props.audio.tags.map((tag) => (
      <TagLink key={tag.id} tag={tag} linkType="audio" />
    ));
    return (
      <>
        <h6>
          <FormattedMessage
            id="countables.tags"
            values={{ count: props.audio.tags.length }}
          />
        </h6>
        {tags}
      </>
    );
  }

  function renderPerformers() {
    if (props.audio.performers.length === 0) return;
    const performers = sortPerformers(props.audio.performers);
    const cards = performers.map((performer) => (
      <PerformerCard
        key={performer.id}
        performer={performer}
        ageFromDate={props.audio.date ?? undefined}
      />
    ));

    return (
      <>
        <h6>
          <FormattedMessage
            id="countables.performers"
            values={{ count: props.audio.performers.length }}
          />
        </h6>
        <div className="row justify-content-center audio-performers">
          {cards}
        </div>
      </>
    );
  }

  const audioDetailsWidth = props.audio.studio ? "col-9" : "col-12";

  return (
    <>
      <div className="row">
        <div className={`${audioDetailsWidth} col-12 audio-details`}>
          <h6>
            <FormattedMessage id="created_at" />:{" "}
            {TextUtils.formatDateTime(intl, props.audio.created_at)}{" "}
          </h6>
          <h6>
            <FormattedMessage id="updated_at" />:{" "}
            {TextUtils.formatDateTime(intl, props.audio.updated_at)}{" "}
          </h6>
          {props.audio.code && (
            <h6>
              <FormattedMessage id="audio_code" />: {props.audio.code}{" "}
            </h6>
          )}
        </div>
      </div>
      <div className="row">
        <div className="col-12">
          {renderDetails()}
          {renderTags()}
          {renderPerformers()}
          <CustomFields values={props.audio.custom_fields} fullWidth />
        </div>
      </div>
    </>
  );
};

export default AudioDetailPanel;
