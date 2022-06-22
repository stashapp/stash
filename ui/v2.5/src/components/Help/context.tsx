import React, {
  lazy,
  PropsWithChildren,
  Suspense,
  useEffect,
  useState,
} from "react";

const Manual = lazy(() => import("./Manual"));

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

  useEffect(() => {
    if (manualLink) setManualLink(undefined);
  }, [manualLink]);

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

export const ManualLink: React.FC<PropsWithChildren<IManualLink>> = ({
  tab,
  children,
}) => {
  const { openManual } = React.useContext(ManualStateContext);

  return (
    <a
      href={`/help/${tab}.md`}
      onClick={(e) => {
        openManual(`${tab}.md`);
        e.preventDefault();
      }}
    >
      {children}
    </a>
  );
};
