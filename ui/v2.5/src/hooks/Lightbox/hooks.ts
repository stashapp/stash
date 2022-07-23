import { useCallback, useContext, useEffect } from "react";
import * as GQL from "src/core/generated-graphql";
import { LightboxContext, IState } from "./context";

export const useLightbox = (state: Partial<Omit<IState, "isVisible">>) => {
  const { setLightboxState } = useContext(LightboxContext);

  useEffect(() => {
    setLightboxState({
      images: state.images,
      showNavigation: state.showNavigation,
      pageCallback: state.pageCallback,
      pageHeader: state.pageHeader,
      slideshowEnabled: state.slideshowEnabled,
      onClose: state.onClose,
    });
  }, [
    setLightboxState,
    state.images,
    state.showNavigation,
    state.pageCallback,
    state.pageHeader,
    state.slideshowEnabled,
    state.onClose,
  ]);

  const show = useCallback(
    (index?: number, slideshowEnabled = false) => {
      setLightboxState({
        initialIndex: index,
        isVisible: true,
        slideshowEnabled,
      });
    },
    [setLightboxState]
  );
  return show;
};

export const useGalleryLightbox = (id: string) => {
  const { setLightboxState } = useContext(LightboxContext);
  const [fetchGallery, { data }] = GQL.useFindGalleryLazyQuery({
    variables: { id },
  });

  useEffect(() => {
    if (data)
      setLightboxState({
        images: data.findGallery?.images ?? [],
        isLoading: false,
        isVisible: true,
      });
  }, [setLightboxState, data]);

  const show = () => {
    if (data)
      setLightboxState({
        isLoading: false,
        isVisible: true,
        images: data.findGallery?.images ?? [],
        pageCallback: undefined,
        pageHeader: undefined,
      });
    else {
      setLightboxState({
        isLoading: true,
        isVisible: true,
        pageCallback: undefined,
        pageHeader: undefined,
      });
      fetchGallery();
    }
  };

  return show;
};
