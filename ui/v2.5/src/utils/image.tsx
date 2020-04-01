import React, { useEffect } from "react";

const readImage = (file: File, onLoadEnd: (this: FileReader) => void) => {
  const reader: FileReader = new FileReader();
  reader.onloadend = onLoadEnd;
  reader.readAsDataURL(file);
};

const pasteImage = (
  event: ClipboardEvent,
  onLoadEnd: (this: FileReader) => void
) => {
  const files = event?.clipboardData?.files;
  if (!files?.length) return;

  const file = files[0];
  readImage(file, onLoadEnd);
};

const onImageChange = (
  event: React.FormEvent<HTMLInputElement>,
  onLoadEnd: (this: FileReader) => void
) => {
  const file = event?.currentTarget?.files?.[0];
  if (file) readImage(file, onLoadEnd);
};

const usePasteImage = (onLoadEnd: (this: FileReader) => void) => {
  useEffect(() => {
    const paste = (event: ClipboardEvent) => pasteImage(event, onLoadEnd);
    document.addEventListener("paste", paste);

    return () => document.removeEventListener("paste", paste);
  });
};

const Image = {
  onImageChange,
  usePasteImage,
};
export default Image;
