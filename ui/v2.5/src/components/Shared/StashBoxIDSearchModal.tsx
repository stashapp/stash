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
import { stashBoxPerformerQuery } from "src/core/StashService";
import { useToast } from "src/hooks/Toast";
import { stringToGender } from "src/utils/gender";

const CLASSNAME = "StashBoxIDSearchModal";
const CLASSNAME_LIST = `${CLASSNAME}-list`;
const CLASSNAME_LIST_CONTAINER = `${CLASSNAME_LIST}-container`;

interface IProps {
  stashBoxes: GQL.StashBox[];
  excludedStashBoxEndpoints?: string[];
  onSelectItem: (item?: GQL.StashIdInput) => void;
}

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

// Main Modal Component
export const StashBoxIDSearchModal: React.FC<IProps> = ({
  stashBoxes,
  excludedStashBoxEndpoints = [],
  onSelectItem,
}) => {
  const intl = useIntl();
  const Toast = useToast();
  const inputRef = useRef<HTMLInputElement>(null);

  const [selectedStashBox, setSelectedStashBox] = useState<GQL.StashBox | null>(
    null
  );
  const [query, setQuery] = useState<string>("");
  const [results, setResults] = useState<GQL.ScrapedPerformerDataFragment[]>(
    []
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
      const queryData = await stashBoxPerformerQuery(
        query,
        selectedStashBox.endpoint
      );
      setResults(queryData.data?.scrapeSinglePerformer ?? []);
    } catch (error) {
      Toast.error(error);
    } finally {
      setLoading(false);
    }
  }, [query, selectedStashBox, Toast]);

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

  function renderResults() {
    if (results.length === 0) {
      return null;
    }

    return (
      <div className={CLASSNAME_LIST_CONTAINER}>
        <div className="mt-1 mb-2">
          <FormattedMessage
            id="dialogs.performers_found"
            values={{ count: results.length }}
          />
        </div>
        <ul className={CLASSNAME_LIST} style={{ listStyleType: "none" }}>
          {results.map((item, i) => (
            <li key={i} onClick={() => handleItemClick(item)}>
              <PerformerSearchResult performer={item} />
            </li>
          ))}
        </ul>
      </div>
    );
  }

  return (
    <ModalComponent
      show
      onHide={handleClose}
      header={intl.formatMessage(
        { id: "stashbox_search.header" },
        { entityType: "Performer" }
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
              { entityType: "Performer" }
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
        ) : results.length > 0 ? (
          renderResults()
        ) : (
          query !== "" &&
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

export default StashBoxIDSearchModal;
