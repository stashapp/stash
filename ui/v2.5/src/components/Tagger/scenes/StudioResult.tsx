import React, { useMemo } from "react";
import { Button, ButtonGroup } from "react-bootstrap";
import { FormattedMessage } from "react-intl";
import cx from "classnames";

import { StudioSelect, SelectObject } from "src/components/Shared/Select";
import * as GQL from "src/core/generated-graphql";

import { OptionalField } from "../IncludeButton";
import { getStashboxBase } from "src/utils/stashbox";
import { ExternalLink } from "src/components/Shared/ExternalLink";
import { Link } from "react-router-dom";
import { LinkButton } from "../LinkButton";

const StudioLink: React.FC<{
  studio: GQL.ScrapedStudio | GQL.SlimStudioDataFragment;
  url: string | undefined;
  internal?: boolean;
}> = ({ studio, url, internal = false }) => {
  const name = useMemo(() => {
    if (!url) return studio.name;

    return internal ? (
      <Link to={url} target="_blank">
        {studio.name}
      </Link>
    ) : (
      <ExternalLink href={url}>{studio.name}</ExternalLink>
    );
  }, [url, studio.name, internal]);

  return <span>{name}</span>;
};

interface IStudioResultProps {
  studio: GQL.ScrapedStudio;
  selectedID: string | undefined;
  setSelectedID: (id: string | undefined) => void;
  onCreate: () => void;
  onLink?: () => Promise<void>;
  endpoint?: string;
}

const StudioResult: React.FC<IStudioResultProps> = ({
  studio,
  selectedID,
  setSelectedID,
  onCreate,
  onLink,
  endpoint,
}) => {
  const { data: studioData, loading: stashLoading } = GQL.useFindStudioQuery({
    variables: { id: studio.stored_id ?? "" },
    skip: !studio.stored_id,
  });

  const matchedStudio = studioData?.findStudio;
  const matchedStashID = matchedStudio?.stash_ids.some(
    (stashID) => stashID.endpoint === endpoint && stashID.stash_id
  );

  const stashboxStudioPrefix = endpoint
    ? `${getStashboxBase(endpoint)}studios/`
    : undefined;
  const studioURLPrefix = "/studios/";

  const handleSelect = (studios: SelectObject[]) => {
    if (studios.length) {
      setSelectedID(studios[0].id);
    } else {
      setSelectedID(undefined);
    }
  };

  const handleSkip = () => {
    setSelectedID(undefined);
  };

  if (stashLoading) return <div>Loading studio</div>;

  if (matchedStudio && matchedStashID) {
    return (
      <div className="row no-gutters my-2">
        <div className="entity-name">
          <FormattedMessage id="countables.studios" values={{ count: 1 }} />:
          <b className="ml-2">
            <StudioLink
              studio={studio}
              url={`${stashboxStudioPrefix}${studio.remote_site_id}`}
            />
          </b>
        </div>
        <span className="ml-auto">
          <OptionalField
            exclude={selectedID === undefined}
            setExclude={(v) =>
              v ? handleSkip() : setSelectedID(matchedStudio.id)
            }
          >
            <div>
              <span className="mr-2">
                <FormattedMessage id="component_tagger.verb_matched" />:
              </span>
              <b className="col-3 text-right">
                <StudioLink
                  studio={matchedStudio}
                  url={`${studioURLPrefix}${matchedStudio.id}`}
                  internal
                />
              </b>
            </div>
          </OptionalField>
        </span>
      </div>
    );
  }

  const selectedSource = !selectedID ? "skip" : "existing";

  const safeBuildStudioScraperLink = (id: string | null | undefined) => {
    return stashboxStudioPrefix && id
      ? `${stashboxStudioPrefix}${id}`
      : undefined;
  };

  return (
    <div className="row no-gutters align-items-center mt-2">
      <div className="entity-name">
        <FormattedMessage id="countables.studios" values={{ count: 1 }} />:
        <b className="ml-2">
          <StudioLink
            studio={studio}
            url={safeBuildStudioScraperLink(studio.remote_site_id)}
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
        <StudioSelect
          ids={selectedID ? [selectedID] : []}
          onSelect={handleSelect}
          className={cx("studio-select", {
            "studio-select-active": selectedSource === "existing",
          })}
          isClearable={false}
        />
        {endpoint && onLink && (
          <LinkButton disabled={selectedID === undefined} onLink={onLink} />
        )}
      </ButtonGroup>
    </div>
  );
};

export default StudioResult;
