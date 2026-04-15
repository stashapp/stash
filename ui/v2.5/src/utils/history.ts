import { useHistory } from "react-router-dom";

type History = ReturnType<typeof useHistory>;

export function goBackOrReplace(history: History, defaultPath: string) {
  if (history.length > 1) {
    history.goBack();
  } else {
    history.replace(defaultPath);
  }
}
