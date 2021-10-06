import React, { useState, useEffect } from "react";
import cx from "classnames";

import * as GQL from "src/core/generated-graphql";
import SceneScraperSceneEditor from "./SceneScraperSceneEditor";

export interface ISceneSearchResults {
  target: GQL.SlimSceneDataFragment;
  scenes: GQL.ScrapedSceneDataFragment[];
}

export const SceneSearchResults: React.FC<ISceneSearchResults> = ({
  target,
  scenes,
}) => {
  const [selectedResult, setSelectedResult] = useState<number | undefined>();

  useEffect(() => {
    if (!scenes) {
      setSelectedResult(undefined);
    }
  }, [scenes]);

  function getClassName(i: number) {
    return cx("row mx-0 mt-2 search-result", {
      "selected-result active": i === selectedResult,
    });
  }

  return (
    <ul>
      {scenes.map((s, i) => (
        // eslint-disable-next-line jsx-a11y/click-events-have-key-events, jsx-a11y/no-noninteractive-element-interactions, react/no-array-index-key
        <li
          // eslint-disable-next-line react/no-array-index-key
          key={i}
          onClick={() => setSelectedResult(i)}
          className={getClassName(i)}
        >
          {/* <SceneSearchResult scene={s} /> */}
          <SceneScraperSceneEditor
            index={i}
            isActive={i === selectedResult}
            scene={s}
            stashScene={target}
          />
        </li>
      ))}
    </ul>
  );
};
