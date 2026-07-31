import React, {
  Suspense,
  useCallback,
  useEffect,
  useRef,
  useState,
} from "react";
import { lazyComponent } from "src/utils/lazyComponent";
import { ILightboxImage, IChapter } from "./types";

const LightboxComponent = lazyComponent(() => import("./Lightbox"));
const LIGHTBOX_HISTORY_KEY = "stashLightbox";

export type LightboxHideReason = "dismiss" | "navigate";

export interface IState {
  images: ILightboxImage[];
  isVisible: boolean;
  isLoading: boolean;
  showNavigation: boolean;
  initialIndex?: number;
  pageCallback?: (props: { direction?: number; page?: number }) => void;
  chapters?: IChapter[];
  page?: number;
  pages?: number;
  pageSize?: number;
  totalCount?: number;
  slideshowEnabled: boolean;
  slideshowAutostart?: boolean;
  onClose?: () => void;
}
interface IContext {
  lightboxState: IState;
  setLightboxState: (state: Partial<IState>) => void;
}

interface ILightboxHistoryState {
  [LIGHTBOX_HISTORY_KEY]?: {
    id?: number;
  };
}

const getLightboxHistoryID = (state: unknown) => {
  if (!state || typeof state !== "object") return;

  const lightboxState = (state as ILightboxHistoryState)[LIGHTBOX_HISTORY_KEY];
  return typeof lightboxState?.id === "number" ? lightboxState.id : undefined;
};

export const LightboxContext = React.createContext<IContext | null>(null);

export function useLightboxContext() {
  const context = React.useContext(LightboxContext);
  if (!context) {
    throw new Error(
      "useLightboxContext must be used within a LightboxProvider"
    );
  }
  return context;
}

export const LightboxProvider: React.FC = ({ children }) => {
  const [lightboxState, setLightboxState] = useState<IState>({
    images: [],
    isVisible: false,
    isLoading: false,
    showNavigation: true,
    slideshowEnabled: false,
  });

  const activeHistoryID = useRef<number>();
  const nextHistoryID = useRef(0);
  const isDismissingRef = useRef(false);
  const isVisibleRef = useRef(lightboxState.isVisible);
  const onCloseRef = useRef<(() => void) | undefined>();

  isVisibleRef.current = lightboxState.isVisible;
  onCloseRef.current = lightboxState.onClose;

  const isCurrentLightboxHistoryEntry = useCallback(() => {
    return (
      activeHistoryID.current !== undefined &&
      getLightboxHistoryID(history.state) === activeHistoryID.current
    );
  }, []);

  const pushLightboxHistory = useCallback(() => {
    const id = nextHistoryID.current + 1;
    nextHistoryID.current = id;
    activeHistoryID.current = id;
    isDismissingRef.current = false;

    history.pushState(
      {
        ...(typeof history.state === "object" && history.state !== null
          ? history.state
          : {}),
        [LIGHTBOX_HISTORY_KEY]: { id },
      },
      "",
      window.location.href
    );
  }, []);

  const clearCurrentLightboxHistory = useCallback(() => {
    if (!isCurrentLightboxHistoryEntry()) return;

    const currentState =
      typeof history.state === "object" && history.state !== null
        ? { ...history.state }
        : {};
    delete (currentState as ILightboxHistoryState)[LIGHTBOX_HISTORY_KEY];
    history.replaceState(currentState, "", window.location.href);
    activeHistoryID.current = undefined;
    isDismissingRef.current = false;
  }, [isCurrentLightboxHistoryEntry]);

  const closeLightbox = useCallback(() => {
    if (!isVisibleRef.current) return;

    isDismissingRef.current = false;
    isVisibleRef.current = false;
    setLightboxState((currentState: IState) =>
      currentState.isVisible
        ? {
            ...currentState,
            isVisible: false,
            slideshowAutostart: false,
          }
        : currentState
    );
    onCloseRef.current?.();
  }, []);

  const setPartialState = useCallback(
    (state: Partial<IState>) => {
      if (state.isVisible === true && !isVisibleRef.current) {
        pushLightboxHistory();
        isVisibleRef.current = true;
      } else if (state.isVisible === false) {
        isDismissingRef.current = false;
        isVisibleRef.current = false;
      }

      setLightboxState((currentState: IState) => ({
        ...currentState,
        ...state,
      }));
    },
    [pushLightboxHistory]
  );

  const onHide = useCallback(
    (reason: LightboxHideReason = "dismiss") => {
      if (reason === "navigate") {
        clearCurrentLightboxHistory();
        closeLightbox();
        return;
      }

      if (isCurrentLightboxHistoryEntry()) {
        if (isDismissingRef.current) return;

        isDismissingRef.current = true;
        history.back();
        return;
      }

      closeLightbox();
    },
    [clearCurrentLightboxHistory, closeLightbox, isCurrentLightboxHistoryEntry]
  );

  useEffect(() => {
    const handlePopState = (event: PopStateEvent) => {
      const historyID = getLightboxHistoryID(event.state);
      if (
        activeHistoryID.current !== undefined &&
        historyID === activeHistoryID.current
      ) {
        if (!isVisibleRef.current) {
          isDismissingRef.current = false;
          isVisibleRef.current = true;
          setLightboxState((currentState: IState) => ({
            ...currentState,
            isVisible: true,
          }));
        }
        return;
      }

      closeLightbox();
    };

    window.addEventListener("popstate", handlePopState);
    return () => {
      window.removeEventListener("popstate", handlePopState);
    };
  }, [closeLightbox]);

  const onDeleteImage = useCallback((id: string) => {
    setLightboxState((s) => ({
      ...s,
      images: s.images.filter((img) => img.id !== id),
    }));
  }, []);

  return (
    <LightboxContext.Provider
      value={{ lightboxState, setLightboxState: setPartialState }}
    >
      {children}
      <Suspense fallback={null}>
        {lightboxState.isVisible && (
          <LightboxComponent
            {...lightboxState}
            hide={onHide}
            onDeleteImage={onDeleteImage}
          />
        )}
      </Suspense>
    </LightboxContext.Provider>
  );
};
