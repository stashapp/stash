import React, { useEffect, useState } from "react";
import { Button, ButtonGroup } from "react-bootstrap";
import { FormattedMessage } from "react-intl";

import * as GQL from "src/core/generated-graphql";
import { Icon } from "src/components/Shared/Icon";
import { OperationButton } from "src/components/Shared/OperationButton";
import { OptionalField } from "../IncludeButton";
import { faSave } from "@fortawesome/free-solid-svg-icons";
import {
  Performer,
  PerformerSelect,
} from "src/components/Performers/PerformerSelect";
import { getStashboxBase } from "src/utils/stashbox";
import { ExternalLink } from "src/components/Shared/ExternalLink";

interface IPerformerName {
  performer: GQL.ScrapedPerformer | Performer;
  id: string | undefined | null;
  baseURL: string | undefined;
}

const PerformerName: React.FC<IPerformerName> = ({
  performer,
  id,
  baseURL,
}) => {
  const name =
    baseURL && id ? (
      <ExternalLink href={`${baseURL}${id}`}>{performer.name}</ExternalLink>
    ) : (
      performer.name
    );

  return (
    <>
      <span>{name}</span>
      {performer.disambiguation && (
        <span className="performer-disambiguation">
          {` (${performer.disambiguation})`}
        </span>
      )}
    </>
  );
};

interface IPerformerResultProps {
  performer: GQL.ScrapedPerformer;
  selectedID: string | undefined;
  setSelectedID: (id: string | undefined) => void;
  onCreate: () => void;
  onLink?: () => Promise<void>;
  endpoint?: string;
}

const PerformerResult: React.FC<IPerformerResultProps> = ({
  performer,
  selectedID,
  setSelectedID,
  onCreate,
  onLink,
  endpoint,
}) => {
  const { data: performerData, loading: stashLoading } =
    GQL.useFindPerformerQuery({
      variables: { id: performer.stored_id ?? "" },
      skip: !performer.stored_id,
    });

  const matchedPerformer = performerData?.findPerformer;
  const matchedStashID = matchedPerformer?.stash_ids.some(
    (stashID) =>
      stashID.endpoint === endpoint &&
      stashID.stash_id === performer.remote_site_id
  );

  const [selectedPerformer, setSelectedPerformer] = useState<Performer>();

  const stashboxPerformerPrefix = endpoint
    ? `${getStashboxBase(endpoint)}performers/`
    : undefined;
  const performerURLPrefix = "/performers/";

  function selectPerformer(selected: Performer | undefined) {
    setSelectedPerformer(selected);
    setSelectedID(selected?.id);
  }

  useEffect(() => {
    if (
      performerData?.findPerformer &&
      selectedID === performerData?.findPerformer?.id
    ) {
      setSelectedPerformer(performerData.findPerformer);
    }
  }, [performerData?.findPerformer, selectedID]);

  const handleSelect = (performers: Performer[]) => {
    if (performers.length) {
      selectPerformer(performers[0]);
    } else {
      selectPerformer(undefined);
    }
  };

  const handleSkip = () => {
    selectPerformer(undefined);
  };

  if (stashLoading) return <div>Loading performer</div>;

  if (matchedPerformer && matchedStashID) {
    return (
      <div className="row no-gutters my-2">
        <div className="entity-name">
          <FormattedMessage id="countables.performers" values={{ count: 1 }} />:
          <b className="ml-2">
            <PerformerName
              performer={performer}
              id={performer.remote_site_id}
              baseURL={stashboxPerformerPrefix}
            />
          </b>
        </div>
        <span className="ml-auto">
          <OptionalField
            exclude={selectedID === undefined}
            setExclude={(v) =>
              v ? handleSkip() : setSelectedID(matchedPerformer.id)
            }
          >
            <div>
              <span className="mr-2">
                <FormattedMessage id="component_tagger.verb_matched" />:
              </span>
              <b className="col-3 text-right">
                <PerformerName
                  performer={matchedPerformer}
                  id={matchedPerformer.id}
                  baseURL={performerURLPrefix}
                />
              </b>
            </div>
          </OptionalField>
        </span>
      </div>
    );
  }

  function maybeRenderLinkButton() {
    if (endpoint && onLink) {
      return (
        <OperationButton
          variant="secondary"
          disabled={selectedID === undefined}
          operation={onLink}
          hideChildrenWhenLoading
        >
          <Icon icon={faSave} />
        </OperationButton>
      );
    }
  }

  const selectedSource = !selectedID ? "skip" : "existing";

  return (
    <div className="row no-gutters align-items-center mt-2">
      <div className="entity-name">
        <FormattedMessage id="countables.performers" values={{ count: 1 }} />:
        <b className="ml-2">
          <PerformerName
            performer={performer}
            id={performer.remote_site_id}
            baseURL={stashboxPerformerPrefix}
          />
        </b>
      </div>
      <ButtonGroup>
        <Button variant="secondary" onClick={() => onCreate()}>
          <FormattedMessage id="actions.create" />
        </Button>
        <Button
          variant={selectedSource === "skip" ? "primary" : "secondary"}
          onClick={() => handleSkip()}
        >
          <FormattedMessage id="actions.skip" />
        </Button>
        <PerformerSelect
          values={selectedPerformer ? [selectedPerformer] : []}
          onSelect={handleSelect}
          active={selectedSource === "existing"}
          isClearable={false}
        />
        {maybeRenderLinkButton()}
      </ButtonGroup>
    </div>
  );
};

export default PerformerResult;
