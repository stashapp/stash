import React from "react";
import { StashService } from "src/core/StashService";
import { LoadingIndicator } from "src/components/Shared";

export const Stats: React.FC = () => {
  const { data, error, loading } = StashService.useStats();

  if (loading || !data) return <LoadingIndicator message="Loading..." />;

  if (error) return <span>error.message</span>;

  return (
    <div className="mt-5">
      <div className="col col-sm-8 m-sm-auto">
        <nav className="col col-sm-8 m-sm-auto row">
          <div className="flex-grow-1">
            <div>
              <p className="heading">Scenes</p>
              <p className="title">{data.stats.scene_count}</p>
            </div>
          </div>
          <div className="flex-grow-1">
            <div>
              <p className="heading">Galleries</p>
              <p className="title">{data.stats.gallery_count}</p>
            </div>
          </div>
          <div className="flex-grow-1">
            <div>
              <p className="heading">Performers</p>
              <p className="title">{data.stats.performer_count}</p>
            </div>
          </div>
          <div className="flex-grow-1">
            <div>
              <p className="heading">Studios</p>
              <p className="title">{data.stats.studio_count}</p>
            </div>
          </div>
          <div className="flex-grow-1">
            <div>
              <p className="heading">Tags</p>
              <p className="title">{data.stats.tag_count}</p>
            </div>
          </div>
        </nav>

        <h5>Notes</h5>
        <em>
          This is still an early version, some things are still a work in
          progress.
        </em>
      </div>
    </div>
  );
};
