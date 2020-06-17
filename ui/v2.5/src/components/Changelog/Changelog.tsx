import React from "react";
import { useChangelogStorage } from "src/hooks";
import Version from "./Version";
import { V010, V011, V020, V021, V030 } from "./versions";

const Changelog: React.FC = () => {
  const [{ data, loading }, setOpenState] = useChangelogStorage();

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
        version="v0.3.0"
        openState={openState}
        setOpenState={setVersionOpenState}
        defaultOpen
      >
        <V030 />
      </Version>
      <Version
        version="v0.2.1"
        date="2020-06-10"
        openState={openState}
        setOpenState={setVersionOpenState}
      >
        <V021 />
      </Version>
      <Version
        version="v0.2.0"
        date="2020-06-06"
        openState={openState}
        setOpenState={setVersionOpenState}
      >
        <V020 />
      </Version>
      <Version
        version="v0.1.1"
        date="2020-02-25"
        openState={openState}
        setOpenState={setVersionOpenState}
      >
        <V011 />
      </Version>
      <Version
        version="v0.1.0"
        date="2020-02-24"
        openState={openState}
        setOpenState={setVersionOpenState}
      >
        <V010 />
      </Version>
    </>
  );
};

export default Changelog;
