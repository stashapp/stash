import React, { useCallback, useEffect, useRef, useState } from "react";
import { Form, Button, Row, Col, Badge, InputGroup } from "react-bootstrap";
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
  stashBoxSceneQuery,
  stashBoxStudioQuery,
  stashBoxTagQuery,
} from "src/core/StashService";
import { useToast } from "src/hooks/Toast";
import { stringToGender } from "src/utils/gender";

type SearchResultItem =
  | GQL.ScrapedPerformerDataFragment
  | GQL.ScrapedSceneDataFragment
  | GQL.ScrapedStudioDataFragment
  | GQL.ScrapedSceneTagDataFragment;

export type StashBoxEntityType = "performer" | "scene" | "studio" | "tag";

interface IProps {
  entityType: StashBoxEntityType;
  stashBoxes: GQL.StashBox[];
  excludedStashBoxEndpoints?: string[];
  onSelectItem: (item?: GQL.StashIdInput) => void;
  initialQuery?: string;
}

const CLASSNAME = "StashBoxIDSearchModal";
const CLASSNAME_LIST = `${CLASSNAME}-list`;
const CLASSNAME_LIST_CONTAINER = `${CLASSNAME_LIST}-container`;

interface IHasRemoteSiteID {
  remote_site_id?: string | null;
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
    <div className="mt-3 search-item" style={{ cursor: "pointer" }}>
      <PerformerSearchResultDetails performer={performer} />
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
          <h4 className="scene-title">
            <span>{scene.title}</span>
            {scene.code && (
              <span className="scene-code">{` (${scene.code})`}</span>
            )}
          </h4>
          <h5 className="scene-details">
            {scene.studio?.name && <span>{scene.studio.name}</span>}
            {scene.date && (
              <span className="scene-date">{` • ${scene.date}`}</span>
            )}
          </h5>
          {scene.performers && scene.performers.length > 0 && (
            <div className="scene-performers">
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
    <div className="mt-3 search-item" style={{ cursor: "pointer" }}>
      <SceneSearchResultDetails scene={scene} />
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
          <h4 className="studio-name">
            <span>{studio.name}</span>
          </h4>
          {studio.parent?.name && (
            <h5 className="studio-parent">
              <span>{studio.parent.name}</span>
            </h5>
          )}
          {studio.urls && studio.urls.length > 0 && (
            <div className="studio-url text-muted small">{studio.urls[0]}</div>
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
    <div className="mt-3 search-item" style={{ cursor: "pointer" }}>
      <StudioSearchResultDetails studio={studio} />
    </div>
  );
};

// Tag Result Component
interface ITagResultProps {
  tag: GQL.ScrapedSceneTagDataFragment;
}

export const TagSearchResult: React.FC<ITagResultProps> = ({ tag }) => {
  return (
    <div className="mt-3 search-item" style={{ cursor: "pointer" }}>
      <div className="tag-result">
        <Row>
          <div className="col flex-column">
            <h4 className="tag-name">
              <span>{tag.name}</span>
            </h4>
          </div>
        </Row>
      </div>
    </div>
  );
};

// Helper to get entity type message id for i18n
function getEntityTypeMessageId(entityType: StashBoxEntityType): string {
  switch (entityType) {
    case "performer":
      return "performer";
    case "scene":
      return "scene";
    case "studio":
      return "studio";
    case "tag":
      return "tag";
  }
}

// Helper to get the "found" message id based on entity type
function getFoundMessageId(entityType: StashBoxEntityType): string {
  switch (entityType) {
    case "performer":
      return "dialogs.performers_found";
    case "scene":
      return "dialogs.scenes_found";
    case "studio":
      return "dialogs.studios_found";
    case "tag":
      return "dialogs.tags_found";
  }
}

// Main Modal Component
export const StashBoxIDSearchModal: React.FC<IProps> = ({
  entityType,
  stashBoxes,
  excludedStashBoxEndpoints = [],
  onSelectItem,
  initialQuery = "",
}) => {
  const intl = useIntl();
  const Toast = useToast();
  const inputRef = useRef<HTMLInputElement>(null);

  const [selectedStashBox, setSelectedStashBox] = useState<GQL.StashBox | null>(
    null
  );
  const [query, setQuery] = useState<string>(initialQuery);
  const [results, setResults] = useState<SearchResultItem[] | undefined>(
    undefined
  );
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (stashBoxes.length > 0) {
      setSelectedStashBox(stashBoxes[0]);
    }
  }, [stashBoxes]);

  useEffect(() => inputRef.current?.focus(), []);

  const doSearch = useCallback(async () => {
    if (!selectedStashBox || !query) {
      return;
    }

    setLoading(true);
    setResults([]);

    try {
      switch (entityType) {
        case "performer": {
          const queryData = await stashBoxPerformerQuery(
            query,
            selectedStashBox.endpoint
          );
          setResults(queryData.data?.scrapeSinglePerformer ?? []);
          break;
        }
        case "scene": {
          const queryData = await stashBoxSceneQuery(
            query,
            selectedStashBox.endpoint
          );
          setResults(queryData.data?.scrapeSingleScene ?? []);
          break;
        }
        case "studio": {
          const queryData = await stashBoxStudioQuery(
            query,
            selectedStashBox.endpoint
          );
          setResults(queryData.data?.scrapeSingleStudio ?? []);
          break;
        }
        case "tag": {
          const queryData = await stashBoxTagQuery(
            query,
            selectedStashBox.endpoint
          );
          setResults(queryData.data?.scrapeSingleTag ?? []);
          break;
        }
      }
    } catch (error) {
      Toast.error(error);
    } finally {
      setLoading(false);
    }
  }, [query, selectedStashBox, Toast, entityType]);

  function handleItemClick(item: IHasRemoteSiteID) {
    if (selectedStashBox && item.remote_site_id) {
      onSelectItem({
        endpoint: selectedStashBox.endpoint,
        stash_id: item.remote_site_id,
      });
    } else {
      onSelectItem(undefined);
    }
  }

  function handleClose() {
    onSelectItem(undefined);
  }

  function renderResultItem(item: SearchResultItem) {
    switch (entityType) {
      case "performer":
        return (
          <PerformerSearchResult
            performer={item as GQL.ScrapedPerformerDataFragment}
          />
        );
      case "scene":
        return (
          <SceneSearchResult scene={item as GQL.ScrapedSceneDataFragment} />
        );
      case "studio":
        return (
          <StudioSearchResult studio={item as GQL.ScrapedStudioDataFragment} />
        );
      case "tag":
        return (
          <TagSearchResult tag={item as GQL.ScrapedSceneTagDataFragment} />
        );
    }
  }

  function renderResults() {
    if (!results || results.length === 0) {
      return null;
    }

    return (
      <div className={CLASSNAME_LIST_CONTAINER}>
        <div className="mt-1 mb-2">
          <FormattedMessage
            id={getFoundMessageId(entityType)}
            values={{ count: results.length }}
          />
        </div>
        <ul className={CLASSNAME_LIST} style={{ listStyleType: "none" }}>
          {results.map((item, i) => (
            <li key={i} onClick={() => handleItemClick(item)}>
              {renderResultItem(item)}
            </li>
          ))}
        </ul>
      </div>
    );
  }

  const entityTypeDisplayName = intl.formatMessage({
    id: getEntityTypeMessageId(entityType),
  });

  return (
    <ModalComponent
      show
      onHide={handleClose}
      header={intl.formatMessage(
        { id: "stashbox_search.header" },
        { entityType: entityTypeDisplayName }
      )}
      accept={{
        text: intl.formatMessage({ id: "actions.cancel" }),
        onClick: handleClose,
        variant: "secondary",
      }}
    >
      <div className={CLASSNAME}>
        <Form.Group className="d-flex align-items-center mb-3">
          <Form.Label className="mb-0 mr-2" style={{ flexShrink: 0 }}>
            <FormattedMessage id="stashbox.source" />
          </Form.Label>
          <Form.Control
            as="select"
            className="input-control"
            style={{ flex: "0 1 auto" }}
            value={selectedStashBox?.endpoint ?? ""}
            onChange={(e) => {
              const box = stashBoxes.find(
                (b) => b.endpoint === e.currentTarget.value
              );
              if (box) {
                setSelectedStashBox(box);
              }
            }}
          >
            {stashBoxes.map((box, index) => (
              <option key={box.endpoint} value={box.endpoint}>
                {stashboxDisplayName(box.name, index)}
              </option>
            ))}
          </Form.Control>
        </Form.Group>

        {selectedStashBox &&
          excludedStashBoxEndpoints.includes(selectedStashBox.endpoint) && (
            <span className="saved-filter-overwrite-warning mb-3 d-block">
              <FormattedMessage id="dialogs.stashid_exists_warning" />
            </span>
          )}

        <InputGroup>
          <Form.Control
            onChange={(e) => setQuery(e.currentTarget.value)}
            value={query}
            placeholder={intl.formatMessage(
              { id: "stashbox_search.placeholder_name_or_id" },
              { entityType: entityTypeDisplayName }
            )}
            className="text-input"
            ref={inputRef}
            onKeyPress={(e: React.KeyboardEvent<HTMLInputElement>) =>
              e.key === "Enter" && doSearch()
            }
          />
          <InputGroup.Append>
            <Button
              onClick={doSearch}
              variant="primary"
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
        ) : results && results.length > 0 ? (
          renderResults()
        ) : (
          results !== undefined &&
          results.length === 0 && (
            <h5 className="text-center mt-4">
              <FormattedMessage id="stashbox_search.no_results" />
            </h5>
          )
        )}
      </div>
    </ModalComponent>
  );
};

export default StashBoxIDSearchModal;
