import React, { useEffect, useRef, useState } from "react";
import { Button, Form } from "react-bootstrap";
import { useIntl } from "react-intl";

import * as GQL from "src/core/generated-graphql";
import { ModalComponent } from "src/components/Shared/Modal";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { stashboxDisplayName } from "src/utils/stashbox";
import { useDebouncedSetState } from "src/hooks/debounce";

const CLASSNAME = "PerformerScrapeModal";
const CLASSNAME_LIST = `${CLASSNAME}-list`;

export interface IStashBox extends GQL.StashBox {
  index: number;
}

interface IProps {
  instance: IStashBox;
  onHide: () => void;
  onSelectPerformer: (performer: GQL.ScrapedPerformer) => void;
  name?: string;
}
const PerformerStashBoxModal: React.FC<IProps> = ({
  instance,
  name,
  onHide,
  onSelectPerformer,
}) => {
  const intl = useIntl();
  const inputRef = useRef<HTMLInputElement>(null);
  const [query, setQuery] = useState<string>(name ?? "");
  const { data, loading } = GQL.useScrapeSinglePerformerQuery({
    variables: {
      source: {
        stash_box_index: instance.index,
      },
      input: {
        query,
      },
    },
    skip: query === "",
  });

  const performers = data?.scrapeSinglePerformer ?? [];

  const onInputChange = useDebouncedSetState(setQuery, 500);

  useEffect(() => inputRef.current?.focus(), []);

  return (
    <ModalComponent
      show
      onHide={onHide}
      header={`Scrape performer from ${stashboxDisplayName(
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
          placeholder="Performer name..."
          className="text-input mb-4"
          ref={inputRef}
        />
        {loading ? (
          <div className="m-4 text-center">
            <LoadingIndicator inline />
          </div>
        ) : performers.length > 0 ? (
          <ul className={CLASSNAME_LIST}>
            {performers.map((p) => (
              <li key={p.remote_site_id}>
                <Button variant="link" onClick={() => onSelectPerformer(p)}>
                  {p.name}
                  {p.disambiguation && ` (${p.disambiguation})`}
                </Button>
              </li>
            ))}
          </ul>
        ) : (
          query !== "" && <h5 className="text-center">No results found.</h5>
        )}
      </div>
    </ModalComponent>
  );
};

export default PerformerStashBoxModal;
