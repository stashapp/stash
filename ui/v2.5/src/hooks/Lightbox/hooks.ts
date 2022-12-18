import { useCallback, useContext, useEffect, useMemo, useState } from "react";
import * as GQL from "src/core/generated-graphql";
import { LightboxContext, IState } from "./context";
import { IChapter } from "./types";

export const useLightbox = (
  state: Partial<Omit<IState, "isVisible">>,
  chapters: IChapter[] = []
) => {
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
        chapters: chapters,
      });
    },
    [setLightboxState, chapters]
  );
  return show;
};

export const useGalleryLightbox = (id: string, chapters: IChapter[] = []) => {
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
      if (direction < 0) {
        if (page === 1) {
          setPage(pages);
        } else {
          setPage(page + direction);
        }
      } else if (direction > 0) {
        if (page === pages) {
          // return to the first page
          setPage(1);
        } else {
          setPage(page + direction);
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

  const show = (index: number = 0) => {
    if (index > 40) {
      setPage(Math.floor(index / 40) + 1);
      index = index % 40;
    }
    if (data)
      setLightboxState({
        isLoading: false,
        isVisible: true,
        initialIndex: index,
        images: data.findImages?.images ?? [],
        pageCallback: pages > 1 ? handleLightBoxPage : undefined,
        pageHeader: `Page ${page} / ${pages}`,
        chapters: chapters,
      });
    else {
      setLightboxState({
        isLoading: true,
        isVisible: true,
        initialIndex: index,
        pageCallback: undefined,
        pageHeader: undefined,
        chapters: chapters,
      });
      fetchGallery();
    }
  };

  return show;
};
