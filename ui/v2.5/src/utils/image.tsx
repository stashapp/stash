import React, { useCallback, useEffect, useState } from "react";
import Jimp from "jimp";

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
  const [isEncoding, setIsEncoding] = useState(false);

  const encodeImage = useCallback(
    (data: string) => {
      setIsEncoding(true);
      Jimp.read(data).then((image) =>
        image.quality(75).getBase64(Jimp.MIME_JPEG, (err, buffer) => {
          setIsEncoding(false);
          onLoadEnd(err ? "" : buffer);
        })
      );
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

  return isEncoding;
};

const Image = {
  onImageChange,
  usePasteImage,
};
export default Image;
