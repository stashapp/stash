import React, { useCallback, useEffect, useRef, useState } from "react";
import { Badge, Button, Col, Form, InputGroup, Row } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";

import * as GQL from "src/core/generated-graphql";
import { ModalComponent } from "src/components/Shared/Modal";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { TruncatedText } from "src/components/Shared/TruncatedText";
import { Icon } from "src/components/Shared/Icon";
import { queryScrapeSceneQuery } from "src/core/StashService";
import { useToast } from "src/hooks/Toast";
import { faSearch } from "@fortawesome/free-solid-svg-icons";

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
              <Badge
                className="tag-item"
                variant="secondary"
                key={performer.name}
              >
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
              <Badge className="tag-item" variant="secondary" key={tag.name}>
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
  return (
    <div className="mt-3 search-item">
      <div className="row">
        <SceneSearchResultDetails scene={scene} />
      </div>
    </div>
  );
};

interface IProps {
  scraper: GQL.ScraperSourceInput;
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
  const [scenes, setScenes] = useState<GQL.ScrapedScene[] | undefined>();
  const [error, setError] = useState<Error | undefined>();

  const doQuery = useCallback(
    async (input: string) => {
      if (!input) return;

      setLoading(true);
      try {
        const r = await queryScrapeSceneQuery(scraper, input);
        setScenes(r.data.scrapeSingleScene);
      } catch (err) {
        if (err instanceof Error) setError(err);
      } finally {
        setLoading(false);
      }
    },
    [scraper]
  );

  useEffect(() => inputRef.current?.focus(), []);
  useEffect(() => {
    if (error) {
      Toast.error(error);
      setError(undefined);
    }
  }, [error, Toast]);

  function renderResults() {
    if (!scenes) {
      return;
    }

    return (
      <div className={CLASSNAME_LIST_CONTAINER}>
        <div className="mt-1">
          <FormattedMessage
            id="dialogs.scenes_found"
            values={{ count: scenes.length }}
          />
        </div>
        <ul className={CLASSNAME_LIST}>
          {scenes.map((s, i) => (
            // eslint-disable-next-line jsx-a11y/click-events-have-key-events, jsx-a11y/no-noninteractive-element-interactions, react/no-array-index-key
            <li key={i} onClick={() => onSelectScene(s)}>
              <SceneSearchResult scene={s} />
            </li>
          ))}
        </ul>
      </div>
    );
  }

  return (
    <ModalComponent
      show
      onHide={onHide}
      modalProps={{ size: "lg", dialogClassName: "scrape-query-dialog" }}
      header={intl.formatMessage(
        { id: "dialogs.scrape_entity_query" },
        { entity_type: intl.formatMessage({ id: "scene" }) }
      )}
      accept={{
        text: intl.formatMessage({ id: "actions.cancel" }),
        onClick: onHide,
        variant: "secondary",
      }}
    >
      <div className={CLASSNAME}>
        <InputGroup>
          <Form.Control
            defaultValue={name ?? ""}
            placeholder={`${intl.formatMessage({ id: "name" })}...`}
            className="text-input"
            ref={inputRef}
            onKeyPress={(e: React.KeyboardEvent<HTMLInputElement>) =>
              e.key === "Enter" && doQuery(inputRef.current?.value ?? "")
            }
          />
          <InputGroup.Append>
            <Button
              onClick={() => {
                doQuery(inputRef.current?.value ?? "");
              }}
              variant="primary"
              title={intl.formatMessage({ id: "actions.search" })}
            >
              <Icon icon={faSearch} />
            </Button>
          </InputGroup.Append>
        </InputGroup>

        {loading ? (
          <div className="m-4 text-center">
            <LoadingIndicator inline />
          </div>
        ) : (
          renderResults()
        )}
      </div>
    </ModalComponent>
  );
};

export default SceneQueryModal;
