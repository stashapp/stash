import React, { useEffect, useRef, useState } from "react";
import { useIntl } from "react-intl";
import { Button, Form } from "react-bootstrap";

import * as GQL from "src/core/generated-graphql";

import { useDebouncedSetState } from "src/hooks/debounce";
import { stashboxDisplayName } from "src/utils/stashbox";

import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { ModalComponent } from "src/components/Shared/Modal";

const CLASSNAME = "StudioScrapeModal";
const CLASSNAME_LIST = `${CLASSNAME}-list`;

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
        stash_box_index: instance.index,
      },
      input: {
        query,
      },
    },
    skip: query === "",
  });

  const studios = data?.scrapeSingleStudio ?? [];

  const onInputChange = useDebouncedSetState(setQuery, 500);

  useEffect(() => inputRef.current?.focus(), []);

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
          placeholder="Studio name..."
          className="text-input mb-4"
          ref={inputRef}
        />
        {loading ? (
          <div className="m-4 text-center">
            <LoadingIndicator inline />
          </div>
        ) : studios.length > 0 ? (
          <ul className={CLASSNAME_LIST}>
            {studios.map((s) => (
              <li key={s.remote_site_id}>
                <Button variant="link" onClick={() => onSelectStudio(s)}>
                  {s.images && s.images?.length > 0 ? (
                    <img
                      width={64}
                      src={s.images[0]}
                      alt={s.name}
                      className="rounded-circle"
                    />
                  ): null}&nbsp;&nbsp;&nbsp;&nbsp;
                  {s.name}
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

export default StudioStashBoxModal;
