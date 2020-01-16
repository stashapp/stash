import React, { useEffect } from "react";

export class ImageUtils {

  private static readImage(file: File, onLoadEnd: (this: FileReader) => any) {
    const reader: FileReader = new FileReader();
    
    reader.onloadend = onLoadEnd;
    reader.readAsDataURL(file);
  }

  public static onImageChange(event: React.FormEvent<HTMLInputElement>, onLoadEnd: (this: FileReader) => any) {
    const file: File = (event.target as any).files[0];
    ImageUtils.readImage(file, onLoadEnd);
  }
    
  public static pasteImage(e : any, onLoadEnd: (this: FileReader) => any) {
    if (e.clipboardData.files.length === 0) {
      return;
    }
    
    const file: File = e.clipboardData.files[0];
    ImageUtils.readImage(file, onLoadEnd);
  }

  public static addPasteImageHook(onLoadEnd: (this: FileReader) => any) {
    useEffect(() => {
      const pasteImage = (e: any) => { ImageUtils.pasteImage(e, onLoadEnd) }
      window.addEventListener("paste", pasteImage);
    
      return () => window.removeEventListener("paste", pasteImage);
    });
  }
}