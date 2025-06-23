import { useCallback, useEffect, useMemo, useState } from "react";
import * as GQL from "src/core/generated-graphql";
import { IState, useLightboxContext } from "./context";
import { IChapter } from "./types";

export const useLightbox = (
  state: Partial<Omit<IState, "isVisible">> = {},
  chapters: IChapter[] = []
) => {
  const { setLightboxState } = useLightboxContext();

  useEffect(() => {
    setLightboxState({
      images: state.images,
      showNavigation: state.showNavigation,
      pageCallback: state.pageCallback,
      page: state.page,
      pages: state.pages,
      pageSize: state.pageSize,
      slideshowEnabled: state.slideshowEnabled,
      onClose: state.onClose,
    });
  }, [
    setLightboxState,
    state.images,
    state.showNavigation,
    state.pageCallback,
    state.page,
    state.pages,
    state.pageSize,
    state.slideshowEnabled,
    state.onClose,
  ]);

  const show = useCallback(
    (props: Partial<IState>) => {
      setLightboxState({
        ...props,
        isVisible: true,
        page: props.page ?? state.page,
        pages: props.pages ?? state.pages,
        pageSize: props.pageSize ?? state.pageSize,
        chapters: chapters,
      });
    },
    [setLightboxState, state.page, state.pages, state.pageSize, chapters]
  );
  return show;
};

export const useGalleryLightbox = (id: string, chapters: IChapter[] = []) => {
  const { setLightboxState } = useLightboxContext();

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
    (props: { direction?: number; page?: number }) => {
      const { direction, page: newPage } = props;

      if (direction !== undefined) {
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
      } else if (newPage !== undefined) {
        setPage(newPage);
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
        page,
        pages,
      });
  }, [setLightboxState, data, handleLightBoxPage, page, pages]);

  const show = (index: number = 0) => {
    if (index > pageSize) {
      setPage(Math.floor(index / pageSize) + 1);
      index = index % pageSize;
    } else {
      setPage(1);
    }
    if (data)
      setLightboxState({
        isLoading: false,
        isVisible: true,
        initialIndex: index,
        images: data.findImages?.images ?? [],
        pageCallback: pages > 1 ? handleLightBoxPage : undefined,
        page,
        pages,
        pageSize,
        chapters: chapters,
      });
    else {
      setLightboxState({
        images: [],
        isLoading: true,
        isVisible: true,
        initialIndex: index,
        pageCallback: undefined,
        page: undefined,
        pageSize,
        chapters: chapters,
      });
      fetchGallery();
    }
  };

  return show;
};
