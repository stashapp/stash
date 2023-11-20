import React, {
  CSSProperties,
  useEffect,
  useRef,
  useState,
  useCallback,
} from "react";
import { Button } from "react-bootstrap";
import * as GQL from "src/core/generated-graphql";
import TextUtils from "src/utils/text";
import { Icon } from "src/components/Shared/Icon";
import {
  faChevronRight,
  faChevronLeft,
} from "@fortawesome/free-solid-svg-icons";
import { useSpriteInfo } from "src/hooks/sprite";

interface IScenePlayerScrubberProps {
  file: GQL.VideoFileDataFragment;
  scene: GQL.SceneDataFragment;
  time: number;
  onSeek: (seconds: number) => void;
  onScroll: () => void;
}

interface ISceneSpriteItem {
  style: CSSProperties;
  time: string;
}

export const ScenePlayerScrubber: React.FC<IScenePlayerScrubberProps> = ({
  file,
  scene,
  time,
  onSeek,
  onScroll,
}) => {
  const contentEl = useRef<HTMLDivElement>(null);
  const indicatorEl = useRef<HTMLDivElement>(null);
  const sliderEl = useRef<HTMLDivElement>(null);
  const mouseDown = useRef(false);
  const lastMouseEvent = useRef<MouseEvent | null>(null);
  const startMouseEvent = useRef<MouseEvent | null>(null);
  const velocity = useRef(0);

  const prevTime = useRef(NaN);
  const _width = useRef(0);
  const [width, setWidth] = useState(0);
  const [scrubWidth, setScrubWidth] = useState(0);
  const position = useRef(0);
  const setPosition = useCallback(
    (value: number, seek: boolean) => {
      if (!scrubWidth) return;

      const slider = sliderEl.current!;
      const indicator = indicatorEl.current!;

      const midpointOffset = slider.clientWidth / 2;

      let newPosition: number;
      let percentage: number;
      if (value >= midpointOffset) {
        percentage = 0;
        newPosition = midpointOffset;
      } else if (value <= midpointOffset - scrubWidth) {
        percentage = 1;
        newPosition = midpointOffset - scrubWidth;
      } else {
        percentage = (midpointOffset - value) / scrubWidth;
        newPosition = value;
      }

      slider.style.transform = `translateX(${newPosition}px)`;
      indicator.style.transform = `translateX(${percentage * 100}%)`;

      position.current = newPosition;

      if (seek) {
        onSeek(percentage * (file.duration || 0));
      }
    },
    [onSeek, file.duration, scrubWidth]
  );

  const spriteInfo = useSpriteInfo(scene.paths.vtt ?? undefined);
  const [spriteItems, setSpriteItems] = useState<ISceneSpriteItem[]>();

  useEffect(() => {
    if (!spriteInfo) return;
    let totalWidth = 0;
    const newSprites = spriteInfo?.map((sprite, index) => {
      totalWidth += sprite.w;
      const left = sprite.w * index;
      const style = {
        width: `${sprite.w}px`,
        height: `${sprite.h}px`,
        backgroundPosition: `${-sprite.x}px ${-sprite.y}px`,
        backgroundImage: `url(${sprite.url})`,
        left: `${left}px`,
      };
      const start = TextUtils.secondsToTimestamp(sprite.start);
      const end = TextUtils.secondsToTimestamp(sprite.end);
      return {
        style,
        time: `${start} - ${end}`,
      };
    });
    setScrubWidth(totalWidth);
    setSpriteItems(newSprites);
  }, [spriteInfo]);

  useEffect(() => {
    const onResize = (entries: ResizeObserverEntry[]) => {
      const newWidth = entries[0].target.clientWidth;
      if (_width.current != newWidth) {
        // set prevTime to NaN to not use a transition when updating the slider position
        prevTime.current = NaN;
        _width.current = newWidth;
        setWidth(newWidth);
      }
    };

    const content = contentEl.current!;
    const resizeObserver = new ResizeObserver(onResize);
    resizeObserver.observe(content);

    return () => {
      resizeObserver.unobserve(content);
    };
  }, []);

  function setLinearTransition() {
    const slider = sliderEl.current!;
    slider.style.transition = "500ms linear";
  }

  function setEaseOutTransition() {
    const slider = sliderEl.current!;
    slider.style.transition = "333ms ease-out";
  }

  function clearTransition() {
    const slider = sliderEl.current!;
    slider.style.transition = "";
  }

  // Update slider position when player time changes
  useEffect(() => {
    if (!scrubWidth || !width) return;

    const duration = Number(file.duration);
    const percentage = time / duration;
    const newPosition = width / 2 - percentage * scrubWidth;

    // Ignore position changes of < 1px
    if (Math.abs(newPosition - position.current) < 1) return;

    const delta = Math.abs(time - prevTime.current);
    if (isNaN(delta)) {
      // Don't use a transition on initial time change or after resize
      clearTransition();
    } else if (delta <= 1) {
      // If time changed by < 1s, use linear transition instead of ease-out
      setLinearTransition();
    } else {
      setEaseOutTransition();
    }
    prevTime.current = time;

    setPosition(newPosition, false);
  }, [file.duration, setPosition, time, width, scrubWidth]);

  const onMouseUp = useCallback(
    (event: MouseEvent) => {
      if (!mouseDown.current) return;
      const slider = sliderEl.current!;

      mouseDown.current = false;

      contentEl.current!.classList.remove("dragging");

      let newPosition = position.current;
      const midpointOffset = slider.clientWidth / 2;
      const delta = Math.abs(event.clientX - startMouseEvent.current!.clientX);
      if (delta < 1 && event.target instanceof HTMLDivElement) {
        const { target } = event;

        if (target.hasAttribute("data-sprite-item-id")) {
          newPosition = midpointOffset - (target.offsetLeft + event.offsetX);
        }

        if (target.hasAttribute("data-marker-id")) {
          newPosition = midpointOffset - target.offsetLeft;
        }
      }
      if (Math.abs(velocity.current) > 25) {
        newPosition = position.current + velocity.current * 10;
        velocity.current = 0;
      }

      setEaseOutTransition();
      setPosition(newPosition, true);
    },
    [setPosition]
  );

  const onMouseDown = useCallback((event: MouseEvent) => {
    // Only if left mouse button pressed
    if (event.button !== 0) return;

    event.preventDefault();

    mouseDown.current = true;
    lastMouseEvent.current = event;
    startMouseEvent.current = event;
    velocity.current = 0;
  }, []);

  const onMouseMove = useCallback(
    (event: MouseEvent) => {
      if (!mouseDown.current) return;

      // negative dragging right (past), positive left (future)
      const delta = event.clientX - lastMouseEvent.current!.clientX;

      if (lastMouseEvent.current === startMouseEvent.current) {
        // this is the first mousemove event after mousedown

        // #4295: a mousemove with delta 0 can be sent when just clicking
        // ignore such an event to prevent pausing the player
        if (delta === 0) return;

        onScroll();
      }

      contentEl.current!.classList.add("dragging");

      const movement = event.movementX;
      velocity.current = movement;

      clearTransition();
      setPosition(position.current + delta, false);
      lastMouseEvent.current = event;
    },
    [onScroll, setPosition]
  );

  useEffect(() => {
    const content = contentEl.current!;

    content.addEventListener("mousedown", onMouseDown, false);
    content.addEventListener("mousemove", onMouseMove, false);
    window.addEventListener("mouseup", onMouseUp, false);

    return () => {
      content.removeEventListener("mousedown", onMouseDown);
      content.removeEventListener("mousemove", onMouseMove);
      window.removeEventListener("mouseup", onMouseUp);
    };
  }, [onMouseDown, onMouseMove, onMouseUp]);

  function goBack() {
    const slider = sliderEl.current!;
    const newPosition = position.current + slider.clientWidth;
    setEaseOutTransition();
    setPosition(newPosition, true);
  }

  function goForward() {
    const slider = sliderEl.current!;
    const newPosition = position.current - slider.clientWidth;
    setEaseOutTransition();
    setPosition(newPosition, true);
  }

  function renderTags() {
    if (!spriteItems) return;

    return scene.scene_markers.map((marker, index) => {
      const { duration } = file;
      const left = (scrubWidth * marker.seconds) / duration;
      const style = { left: `${left}px` };

      return (
        <div
          key={index}
          className="scrubber-tag"
          style={style}
          data-marker-id={index}
        >
          {marker.title || marker.primary_tag.name}
        </div>
      );
    });
  }

  function renderSprites() {
    if (!scene.paths.vtt) return;

    return spriteItems?.map((sprite, index) => {
      return (
        <div
          key={index}
          className="scrubber-item"
          style={sprite.style}
          data-sprite-item-id={index}
        >
          <span className="scrubber-item-time">{sprite.time}</span>
        </div>
      );
    });
  }

  return (
    <div className="scrubber-wrapper">
      <Button
        className="scrubber-button"
        id="scrubber-back"
        onClick={() => goBack()}
      >
        <Icon className="fa-fw" icon={faChevronLeft} />
      </Button>
      <div ref={contentEl} className="scrubber-content">
        <div className="scrubber-tags-background" />
        <div
          className="scrubber-heatmap"
          style={{
            backgroundImage: scene.paths.interactive_heatmap
              ? `url(${scene.paths.interactive_heatmap})`
              : undefined,
          }}
        />
        <div ref={indicatorEl} id="scrubber-position-indicator" />
        <div id="scrubber-current-position" />
        <div className="scrubber-viewport">
          <div ref={sliderEl} className="scrubber-slider">
            <div className="scrubber-tags">{renderTags()}</div>
            {renderSprites()}
          </div>
        </div>
      </div>
      <Button
        className="scrubber-button"
        id="scrubber-forward"
        onClick={() => goForward()}
      >
        <Icon className="fa-fw" icon={faChevronRight} />
      </Button>
    </div>
  );
};
