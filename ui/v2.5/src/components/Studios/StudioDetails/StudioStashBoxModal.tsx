import React, { useEffect, useRef, useState } from "react";
import { Form, Row } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";

import * as GQL from "src/core/generated-graphql";
import { ModalComponent } from "src/components/Shared/Modal";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { stashboxDisplayName } from "src/utils/stashbox";
import { useDebounce } from "src/hooks/debounce";
import { TruncatedText } from "src/components/Shared/TruncatedText";

const CLASSNAME = "StudioScrapeModal";
const CLASSNAME_LIST = `${CLASSNAME}-list`;
const CLASSNAME_LIST_CONTAINER = `${CLASSNAME_LIST}-container`;

interface IStudioSearchResultDetailsProps {
  studio: GQL.ScrapedStudioDataFragment;
}

const StudioSearchResultDetails: React.FC<IStudioSearchResultDetailsProps> = ({
  studio,
}) => {
  function renderImage() {
    if (studio.image) {
      return (
        <div className="scene-image-container">
          <img
            src={studio.image}
            alt=""
            className="align-self-center scene-image"
          />
        </div>
      );
    }
  }

  return (
    <div className="studio-result">
      <Row>
        {renderImage()}
        <div className="col flex-column">
          <h4 className="studio-name">
            <span>{studio.name}</span>
          </h4>
          {studio.parent?.name && (
            <h5 className="studio-parent text-muted">
              <span>{studio.parent.name}</span>
            </h5>
          )}
          {studio.urls && studio.urls.length > 0 && (
            <div className="studio-url text-muted small">{studio.urls[0]}</div>
          )}
        </div>
      </Row>
      {studio.details && (
        <Row>
          <div className="col">
            <TruncatedText text={studio.details} lineCount={3} />
          </div>
        </Row>
      )}
    </div>
  );
};

export interface IStudioSearchResult {
  studio: GQL.ScrapedStudioDataFragment;
}

export const StudioSearchResult: React.FC<IStudioSearchResult> = ({
  studio,
}) => {
  return (
    <div className="mt-3 search-item">
      <StudioSearchResultDetails studio={studio} />
    </div>
  );
};

export interface IStashBox extends GQL.StashBox {
  index: number;
}

interface IProps {
  instance: IStashBox;
  onHide: () => void;
  onSelectStudio: (studio: GQL.ScrapedStudio) => void;
  name?: string;
}

const StudioStashBoxModal: React.FC<IProps> = ({
  instance,
  name,
  onHide,
  onSelectStudio,
}) => {
  const intl = useIntl();
  const inputRef = useRef<HTMLInputElement>(null);
  const [query, setQuery] = useState<string>(name ?? "");
  const { data, loading } = GQL.useScrapeSingleStudioQuery({
    variables: {
      source: {
        stash_box_endpoint: instance.endpoint,
      },
      input: {
        query,
      },
    },
    skip: query === "",
  });

  const studios = data?.scrapeSingleStudio ?? [];

  const onInputChange = useDebounce(setQuery, 500);

  useEffect(() => inputRef.current?.focus(), []);

  function renderResults() {
    if (!studios) {
      return;
    }

    return (
      <div className={CLASSNAME_LIST_CONTAINER}>
        <div className="mt-1">
          <FormattedMessage
            id="dialogs.studios_found"
            values={{ count: studios.length }}
          />
        </div>
        <ul className={CLASSNAME_LIST}>
          {studios.map((s, i) => (
            // eslint-disable-next-line jsx-a11y/click-events-have-key-events, jsx-a11y/no-noninteractive-element-interactions, react/no-array-index-key
            <li key={i} onClick={() => onSelectStudio(s)}>
              <StudioSearchResult studio={s} />
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
      header={`Scrape studio from ${stashboxDisplayName(
        instance.name,
        instance.index
      )}`}
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
          placeholder={intl.formatMessage({ id: "studio_name" }) + "..."}
          className="text-input mb-4"
          ref={inputRef}
        />
        {loading ? (
          <div className="m-4 text-center">
            <LoadingIndicator inline />
          </div>
        ) : studios.length > 0 ? (
          renderResults()
        ) : (
          query !== "" && (
            <h5 className="text-center">
              <FormattedMessage id="stashbox_search.no_results" />
            </h5>
          )
        )}
      </div>
    </ModalComponent>
  );
};

export default StudioStashBoxModal;
