import React, { useCallback } from "react";

export function useModal() {
  const [modal, setModal] = React.useState<React.ReactNode>();

  const closeModal = useCallback(() => setModal(undefined), []);
  const showModal = useCallback((m: React.ReactNode) => setModal(m), []);

  return { modal, closeModal, showModal };
}

export interface IModalContextState {
  modal: React.ReactNode;
  closeModal: () => void;
  showModal: (m: React.ReactNode) => void;
}

export const ModalStateContext = React.createContext<IModalContextState | null>(
  null
);

export const useModalContext = () => {
  const context = React.useContext(ModalStateContext);

  if (context === null) {
    throw new Error("useModalContext must be used within a ModalContext");
  }

  return context;
};

// ModalContext is a provider that allows you to show a modal anywhere in the app
// by calling showModal(modal) and close it by calling closeModal()
export const ModalContext: React.FC = ({ children }) => {
  const modalState = useModal();

  return (
    <ModalStateContext.Provider value={modalState}>
      {modalState.modal}
      {children}
    </ModalStateContext.Provider>
  );
};
