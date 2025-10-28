import React, { useEffect, useRef, useState } from "react";
import { Button, Form } from "react-bootstrap";
import { useIntl } from "react-intl";

import * as GQL from "src/core/generated-graphql";
import { ModalComponent } from "src/components/Shared/Modal";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { useScrapePerformerList } from "src/core/StashService";
import { useDebounce } from "src/hooks/debounce";

const CLASSNAME = "PerformerScrapeModal";
const CLASSNAME_LIST = `${CLASSNAME}-list`;

interface IProps {
  scraper: GQL.Scraper;
  onHide: () => void;
  onSelectPerformer: (
    performer: GQL.ScrapedPerformerDataFragment,
    scraper: GQL.Scraper
  ) => void;
  name?: string;
}
const PerformerScrapeModal: React.FC<IProps> = ({
  scraper,
  name,
  onHide,
  onSelectPerformer,
}) => {
  const intl = useIntl();
  const inputRef = useRef<HTMLInputElement>(null);
  const [query, setQuery] = useState<string>(name ?? "");
  const { data, loading } = useScrapePerformerList(scraper.id, query);

  const performers = data?.scrapeSinglePerformer ?? [];

  const onInputChange = useDebounce(setQuery, 500);

  useEffect(() => inputRef.current?.focus(), []);

  return (
    <ModalComponent
      show
      onHide={onHide}
      header={`Scrape performer from ${scraper.name}`}
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
        ) : (
          <ul className={CLASSNAME_LIST}>
            {performers.map((p, i) => (
              <li key={i}>
                <Button
                  variant="link"
                  onClick={() => onSelectPerformer(p, scraper)}
                >
                  {p.name}
                  {p.disambiguation && ` (${p.disambiguation})`}
                </Button>
              </li>
            ))}
          </ul>
        )}
      </div>
    </ModalComponent>
  );
};

export default PerformerScrapeModal;
