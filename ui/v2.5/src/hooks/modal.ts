import React from "react";

export function useModal() {
  const [modal, setModal] = React.useState<React.ReactNode>();

  const closeModal = () => setModal(undefined);
  const showModal = (m: React.ReactNode) => setModal(m);

  return { modal, closeModal, showModal };
}
