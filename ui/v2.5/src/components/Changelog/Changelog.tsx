import React from "react";
import { useChangelogStorage } from "src/hooks";
import Version from "./Version";
import V010 from "./versions/v010.md";
import V011 from "./versions/v011.md";
import V020 from "./versions/v020.md";
import V021 from "./versions/v021.md";
import V030 from "./versions/v030.md";
import V040 from "./versions/v040.md";
import { MarkdownPage } from "../Shared/MarkdownPage";

const Changelog: React.FC = () => {
  const [{ data, loading }, setOpenState] = useChangelogStorage();

  const stashVersion = process.env.REACT_APP_STASH_VERSION;
  const buildTime = process.env.REACT_APP_DATE;

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

  return (
    <>
      <h1 className="mb-4">Changelog:</h1>
      <Version
        version={stashVersion || "v0.4.0"}
        date={buildDate}
        openState={openState}
        setOpenState={setVersionOpenState}
        defaultOpen
      >
        <MarkdownPage page={V040} />
      </Version>
      <Version
        version="v0.3.0"
        date="2020-09-02"
        openState={openState}
        setOpenState={setVersionOpenState}
      >
        <MarkdownPage page={V030} />
      </Version>
      <Version
        version="v0.2.1"
        date="2020-06-10"
        openState={openState}
        setOpenState={setVersionOpenState}
      >
        <MarkdownPage page={V021} />
      </Version>
      <Version
        version="v0.2.0"
        date="2020-06-06"
        openState={openState}
        setOpenState={setVersionOpenState}
      >
        <MarkdownPage page={V020} />
      </Version>
      <Version
        version="v0.1.1"
        date="2020-02-25"
        openState={openState}
        setOpenState={setVersionOpenState}
      >
        <MarkdownPage page={V011} />
      </Version>
      <Version
        version="v0.1.0"
        date="2020-02-24"
        openState={openState}
        setOpenState={setVersionOpenState}
      >
        <MarkdownPage page={V010} />
      </Version>
    </>
  );
};

export default Changelog;
