import React, {
  useState,
  useContext,
  createContext,
  useMemo,
} from "react";
import { Toast } from "react-bootstrap";
import { errorToString } from "src/utils";

export interface IToast {
  content: React.ReactNode | string;
  delay?: number;
  variant?: "success" | "danger" | "warning";
  priority?: number; // higher is more important
}

interface IActiveToast extends IToast {
  id: number;
}

// errors are always more important than regular toasts
const errorPriority = 100;
// errors should stay on screen longer
const errorDelay = 5000;

let toastID = 0;

const ToastContext = createContext<(item: IToast) => void>(() => {});

export const ToastProvider: React.FC = ({ children }) => {
  const [toast, setToast] = useState<IActiveToast>();

  const toastItem = useMemo(() => {
    if (!toast) return null;

    return (
      <Toast
        autohide
        key={toast.id}
        onClose={() => setToast(undefined)}
        className={toast.variant ?? "success"}
        delay={toast.delay ?? 3000}
      >
        <Toast.Header>
          <span className="mr-auto">{toast.content}</span>
        </Toast.Header>
      </Toast>
    );
  }, [toast]);

  function addToast(item: IToast) {
    if (!toast || (item.priority ?? 0) >= (toast.priority ?? 0)) {
      setToast({ ...item, id: toastID++ });
    }
  }

  return (
    <ToastContext.Provider value={addToast}>
      {children}
      <div className="toast-container row">{toastItem}</div>
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
          content: message,
          priority: errorPriority,
          delay: errorDelay,
        });
      },
    }),
    [addToast]
  );
};
