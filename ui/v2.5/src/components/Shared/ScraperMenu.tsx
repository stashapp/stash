import React from "react";
import { Dropdown, DropdownButton } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import { Icon } from "./Icon";
import { stashboxDisplayName } from "src/utils/stashbox";
import { ScraperSourceInput, StashBox } from "src/core/generated-graphql";
import { faSyncAlt } from "@fortawesome/free-solid-svg-icons";

export const ScraperMenu: React.FC<{
  stashBoxes: StashBox[];
  fragmentScrapers: { id: string; name: string }[];
  onScrapeClicked: (s: ScraperSourceInput) => void;
  onReloadScrapers: () => void;
}> = ({ stashBoxes, fragmentScrapers, onScrapeClicked, onReloadScrapers }) => {
  const intl = useIntl();

  return (
    <DropdownButton
      className="scraper-menu"
      title={intl.formatMessage({ id: "actions.scrape_with" })}
    >
      {stashBoxes.map((s, index) => (
        <Dropdown.Item
          key={s.endpoint}
          onClick={() =>
            onScrapeClicked({
              stash_box_endpoint: s.endpoint,
            })
          }
        >
          {stashboxDisplayName(s.name, index)}
        </Dropdown.Item>
      ))}
      {fragmentScrapers.map((s) => (
        <Dropdown.Item
          key={s.name}
          onClick={() => onScrapeClicked({ scraper_id: s.id })}
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
    </DropdownButton>
  );
};
