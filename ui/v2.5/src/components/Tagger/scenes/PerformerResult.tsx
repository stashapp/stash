import React from "react";
import { Button, ButtonGroup } from "react-bootstrap";
import { FormattedMessage } from "react-intl";
import cx from "classnames";

import * as GQL from "src/core/generated-graphql";
import {
  Icon,
  OperationButton,
  PerformerSelect,
  SelectObject,
} from "src/components/Shared";
import { OptionalField } from "../IncludeButton";
import { faSave } from "@fortawesome/free-solid-svg-icons";

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
  const {
    data: performerData,
    loading: stashLoading,
  } = GQL.useFindPerformerQuery({
    variables: { id: performer.stored_id ?? "" },
    skip: !performer.stored_id,
  });

  const matchedPerformer = performerData?.findPerformer;
  const matchedStashID = matchedPerformer?.stash_ids.some(
    (stashID) => stashID.endpoint === endpoint && stashID.stash_id
  );

  const handlePerformerSelect = (performers: SelectObject[]) => {
    if (performers.length) {
      setSelectedID(performers[0].id);
    } else {
      setSelectedID(undefined);
    }
  };

  const handlePerformerSkip = () => {
    setSelectedID(undefined);
  };

  if (stashLoading) return <div>Loading performer</div>;

  if (matchedPerformer && matchedStashID) {
    return (
      <div className="row no-gutters my-2">
        <div className="entity-name">
          <FormattedMessage id="countables.performers" values={{ count: 1 }} />:
          <b className="ml-2">{performer.name}</b>
        </div>
        <span className="ml-auto">
          <OptionalField
            exclude={selectedID === undefined}
            setExclude={(v) =>
              v ? handlePerformerSkip() : setSelectedID(matchedPerformer.id)
            }
          >
            <div>
              <span className="mr-2">
                <FormattedMessage id="component_tagger.verb_matched" />:
              </span>
              <b className="col-3 text-right">{matchedPerformer.name}</b>
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
        <b className="ml-2">{performer.name}</b>
      </div>
      <ButtonGroup>
        <Button variant="secondary" onClick={() => onCreate()}>
          <FormattedMessage id="actions.create" />
        </Button>
        <Button
          variant={selectedSource === "skip" ? "primary" : "secondary"}
          onClick={() => handlePerformerSkip()}
        >
          <FormattedMessage id="actions.skip" />
        </Button>
        <PerformerSelect
          ids={selectedID ? [selectedID] : []}
          onSelect={handlePerformerSelect}
          className={cx("performer-select", {
            "performer-select-active": selectedSource === "existing",
          })}
          isClearable={false}
        />
        {maybeRenderLinkButton()}
      </ButtonGroup>
    </div>
  );
};

export default PerformerResult;
