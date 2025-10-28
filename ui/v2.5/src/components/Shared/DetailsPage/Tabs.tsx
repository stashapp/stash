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

  const setTabKey = useCallback(
    (newTabKey: string | null) => {
      if (!newTabKey) newTabKey = defaultTabKey;
      if (newTabKey === tabKey) return;

      if (validTabs.includes(newTabKey)) {
        history.replace(`${baseURL}/${newTabKey}`);
      }
    },
    [defaultTabKey, validTabs, tabKey, history, baseURL]
  );

  useEffect(() => {
    if (!tabKey) {
      setTabKey(defaultTabKey);
    }
  }, [setTabKey, defaultTabKey, tabKey]);

  return { setTabKey };
}
