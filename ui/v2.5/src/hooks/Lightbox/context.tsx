import React, { useCallback, useState } from "react";
import * as GQL from "src/core/generated-graphql";
import { LightboxComponent } from "./Lightbox";

type Image = Pick<GQL.Image, "paths">;

export interface IState {
  images: Image[];
  isVisible: boolean;
  isLoading: boolean;
  showNavigation: boolean;
  initialIndex?: number;
  pageCallback?: (direction: number) => boolean;
  pageHeader?: string;
}
interface IContext {
  setLightboxState: (state: Partial<IState>) => void;
}

export const LightboxContext = React.createContext<IContext>({
  setLightboxState: () => {},
});
const Lightbox: React.FC = ({ children }) => {
  const [lightboxState, setLightboxState] = useState<IState>({
    images: [],
    isVisible: false,
    isLoading: false,
    showNavigation: true,
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

  return (
    <LightboxContext.Provider value={{ setLightboxState: setPartialState }}>
      {children}
      {lightboxState.isVisible && (
        <LightboxComponent
          {...lightboxState}
          hide={() => setLightboxState({ ...lightboxState, isVisible: false })}
        />
      )}
    </LightboxContext.Provider>
  );
};

export default Lightbox;
