import React, { useState, useEffect } from "react";
import { Badge, Col, Row } from "react-bootstrap";
import cx from "classnames";

import * as GQL from "src/core/generated-graphql";
import { TruncatedText } from "src/components/Shared";
import SceneScraperSceneEditor from "./SceneScraperSceneEditor";

interface ISceneSearchResultDetailsProps {
  scene: GQL.ScrapedSceneDataFragment;
}

const SceneSearchResultDetails: React.FC<ISceneSearchResultDetailsProps> = ({
  scene,
}) => {
  function renderPerformers() {
    if (scene.performers) {
      return (
        <Row>
          <Col>
            {scene.performers?.map((performer) => (
              <Badge className="tag-item" variant="secondary">
                {performer.name}
              </Badge>
            ))}
          </Col>
        </Row>
      );
    }
  }

  function renderTags() {
    if (scene.tags) {
      return (
        <Row>
          <Col>
            {scene.tags?.map((tag) => (
              <Badge className="tag-item" variant="secondary">
                {tag.name}
              </Badge>
            ))}
          </Col>
        </Row>
      );
    }
  }

  function renderImage() {
    if (scene.image) {
      return (
        <div className="scene-image-container">
          <img
            src={scene.image}
            alt=""
            className="align-self-center scene-image"
          />
        </div>
      );
    }
  }

  return (
    <div className="scene-details">
      <Row>
        {renderImage()}
        <div className="col flex-column">
          <h4>{scene.title}</h4>
          <h5>
            {scene.studio?.name}
            {scene.studio?.name && scene.date && ` • `}
            {scene.date}
          </h5>
        </div>
      </Row>
      <Row>
        <Col>
          <TruncatedText text={scene.details ?? ""} lineCount={3} />
        </Col>
      </Row>
      {renderPerformers()}
      {renderTags()}
    </div>
  );
};

interface ISceneSearchResult {
  scene: GQL.ScrapedSceneDataFragment;
}

// TODO - decide if we want to keep this
// eslint-disable-next-line @typescript-eslint/no-unused-vars
const SceneSearchResult: React.FC<ISceneSearchResult> = ({ scene }) => {
  return (
    <div className="mt-3 search-item">
      <div className="row">
        <SceneSearchResultDetails scene={scene} />
      </div>
    </div>
  );
};

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
