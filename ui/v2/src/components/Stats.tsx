import { H1, Spinner } from "@blueprintjs/core";
import React, { FunctionComponent } from "react";
import { StashService } from "../core/StashService";

export const Stats: FunctionComponent = () => {
  const { data, error, loading } = StashService.useStats();

  function renderStats() {
    if (!data || !data.stats) { return; }
    return (
      <nav id="details-container" className="level">
        <div className="level-item has-text-centered">
          <div>
            <p className="heading">Scenes</p>
            <p className="title">{data.stats.scene_count}</p>
          </div>
        </div>
        <div className="level-item has-text-centered">
          <div>
            <p className="heading">Galleries</p>
            <p className="title">{data.stats.gallery_count}</p>
          </div>
        </div>
        <div className="level-item has-text-centered">
          <div>
            <p className="heading">Performers</p>
            <p className="title">{data.stats.performer_count}</p>
          </div>
        </div>
        <div className="level-item has-text-centered">
          <div>
            <p className="heading">Studios</p>
            <p className="title">{data.stats.studio_count}</p>
          </div>
        </div>
        <div className="level-item has-text-centered">
          <div>
            <p className="heading">Dvds</p>
            <p className="title">{data.stats.dvd_count}</p>
          </div>
        </div>

        <div className="level-item has-text-centered">
          <div>
            <p className="heading">Tags</p>
            <p className="title">{data.stats.tag_count}</p>
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

      <h3>Notes</h3>
      <pre>
        {`
        This is still an early version, some things are still a work in progress.
        `}
      </pre>
    </div>
  );
};
