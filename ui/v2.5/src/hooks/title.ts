import { MessageDescriptor, useIntl } from "react-intl";

export const TITLE = "Stash";
export const TITLE_SEPARATOR = " | ";

export function useTitleProps(...messages: (string | MessageDescriptor)[]) {
  const intl = useIntl();

  const parts = messages.map((msg) => {
    if (typeof msg === "object") {
      return intl.formatMessage(msg);
    } else {
      return msg;
    }
  });

  return makeTitleProps(...parts);
}

export function makeTitleProps(...parts: string[]) {
  const title = [...parts, TITLE].join(TITLE_SEPARATOR);
  return {
    titleTemplate: `%s | ${title}`,
    defaultTitle: title,
  };
}
