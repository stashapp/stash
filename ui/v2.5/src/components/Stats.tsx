import React from "react";
import { useStats } from "src/core/StashService";
import { FormattedMessage, FormattedNumber } from "react-intl";
import { LoadingIndicator } from "src/components/Shared";
import Changelog from "src/components/Changelog/Changelog";
import { TextUtils } from "src/utils";

export const Stats: React.FC = () => {
  const { data, error, loading } = useStats();

  if (error) return <span>{error.message}</span>;
  if (loading || !data) return <LoadingIndicator />;

  const scenesSize = TextUtils.fileSize(data.stats.scenes_size);
  const imagesSize = TextUtils.fileSize(data.stats.images_size);

  return (
    <div className="mt-5">
      <div className="col col-sm-8 m-sm-auto row stats">
        <div className="stats-element">
          <p className="title">
            <FormattedNumber value={Math.floor(scenesSize.size)} />
            {` ${TextUtils.formatFileSizeUnit(scenesSize.unit)}`}
          </p>
          <p className="heading">
            <FormattedMessage id="scenes-size" defaultMessage="Scenes size" />
          </p>
        </div>
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
            <FormattedNumber value={Math.floor(imagesSize.size)} />
            {` ${TextUtils.formatFileSizeUnit(imagesSize.unit)}`}
          </p>
          <p className="heading">
            <FormattedMessage id="images-size" defaultMessage="Images size" />
          </p>
        </div>
        <div className="stats-element">
          <p className="title">
            <FormattedNumber value={data.stats.image_count} />
          </p>
          <p className="heading">
            <FormattedMessage id="images" defaultMessage="Images" />
          </p>
        </div>
      </div>
      <div className="col col-sm-8 m-sm-auto row stats">
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
      <div className="changelog col col-sm-8 mx-sm-auto">
        <Changelog />
      </div>
    </div>
  );
};
