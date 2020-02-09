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
      <div className="col col-sm-8 m-sm-auto">
        <nav className="col col-sm-8 m-sm-auto row">
          <div className="flex-grow-1">
            <div>
              <p className="heading">
                <FormattedMessage id="scenes" defaultMessage="Scenes" />
              </p>
              <p className="title">
                <FormattedNumber value={data.stats.scene_count} />
              </p>
            </div>
          </div>
          <div className="flex-grow-1">
            <div>
              <p className="heading">
                <FormattedMessage id="galleries" defaultMessage="Galleries" />
              </p>
              <p className="title">
                <FormattedNumber value={data.stats.gallery_count} />
              </p>
            </div>
          </div>
          <div className="flex-grow-1">
            <div>
              <p className="heading">
                <FormattedMessage id="performers" defaultMessage="Performers" />
              </p>
              <p className="title">
                <FormattedNumber value={data.stats.performer_count} />
              </p>
            </div>
          </div>
          <div className="flex-grow-1">
            <div>
              <p className="heading">
                <FormattedMessage id="studios" defaultMessage="Studios" />
              </p>
              <p className="title">
                <FormattedNumber value={data.stats.studio_count} />
              </p>
            </div>
          </div>
          <div className="flex-grow-1">
            <div>
              <p className="heading">
                <FormattedMessage id="tags" defaultMessage="Tags" />
              </p>
              <p className="title">
                <FormattedNumber value={data.stats.tag_count} />
              </p>
            </div>
          </div>
        </nav>

        <h5>
          <FormattedMessage id="stats.notes" defaultMessage="Notes" />
        </h5>
        <em>
          <FormattedMessage
            id="stats.warning"
            defaultMessage="This is still an early version, some things are still a work in progress."
          />
        </em>
      </div>
    </div>
  );
};
