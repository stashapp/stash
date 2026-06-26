import { FormattedMessage } from "react-intl";
import { Counter } from "../Counter";
import { useCallback, useEffect } from "react";
import { useHistory } from "react-router-dom";
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
  const activeTabKey =
    tabKey && tabKey !== "default" && validTabs.includes(tabKey)
      ? tabKey
      : defaultTabKey;

  const setTabKey = useCallback(
    (newTabKey: string | null) => {
      if (
        !newTabKey ||
        newTabKey === "default" ||
        !validTabs.includes(newTabKey)
      ) {
        newTabKey = defaultTabKey;
      }
      if (newTabKey === activeTabKey) return;

      history.replace(`${baseURL}/${newTabKey}`);
    },
    [activeTabKey, defaultTabKey, validTabs, history, baseURL]
  );

  useEffect(() => {
    if (tabKey !== activeTabKey) {
      history.replace(`${baseURL}/${activeTabKey}`);
    }
  }, [activeTabKey, baseURL, history, tabKey]);

  return { activeTabKey, setTabKey };
}
