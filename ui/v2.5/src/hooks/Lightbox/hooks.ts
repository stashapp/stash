import { useCallback, useEffect, useRef } from "react";
import * as GQL from "src/core/generated-graphql";
import { useToast } from "../Toast";
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
      totalCount: state.totalCount,
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
    state.totalCount,
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
        totalCount: props.totalCount ?? state.totalCount,
        chapters: chapters,
      });
    },
    [
      setLightboxState,
      state.page,
      state.pages,
      state.pageSize,
      state.totalCount,
      chapters,
    ]
  );
  return show;
};

// Drives the lightbox for a gallery chosen at call time. A single lazy query is
// used and variables are passed to the execute function, so it can be
// instantiated once per list (rather than once per card) and the gallery to
// open is supplied to show() rather than baked into the hook.
interface ILightboxGallery {
  id: string;
  chapters?: IChapter[];
}

export const useGalleriesLightbox = () => {
  const { setLightboxState } = useLightboxContext();
  const Toast = useToast();

  const pageSize = 40;
  const [fetchImages, { data }] = GQL.useFindImagesLazyQuery();

  // details of the gallery currently shown in the lightbox, used for paging
  const active = useRef<{
    id: string;
    chapters: IChapter[];
    slideshowEnabled: boolean;
    slideshowAutostart: boolean;
    page: number;
    pages: number;
  }>();

  // Keep the lightbox images in sync with the live query result. The lazy query
  // stays subscribed after execution, so its data updates whenever the Apollo
  // cache changes (e.g. rating an image from inside the lightbox evicts and
  // refetches findImages). Push those updates into the lightbox state so edits
  // are reflected immediately, mirroring useLightbox above. Without this the
  // images are a one-shot snapshot taken at fetch time and never refresh.
  useEffect(() => {
    if (!active.current) return;
    const images = data?.findImages?.images;
    if (!images) return;
    // totalCount drives the lightbox position counter; keep it in sync with the
    // live query so it stays correct after edits that refetch (e.g. deletion).
    setLightboxState({ images, totalCount: data?.findImages.count });
  }, [data, setLightboxState]);

  function loadPage(page: number) {
    const current = active.current;
    if (!current) return;

    fetchImages({
      variables: {
        filter: { page, per_page: pageSize, sort: "path" },
        image_filter: {
          galleries: {
            modifier: GQL.CriterionModifier.Includes,
            value: [current.id],
          },
        },
      },
    })
      .then((result) => {
        // ignore if a different gallery was opened in the meantime
        if (active.current !== current) return;

        const totalCount = result.data?.findImages.count ?? 0;
        const pages = Math.ceil(totalCount / pageSize);
        current.page = page;
        current.pages = pages;

        setLightboxState({
          isLoading: false,
          isVisible: true,
          images: result.data?.findImages?.images ?? [],
          pageCallback: pages > 1 ? handleLightBoxPage : undefined,
          page,
          pages,
          pageSize,
          totalCount,
          chapters: current.chapters,
          slideshowEnabled: current.slideshowEnabled,
          slideshowAutostart: current.slideshowAutostart,
        });
      })
      .catch((e) => {
        // A failed page load leaves the lightbox on its loader: the page and
        // images props never change, so an in-flight switch cannot settle.
        // Surface the error so the stall is at least explained.
        Toast.error(e);
      });
  }

  function handleLightBoxPage(props: { direction?: number; page?: number }) {
    const current = active.current;
    if (!current) return;

    const { direction, page: newPage } = props;
    let target = current.page;

    if (direction !== undefined) {
      if (direction < 0) {
        target = current.page === 1 ? current.pages : current.page + direction;
      } else if (direction > 0) {
        target = current.page === current.pages ? 1 : current.page + direction;
      }
    } else if (newPage !== undefined) {
      target = newPage;
    }

    loadPage(target);
  }

  function show(
    gallery: ILightboxGallery,
    index: number = 0,
    slideshowEnabled: boolean = false,
    slideshowAutostart: boolean = false
  ) {
    const chapters = gallery.chapters ?? [];
    const startPage = Math.floor(index / pageSize) + 1;
    const initialIndex = index % pageSize;

    active.current = {
      id: gallery.id,
      chapters,
      slideshowEnabled,
      slideshowAutostart,
      page: startPage,
      pages: 1,
    };

    setLightboxState({
      images: [],
      isLoading: true,
      isVisible: true,
      initialIndex,
      showNavigation: false,
      pageCallback: undefined,
      page: undefined,
      pageSize,
      chapters,
      slideshowEnabled,
      slideshowAutostart,
    });

    loadPage(startPage);
  }

  return show;
};

// Bound to a single gallery, for callers that always open the same one (the
// gallery detail page, wall cards, and the plugin API). Thin wrapper over
// useGalleriesLightbox so the paging/fetch logic lives in one place.
export const useGalleryLightbox = (id: string, chapters: IChapter[] = []) => {
  const show = useGalleriesLightbox();
  return (index: number = 0) => show({ id, chapters }, index);
};
