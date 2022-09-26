import React, { useCallback, useEffect, useRef, useState } from "react";
import { Form, InputGroup, Row, Col, Button, Badge } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";

import * as GQL from "src/core/generated-graphql";
import { Icon } from "src/components/Shared/Icon";
import { ModalComponent } from "src/components/Shared/Modal";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { TruncatedText } from "src/components/Shared/TruncatedText";
import { genderToString, stringToGender } from "src/utils/gender";
import { queryScrapePerformerQuery } from "src/core/StashService";
import { useToast } from "src/hooks/Toast";
import { faSearch } from "@fortawesome/free-solid-svg-icons";
import TextUtils from "src/utils/text";

interface IPerformerSearchResultDetailsProps {
  performer: GQL.ScrapedPerformerDataFragment;
}

const PerformerSearchResultDetails: React.FC<IPerformerSearchResultDetailsProps> = ({
  performer,
}) => {
  function renderImage() {
    if (performer.images && performer.images.length > 0) {
      return (
        <div className="scene-image-container">
          <img
            src={performer.images[0]}
            alt=""
            className="align-self-center scene-image"
          />
        </div>
      );
    }
  }

  function calculateAge() {
    if (performer?.birthdate) {
      // calculate the age from birthdate. In future, this should probably be
      // provided by the server
      return TextUtils.age(performer.birthdate, performer.death_date);
    }
  }

  function renderTags() {
    if (performer.tags) {
      return (
        <Row>
          <Col>
            {performer.tags?.map((tag) => (
              <Badge
                className="tag-item"
                variant="secondary"
                key={tag.stored_id}
              >
                {tag.name}
              </Badge>
            ))}
          </Col>
        </Row>
      );
    }
  }

  let calculated_age = calculateAge();

  return (
    <div className="scene-details">
      <Row>
        {renderImage()}
        <div className="col flex-column">
          <h4>{performer.name}{performer.disambiguation && ` (${performer.disambiguation})`}</h4>
          <h5>
            {performer.gender &&
              genderToString(stringToGender(performer.gender, true))}
            {performer.gender && calculated_age && ` • `}
            {calculated_age}
            {calculated_age && " "}
            {calculated_age && <FormattedMessage id="years_old" />}
          </h5>
        </div>
      </Row>
      <Row>
        <Col>
          <TruncatedText text={performer.details ?? ""} lineCount={3} />
        </Col>
      </Row>
      {renderTags()}
    </div>
  );
};

export interface IPerformerSearchResult {
  performer: GQL.ScrapedPerformerDataFragment;
}

export const PerformerSearchResult: React.FC<IPerformerSearchResult> = ({
  performer,
}) => {
  return (
    <div className="mt-3 search-item">
      <PerformerSearchResultDetails performer={performer} />
    </div>
  );
};

interface IProps {
  scraper: GQL.ScraperSourceInput;
  onHide: () => void;
  onSelectPerformer: (performer: GQL.ScrapedPerformerDataFragment) => void;
  name?: string;
}

const PerformerScrapeModal: React.FC<IProps> = ({
  scraper,
  name,
  onHide,
  onSelectPerformer,
}) => {
  const CLASSNAME = "PerformerScrapeModal";
  const CLASSNAME_LIST = `${CLASSNAME}-list`;
  const CLASSNAME_LIST_CONTAINER = `${CLASSNAME_LIST}-container`;

  const intl = useIntl();
  const Toast = useToast();

  const inputRef = useRef<HTMLInputElement>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [performers, setPerformers] = useState<
    GQL.ScrapedPerformer[] | undefined
  >();
  const [error, setError] = useState<Error | undefined>();

  const doQuery = useCallback(
    async (input: string) => {
      if (!input) return;

      setLoading(true);
      try {
        const r = await queryScrapePerformerQuery(scraper, input);
        setPerformers(r.data?.scrapeSinglePerformer);
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
    doQuery(name ?? "");
  }, [doQuery, name]);
  useEffect(() => {
    if (error) {
      Toast.error(error);
      setError(undefined);
    }
  }, [error, Toast]);

  function renderResults() {
    if (!performers) {
      return;
    }

    return (
      <div className={CLASSNAME_LIST_CONTAINER}>
        <div className="mt-1">
          <FormattedMessage
            id="dialogs.performers_found"
            values={{ count: performers.length }}
          />
        </div>
        <ul className={CLASSNAME_LIST}>
          {performers.map((p, i) => (
            // eslint-disable-next-line jsx-a11y/click-events-have-key-events, jsx-a11y/no-noninteractive-element-interactions, react/no-array-index-key
            <li key={i} onClick={() => onSelectPerformer(p)}>
              <PerformerSearchResult performer={p} />
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
        { entity_type: intl.formatMessage({ id: "performer" }) }
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

export default PerformerScrapeModal;
