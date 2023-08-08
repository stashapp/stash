import React, { useEffect, useState, useContext, createContext } from "react";
import { Link } from "react-router-dom";
import { Toast } from "react-bootstrap";

interface IToast {
  header?: string;
  content: React.ReactNode | string;
  delay?: number;
  variant?: "success" | "danger" | "warning";
}
interface IActiveToast extends IToast {
  id: number;
}

let toastID = 0;
const ToastContext = createContext<(item: IToast) => void>(() => {});

export const ToastProvider: React.FC = ({ children }) => {
  const [toasts, setToasts] = useState<IActiveToast[]>([]);

  const removeToast = (id: number) =>
    setToasts(toasts.filter((item) => item.id !== id));

  const toastItems = toasts.map((toast) => (
    <Toast
      autohide
      key={toast.id}
      onClose={() => removeToast(toast.id)}
      className={toast.variant ?? "success"}
      delay={toast.delay ?? 3000}
    >
      <Toast.Header>
        <span className="mr-auto">{toast.header}</span>
      </Toast.Header>
      <Toast.Body>{toast.content}</Toast.Body>
    </Toast>
  ));

  const addToast = (toast: IToast) =>
    setToasts([...toasts, { ...toast, id: toastID++ }]);

  return (
    <ToastContext.Provider value={addToast}>
      {children}
      <div className="toast-container row">{toastItems}</div>
    </ToastContext.Provider>
  );
};

function createHookObject(toastFunc: (toast: IToast) => void) {
  return {
    success: toastFunc,
    error: (error: unknown, delay: number = 10000) => {
      /* eslint-disable @typescript-eslint/no-explicit-any, no-console */
      let message: string;
      if (error instanceof Error) {
        message = error.message ?? error.toString();
      } else if ((error as any).toString) {
        message = (error as any).toString();
      } else {
        console.error(error);
        message = "Unknown error, check the logs for more information";
      }

      console.error(message);
      let content = maybeAddLinkToErrMessage(message);
      toastFunc({
        variant: "danger",
        header: "Error",
        content: content,
        delay: delay,
      });
      /* eslint-enable @typescript-eslint/no-explicit-any, no-console */
    },
  };

  function maybeAddLinkToErrMessage(error: string): React.ReactNode | string {
    // This function is a hack to enable the text to be a direct link in the Toast notification
    if (error.includes("check the logs")) {
      const msgParts = error.split("check the logs");
      return (
        <>
          {msgParts[0]}
          <Link to={`/settings?tab=logs`}>check the logs</Link>
          {msgParts.length > 1 ? msgParts[1] : ""}
        </>
      );
    }
    return error;
  }
}

export const useToast = () => {
  const setToast = useContext(ToastContext);
  const [hookObject, setHookObject] = useState(createHookObject(setToast));
  useEffect(() => setHookObject(createHookObject(setToast)), [setToast]);

  return hookObject;
};
