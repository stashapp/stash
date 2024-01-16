import React, {
  useState,
  useContext,
  createContext,
  useCallback,
  useMemo,
} from "react";
import { Toast } from "react-bootstrap";
import { errorToString } from "src/utils";

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
    setToasts((prev) => prev.filter((item) => item.id !== id));

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

  const addToast = useCallback((toast: IToast) => {
    setToasts((prev) => [...prev, { ...toast, id: toastID++ }]);
  }, []);

  return (
    <ToastContext.Provider value={addToast}>
      {children}
      <div className="toast-container row">{toastItems}</div>
    </ToastContext.Provider>
  );
};

export const useToast = () => {
  const addToast = useContext(ToastContext);

  return useMemo(
    () => ({
      toast: addToast,
      success(message: React.ReactNode | string) {
        addToast({
          content: message,
        });
      },
      error(error: unknown) {
        const message = errorToString(error);

        console.error(error);
        addToast({
          variant: "danger",
          header: "Error",
          content: message,
        });
      },
    }),
    [addToast]
  );
};
