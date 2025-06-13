import { FormattedMessage } from "react-intl";
import { Counter } from "../Counter";
import { useCallback, useEffect } from "react";
import { useHistory, useLocation } from "react-router-dom";
import { PatchComponent } from "src/patch";

export const TabTitleCounter: React.FC<{
  messageID: string;
  count: number;
  abbreviateCounter: boolean;
}> = PatchComponent(
  "TabTitleCounter",
  ({ messageID, count, abbreviateCounter }) => {
    return (
      <>
        <FormattedMessage id={messageID} />
        <Counter count={count} abbreviateCounter={abbreviateCounter} hideZero />
      </>
    );
  }
);

export function useTabKey(props: {
  tabKey: string | undefined;
  validTabs: readonly string[];
  defaultTabKey: string;
  baseURL: string;
}) {
  const { tabKey, validTabs, defaultTabKey, baseURL } = props;

  const history = useHistory();
  const location = useLocation();

  const setTabKey = useCallback(
    (newTabKey: string | null) => {
      if (!newTabKey) newTabKey = defaultTabKey;
      if (newTabKey === tabKey) return;

      if (validTabs.includes(newTabKey)) {
        const params = new URLSearchParams(location.search);
        const returnTo = params.get("returnTo");
        const newSearch = returnTo
          ? `?returnTo=${encodeURIComponent(returnTo)}`
          : "";

        history.replace({
          pathname: `${baseURL}/${newTabKey}`,
          search: newSearch,
        });
      }
    },
    [defaultTabKey, validTabs, tabKey, history, baseURL, location.search]
  );

  useEffect(() => {
    if (!tabKey) {
      setTabKey(defaultTabKey);
    }
  }, [setTabKey, defaultTabKey, tabKey]);

  return { setTabKey };
}
