import React, { useCallback, useEffect, useRef, useState } from 'react';
import * as GQL from 'src/core/generated-graphql';
import { Button } from 'react-bootstrap';
import cx from 'classnames';
import Mousetrap from "mousetrap";
import { debounce } from 'lodash';

import { Icon, LoadingIndicator } from 'src/components/Shared';

const CLASSNAME = 'Lightbox';
const CLASSNAME_HEADER = `${CLASSNAME}-header`;
const CLASSNAME_DISPLAY = `${CLASSNAME}-display`;
const CLASSNAME_CAROUSEL = `${CLASSNAME}-carousel`;
const CLASSNAME_INSTANT = `${CLASSNAME_CAROUSEL}-instant`;
const CLASSNAME_IMAGE = `${CLASSNAME_CAROUSEL}-image`;
const CLASSNAME_SELECTED = `${CLASSNAME_IMAGE}-selected`;
const CLASSNAME_PRELOAD = `${CLASSNAME_IMAGE}-preload`;
const CLASSNAME_NAV = `${CLASSNAME}-nav`;
const CLASSNAME_NAVIMAGE = `${CLASSNAME}-nav-image`;
const CLASSNAME_NAVSELECTED = `${CLASSNAME}-nav-selected`;

type Image = Pick<GQL.Image, 'paths'>;
type LightboxHookResult = [
  (index?: number) => void,
  React.ReactNode
];

export const useLightbox = (images: Image[], showNavigation = true): LightboxHookResult => {
  const [isVisible, setVisible] = useState(false);
  const [index, setIndex] = useState(0);
  const [instantTransition, setInstantTransition] = useState(false);
  const [isFullscreen, setFullscreen] = useState(false);
  const containerRef = useRef<HTMLDivElement | null>(null);

  const disableInstantTransition = debounce(() => setInstantTransition(false), 400);
  const setInstant = useCallback(() => {
    setInstantTransition(true);
    disableInstantTransition();
  }, [disableInstantTransition]);

  const selectIndex = (e: React.MouseEvent, i: number) => {
    setIndex(i);
    e.stopPropagation();
  }

  const exitFullscreen = () => document.exitFullscreen().then(() => setFullscreen(false));

  const close = useCallback(() => {
    if (!isFullscreen) {
      setVisible(false)
      document.body.style.overflow = 'auto';
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      (Mousetrap as any).unpause();
    }
    else
      exitFullscreen();
  }, [isFullscreen]);

  const handleClose = (e: React.MouseEvent<HTMLDivElement>) => {
    if ((e.target as Node).nodeName === 'DIV')
      close();
  }

  const handleLeft = useCallback(() => setIndex(i => i === 0 ? images.length - 1 : i - 1), [images]);
  const handleRight = useCallback(() => setIndex(i => i === images.length - 1 ? 0 : i + 1), [images]);

  const handleKey = useCallback((e: KeyboardEvent) => {
    if (e.repeat && (e.key === "ArrowRight" || e.key === "ArrowLeft"))
      setInstant();
    if (e.key === "ArrowLeft")
      handleLeft();
    else if (e.key === "ArrowRight")
      handleRight();
    else if (e.key === "Escape")
      close();
  }, [setInstant, handleLeft, handleRight, close]);

  const show = (showIndex: number = 0) => {
    setIndex(showIndex);
    setVisible(true);
    document.body.style.overflow = 'hidden';
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    (Mousetrap as any).pause();
  }

  useEffect(() => {
    if (isVisible)
      document.addEventListener('keydown', handleKey);
    return () => document.removeEventListener('keydown', handleKey);
  }, [isVisible, handleKey]);

  const toggleFullscreen = useCallback(() => {
    if (!isFullscreen)
      containerRef.current?.requestFullscreen().then(() => setFullscreen(true));
    else
      exitFullscreen();
  }, [isFullscreen]);

  const shouldPreload = (i: number) => {
    const distance = Math.abs(index - i);
    return (distance <= 2) || (distance >= (images.length - 2));
  }

  const navItems = images.map((image, i) => (
    <img src={image.paths.thumbnail ?? ''} alt="" className={cx(CLASSNAME_NAVIMAGE, { [CLASSNAME_NAVSELECTED]: i === index})} onClick={(e: React.MouseEvent) => selectIndex(e, i)} role="presentation" loading="lazy" />
  ));

  const element = isVisible ? (
    <div className={CLASSNAME} role="presentation" ref={containerRef} onClick={handleClose}>
      { images.length > 0 ? (
        <>
          <div className={CLASSNAME_HEADER}>
            <Button variant="link" onClick={toggleFullscreen} title="Toggle Fullscreen">
              <Icon icon="expand" />
            </Button>
            <Button variant="link" onClick={() => close()} title="Close Lightbox">
              <Icon icon="times" />
            </Button>
          </div>
          <div className={CLASSNAME_DISPLAY}>
            <Button variant="link" onClick={handleLeft}>
              <Icon icon="chevron-left" />
            </Button>

            <div className={cx(CLASSNAME_CAROUSEL, { [CLASSNAME_INSTANT]: instantTransition })} style={{ left: `${index * -100}vw` }}>
              { images.map((image, i) => {
                const preload = shouldPreload(i);
                return (
                  <div className={cx(CLASSNAME_IMAGE, { [CLASSNAME_SELECTED]: i === index, [CLASSNAME_PRELOAD]: preload })}>
                    <img src={image.paths.image ?? ''} alt="" loading={preload ? 'eager' : 'lazy'} />
                  </div>
                );
              })}
            </div>

            <Button variant="link" onClick={handleRight}>
              <Icon icon="chevron-right" />
            </Button>
          </div>
          { showNavigation && !isFullscreen && images.length > 1 && (
            <div className={CLASSNAME_NAV}>
              { navItems }
            </div>
          )}
        </>
      )  : <LoadingIndicator /> }
    </div>
  ) : <></>;

  return [show, element];
}

export const useGalleryLightbox = (id: string) => {
  const [fetchGallery, { data }] = GQL.useFindGalleryLazyQuery({ variables: { id } });
  const [showGallery, container] = useLightbox(data?.findGallery?.images ?? [])

  const show = () => {
    fetchGallery();
    showGallery();
    return false;
  }

  return [show, container] as const;
}
