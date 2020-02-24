import { Spinner } from "@blueprintjs/core";
import React, { FunctionComponent } from "react";
import { StashService } from "../core/StashService";

export const Stats: FunctionComponent = () => {
  const { data, error, loading } = StashService.useStats();

  function renderStats() {
    if (!data || !data.stats) { return; }
    return (
      <nav id="details-container" className="level stats">
        <div className="level-item has-text-centered">
          <div>
            <p className="title">{data.stats.scene_count}</p>
            <p className="heading">Scenes</p>
          </div>
        </div>
        <div className="level-item has-text-centered">
          <div>
            <p className="title">{data.stats.gallery_count}</p>
            <p className="heading">Galleries</p>
          </div>
        </div>
        <div className="level-item has-text-centered">
          <div>
            <p className="title">{data.stats.performer_count}</p>
            <p className="heading">Performers</p>
          </div>
        </div>
        <div className="level-item has-text-centered">
          <div>
            <p className="title">{data.stats.studio_count}</p>
            <p className="heading">Studios</p>
          </div>
        </div>
        <div className="level-item has-text-centered">
          <div>
            <p className="title">{data.stats.tag_count}</p>
            <p className="heading">Tags</p>
          </div>
        </div>
      </nav>
    );
  }

  return (
    <div id="details-container">
      {!data || loading ? <Spinner size={Spinner.SIZE_LARGE} /> : undefined}
      {!!error ? <span>error.message</span> : undefined}
      {renderStats()}
    </div>
  );
};
