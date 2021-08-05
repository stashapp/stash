import React, { useCallback, useEffect, useRef, useState } from "react";
import { debounce } from "lodash";
import { Badge, Col, Form, Row } from "react-bootstrap";
import { useIntl } from "react-intl";

import * as GQL from "src/core/generated-graphql";
import { Modal, LoadingIndicator, TruncatedText } from "src/components/Shared";
import {
  queryScrapeSceneQuery,
  queryStashBoxSceneQuery,
} from "src/core/StashService";
import { IStashBox } from "src/components/Performers/PerformerDetails/PerformerStashBoxModal";
import { useToast } from "src/hooks";

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

export interface ISceneSearchResult {
  scene: GQL.ScrapedSceneDataFragment;
}

export const SceneSearchResult: React.FC<ISceneSearchResult> = ({ scene }) => {
  // const width = scene.file.width ? scene.file.width : 0;
  // const height = scene.file.height ? scene.file.height : 0;
  // const isPortrait = height > width;

  return (
    <div className="mt-3 search-item">
      <div className="row">
        <SceneSearchResultDetails scene={scene} />
      </div>
    </div>
  );
};

interface IProps {
  scraper: GQL.Scraper | IStashBox;
  onHide: () => void;
  onSelectScene: (scene: GQL.ScrapedSceneDataFragment) => void;
  name?: string;
}
export const SceneQueryModal: React.FC<IProps> = ({
  scraper,
  name,
  onHide,
  onSelectScene,
}) => {
  const CLASSNAME = "SceneScrapeModal";
  const CLASSNAME_LIST = `${CLASSNAME}-list`;
  const CLASSNAME_LIST_CONTAINER = `${CLASSNAME_LIST}-container`;

  const intl = useIntl();
  const Toast = useToast();

  const inputRef = useRef<HTMLInputElement>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [scenes, setScenes] = useState<GQL.ScrapedScene[]>([]);

  const doQuery = useCallback(
    async (input: string) => {
      if (!input) return;

      setLoading(true);
      try {
        if ((scraper as IStashBox).index !== undefined) {
          const r = await queryStashBoxSceneQuery(
            (scraper as IStashBox).index,
            input
          );
          setScenes(r.data.queryStashBoxScene);
        } else {
          const r = await queryScrapeSceneQuery(
            (scraper as GQL.Scraper).id,
            input
          );
          setScenes(r.data.scrapeSceneQuery);
        }
      } catch (err) {
        Toast.error(err);
      } finally {
        setLoading(false);
      }
    },
    [Toast, scraper]
  );

  const onInputChange = debounce((input: string) => {
    doQuery(input);
  }, 500);

  useEffect(() => inputRef.current?.focus(), []);

  useEffect(() => {
    if (name) doQuery(name);
  }, [name, doQuery]);

  // TODO - message for header, placeholder
  return (
    <Modal
      show
      onHide={onHide}
      modalProps={{ size: "lg", dialogClassName: "scrape-query-dialog" }}
      header={`Scrape scene from ${scraper.name}`}
      accept={{
        text: intl.formatMessage({ id: "actions.cancel" }),
        onClick: onHide,
        variant: "secondary",
      }}
    >
      <div className={CLASSNAME}>
        <Form.Control
          onChange={(e) => onInputChange(e.currentTarget.value)}
          defaultValue={name ?? ""}
          placeholder="Scene name..."
          className="text-input mb-4"
          ref={inputRef}
        />
        {loading ? (
          <div className="m-4 text-center">
            <LoadingIndicator inline />
          </div>
        ) : (
          <div className={CLASSNAME_LIST_CONTAINER}>
            <ul className={CLASSNAME_LIST}>
              {scenes.map((s, i) => (
                // eslint-disable-next-line jsx-a11y/click-events-have-key-events, jsx-a11y/no-noninteractive-element-interactions, react/no-array-index-key
                <li key={i} onClick={() => onSelectScene(s)}>
                  <SceneSearchResult scene={s} />
                </li>
              ))}
            </ul>
          </div>
        )}
      </div>
    </Modal>
  );
};
