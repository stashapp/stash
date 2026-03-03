import React, { useCallback, useEffect } from "react";

const blobToDataURL = (blob: Blob): Promise<string> =>
  new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onloadend = () => resolve(reader.result as string);
    reader.onerror = reject;
    reader.readAsDataURL(blob);
  });

const readImage = (file: File, onLoadEnd: (imageData: string) => void) => {
  // only proceed if no error encountered
  blobToDataURL(file).then(onLoadEnd).catch(() => {});
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
  return blobToDataURL(blob);
};

const readClipboardImage = async (): Promise<string | null> => {
  const items = await navigator.clipboard.read();
  for (const item of items) {
    const imageType = item.types.find((t) => t.startsWith("image/"));
    if (imageType) {
      const blob = await item.getType(imageType);
      return blobToDataURL(blob);
    }
  }
  return null;
};

const ImageUtils = {
  onImageChange,
  usePasteImage,
  imageToDataURL,
  readClipboardImage,
};

export default ImageUtils;
