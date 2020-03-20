import React from "react";
import { StashService } from "src/core/StashService";
import { FormattedMessage, FormattedNumber } from "react-intl";
import { LoadingIndicator } from "src/components/Shared";

export const Stats: React.FC = () => {
  const { data, error, loading } = StashService.useStats();

  if (loading || !data) return <LoadingIndicator />;

  if (error) return <span>error.message</span>;

  return (
    <div className="mt-5">
      <div className="col col-sm-8 m-sm-auto row stats">
        <div className="stats-element">
          <p className="title">
            <FormattedNumber value={data.stats.scene_count} />
          </p>
          <p className="heading">
            <FormattedMessage id="scenes" defaultMessage="Scenes" />
          </p>
        </div>
        <div className="stats-element">
          <p className="title">
            <FormattedNumber value={data.stats.movie_count} />
          </p>
          <p className="heading">
            <FormattedMessage id="movies" defaultMessage="Movies" />
          </p>
        </div>
        <div className="stats-element">
          <p className="title">
            <FormattedNumber value={data.stats.gallery_count} />
          </p>
          <p className="heading">
            <FormattedMessage id="galleries" defaultMessage="Galleries" />
          </p>
        </div>
        <div className="stats-element">
          <p className="title">
            <FormattedNumber value={data.stats.performer_count} />
          </p>
          <p className="heading">
            <FormattedMessage id="performers" defaultMessage="Performers" />
          </p>
        </div>
        <div className="stats-element">
          <p className="title">
            <FormattedNumber value={data.stats.studio_count} />
          </p>
          <p className="heading">
            <FormattedMessage id="studios" defaultMessage="Studios" />
          </p>
        </div>
        <div className="stats-element">
          <p className="title">
            <FormattedNumber value={data.stats.tag_count} />
          </p>
          <p className="heading">
            <FormattedMessage id="tags" defaultMessage="Tags" />
          </p>
        </div>
      </div>
    </div>
  );
};
