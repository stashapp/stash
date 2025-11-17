import React, { useCallback, useEffect, useRef, useState } from "react";
import {
  Form,
  Button,
  Dropdown,
  Row,
  Col,
  Badge,
  InputGroup,
} from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import { faSearch } from "@fortawesome/free-solid-svg-icons";
import * as GQL from "src/core/generated-graphql";
import { ModalComponent } from "src/components/Shared/Modal";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { stashboxDisplayName } from "src/utils/stashbox";
import { TruncatedText } from "src/components/Shared/TruncatedText";
import TextUtils from "src/utils/text";
import GenderIcon from "src/components/Performers/GenderIcon";
import { CountryFlag } from "src/components/Shared/CountryFlag";
import { Icon } from "src/components/Shared/Icon";
import {
  stashBoxPerformerQuery,
  stashBoxStudioQuery,
  stashBoxSceneQuery,
} from "src/core/StashService";
import { useToast } from "src/hooks/Toast";
import { stringToGender } from "src/utils/gender";

const CLASSNAME = "StashBoxSearchModal";
const CLASSNAME_LIST = `${CLASSNAME}-list`;
const CLASSNAME_LIST_CONTAINER = `${CLASSNAME_LIST}-container`;

export type EntityType = "performer" | "studio" | "scene";

interface IStashBox extends GQL.StashBox {
  index: number;
}

interface IProps {
  entityType: EntityType;
  stashBoxes: GQL.StashBox[];
  onHide: () => void;
  onSelectItem: (
    item: GQL.ScrapedPerformer | GQL.ScrapedStudio | GQL.ScrapedScene,
    endpoint: string
  ) => void;
}

// Shared component for rendering images
const SearchResultImage: React.FC<{ imageUrl?: string | null }> = ({
  imageUrl,
}) => {
  if (!imageUrl) return null;

  return (
    <div className="scene-image-container">
      <img src={imageUrl} alt="" className="align-self-center scene-image" />
    </div>
  );
};

// Shared component for rendering tags
const SearchResultTags: React.FC<{
  tags?: GQL.ScrapedTag[] | null;
}> = ({ tags }) => {
  if (!tags || tags.length === 0) return null;

  return (
    <Row>
      <Col>
        {tags.map((tag) => (
          <Badge className="tag-item" variant="secondary" key={tag.stored_id}>
            {tag.name}
          </Badge>
        ))}
      </Col>
    </Row>
  );
};

// Performer Result Component
interface IPerformerResultProps {
  performer: GQL.ScrapedPerformerDataFragment;
}

const PerformerSearchResultDetails: React.FC<IPerformerResultProps> = ({
  performer,
}) => {
  const age = performer?.birthdate
    ? TextUtils.age(performer.birthdate, performer.death_date)
    : undefined;

  return (
    <div className="performer-result">
      <Row>
        <SearchResultImage imageUrl={performer.images?.[0]} />
        <div className="col flex-column">
          <h4 className="performer-name">
            <span>{performer.name}</span>
            {performer.disambiguation && (
              <span className="performer-disambiguation">
                {` (${performer.disambiguation})`}
              </span>
            )}
          </h4>
          <h5 className="performer-details">
            {performer.gender && (
              <span>
                <GenderIcon
                  className="gender-icon"
                  gender={stringToGender(performer.gender, true)}
                />
              </span>
            )}
            {age && (
              <span>
                {`${age} `}
                <FormattedMessage id="years_old" />
              </span>
            )}
          </h5>
          {performer.country && (
            <span>
              <CountryFlag
                className="performer-result__country-flag"
                country={performer.country}
              />
            </span>
          )}
        </div>
      </Row>
      <Row>
        <Col>
          <TruncatedText text={performer.details ?? ""} lineCount={3} />
        </Col>
      </Row>
      <SearchResultTags tags={performer.tags} />
    </div>
  );
};

export const PerformerSearchResult: React.FC<IPerformerResultProps> = ({
  performer,
}) => {
  return (
    <div className="mt-3 search-item">
      <PerformerSearchResultDetails performer={performer} />
    </div>
  );
};

// Studio Result Component
interface IStudioResultProps {
  studio: GQL.ScrapedStudioDataFragment;
}

const StudioSearchResultDetails: React.FC<IStudioResultProps> = ({
  studio,
}) => {
  return (
    <div className="studio-result">
      <Row>
        <SearchResultImage imageUrl={studio.image} />
        <div className="col flex-column">
          <h4>{studio.name}</h4>
          {studio.urls && studio.urls.length > 0 && (
            <p>
              <a
                href={studio.urls[0]}
                target="_blank"
                rel="noopener noreferrer"
              >
                {studio.urls[0]}
              </a>
            </p>
          )}
        </div>
      </Row>
    </div>
  );
};

export const StudioSearchResult: React.FC<IStudioResultProps> = ({
  studio,
}) => {
  return (
    <div className="mt-3 search-item">
      <StudioSearchResultDetails studio={studio} />
    </div>
  );
};

// Scene Result Component
interface ISceneResultProps {
  scene: GQL.ScrapedSceneDataFragment;
}

const SceneSearchResultDetails: React.FC<ISceneResultProps> = ({ scene }) => {
  return (
    <div className="scene-result">
      <Row>
        <SearchResultImage imageUrl={scene.image} />
        <div className="col flex-column">
          <h4>{scene.title}</h4>
          {scene.studio?.name && (
            <div>
              <strong>
                <FormattedMessage id="studio" />:{" "}
              </strong>
              {scene.studio.name}
            </div>
          )}
          {scene.date && (
            <div>
              <strong>
                <FormattedMessage id="date" />:{" "}
              </strong>
              {scene.date}
            </div>
          )}
          {scene.performers && scene.performers.length > 0 && (
            <div>
              <strong>
                <FormattedMessage id="performers" />:{" "}
              </strong>
              {scene.performers.map((p) => p.name).join(", ")}
            </div>
          )}
        </div>
      </Row>
      <Row>
        <Col>
          <TruncatedText text={scene.details ?? ""} lineCount={3} />
        </Col>
      </Row>
      <SearchResultTags tags={scene.tags} />
    </div>
  );
};

export const SceneSearchResult: React.FC<ISceneResultProps> = ({ scene }) => {
  return (
    <div className="mt-3 search-item">
      <SceneSearchResultDetails scene={scene} />
    </div>
  );
};

// Main Modal Component
export const StashBoxSearchModal: React.FC<IProps> = ({
  entityType,
  stashBoxes,
  onHide,
  onSelectItem,
}) => {
  const intl = useIntl();
  const Toast = useToast();
  const inputRef = useRef<HTMLInputElement>(null);

  const [selectedStashBox, setSelectedStashBox] = useState<IStashBox | null>(
    null
  );
  const [results, setResults] = useState<
    | GQL.ScrapedPerformerDataFragment[]
    | GQL.ScrapedStudioDataFragment[]
    | GQL.ScrapedSceneDataFragment[]
  >([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (stashBoxes.length > 0) {
      setSelectedStashBox({ ...stashBoxes[0], index: 0 });
    }
  }, [stashBoxes]);

  useEffect(() => inputRef.current?.focus(), []);

  const doSearch = useCallback(async () => {
    const query = inputRef.current?.value;
    if (!selectedStashBox || !query) {
      return;
    }

    setLoading(true);
    setResults([]);

    try {
      let queryData;

      switch (entityType) {
        case "performer":
          queryData = await stashBoxPerformerQuery(
            query,
            selectedStashBox.endpoint
          );
          setResults(queryData.data?.scrapeSinglePerformer ?? []);
          break;
        case "studio":
          queryData = await stashBoxStudioQuery(
            query,
            selectedStashBox.endpoint
          );
          setResults(queryData.data?.scrapeSingleStudio ?? []);
          break;
        case "scene":
          queryData = await stashBoxSceneQuery(
            query,
            selectedStashBox.endpoint
          );
          setResults(queryData.data?.scrapeSingleScene ?? []);
          break;
      }
    } catch (error) {
      Toast.error(error);
    } finally {
      setLoading(false);
    }
  }, [selectedStashBox, entityType, Toast]);

  function handleItemClick(
    item:
      | GQL.ScrapedPerformerDataFragment
      | GQL.ScrapedStudioDataFragment
      | GQL.ScrapedSceneDataFragment
  ) {
    if (selectedStashBox) {
      onSelectItem(item, selectedStashBox.endpoint);
    }
  }

  function renderResults() {
    if (results.length === 0) {
      return null;
    }

    return (
      <div className={CLASSNAME_LIST_CONTAINER}>
        <div className="mt-1 mb-2">
          <FormattedMessage
            id={`dialogs.${entityType}s_found`}
            values={{ count: results.length }}
          />
        </div>
        <ul className={CLASSNAME_LIST}>
          {results.map((item, i) => (
            <li key={i} onClick={() => handleItemClick(item)}>
              {entityType === "performer" && (
                <PerformerSearchResult
                  performer={item as GQL.ScrapedPerformerDataFragment}
                />
              )}
              {entityType === "studio" && (
                <StudioSearchResult
                  studio={item as GQL.ScrapedStudioDataFragment}
                />
              )}
              {entityType === "scene" && (
                <SceneSearchResult
                  scene={item as GQL.ScrapedSceneDataFragment}
                />
              )}
            </li>
          ))}
        </ul>
      </div>
    );
  }

  const entityTypeLabel =
    entityType.charAt(0).toUpperCase() + entityType.slice(1);

  return (
    <ModalComponent
      show
      onHide={onHide}
      header={intl.formatMessage(
        { id: "stashbox_search.header" },
        { entityType: entityTypeLabel }
      )}
      accept={{
        text: intl.formatMessage({ id: "actions.cancel" }),
        onClick: onHide,
        variant: "secondary",
      }}
    >
      <div className={CLASSNAME}>
        <Form.Group className="d-flex align-items-center mb-3">
          <Form.Label className="mb-0 mr-2">
            <FormattedMessage id="stashbox_instance" />
          </Form.Label>
          <Dropdown>
            <Dropdown.Toggle variant="secondary">
              {selectedStashBox
                ? stashboxDisplayName(
                    selectedStashBox.name,
                    selectedStashBox.index
                  )
                : intl.formatMessage({ id: "stashbox_search.select_stashbox" })}
            </Dropdown.Toggle>
            <Dropdown.Menu>
              {stashBoxes.map((box, index) => (
                <Dropdown.Item
                  key={box.endpoint}
                  onClick={() => setSelectedStashBox({ ...box, index })}
                >
                  {stashboxDisplayName(box.name, index)}
                </Dropdown.Item>
              ))}
            </Dropdown.Menu>
          </Dropdown>
        </Form.Group>

        <InputGroup className="mb-3">
          <Form.Control
            ref={inputRef}
            type="text"
            placeholder={intl.formatMessage(
              { id: "stashbox_search.placeholder_name_or_id" },
              { entityType: entityTypeLabel }
            )}
            onKeyPress={(e: React.KeyboardEvent<HTMLInputElement>) =>
              e.key === "Enter" && doSearch()
            }
          />
          <InputGroup.Append>
            <Button
              variant="primary"
              onClick={doSearch}
              disabled={!selectedStashBox}
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
        ) : results.length > 0 ? (
          renderResults()
        ) : (
          inputRef.current?.value &&
          !loading && (
            <h5 className="text-center">
              <FormattedMessage id="stashbox_search.no_results" />
            </h5>
          )
        )}
      </div>
    </ModalComponent>
  );
};

export default StashBoxSearchModal;
