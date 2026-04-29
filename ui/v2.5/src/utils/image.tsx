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
  blobToDataURL(file)
    .then(onLoadEnd)
    .catch(() => {});
};

const onImageChange = (
  event: React.FormEvent<HTMLInputElement>,
  onLoadEnd: (imageData: string) => void
) => {
  const file = event?.currentTarget?.files?.[0];
  if (file) readImage(file, onLoadEnd);
};

const imageToDataURL = async (url: string) => {
  const response = await fetch(url);
  const blob = await response.blob();
  return blobToDataURL(blob);
};

// uses event.clipboardData which works in all contexts including insecure HTTP
const pasteImage = (
  event: ClipboardEvent,
  onLoadEnd: (imageData: string) => void
) => {
  const files = event?.clipboardData?.files;
  if (!files?.length) return;

  if (document.activeElement instanceof HTMLInputElement) {
    // don't interfere with pasting text into inputs
    return;
  }

  const file = Array.from(files).find((f) => f.type.startsWith("image/"));
  if (file) readImage(file, onLoadEnd);
};

// uses Clipboard API which requires secure context (HTTPS or localhost)
const readClipboardImage = async (): Promise<string | null> => {
  if (!window.isSecureContext) {
    return null;
  }

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

const ImageUtils = {
  onImageChange,
  usePasteImage,
  imageToDataURL,
  readClipboardImage,
};

export default ImageUtils;
