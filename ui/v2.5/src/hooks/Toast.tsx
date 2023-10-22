import React, { useEffect, useState, useContext, createContext } from "react";
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
    error: (error: unknown) => {
      /* eslint-disable no-console */
      let message: string;
      if (error instanceof Error) {
        message = error.message ?? error.toString();
      } else if (error instanceof String) {
        message = error.toString();
      } else {
        console.error(error);
        toastFunc({
          variant: "danger",
          header: "Error",
          content: "Unknown error",
        });
        return;
      }

      console.error(message);
      toastFunc({
        variant: "danger",
        header: "Error",
        content: message,
      });
      /* eslint-enable no-console */
    },
  };
}

export const useToast = () => {
  const setToast = useContext(ToastContext);
  const [hookObject, setHookObject] = useState(createHookObject(setToast));
  useEffect(() => setHookObject(createHookObject(setToast)), [setToast]);

  return hookObject;
};
