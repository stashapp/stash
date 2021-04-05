import React, { useState } from "react";
import { Button, Card, Form, InputGroup } from "react-bootstrap";
import { Link } from "react-router-dom";
import { HashLink } from "react-router-hash-link";
import { useLocalForage } from "src/hooks";

import * as GQL from "src/core/generated-graphql";
import { LoadingIndicator } from "src/components/Shared";
import {
  stashBoxPerformerQuery,
  useConfiguration,
} from "src/core/StashService";
import { Manual } from "src/components/Help/Manual";

import StashSearchResult from "./StashSearchResult";
import PerformerConfig from "./Config";
import {
  LOCAL_FORAGE_KEY,
  ITaggerConfig,
  initialConfig,
} from "../constants";
import {
  IStashBoxPerformer,
  selectPerformers,
} from "../utils";

const CLASSNAME = 'PerformerTagger';

interface IPerformerTaggerListProps {
  performers: GQL.PerformerDataFragment[];
  selectedEndpoint: { endpoint: string; index: number };
  config: ITaggerConfig;
}

const PerformerTaggerList: React.FC<IPerformerTaggerListProps> = ({
  performers,
  selectedEndpoint,
}) => {
  const [loading, setLoading] = useState(false);
  const [searchResults, setSearchResults] = useState<
    Record<string, IStashBoxPerformer[]>
  >({});
  const [searchErrors, setSearchErrors] = useState<
    Record<string, string | undefined>
  >({});
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const [taggedPerformers, setTaggedPerformers] = useState<
    Record<string, Partial<GQL.SlimPerformerDataFragment>>
  >({});
  const [queries, setQueries] = useState<Record<string, string>>({});

  const doBoxSearch = (performerID: string, searchVal: string) => {
    stashBoxPerformerQuery(searchVal, selectedEndpoint.index)
      .then((queryData) => {
        const s = selectPerformers(queryData.data?.queryStashBoxPerformer?.[0].results ?? []);
        setSearchResults({
          ...searchResults,
          [performerID]: s,
        });
        setSearchErrors({
          ...searchErrors,
          [performerID]: undefined,
        });
        setLoading(false);
      })
      .catch(() => {
        setLoading(false);
        // Destructure to remove existing result
        const { [performerID]: unassign, ...results } = searchResults;
        setSearchResults(results);
        setSearchErrors({
          ...searchErrors,
          [performerID]: "Network Error",
        });
      });

    setLoading(true);
  };

  const handleTaggedPerformer = (performer: Pick<GQL.SlimPerformerDataFragment, 'id'> & Partial<Omit<GQL.SlimPerformerDataFragment, 'id'>>) => {
    setTaggedPerformers({
      ...taggedPerformers,
      [performer.id]: performer,
    });
  };

  const renderPerformers = () =>
    performers.map((performer) => {
      const isTagged = taggedPerformers[performer.id];
      const hasStashIDs = performer.stash_ids.length > 0;

      let mainContent;
      if (!isTagged && hasStashIDs) {
        mainContent = (
          <div className="text-left">
            <h5 className="text-bold">Performer already tagged</h5>
          </div>
        );
      } else if (!isTagged && !hasStashIDs) {
        mainContent = (
          <InputGroup>
            <InputGroup.Prepend>
              <InputGroup.Text>Query</InputGroup.Text>
            </InputGroup.Prepend>
            <Form.Control
              className="text-input"
              defaultValue={performer.name ?? ''}
              onChange={e => setQueries({ ...queries, [performer.id]: e.currentTarget.value })}
              onKeyPress={(e: React.KeyboardEvent<HTMLInputElement>) =>
                e.key === "Enter" && doBoxSearch(
                  performer.id,
                  queries[performer.id] ?? performer.name ?? '',
                )
              }
            />
            <InputGroup.Append>
              <Button
                disabled={loading}
                onClick={() =>
                  doBoxSearch(
                    performer.id,
                    queries[performer.id] ?? performer.name ?? '',
                  )
                }
              >
                Search
              </Button>
            </InputGroup.Append>
          </InputGroup>
        );
      } else if (isTagged) {
        mainContent = (
          <div className="d-flex flex-column text-left">
            <h5>Performer successfully tagged:</h5>
            <h6>
              <Link className="bold" to={`/performers/${performer.id}`}>
                {taggedPerformers[performer.id].name}
              </Link>
            </h6>
          </div>
        );
      }

      let subContent;
      if (performer.stash_ids.length > 0) {
        const stashLinks = performer.stash_ids.map((stashID) => {
          const base = stashID.endpoint.match(/https?:\/\/.*?\//)?.[0];
          const link = base ? (
            <a
              className="small d-block"
              href={`${base}performers/${stashID.stash_id}`}
              target="_blank"
              rel="noopener noreferrer"
            >
              {stashID.stash_id}
            </a>
          ) : (
            <div className="small">{stashID.stash_id}</div>
          );

          return link;
        });
        subContent = <>{stashLinks}</>;
      } else if (searchErrors[performer.id]) {
        subContent = (
          <div className="text-danger font-weight-bold">
            {searchErrors[performer.id]}
          </div>
        );
      } else if (searchResults[performer.id]?.length === 0) {
        subContent = (
          <div className="text-danger font-weight-bold">No results found.</div>
        );
      }

      let searchResult;
      if (searchResults[performer.id]?.length > 0 && !isTagged) {
        searchResult = (
          <StashSearchResult
            key={performer.id}
            stashboxPerformers={searchResults[performer.id]}
            performer={performer}
            endpoint={selectedEndpoint.endpoint}
            onPerformerTagged={handleTaggedPerformer}
          />
        );
      }

      return (
        <div key={performer.id} className={`${CLASSNAME}-performer`}>
          <Card className="performer-card p-0 m-0">
            <img src={performer.image_path ?? ''} alt="" />
          </Card>
          <div className="flex-grow-1 ml-3">
            <Link to={`/performers/${performer.id}`} className={`${CLASSNAME}-header`}>
              <h2>{performer.name}</h2>
            </Link>
            {mainContent}
            <div className="sub-content text-left">{subContent}</div>
            {searchResult}
          </div>
        </div>
      );
    });

  return (
    <Card className={CLASSNAME}>
      {renderPerformers()}
    </Card>
  );
};

interface ITaggerProps {
  performers: GQL.PerformerDataFragment[];
}

export const PerformerTagger: React.FC<ITaggerProps> = ({ performers }) => {
  const stashConfig = useConfiguration();
  const [{ data: config }, setConfig] = useLocalForage<ITaggerConfig>(
    LOCAL_FORAGE_KEY,
    initialConfig
  );
  const [showConfig, setShowConfig] = useState(false);
  const [showManual, setShowManual] = useState(false);

  if (!config) return <LoadingIndicator />;

  const savedEndpointIndex =
    stashConfig.data?.configuration.general.stashBoxes.findIndex(
      (s) => s.endpoint === config.selectedEndpoint
    ) ?? -1;
  const selectedEndpointIndex =
    savedEndpointIndex === -1 &&
    stashConfig.data?.configuration.general.stashBoxes.length
      ? 0
      : savedEndpointIndex;
  const selectedEndpoint =
    stashConfig.data?.configuration.general.stashBoxes[selectedEndpointIndex];

  return (
    <>
      <Manual
        show={showManual}
        onClose={() => setShowManual(false)}
        defaultActiveTab="Tagger.md"
      />
      <div className="tagger-container mx-md-auto">
        {selectedEndpointIndex !== -1 && selectedEndpoint ? (
          <>
            <div className="row mb-2 no-gutters">
              <Button onClick={() => setShowConfig(!showConfig)} variant="link">
                {showConfig ? "Hide" : "Show"} Configuration
              </Button>
              <Button
                className="ml-auto"
                onClick={() => setShowManual(true)}
                title="Help"
                variant="link"
              >
                Help
              </Button>
            </div>

            <PerformerConfig config={config} setConfig={setConfig} show={showConfig} />
            <PerformerTaggerList
              performers={performers}
              config={config}
              selectedEndpoint={{
                endpoint: selectedEndpoint.endpoint,
                index: selectedEndpointIndex,
              }}
            />
          </>
        ) : (
          <div className="my-4">
            <h3 className="text-center mt-4">
              To use the performer tagger a stash-box instance needs to be
              configured.
            </h3>
            <h5 className="text-center">
              Please see{" "}
              <HashLink
                to="/settings?tab=configuration#stashbox"
                scroll={(el) =>
                  el.scrollIntoView({ behavior: "smooth", block: "center" })
                }
              >
                Settings.
              </HashLink>
            </h5>
          </div>
        )}
      </div>
    </>
  );
};
