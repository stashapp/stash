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

  const setPartialState = useCallback(
    (state: Partial<IState>) => {
      setLightboxState((currentState: IState) => ({
        ...currentState,
        ...state,
      }));
    },
    [setLightboxState]
  );

  const onHide = () => {
    setLightboxState({ ...lightboxState, isVisible: false });
    if (lightboxState.onClose) {
      lightboxState.onClose();
    }
  };

  return (
    <LightboxContext.Provider
      value={{ lightboxState, setLightboxState: setPartialState }}
    >
      {children}
      <Suspense fallback={<></>}>
        {lightboxState.isVisible && (
          <LightboxComponent {...lightboxState} hide={onHide} />
        )}
      </Suspense>
    </LightboxContext.Provider>
  );
};
