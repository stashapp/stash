import React from "react";
import { useChangelogStorage } from "src/hooks/LocalForage";
import Version from "./Version";
import V010 from "src/docs/en/Changelog/v010.md";
import V011 from "src/docs/en/Changelog/v011.md";
import V020 from "src/docs/en/Changelog/v020.md";
import V021 from "src/docs/en/Changelog/v021.md";
import V030 from "src/docs/en/Changelog/v030.md";
import V040 from "src/docs/en/Changelog/v040.md";
import V050 from "src/docs/en/Changelog/v050.md";
import V060 from "src/docs/en/Changelog/v060.md";
import V070 from "src/docs/en/Changelog/v070.md";
import V080 from "src/docs/en/Changelog/v080.md";
import V090 from "src/docs/en/Changelog/v090.md";
import V0100 from "src/docs/en/Changelog/v0100.md";
import V0110 from "src/docs/en/Changelog/v0110.md";
import V0120 from "src/docs/en/Changelog/v0120.md";
import V0130 from "src/docs/en/Changelog/v0130.md";
import V0131 from "src/docs/en/Changelog/v0131.md";
import V0140 from "src/docs/en/Changelog/v0140.md";
import V0150 from "src/docs/en/Changelog/v0150.md";
import V0160 from "src/docs/en/Changelog/v0160.md";
import V0161 from "src/docs/en/Changelog/v0161.md";
import V0170 from "src/docs/en/Changelog/v0170.md";
import V0180 from "src/docs/en/Changelog/v0180.md";
import V0190 from "src/docs/en/Changelog/v0190.md";
import V0200 from "src/docs/en/Changelog/v0200.md";
import V0210 from "src/docs/en/Changelog/v0210.md";
import V0220 from "src/docs/en/Changelog/v0220.md";
import V0230 from "src/docs/en/Changelog/v0230.md";
import V0240 from "src/docs/en/Changelog/v0240.md";
import V0250 from "src/docs/en/Changelog/v0250.md";
import V0260 from "src/docs/en/Changelog/v0260.md";
import V0270 from "src/docs/en/Changelog/v0270.md";
import V0280 from "src/docs/en/Changelog/v0280.md";
import V0290 from "src/docs/en/Changelog/v0290.md";
import V0300 from "src/docs/en/Changelog/v0300.md";

import V0290ReleaseNotes from "src/docs/en/ReleaseNotes/v0290.md";

import { MarkdownPage } from "../Shared/MarkdownPage";
import { FormattedMessage } from "react-intl";

const Changelog: React.FC = () => {
  const [{ data, loading }, setOpenState] = useChangelogStorage();

  const stashVersion = import.meta.env.VITE_APP_STASH_VERSION;
  const buildTime = import.meta.env.VITE_APP_DATE;

  let buildDate;
  if (buildTime) {
    buildDate = buildTime.substring(0, buildTime.indexOf(" "));
  }

  if (loading) return <></>;

  const openState = data?.versions ?? {};

  const setVersionOpenState = (key: string, state: boolean) =>
    setOpenState({
      versions: {
        ...openState,
        [key]: state,
      },
    });

  interface IStashRelease {
    version: string;
    date?: string;
    page: string;
    defaultOpen?: boolean;
    releaseNotes?: string;
  }

  // after new release:
  // add entry to releases, using the current* fields
  // then update the current fields.
  const currentVersion = stashVersion || "v0.30.0";
  const currentDate = buildDate;
  const currentPage = V0300;

  const releases: IStashRelease[] = [
    {
      version: currentVersion,
      date: currentDate,
      page: currentPage,
      defaultOpen: true,
    },
    {
      version: "v0.29.3",
      date: "2025-11-06",
      page: V0290,
      releaseNotes: V0290ReleaseNotes,
    },
    {
      version: "v0.28.1",
      date: "2025-03-20",
      page: V0280,
    },
    {
      version: "v0.27.2",
      date: "2024-10-16",
      page: V0270,
    },
    {
      version: "v0.26.2",
      date: "2024-06-27",
      page: V0260,
    },
    {
      version: "v0.25.1",
      date: "2024-03-13",
      page: V0250,
    },
    {
      version: "v0.24.3",
      date: "2024-01-15",
      page: V0240,
    },
    {
      version: "v0.23.1",
      date: "2023-10-14",
      page: V0230,
    },
    {
      version: "v0.22.1",
      date: "2023-08-21",
      page: V0220,
    },
    {
      version: "v0.21.0",
      date: "2023-06-13",
      page: V0210,
    },
    {
      version: "v0.20.2",
      date: "2023-04-08",
      page: V0200,
    },
    {
      version: "v0.19.1",
      date: "2023-02-21",
      page: V0190,
    },
    {
      version: "v0.18.0",
      date: "2022-11-30",
      page: V0180,
    },
    {
      version: "v0.17.2",
      date: "2022-10-25",
      page: V0170,
    },
    {
      version: "v0.16.1",
      date: "2022-07-26",
      page: V0161,
    },
    {
      version: "v0.16.0",
      date: "2022-07-05",
      page: V0160,
    },
    {
      version: "v0.15.0",
      date: "2022-05-18",
      page: V0150,
    },
    {
      version: "v0.14.0",
      date: "2022-04-11",
      page: V0140,
    },
    {
      version: "v0.13.1",
      date: "2022-03-16",
      page: V0131,
    },
    {
      version: "v0.13.0",
      date: "2022-03-08",
      page: V0130,
    },
    {
      version: "v0.12.0",
      date: "2021-12-29",
      page: V0120,
    },
    {
      version: "v0.11.0",
      date: "2021-11-16",
      page: V0110,
    },
    {
      version: "v0.10.0",
      date: "2021-10-11",
      page: V0100,
    },
    {
      version: "v0.9.0",
      date: "2021-09-06",
      page: V090,
    },
    {
      version: "v0.8.0",
      date: "2021-07-02",
      page: V080,
    },
    {
      version: "v0.7.0",
      date: "2021-05-15",
      page: V070,
    },
    {
      version: "v0.6.0",
      date: "2021-03-29",
      page: V060,
    },
    {
      version: "v0.5.0",
      date: "2021-02-23",
      page: V050,
    },
    {
      version: "v0.4.0",
      date: "2020-11-24",
      page: V040,
    },
    {
      version: "v0.3.0",
      date: "2020-09-02",
      page: V030,
    },
    {
      version: "v0.2.1",
      date: "2020-06-10",
      page: V021,
    },
    {
      version: "v0.2.0",
      date: "2020-06-06",
      page: V020,
    },
    {
      version: "v0.1.1",
      date: "2020-02-25",
      page: V011,
    },
    {
      version: "v0.1.0",
      date: "2020-02-24",
      page: V010,
    },
  ];

  return (
    <div className="changelog">
      <h1 className="mb-4">Changelog:</h1>
      {releases.map((r) => (
        <Version
          key={r.version}
          version={r.version}
          date={r.date}
          openState={openState}
          setOpenState={setVersionOpenState}
          defaultOpen={r.defaultOpen}
        >
          {r.releaseNotes && (
            <div>
              <h3 className="mt-0">
                <FormattedMessage id="release_notes" />
              </h3>
              <MarkdownPage page={r.releaseNotes} />
              <hr />
            </div>
          )}
          <MarkdownPage page={r.page} />
        </Version>
      ))}
    </div>
  );
};

export default Changelog;
