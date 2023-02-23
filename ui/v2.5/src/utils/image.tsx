import React, { useCallback, useEffect } from "react";

const readImage = (file: File, onLoadEnd: (imageData: string) => void) => {
  const reader: FileReader = new FileReader();
  reader.onloadend = () => {
    // only proceed if no error encountered
    if (!reader.error) {
      onLoadEnd(reader.result as string);
    }
  };
  reader.readAsDataURL(file);
};

const pasteImage = (
  event: ClipboardEvent,
  onLoadEnd: (imageData: string) => void
) => {
  const files = event?.clipboardData?.files;
  if (!files?.length) return;

  const file = files[0];
  readImage(file, onLoadEnd);
};

const onImageChange = (
  event: React.FormEvent<HTMLInputElement>,
  onLoadEnd: (imageData: string) => void
) => {
  const file = event?.currentTarget?.files?.[0];
  if (file) readImage(file, onLoadEnd);
};

const usePasteImage = (
  onLoadEnd: (imageData: string) => void,
  isActive: boolean = true
) => {
  const encodeImage = useCallback(
    (data: string) => {
      onLoadEnd(data);
    },
    [onLoadEnd]
  );

  useEffect(() => {
    const paste = (event: ClipboardEvent) => pasteImage(event, encodeImage);
    if (isActive) {
      document.addEventListener("paste", paste);
    }

    return () => document.removeEventListener("paste", paste);
  }, [isActive, encodeImage]);

  return false;
};

const imageToDataURL = async (url: string) => {
  const response = await fetch(url);
  const blob = await response.blob();
  return new Promise<string>((resolve, reject) => {
    const reader = new FileReader();
    reader.onloadend = () => {
      resolve(reader.result as string);
    };
    reader.onerror = reject;
    reader.readAsDataURL(blob);
  });
};

const ImageUtils = {
  onImageChange,
  usePasteImage,
  imageToDataURL,
};

export default ImageUtils;
