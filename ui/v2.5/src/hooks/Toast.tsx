import React, { useState, useContext, createContext } from 'react';
import { Toast } from 'react-bootstrap';

interface IToast {
  header?: string;
  content: JSX.Element|string;
  delay?: number;
  variant?: 'success'|'danger'|'warning'|'info';
}
interface IActiveToast extends IToast {
  id: number;
}

let toastID = 0;
const ToastContext = createContext<(item:IToast) => void>(() => {});

export const ToastProvider: React.FC  = ({children}) => {
    const [toasts, setToasts] = useState<IActiveToast[]>([]);

  const removeToast = (id:number) => (
    setToasts(toasts.filter(item => item.id !== id))
  );

  const toastItems = toasts.map(toast => (
    <Toast
      autohide
      key={toast.id}
      onClose={() => removeToast(toast.id)}
      className={toast.variant ?? 'success'}
      delay={toast.delay ?? 5000}
    >
      <Toast.Header>
        <span className="mr-auto">
          { toast.header ?? 'Stash' }
        </span>
      </Toast.Header>
      <Toast.Body>{toast.content}</Toast.Body>
    </Toast>
  ));

  const addToast = (toast:IToast) => (
    setToasts([...toasts, { ...toast, id: toastID++ }])
  );

  return (
    <ToastContext.Provider value={addToast}>
      {children}
      <div className="toast-container row">
        { toastItems }
      </div>
    </ToastContext.Provider>
  )
}

const useToasts = () => {
  const setToast = useContext(ToastContext);
  return {
    success: setToast,
    error: (error: Error) => {
      console.error(error.message);
      setToast({
        variant: 'danger',
        header: 'Error',
        content: error.message ?? error.toString()
      });
    }
  };
}

export default useToasts;
