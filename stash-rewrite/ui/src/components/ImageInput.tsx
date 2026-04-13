import { useRef, useState, useCallback } from "react";
import { useMutation } from "@tanstack/react-query";
import { Upload, X, Image as ImageIcon, Link, Clipboard } from "lucide-react";

interface ImageInputProps {
  currentImageUrl?: string;
  onUpload: (file: File) => Promise<unknown>;
  onDelete?: () => Promise<unknown>;
  onSuccess?: () => void;
  label?: string;
  className?: string;
  aspectRatio?: string; // e.g. "2/3", "16/9", "1/1"
}

export function ImageInput({
  currentImageUrl,
  onUpload,
  onDelete,
  onSuccess,
  label = "Image",
  className = "",
  aspectRatio = "2/3",
}: ImageInputProps) {
  const fileRef = useRef<HTMLInputElement>(null);
  const [preview, setPreview] = useState<string | null>(null);
  const [showUrlInput, setShowUrlInput] = useState(false);
  const [urlValue, setUrlValue] = useState("");
  const [imgError, setImgError] = useState(false);

  const uploadMut = useMutation({
    mutationFn: onUpload,
    onSuccess: () => {
      setPreview(null);
      onSuccess?.();
    },
  });

  const deleteMut = useMutation({
    mutationFn: onDelete ?? (() => Promise.resolve()),
    onSuccess: () => {
      setPreview(null);
      setImgError(false);
      onSuccess?.();
    },
  });

  const handleFile = useCallback(
    (file: File) => {
      if (!file.type.startsWith("image/")) return;
      const url = URL.createObjectURL(file);
      setPreview(url);
      setImgError(false);
      uploadMut.mutate(file);
    },
    [uploadMut],
  );

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) handleFile(file);
    e.target.value = "";
  };

  const handleDrop = useCallback(
    (e: React.DragEvent) => {
      e.preventDefault();
      const file = e.dataTransfer.files[0];
      if (file) handleFile(file);
    },
    [handleFile],
  );

  const handlePaste = useCallback(async () => {
    try {
      const items = await navigator.clipboard.read();
      for (const item of items) {
        const imageType = item.types.find((t) => t.startsWith("image/"));
        if (imageType) {
          const blob = await item.getType(imageType);
          const file = new File([blob], "clipboard.png", { type: imageType });
          handleFile(file);
          return;
        }
      }
    } catch {
      // Clipboard API not available or no image data
    }
  }, [handleFile]);

  const handleUrlSubmit = useCallback(async () => {
    if (!urlValue.trim()) return;
    try {
      const res = await fetch(urlValue.trim());
      const blob = await res.blob();
      const ext = blob.type.split("/")[1] || "jpg";
      const file = new File([blob], `url-image.${ext}`, { type: blob.type });
      handleFile(file);
      setShowUrlInput(false);
      setUrlValue("");
    } catch {
      // Fetch failed
    }
  }, [urlValue, handleFile]);

  const displayUrl = preview || (currentImageUrl && !imgError ? currentImageUrl : null);
  const isLoading = uploadMut.isPending || deleteMut.isPending;

  return (
    <div className={`space-y-2 ${className}`}>
      <label className="block text-sm font-medium text-gray-400">{label}</label>

      <div
        className="relative rounded-lg overflow-hidden bg-gray-800 border border-gray-700 border-dashed hover:border-gray-600 transition-colors cursor-pointer"
        style={{ aspectRatio }}
        onDrop={handleDrop}
        onDragOver={(e) => e.preventDefault()}
        onClick={() => !isLoading && fileRef.current?.click()}
      >
        {displayUrl ? (
          <>
            <img
              src={displayUrl}
              alt={label}
              className="w-full h-full object-cover"
              onError={() => setImgError(true)}
            />
            {isLoading && (
              <div className="absolute inset-0 bg-black/50 flex items-center justify-center">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500" />
              </div>
            )}
          </>
        ) : (
          <div className="flex flex-col items-center justify-center w-full h-full text-gray-500 p-4">
            <ImageIcon className="w-8 h-8 mb-2" />
            <span className="text-xs text-center">Click or drag to upload</span>
          </div>
        )}
      </div>

      <input
        ref={fileRef}
        type="file"
        accept="image/jpeg,image/png,image/webp,image/gif"
        onChange={handleFileChange}
        className="hidden"
      />

      {/* Action buttons */}
      <div className="flex gap-1.5 flex-wrap">
        <button
          type="button"
          onClick={(e) => {
            e.stopPropagation();
            fileRef.current?.click();
          }}
          disabled={isLoading}
          className="flex items-center gap-1 px-2 py-1 text-xs bg-gray-700 hover:bg-gray-600 rounded transition-colors disabled:opacity-50"
          title="Upload from file"
        >
          <Upload className="w-3 h-3" /> File
        </button>
        <button
          type="button"
          onClick={(e) => {
            e.stopPropagation();
            setShowUrlInput(!showUrlInput);
          }}
          disabled={isLoading}
          className="flex items-center gap-1 px-2 py-1 text-xs bg-gray-700 hover:bg-gray-600 rounded transition-colors disabled:opacity-50"
          title="Upload from URL"
        >
          <Link className="w-3 h-3" /> URL
        </button>
        <button
          type="button"
          onClick={(e) => {
            e.stopPropagation();
            handlePaste();
          }}
          disabled={isLoading}
          className="flex items-center gap-1 px-2 py-1 text-xs bg-gray-700 hover:bg-gray-600 rounded transition-colors disabled:opacity-50"
          title="Paste from clipboard"
        >
          <Clipboard className="w-3 h-3" /> Paste
        </button>
        {onDelete && displayUrl && (
          <button
            type="button"
            onClick={(e) => {
              e.stopPropagation();
              deleteMut.mutate();
            }}
            disabled={isLoading}
            className="flex items-center gap-1 px-2 py-1 text-xs bg-red-900/50 hover:bg-red-800/50 text-red-400 rounded transition-colors disabled:opacity-50 ml-auto"
            title="Remove image"
          >
            <X className="w-3 h-3" /> Remove
          </button>
        )}
      </div>

      {/* URL input */}
      {showUrlInput && (
        <div className="flex gap-2">
          <input
            type="url"
            value={urlValue}
            onChange={(e) => setUrlValue(e.target.value)}
            placeholder="https://..."
            className="flex-1 bg-gray-800 border border-gray-700 rounded px-2 py-1 text-xs text-gray-200 focus:outline-none focus:border-blue-500"
            onKeyDown={(e) => e.key === "Enter" && handleUrlSubmit()}
          />
          <button
            type="button"
            onClick={handleUrlSubmit}
            className="px-2 py-1 text-xs bg-blue-600 hover:bg-blue-500 rounded transition-colors"
          >
            Load
          </button>
        </div>
      )}

      {uploadMut.error && (
        <p className="text-xs text-red-400">{(uploadMut.error as Error).message}</p>
      )}
    </div>
  );
}
