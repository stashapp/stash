import React, { PropsWithChildren, useMemo, useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { Button } from "react-bootstrap";
import {
  mutateReloadScrapers,
  useListGroupScrapers,
  useListPerformerScrapers,
  useListSceneScrapers,
  useListGalleryScrapers,
  useListImageScrapers,
  useScraperURLConflicts,
} from "src/core/StashService";
import * as GQL from "src/core/generated-graphql";
import { ConflictBadge, findConflict } from "./ScraperURLConflicts";
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
import { ClearableInput } from "../Shared/ClearableInput";
import { Counter } from "../Shared/Counter";

const ScraperTable: React.FC<
  PropsWithChildren<{
    entityType: string;
    count?: number;
  }>
> = ({ entityType, count, children }) => {
  const intl = useIntl();

  const titleEl = useMemo(() => {
    const title = intl.formatMessage(
      { id: "config.scraping.entity_scrapers" },
      { entityType: intl.formatMessage({ id: entityType }) }
    );

    if (count) {
      return (
        <span>
          {title} <Counter count={count} />
        </span>
      );
    }

    return title;
  }, [count, entityType, intl]);

  return (
    <CollapseButton text={titleEl}>
      <table className="scraper-table">
        <thead>
          <tr>
            <th>
              <FormattedMessage id="name" />
            </th>
            <th>
              <FormattedMessage id="config.scraping.supported_types" />
            </th>
            <th>
              <FormattedMessage id="config.scraping.supported_urls" />
            </th>
          </tr>
        </thead>
        <tbody>{children}</tbody>
      </table>
    </CollapseButton>
  );
};

const ScrapeTypeList: React.FC<{
  types: ScrapeType[];
  entityType: string;
}> = ({ types, entityType }) => {
  const intl = useIntl();

  const typeStrings = useMemo(
    () =>
      types.map((t) => {
        switch (t) {
          case ScrapeType.Fragment:
            return intl.formatMessage(
              { id: "config.scraping.entity_metadata" },
              { entityType: intl.formatMessage({ id: entityType }) }
            );
          default:
            return t;
        }
      }),
    [types, entityType, intl]
  );

  return (
    <ul>
      {typeStrings.map((t) => (
        <li key={t}>{t}</li>
      ))}
    </ul>
  );
};

type ScraperURLConflict =
  GQL.ScraperUrlConflictsQuery["scraperURLConflicts"][number];

interface IURLList {
  scraperId: string;
  urls: string[];
  conflicts: ScraperURLConflict[];
  scraperNames: Map<string, string>;
}

const URLList: React.FC<IURLList> = ({
  scraperId,
  urls,
  conflicts,
  scraperNames,
}) => {
  const items = useMemo(() => {
    function linkSite(url: string) {
      const u = new URL(url);
      return `${u.protocol}//${u.host}`;
    }

    const ret = urls
      .slice()
      .sort()
      .map((u) => {
        const sanitised = TextUtils.sanitiseURL(u);
        const siteURL = linkSite(sanitised!);
        const conflict = findConflict(conflicts, scraperId, u);

        return (
          <li key={u}>
            <ExternalLink href={siteURL}>{sanitised}</ExternalLink>
            {conflict && (
              <ConflictBadge
                conflict={conflict}
                mineName={scraperNames.get(scraperId) ?? scraperId}
                otherName={
                  scraperNames.get(conflict.otherID) ?? conflict.otherID
                }
              />
            )}
          </li>
        );
      });

    return ret;
  }, [urls, conflicts, scraperId, scraperNames]);

  return <ul>{items}</ul>;
};

const ScraperTableRow: React.FC<{
  id: string;
  name: string;
  entityType: string;
  supportedScrapes: ScrapeType[];
  urls: string[];
  conflicts: ScraperURLConflict[];
  scraperNames: Map<string, string>;
}> = ({
  id,
  name,
  entityType,
  supportedScrapes,
  urls,
  conflicts,
  scraperNames,
}) => {
  return (
    <tr>
      <td>{name}</td>
      <td>
        <ScrapeTypeList types={supportedScrapes} entityType={entityType} />
      </td>
      <td>
        <URLList
          scraperId={id}
          urls={urls}
          conflicts={conflicts}
          scraperNames={scraperNames}
        />
      </td>
    </tr>
  );
};

function filterScraper(filter: string) {
  return (name: string, urls: string[] | undefined | null) => {
    if (!filter) return true;

    return (
      name.toLowerCase().includes(filter) ||
      urls?.some((url) => url.toLowerCase().includes(filter))
    );
  };
}

const ScrapersSection: React.FC = () => {
  const Toast = useToast();
  const intl = useIntl();

  const [filter, setFilter] = useState("");

  const { data: performerScrapers, loading: loadingPerformers } =
    useListPerformerScrapers();
  const { data: sceneScrapers, loading: loadingScenes } =
    useListSceneScrapers();
  const { data: galleryScrapers, loading: loadingGalleries } =
    useListGalleryScrapers();
  const { data: imageScrapers, loading: loadingImages } =
    useListImageScrapers();
  const { data: groupScrapers, loading: loadingGroups } =
    useListGroupScrapers();
  const { data: conflictsData } = useScraperURLConflicts();

  const scraperNames = useMemo(() => {
    const m = new Map<string, string>();
    for (const list of [
      performerScrapers?.listScrapers,
      sceneScrapers?.listScrapers,
      galleryScrapers?.listScrapers,
      imageScrapers?.listScrapers,
      groupScrapers?.listScrapers,
    ]) {
      list?.forEach((s) => {
        m.set(s.id, s.name);
      });
    }
    return m;
  }, [
    performerScrapers,
    sceneScrapers,
    galleryScrapers,
    imageScrapers,
    groupScrapers,
  ]);

  const conflictsByType = useMemo(() => {
    const conflicts = conflictsData?.scraperURLConflicts ?? [];
    return {
      performers: conflicts.filter(
        (c) => c.type === GQL.ScrapeContentType.Performer
      ),
      scenes: conflicts.filter((c) => c.type === GQL.ScrapeContentType.Scene),
      galleries: conflicts.filter(
        (c) => c.type === GQL.ScrapeContentType.Gallery
      ),
      images: conflicts.filter((c) => c.type === GQL.ScrapeContentType.Image),
      groups: conflicts.filter((c) => c.type === GQL.ScrapeContentType.Group),
    };
  }, [conflictsData]);

  const filteredScrapers = useMemo(() => {
    const filterFn = filterScraper(filter.toLowerCase());
    return {
      performers: performerScrapers?.listScrapers.filter((s) =>
        filterFn(s.name, s.performer?.urls)
      ),
      scenes: sceneScrapers?.listScrapers.filter((s) =>
        filterFn(s.name, s.scene?.urls)
      ),
      galleries: galleryScrapers?.listScrapers.filter((s) =>
        filterFn(s.name, s.gallery?.urls)
      ),
      images: imageScrapers?.listScrapers.filter((s) =>
        filterFn(s.name, s.image?.urls)
      ),
      groups: groupScrapers?.listScrapers.filter((s) =>
        filterFn(s.name, s.group?.urls)
      ),
    };
  }, [
    performerScrapers,
    sceneScrapers,
    galleryScrapers,
    imageScrapers,
    groupScrapers,
    filter,
  ]);

  async function onReloadScrapers() {
    try {
      await mutateReloadScrapers();
    } catch (e) {
      Toast.error(e);
    }
  }

  if (
    loadingScenes ||
    loadingGalleries ||
    loadingPerformers ||
    loadingGroups ||
    loadingImages
  )
    return (
      <SettingSection headingID="config.scraping.scrapers">
        <LoadingIndicator />
      </SettingSection>
    );

  return (
    <SettingSection headingID="config.scraping.scrapers">
      <div className="content scraper-toolbar">
        <ClearableInput
          placeholder={`${intl.formatMessage({ id: "filter" })}...`}
          value={filter}
          setValue={(v) => setFilter(v)}
        />

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
        {!!filteredScrapers.scenes?.length && (
          <ScraperTable
            entityType="scene"
            count={filteredScrapers.scenes?.length}
          >
            {filteredScrapers.scenes?.map((scraper) => (
              <ScraperTableRow
                key={scraper.id}
                id={scraper.id}
                name={scraper.name}
                entityType="scene"
                supportedScrapes={scraper.scene?.supported_scrapes ?? []}
                urls={scraper.scene?.urls ?? []}
                conflicts={conflictsByType.scenes}
                scraperNames={scraperNames}
              />
            ))}
          </ScraperTable>
        )}

        {!!filteredScrapers.galleries?.length && (
          <ScraperTable
            entityType="gallery"
            count={filteredScrapers.galleries?.length}
          >
            {filteredScrapers.galleries?.map((scraper) => (
              <ScraperTableRow
                key={scraper.id}
                id={scraper.id}
                name={scraper.name}
                entityType="gallery"
                supportedScrapes={scraper.gallery?.supported_scrapes ?? []}
                urls={scraper.gallery?.urls ?? []}
                conflicts={conflictsByType.galleries}
                scraperNames={scraperNames}
              />
            ))}
          </ScraperTable>
        )}

        {!!filteredScrapers.images?.length && (
          <ScraperTable
            entityType="image"
            count={filteredScrapers.images?.length}
          >
            {filteredScrapers.images?.map((scraper) => (
              <ScraperTableRow
                key={scraper.id}
                id={scraper.id}
                name={scraper.name}
                entityType="image"
                supportedScrapes={scraper.image?.supported_scrapes ?? []}
                urls={scraper.image?.urls ?? []}
                conflicts={conflictsByType.images}
                scraperNames={scraperNames}
              />
            ))}
          </ScraperTable>
        )}

        {!!filteredScrapers.performers?.length && (
          <ScraperTable
            entityType="performer"
            count={filteredScrapers.performers?.length}
          >
            {filteredScrapers.performers?.map((scraper) => (
              <ScraperTableRow
                key={scraper.id}
                id={scraper.id}
                name={scraper.name}
                entityType="performer"
                supportedScrapes={scraper.performer?.supported_scrapes ?? []}
                urls={scraper.performer?.urls ?? []}
                conflicts={conflictsByType.performers}
                scraperNames={scraperNames}
              />
            ))}
          </ScraperTable>
        )}

        {!!filteredScrapers.groups?.length && (
          <ScraperTable
            entityType="group"
            count={filteredScrapers.groups?.length}
          >
            {filteredScrapers.groups?.map((scraper) => (
              <ScraperTableRow
                key={scraper.id}
                id={scraper.id}
                name={scraper.name}
                entityType="group"
                supportedScrapes={scraper.group?.supported_scrapes ?? []}
                urls={scraper.group?.urls ?? []}
                conflicts={conflictsByType.groups}
                scraperNames={scraperNames}
              />
            ))}
          </ScraperTable>
        )}
      </div>
    </SettingSection>
  );
};

export const SettingsScrapingPanel: React.FC = () => {
  const { general, scraping, loading, error, saveGeneral, saveScraping } =
    useSettings();

  if (error) return <h1>{error.message}</h1>;
  if (loading) return <LoadingIndicator />;

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

      <ScrapersSection />
    </>
  );
};
