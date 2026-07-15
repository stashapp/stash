import React, { Suspense, useCallback, useState } from "react";
import { lazyComponent } from "src/utils/lazyComponent";
import { ILightboxImage, IChapter } from "./types";

const LightboxComponent = lazyComponent(() => import("./Lightbox"));

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

  const setPartialState = useCallback((state: Partial<IState>) => {
    setLightboxState((currentState: IState) => ({
      ...currentState,
      ...state,
    }));
  }, []);

  const onHide = () => {
    // slideshowAutostart is a per-open instruction (set when opening a gallery's
    // lightbox from the galleries page). Clear it on close so it doesn't leak
    // into the next lightbox opened from another entry point.
    setLightboxState({
      ...lightboxState,
      isVisible: false,
      slideshowAutostart: false,
    });
    if (lightboxState.onClose) {
      lightboxState.onClose();
    }
  };

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
