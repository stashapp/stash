import { MessageDescriptor, useIntl } from "react-intl";
import { useConfigurationContext } from "./Config";

export const TITLE = "Stash";
export const TITLE_SEPARATOR = " | ";

export function useTitleProps(...messages: (string | MessageDescriptor)[]) {
  const intl = useIntl();
  const config = useConfigurationContext();
  const title = config.configuration.ui.title || TITLE;

  const parts = messages.map((msg) => {
    if (typeof msg === "object") {
      return intl.formatMessage(msg);
    } else {
      return msg;
    }
  });

  return makeTitleProps(title, ...parts);
}

export function makeTitleProps(title: string, ...parts: string[]) {
  const fullTitle = [...parts, title].join(TITLE_SEPARATOR);
  return {
    titleTemplate: `%s | ${fullTitle}`,
    defaultTitle: fullTitle,
  };
}
