import React, { useMemo, useState } from "react";
import { Dropdown, Button } from "react-bootstrap";
import { useIntl } from "react-intl";
import { Icon } from "./Icon";
import { stashboxDisplayName } from "src/utils/stashbox";
import { ScraperSourceInput, StashBox } from "src/core/generated-graphql";
import { faSyncAlt } from "@fortawesome/free-solid-svg-icons";
import { ClearableInput } from "./ClearableInput";
import useFocus from "src/utils/focus";
import ScreenUtils from "src/utils/screen";

export const ScraperMenu: React.FC<{
  toggle: React.ReactNode;
  variant?: string;
  stashBoxes?: StashBox[];
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
  const [filter, setFilter] = useState("");

  const focusOnOpen = !ScreenUtils.isTouch();
  const focusRef = useFocus();
  const [, setFocus] = focusRef;

  const filteredStashboxes = useMemo(() => {
    if (!stashBoxes) return [];
    if (!filter) return stashBoxes;

    return stashBoxes.filter((s) =>
      s.name.toLowerCase().includes(filter.toLowerCase())
    );
  }, [stashBoxes, filter]);

  const filteredScrapers = useMemo(() => {
    if (!filter) return scrapers;

    return scrapers.filter(
      (s) =>
        s.name.toLowerCase().includes(filter.toLowerCase()) ||
        s.id.toLowerCase().includes(filter.toLowerCase())
    );
  }, [scrapers, filter]);

  return (
    <Dropdown
      className="scraper-menu"
      title={intl.formatMessage({ id: "actions.scrape_query" })}
      onToggle={(v) => {
        if (focusOnOpen && v) setTimeout(() => setFocus(true), 0);
      }}
    >
      <Dropdown.Toggle variant={variant}>{toggle}</Dropdown.Toggle>

      <Dropdown.Menu>
        <div className="scraper-filter-container">
          <ClearableInput
            placeholder={`${intl.formatMessage({ id: "filter" })}...`}
            value={filter}
            setValue={setFilter}
            focus={focusRef}
          />
          <Button
            onClick={onReloadScrapers}
            className="reload-button"
            title={intl.formatMessage({ id: "actions.reload_scrapers" })}
          >
            <Icon icon={faSyncAlt} />
          </Button>
        </div>

        {filteredStashboxes.map((s, index) => (
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

        {filteredStashboxes.length > 0 && filteredScrapers.length > 0 && (
          <Dropdown.Divider />
        )}

        {filteredScrapers.map((s) => (
          <Dropdown.Item
            key={s.name}
            onClick={() => onScraperClicked({ scraper_id: s.id })}
          >
            {s.name}
          </Dropdown.Item>
        ))}
      </Dropdown.Menu>
    </Dropdown>
  );
};
