import React, { Suspense, useCallback, useEffect, useRef, useState } from "react";
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
  slideshowEnabled: boolean;
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

  // Use ref for onClose to avoid stale closures in callbacks
  const onCloseRef = useRef<(() => void) | undefined>();
  onCloseRef.current = lightboxState.onClose;

  const setPartialState = useCallback((state: Partial<IState>) => {
    setLightboxState((currentState: IState) => {
      // Push history entry when lightbox opens so back button closes it
      // instead of navigating away from the current page
      if (state.isVisible && !currentState.isVisible) {
        history.pushState({ lightbox: true }, '');
      }
      return { ...currentState, ...state };
    });
  }, []);

  const onHide = useCallback(() => {
    // User-initiated close (close button, escape, etc.) — navigate back
    // which will trigger popstate and handle the actual state update
    history.back();
  }, []);

  // Close lightbox on browser back/forward navigation
  useEffect(() => {
    const handlePopState = () => {
      setLightboxState((currentState: IState) => {
        if (currentState.isVisible) {
          return { ...currentState, isVisible: false };
        }
        return currentState;
      });
      onCloseRef.current?.();
    };

    window.addEventListener("popstate", handlePopState);
    return () => {
      window.removeEventListener("popstate", handlePopState);
    };
  }, []);

  return (
    <LightboxContext.Provider
      value={{ lightboxState, setLightboxState: setPartialState }}
    >
      {children}
      <Suspense fallback={null}>
        {lightboxState.isVisible && (
          <LightboxComponent {...lightboxState} hide={onHide} />
        )}
      </Suspense>
    </LightboxContext.Provider>
  );
};
