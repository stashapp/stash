import React, { useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { Button } from "react-bootstrap";
import {
  mutateReloadScrapers,
  useListMovieScrapers,
  useListPerformerScrapers,
  useListSceneScrapers,
  useListGalleryScrapers,
} from "src/core/StashService";
import { useToast } from "src/hooks/Toast";
import TextUtils from "src/utils/text";
import { CollapseButton } from "../Shared/CollapseButton";
import { Icon } from "../Shared/Icon";
import { LoadingIndicator } from "../Shared/LoadingIndicator";
import { ScrapeType } from "src/core/generated-graphql";
import { SettingSection } from "./SettingSection";
import { BooleanSetting, StringListSetting, StringSetting } from "./Inputs";
import { useSettings } from "./context";
import { StashBoxSetting } from "./StashBoxConfiguration";
import { faSyncAlt } from "@fortawesome/free-solid-svg-icons";
import {
  AvailableScraperPackages,
  InstalledScraperPackages,
} from "./ScraperPackageManager";
import { ExternalLink } from "../Shared/ExternalLink";

interface IURLList {
  urls: string[];
}

const URLList: React.FC<IURLList> = ({ urls }) => {
  const maxCollapsedItems = 5;
  const [expanded, setExpanded] = useState<boolean>(false);

  function linkSite(url: string) {
    const u = new URL(url);
    return `${u.protocol}//${u.host}`;
  }

  function renderLink(url?: string) {
    if (url) {
      const sanitised = TextUtils.sanitiseURL(url);
      const siteURL = linkSite(sanitised!);

      return <ExternalLink href={siteURL}>{sanitised}</ExternalLink>;
    }
  }

  function getListItems() {
    const items = urls.map((u) => <li key={u}>{renderLink(u)}</li>);

    if (items.length > maxCollapsedItems) {
      if (!expanded) {
        items.length = maxCollapsedItems;
      }

      items.push(
        <li key="expand/collapse">
          <Button onClick={() => setExpanded(!expanded)} variant="link">
            {expanded ? "less" : "more"}
          </Button>
        </li>
      );
    }

    return items;
  }

  return <ul>{getListItems()}</ul>;
};

export const SettingsScrapingPanel: React.FC = () => {
  const Toast = useToast();
  const intl = useIntl();
  const { data: performerScrapers, loading: loadingPerformers } =
    useListPerformerScrapers();
  const { data: sceneScrapers, loading: loadingScenes } =
    useListSceneScrapers();
  const { data: galleryScrapers, loading: loadingGalleries } =
    useListGalleryScrapers();
  const { data: groupScrapers, loading: loadingGroups } =
    useListMovieScrapers();

  const { general, scraping, loading, error, saveGeneral, saveScraping } =
    useSettings();

  async function onReloadScrapers() {
    try {
      await mutateReloadScrapers();
    } catch (e) {
      Toast.error(e);
    }
  }

  function renderPerformerScrapeTypes(types: ScrapeType[]) {
    const typeStrings = types
      .filter((t) => t !== ScrapeType.Fragment)
      .map((t) => {
        switch (t) {
          case ScrapeType.Name:
            return intl.formatMessage({ id: "config.scraping.search_by_name" });
          default:
            return t;
        }
      });

    return (
      <ul>
        {typeStrings.map((t) => (
          <li key={t}>{t}</li>
        ))}
      </ul>
    );
  }

  function renderSceneScrapeTypes(types: ScrapeType[]) {
    const typeStrings = types.map((t) => {
      switch (t) {
        case ScrapeType.Fragment:
          return intl.formatMessage(
            { id: "config.scraping.entity_metadata" },
            { entityType: intl.formatMessage({ id: "scene" }) }
          );
        default:
          return t;
      }
    });

    return (
      <ul>
        {typeStrings.map((t) => (
          <li key={t}>{t}</li>
        ))}
      </ul>
    );
  }

  function renderGalleryScrapeTypes(types: ScrapeType[]) {
    const typeStrings = types.map((t) => {
      switch (t) {
        case ScrapeType.Fragment:
          return intl.formatMessage(
            { id: "config.scraping.entity_metadata" },
            { entityType: intl.formatMessage({ id: "gallery" }) }
          );
        default:
          return t;
      }
    });

    return (
      <ul>
        {typeStrings.map((t) => (
          <li key={t}>{t}</li>
        ))}
      </ul>
    );
  }

  function renderGroupScrapeTypes(types: ScrapeType[]) {
    const typeStrings = types.map((t) => {
      switch (t) {
        case ScrapeType.Fragment:
          return intl.formatMessage(
            { id: "config.scraping.entity_metadata" },
            { entityType: intl.formatMessage({ id: "group" }) }
          );
        default:
          return t;
      }
    });

    return (
      <ul>
        {typeStrings.map((t) => (
          <li key={t}>{t}</li>
        ))}
      </ul>
    );
  }

  function renderURLs(urls: string[]) {
    return <URLList urls={urls} />;
  }

  function renderSceneScrapers() {
    const elements = (sceneScrapers?.listScrapers ?? []).map((scraper) => (
      <tr key={scraper.id}>
        <td>{scraper.name}</td>
        <td>
          {renderSceneScrapeTypes(scraper.scene?.supported_scrapes ?? [])}
        </td>
        <td>{renderURLs(scraper.scene?.urls ?? [])}</td>
      </tr>
    ));

    return renderTable(
      intl.formatMessage(
        { id: "config.scraping.entity_scrapers" },
        { entityType: intl.formatMessage({ id: "scene" }) }
      ),
      elements
    );
  }

  function renderGalleryScrapers() {
    const elements = (galleryScrapers?.listScrapers ?? []).map((scraper) => (
      <tr key={scraper.id}>
        <td>{scraper.name}</td>
        <td>
          {renderGalleryScrapeTypes(scraper.gallery?.supported_scrapes ?? [])}
        </td>
        <td>{renderURLs(scraper.gallery?.urls ?? [])}</td>
      </tr>
    ));

    return renderTable(
      intl.formatMessage(
        { id: "config.scraping.entity_scrapers" },
        { entityType: intl.formatMessage({ id: "gallery" }) }
      ),
      elements
    );
  }

  function renderPerformerScrapers() {
    const elements = (performerScrapers?.listScrapers ?? []).map((scraper) => (
      <tr key={scraper.id}>
        <td>{scraper.name}</td>
        <td>
          {renderPerformerScrapeTypes(
            scraper.performer?.supported_scrapes ?? []
          )}
        </td>
        <td>{renderURLs(scraper.performer?.urls ?? [])}</td>
      </tr>
    ));

    return renderTable(
      intl.formatMessage(
        { id: "config.scraping.entity_scrapers" },
        { entityType: intl.formatMessage({ id: "performer" }) }
      ),
      elements
    );
  }

  function renderGroupScrapers() {
    const elements = (groupScrapers?.listScrapers ?? []).map((scraper) => (
      <tr key={scraper.id}>
        <td>{scraper.name}</td>
        <td>
          {renderGroupScrapeTypes(scraper.movie?.supported_scrapes ?? [])}
        </td>
        <td>{renderURLs(scraper.movie?.urls ?? [])}</td>
      </tr>
    ));

    return renderTable(
      intl.formatMessage(
        { id: "config.scraping.entity_scrapers" },
        { entityType: intl.formatMessage({ id: "group" }) }
      ),
      elements
    );
  }

  function renderTable(title: string, elements: JSX.Element[]) {
    if (elements.length > 0) {
      return (
        <CollapseButton text={title}>
          <table className="scraper-table">
            <thead>
              <tr>
                <th>{intl.formatMessage({ id: "name" })}</th>
                <th>
                  {intl.formatMessage({
                    id: "config.scraping.supported_types",
                  })}
                </th>
                <th>
                  {intl.formatMessage({ id: "config.scraping.supported_urls" })}
                </th>
              </tr>
            </thead>
            <tbody>{elements}</tbody>
          </table>
        </CollapseButton>
      );
    }
  }

  if (error) return <h1>{error.message}</h1>;
  if (
    loading ||
    loadingScenes ||
    loadingGalleries ||
    loadingPerformers ||
    loadingGroups
  )
    return <LoadingIndicator />;

  return (
    <>
      <StashBoxSetting
        value={general.stashBoxes ?? []}
        onChange={(v) => saveGeneral({ stashBoxes: v })}
      />

      <SettingSection headingID="config.general.scraping">
        <StringSetting
          id="scraperUserAgent"
          headingID="config.general.scraper_user_agent"
          subHeadingID="config.general.scraper_user_agent_desc"
          value={scraping.scraperUserAgent ?? undefined}
          onChange={(v) => saveScraping({ scraperUserAgent: v })}
        />

        <StringSetting
          id="scraperCDPPath"
          headingID="config.general.chrome_cdp_path"
          subHeadingID="config.general.chrome_cdp_path_desc"
          value={scraping.scraperCDPPath ?? undefined}
          onChange={(v) => saveScraping({ scraperCDPPath: v })}
        />

        <BooleanSetting
          id="scraper-cert-check"
          headingID="config.general.check_for_insecure_certificates"
          subHeadingID="config.general.check_for_insecure_certificates_desc"
          checked={scraping.scraperCertCheck ?? undefined}
          onChange={(v) => saveScraping({ scraperCertCheck: v })}
        />

        <StringListSetting
          id="excluded-tag-patterns"
          headingID="config.scraping.excluded_tag_patterns_head"
          subHeadingID="config.scraping.excluded_tag_patterns_desc"
          value={scraping.excludeTagPatterns ?? undefined}
          onChange={(v) => saveScraping({ excludeTagPatterns: v })}
        />
      </SettingSection>

      <InstalledScraperPackages />
      <AvailableScraperPackages />

      <SettingSection headingID="config.scraping.scrapers">
        <div className="content">
          <Button onClick={() => onReloadScrapers()}>
            <span className="fa-icon">
              <Icon icon={faSyncAlt} />
            </span>
            <span>
              <FormattedMessage id="actions.reload_scrapers" />
            </span>
          </Button>
        </div>

        <div className="content">
          {renderSceneScrapers()}
          {renderGalleryScrapers()}
          {renderPerformerScrapers()}
          {renderGroupScrapers()}
        </div>
      </SettingSection>
    </>
  );
};
