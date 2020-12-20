import React, { useCallback, useEffect, useRef, useState } from 'react';
import * as GQL from 'src/core/generated-graphql';
import { Button } from 'react-bootstrap';
import cx from 'classnames';
import Mousetrap from "mousetrap";
import { debounce } from 'lodash';

import { Icon, LoadingIndicator } from 'src/components/Shared';

const CLASSNAME = 'Lightbox';
const CLASSNAME_HEADER = `${CLASSNAME}-header`;
const CLASSNAME_INDICATOR = `${CLASSNAME_HEADER}-indicator`;
const CLASSNAME_DISPLAY = `${CLASSNAME}-display`;
const CLASSNAME_CAROUSEL = `${CLASSNAME}-carousel`;
const CLASSNAME_INSTANT = `${CLASSNAME_CAROUSEL}-instant`;
const CLASSNAME_IMAGE = `${CLASSNAME_CAROUSEL}-image`;
const CLASSNAME_NAVBUTTON = `${CLASSNAME}-navbutton`;
const CLASSNAME_NAV = `${CLASSNAME}-nav`;
const CLASSNAME_NAVIMAGE = `${CLASSNAME_NAV}-image`;
const CLASSNAME_NAVSELECTED = `${CLASSNAME_NAV}-selected`;

type Image = Pick<GQL.Image, 'paths'>;
interface IProps {
  images: Image[];
  isVisible: boolean;
  isLoading: boolean;
  initialIndex?: number;
  showNavigation: boolean;
  pageHeader?: string;
  pageCallback?: (direction: number) => boolean;
  hide: () => void;
}

export const LightboxComponent: React.FC<IProps> = ({
  images,
  isVisible,
  isLoading,
  initialIndex = 0,
  showNavigation,
  pageHeader,
  pageCallback,
  hide,
}) => {
  const index = useRef<number | null>(null);
  const [instantTransition, setInstantTransition] = useState(false);
  const [isSwitchingPage, setIsSwitchingPage] = useState(false);
  const [isFullscreen, setFullscreen] = useState(false);
  const containerRef = useRef<HTMLDivElement | null>(null);
  const carouselRef = useRef<HTMLDivElement | null>(null);
  const indicatorRef = useRef<HTMLDivElement | null>(null);
  const navRef = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    setIsSwitchingPage(false);
    if (index.current === -1)
      index.current = images.length - 1;
  }, [images]);

  const disableInstantTransition = debounce(() => setInstantTransition(false), 400);
  const setInstant = useCallback(() => {
    setInstantTransition(true);
    disableInstantTransition();
  }, [disableInstantTransition]);

  const setIndex = useCallback((i: number) => {
    if (images.length < 2)
      return;

    index.current = i;
    if (carouselRef.current)
      carouselRef.current.style.left = `${(i * -100)}vw`;
    if (indicatorRef.current)
      indicatorRef.current.innerHTML = `${i + 1} / ${images.length}`;
    if (navRef.current) {
      const currentThumb = navRef.current.children[i + 1];
      if (currentThumb instanceof HTMLImageElement) {
        const offset = -1 * (currentThumb.offsetLeft - (document.documentElement.clientWidth / 2));
        navRef.current.style.left = `${offset}px`;

        const previouslySelected = navRef.current.getElementsByClassName(CLASSNAME_NAVSELECTED)?.[0]
        if (previouslySelected)
          previouslySelected.className = CLASSNAME_NAVIMAGE;

        currentThumb.className = `${CLASSNAME_NAVIMAGE} ${CLASSNAME_NAVSELECTED}`;
      }
    }
  }, [images]);

  const selectIndex = (e: React.MouseEvent, i: number) => {
    setIndex(i)
    e.stopPropagation();
  }

  useEffect(() => {
    if (isVisible) {
      if (index.current === null)
        setIndex(initialIndex);
      document.body.style.overflow = 'hidden';
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      (Mousetrap as any).pause();
    }
  }, [initialIndex, isVisible, setIndex]);

  const close = useCallback(() => {
    if (!isFullscreen) {
      hide();
      document.body.style.overflow = 'auto';
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      (Mousetrap as any).unpause();
    }
    else
      document.exitFullscreen();
  }, [isFullscreen, hide]);

  const handleClose = (e: React.MouseEvent<HTMLDivElement>) => {
    const { nodeName } = (e.target as Node);
    if (nodeName === "DIV" || nodeName === "PICTURE")
      close();
  }

  const handleLeft = useCallback(() => {
    if (isSwitchingPage || index.current === -1)
      return;

    if (index.current === 0) {
      if (pageCallback) {
        setIsSwitchingPage(true);
        setIndex(-1);
        // Check if calling page wants to swap page
        const repage = pageCallback(-1);
        if (!repage) {
          setIsSwitchingPage(false);
          setIndex(0);
        }
      }
      else
        setIndex(images.length - 1);
    }
    else
      setIndex((index.current ?? 0) - 1);
  }, [images, setIndex, pageCallback, isSwitchingPage]);
  const handleRight = useCallback(() => {
    if (isSwitchingPage)
      return;

    if (index.current === images.length - 1) {
      if (pageCallback) {
        setIsSwitchingPage(true);
        setIndex(0);
        const repage = pageCallback?.(1);
        if (!repage) {
          setIsSwitchingPage(false);
          setIndex(images.length - 1);
        }
      }
      else
        setIndex(0);
    }
    else
      setIndex((index.current ?? 0) + 1);
  }, [images, setIndex, pageCallback, isSwitchingPage]);

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
  const handleFullScreenChange = () => setFullscreen(document.fullscreenElement !== null);

  const handleTouchStart = (ev: React.TouchEvent<HTMLDivElement>) => {
    setInstantTransition(true);

    const el = ev.currentTarget;
    if (ev.touches.length !== 1)
      return;

    const startX = ev.touches[0].clientX;
    let position = 0;

    const resetPosition = () => {
      if (carouselRef.current)
        carouselRef.current.style.left = `${((index.current ?? 0) * -100)}vw`;
    }
    const handleMove = (e: TouchEvent) => {
      position = e.touches[0].clientX;
      if (carouselRef.current)
        carouselRef.current.style.left = `calc(${((index.current ?? 0) * -100)}vw + ${e.touches[0].clientX - startX}px)`;
    }
    const handleEnd = () => {
      const diff = position - startX;
      if (diff <= -50)
        handleRight();
      else if (diff >= 50)
        handleLeft();
      else
        resetPosition();
      // eslint-disable-next-line @typescript-eslint/no-use-before-define
      cleanup();
    }
    const handleCancel = () => {
      // eslint-disable-next-line @typescript-eslint/no-use-before-define
      cleanup();
      resetPosition();
    }
    const cleanup = () => {
      el.removeEventListener('touchmove', handleMove);
      el.removeEventListener('touchend', handleEnd);
      el.removeEventListener('touchcancel', handleCancel);
      setInstantTransition(false);
    }

    el.addEventListener('touchmove', handleMove);
    el.addEventListener('touchend', handleEnd);
    el.addEventListener('touchcancel', handleCancel);
  }

  useEffect(() => {
    if (isVisible) {
      document.addEventListener('keydown', handleKey);
      document.addEventListener('fullscreenchange', handleFullScreenChange);
    }
    return () => {
      document.removeEventListener('keydown', handleKey);
      document.removeEventListener('fullscreenchange', handleFullScreenChange);
    }
  }, [isVisible, handleKey]);

  const toggleFullscreen = useCallback(() => {
    if (!isFullscreen)
      containerRef.current?.requestFullscreen();
    else
      document.exitFullscreen();
  }, [isFullscreen]);

  const navItems = images.map((image, i) => (
    <img src={image.paths.thumbnail ?? ''} alt="" className={cx(CLASSNAME_NAVIMAGE, { [CLASSNAME_NAVSELECTED]: i === index.current})} onClick={(e: React.MouseEvent) => selectIndex(e, i)} role="presentation" loading="lazy" />
  ));

  const currentIndex = index.current === null ? initialIndex : index.current;

  const element = isVisible ? (
    <div className={CLASSNAME} role="presentation" ref={containerRef} onClick={handleClose}>
      { images.length > 0 && !isLoading && !isSwitchingPage ? (
        <>
          <div className={CLASSNAME_HEADER}>
            <div className={CLASSNAME_INDICATOR}>
              <span>{ pageHeader }</span>
              <b ref={indicatorRef}>
                { `${currentIndex + 1} / ${images.length}` }
              </b>
            </div>
            { document.fullscreenEnabled && (
              <Button variant="link" onClick={toggleFullscreen} title="Toggle Fullscreen">
                <Icon icon="expand" />
              </Button>
            )}
            <Button variant="link" onClick={() => close()} title="Close Lightbox">
              <Icon icon="times" />
            </Button>
          </div>
          <div className={CLASSNAME_DISPLAY} onTouchStart={handleTouchStart}>
            { images.length > 1 && (
              <Button variant="link" onClick={handleLeft} className={`${CLASSNAME_NAVBUTTON} d-none d-lg-block`}>
                <Icon icon="chevron-circle-left" />
              </Button>
            )}

            <div
              className={cx(CLASSNAME_CAROUSEL, { [CLASSNAME_INSTANT]: instantTransition })}
              style={{ left: `${currentIndex * -100}vw` }}
              ref={carouselRef}
            >
              { images.map(image => (
                  <div className={CLASSNAME_IMAGE}>
                    <picture>
                      <source srcSet={image.paths.image ?? ''} media="(min-width: 800px)" />
                      <img src={image.paths.thumbnail ?? ''} alt="" />
                    </picture>
                  </div>
              ))}
            </div>

            { images.length > 1 && (
              <Button variant="link" onClick={handleRight} className={`${CLASSNAME_NAVBUTTON} d-none d-lg-block`}>
                <Icon icon="chevron-circle-right" />
              </Button>
            )}
          </div>
          { showNavigation && !isFullscreen && images.length > 1 && (
            <div className={CLASSNAME_NAV} ref={navRef}>
              <Button variant="link" onClick={() => setIndex(images.length - 1)} className={CLASSNAME_NAVBUTTON}>
                <Icon icon="chevron-left" />
              </Button>
              { navItems }
              <Button variant="link" onClick={() => setIndex(0)} className={CLASSNAME_NAVBUTTON}>
                <Icon icon="chevron-right" />
              </Button>
            </div>
          )}
        </>
      )  : <LoadingIndicator /> }
    </div>
  ) : <></>;

  return element;
}
