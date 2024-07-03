import React from "react";
import { Dropdown } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import { Icon } from "./Icon";
import { stashboxDisplayName } from "src/utils/stashbox";
import { ScraperSourceInput, StashBox } from "src/core/generated-graphql";
import { faSyncAlt } from "@fortawesome/free-solid-svg-icons";

export const ScraperMenu: React.FC<{
  toggle: React.ReactNode;
  variant?: string;
  stashBoxes: StashBox[];
  scrapers: { id: string; name: string }[];
  onScraperClicked: (s: ScraperSourceInput) => void;
  onReloadScrapers: () => void;
}> = ({
  toggle,
  variant,
  stashBoxes,
  scrapers,
  onScraperClicked,
  onReloadScrapers,
}) => {
  const intl = useIntl();

  return (
    <Dropdown title={intl.formatMessage({ id: "actions.scrape_query" })}>
      <Dropdown.Toggle variant={variant}>{toggle}</Dropdown.Toggle>

      <Dropdown.Menu>
        {stashBoxes.map((s, index) => (
          <Dropdown.Item
            key={s.endpoint}
            onClick={() =>
              onScraperClicked({
                stash_box_endpoint: s.endpoint,
              })
            }
          >
            {stashboxDisplayName(s.name, index)}
          </Dropdown.Item>
        ))}
        {scrapers.map((s) => (
          <Dropdown.Item
            key={s.name}
            onClick={() => onScraperClicked({ scraper_id: s.id })}
          >
            {s.name}
          </Dropdown.Item>
        ))}
        <Dropdown.Item onClick={() => onReloadScrapers()}>
          <span className="fa-icon">
            <Icon icon={faSyncAlt} />
          </span>
          <span>
            <FormattedMessage id="actions.reload_scrapers" />
          </span>
        </Dropdown.Item>
      </Dropdown.Menu>
    </Dropdown>
  );
};
