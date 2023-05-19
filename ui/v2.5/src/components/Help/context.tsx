import React, { Suspense, useState } from "react";
import { Link } from "react-router-dom";
import { lazyComponent } from "src/utils/lazyComponent";

const Manual = lazyComponent(() => import("./Manual"));

interface IManualContextState {
  openManual: (tab?: string) => void;
}

export const ManualStateContext = React.createContext<IManualContextState>({
  openManual: () => {},
});

export const ManualProvider: React.FC = ({ children }) => {
  const [showManual, setShowManual] = useState(false);
  const [manualLink, setManualLink] = useState<string | undefined>();

  function openManual(tab?: string) {
    setManualLink(tab);
    setShowManual(true);
  }

  return (
    <ManualStateContext.Provider
      value={{
        openManual,
      }}
    >
      <Suspense fallback={<></>}>
        {showManual && (
          <Manual
            show={showManual}
            onClose={() => setShowManual(false)}
            defaultActiveTab={manualLink}
          />
        )}
      </Suspense>
      {children}
    </ManualStateContext.Provider>
  );
};

interface IManualLink {
  tab: string;
}

export const ManualLink: React.FC<IManualLink> = ({ tab, children }) => {
  const { openManual } = React.useContext(ManualStateContext);

  return (
    <Link
      to={`/help/${tab}.md`}
      onClick={(e) => {
        openManual(`${tab}.md`);
        e.preventDefault();
      }}
    >
      {children}
    </Link>
  );
};
