import { useCallback, useContext, useEffect, useMemo, useState } from "react";
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

  const pageSize = 40;
  const [page, setPage] = useState(1);

  const currentFilter = useMemo(() => {
    return {
      page,
      per_page: pageSize,
      sort: "path",
    };
  }, [page]);

  const [fetchGallery, { data }] = GQL.useFindImagesLazyQuery({
    variables: {
      filter: currentFilter,
      image_filter: {
        galleries: {
          modifier: GQL.CriterionModifier.Includes,
          value: [id],
        },
      },
    },
  });

  const pages = useMemo(() => {
    const totalCount = data?.findImages.count ?? 0;
    return Math.ceil(totalCount / pageSize);
  }, [data?.findImages.count]);

  const handleLightBoxPage = useCallback(
    (direction: number) => {
      if (direction === -1) {
        if (page === 1) {
          setPage(pages);
        } else {
          setPage(page - 1);
        }
      } else if (direction === 1) {
        if (page === pages) {
          // return to the first page
          setPage(1);
        } else {
          setPage(page + 1);
        }
      }
    },
    [page, pages]
  );

  useEffect(() => {
    if (data)
      setLightboxState({
        isLoading: false,
        isVisible: true,
        images: data.findImages?.images ?? [],
        pageCallback: pages > 1 ? handleLightBoxPage : undefined,
        pageHeader: `Page ${page} / ${pages}`,
      });
  }, [setLightboxState, data, handleLightBoxPage, page, pages]);

  const show = () => {
    if (data)
      setLightboxState({
        isLoading: false,
        isVisible: true,
        images: data.findImages?.images ?? [],
        pageCallback: pages > 1 ? handleLightBoxPage : undefined,
        pageHeader: `Page ${page} / ${pages}`,
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
