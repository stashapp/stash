/* eslint-disable react/no-array-index-key */

import React, {
  CSSProperties,
  useEffect,
  useRef,
  useState,
  useCallback,
} from "react";
import { Button } from "react-bootstrap";
import axios from "axios";
import * as GQL from "src/core/generated-graphql";
import { TextUtils } from "src/utils";

interface IScenePlayerScrubberProps {
  file: GQL.VideoFileDataFragment;
  scene: GQL.SceneDataFragment;
  position: number;
  onSeek: (seconds: number) => void;
  onScrolled: () => void;
}

interface ISceneSpriteItem {
  start: number;
  end: number;
  x: number;
  y: number;
  w: number;
  h: number;
}

async function fetchSpriteInfo(vttPath: string) {
  const response = await axios.get<string>(vttPath, { responseType: "text" });

  // TODO: This is gnarly
  const lines = response.data.split("\n");
  if (lines.shift() !== "WEBVTT") {
    return;
  }
  if (lines.shift() !== "") {
    return;
  }
  let item: ISceneSpriteItem = { start: 0, end: 0, x: 0, y: 0, w: 0, h: 0 };
  const newSpriteItems: ISceneSpriteItem[] = [];
  while (lines.length) {
    const line = lines.shift();
    if (line !== undefined) {
      if (line.includes("#") && line.includes("=") && line.includes(",")) {
        const size = line.split("#")[1].split("=")[1].split(",");
        item.x = Number(size[0]);
        item.y = Number(size[1]);
        item.w = Number(size[2]);
        item.h = Number(size[3]);

        newSpriteItems.push(item);
        item = { start: 0, end: 0, x: 0, y: 0, w: 0, h: 0 };
      } else if (line.includes(" --> ")) {
        const times = line.split(" --> ");

        const start = times[0].split(":");
        item.start = +start[0] * 60 * 60 + +start[1] * 60 + +start[2];

        const end = times[1].split(":");
        item.end = +end[0] * 60 * 60 + +end[1] * 60 + +end[2];
      }
    }
  }

  return newSpriteItems;
}

export const ScenePlayerScrubber: React.FC<IScenePlayerScrubberProps> = (
  props: IScenePlayerScrubberProps
) => {
  const contentEl = useRef<HTMLDivElement>(null);
  const positionIndicatorEl = useRef<HTMLDivElement>(null);
  const scrubberSliderEl = useRef<HTMLDivElement>(null);
  const mouseDown = useRef(false);
  const lastMouseEvent = useRef<MouseEvent | null>(null);
  const startMouseEvent = useRef<MouseEvent | null>(null);
  const mouseVelocity = useRef(0);

  const touchStart = useRef(false);
  const lastTouchEvent = useRef<TouchEvent | null>(null);
  const startTouchEvent = useRef<TouchEvent | null>(null);
  const touchVelocity = useRef(0);

  const _position = useRef(0);
  const getPosition = useCallback(() => _position.current, []);
  const setPosition = useCallback(
    (newPostion: number, shouldEmit: boolean = true) => {
      if (!scrubberSliderEl.current || !positionIndicatorEl.current) {
        return;
      }
      if (shouldEmit) {
        props.onScrolled();
      }

      const midpointOffset = scrubberSliderEl.current.clientWidth / 2;

      const bounds = getBounds() * -1;
      if (newPostion > midpointOffset) {
        _position.current = midpointOffset;
      } else if (newPostion < bounds - midpointOffset) {
        _position.current = bounds - midpointOffset;
      } else {
        _position.current = newPostion;
      }

      scrubberSliderEl.current.style.transform = `translateX(${_position.current}px)`;

      const indicatorPosition =
        ((newPostion - midpointOffset) / (bounds - midpointOffset * 2)) *
        scrubberSliderEl.current.clientWidth;
      positionIndicatorEl.current.style.transform = `translateX(${indicatorPosition}px)`;
    },
    [props]
  );

  const [spriteItems, setSpriteItems] = useState<ISceneSpriteItem[]>([]);

  useEffect(() => {
    if (!scrubberSliderEl.current) {
      return;
    }
    scrubberSliderEl.current.style.transform = `translateX(${
      scrubberSliderEl.current.clientWidth / 2
    }px)`;
  }, [scrubberSliderEl]);

  useEffect(() => {
    if (!props.scene.paths.vtt) return;
    fetchSpriteInfo(props.scene.paths.vtt).then((sprites) => {
      if (sprites) setSpriteItems(sprites);
    });
  }, [props.scene]);

  useEffect(() => {
    if (!scrubberSliderEl.current) {
      return;
    }
    const duration = Number(props.file.duration);
    const percentage = props.position / duration;
    const position =
      (scrubberSliderEl.current.scrollWidth * percentage -
        scrubberSliderEl.current.clientWidth / 2) *
      -1;
    setPosition(position, false);
  }, [props.position, props.file.duration, setPosition]);

  useEffect(() => {
    window.addEventListener("mouseup", onMouseUp, false);
    window.addEventListener("touchend", onTouchEnd, { passive: true });
    return () => {
      window.removeEventListener("mouseup", onMouseUp);
      window.removeEventListener("touchend", onTouchEnd);
    };
  });

  useEffect(() => {
    if (!contentEl.current) {
      return;
    }
    const el = contentEl.current;
    el.addEventListener("mousedown", onMouseDown, false);
    el.addEventListener("touchstart", onTouchDown, { passive: false });
    return () => {
      if (!el) {
        return;
      }
      el.removeEventListener("mousedown", onMouseDown);
      el.removeEventListener("touchstart", onTouchDown);
    };
  });

  useEffect(() => {
    if (!contentEl.current) {
      return;
    }
    const el = contentEl.current;
    el.addEventListener("mousemove", onMouseMove, false);
    el.addEventListener("touchmove", onTouchMove, { passive: true });
    return () => {
      if (!el) {
        return;
      }
      el.removeEventListener("mousemove", onMouseMove);
      el.removeEventListener("touchmove", onTouchMove);
    };
  });

  function onMouseUp(this: Window, event: MouseEvent) {
    if (!startMouseEvent.current) {
      return;
    }

    mouseDown.current = false;
    handleInputUp(
      event.target,
      event.clientX,
      startMouseEvent.current.clientX,
      event.offsetX,
      mouseVelocity.current
    );
  }

  function onTouchEnd(this: Window, event: TouchEvent) {
    if (!startTouchEvent.current) {
      return;
    }

    const offsetX =
      event.changedTouches[0].pageX -
      (event.target as HTMLDivElement).getBoundingClientRect().left;

    touchStart.current = false;
    handleInputUp(
      event.target,
      event.changedTouches[0].clientX,
      startTouchEvent.current.changedTouches[0].clientX,
      offsetX,
      touchVelocity.current
    );
  }

  function handleInputUp(
    target: EventTarget | null,
    inputX: number,
    startX: number,
    offsetX: number,
    velocity: number
  ) {
    if (!scrubberSliderEl.current) {
      return;
    }

    const delta = Math.abs(inputX - startX);
    if (delta < 1 && target instanceof HTMLDivElement) {
      let seekSeconds: number | undefined;

      const spriteIdString = target.getAttribute("data-sprite-item-id");
      if (spriteIdString != null) {
        const spritePercentage = offsetX / target.clientWidth;
        const offset =
          target.offsetLeft + target.clientWidth * spritePercentage;
        const percentage = offset / scrubberSliderEl.current.scrollWidth;
        seekSeconds = percentage * (props.file.duration || 0);
      }

      const markerIdString = target.getAttribute("data-marker-id");
      if (markerIdString != null) {
        const marker = props.scene.scene_markers[Number(markerIdString)];
        seekSeconds = marker.seconds;
      }

      if (seekSeconds) {
        props.onSeek(seekSeconds);
      }
    } else if (Math.abs(velocity) > 25) {
      const newPosition = getPosition() + velocity * 10;
      setPosition(newPosition);
      velocity = 0;
    }
  }

  function onMouseDown(this: HTMLDivElement, event: MouseEvent) {
    event.preventDefault();
    mouseDown.current = true;
    lastMouseEvent.current = event;
    startMouseEvent.current = event;
    mouseVelocity.current = 0;
  }

  function onTouchDown(this: HTMLDivElement, event: TouchEvent) {
    event.preventDefault();
    touchStart.current = true;
    lastTouchEvent.current = event;
    startTouchEvent.current = event;
    touchVelocity.current = 0;
  }

  function onMouseMove(this: HTMLDivElement, event: MouseEvent) {
    if (!mouseDown.current) {
      return;
    }

    // negative dragging right (past), positive left (future)
    const delta = event.clientX - (lastMouseEvent.current?.clientX ?? 0);

    mouseVelocity.current = event.movementX;
    lastMouseEvent.current = event;

    handleInputMove(delta);
  }

  function onTouchMove(this: HTMLDivElement, event: TouchEvent) {
    if (!touchStart.current) {
      return;
    }

    // negative dragging right (past), positive left (future)
    const delta =
      event.changedTouches[0].clientX -
      (lastTouchEvent.current?.changedTouches[0].clientX ?? 0);

    touchVelocity.current = delta;
    lastTouchEvent.current = event;

    const newPostion = getPosition() + delta;
    setPosition(newPostion);
  }

  function handleInputMove(delta: number) {
    const newPostion = getPosition() + delta;
    setPosition(newPostion);
  }

  function getBounds(): number {
    if (!scrubberSliderEl.current || !positionIndicatorEl.current) {
      return 0;
    }
    return (
      scrubberSliderEl.current.scrollWidth -
      scrubberSliderEl.current.clientWidth
    );
  }

  function goBack() {
    if (!scrubberSliderEl.current) {
      return;
    }
    const newPosition = getPosition() + scrubberSliderEl.current.clientWidth;
    setPosition(newPosition);
  }

  function goForward() {
    if (!scrubberSliderEl.current) {
      return;
    }
    const newPosition = getPosition() - scrubberSliderEl.current.clientWidth;
    setPosition(newPosition);
  }

  function renderTags() {
    function getTagStyle(i: number): CSSProperties {
      if (
        !scrubberSliderEl.current ||
        spriteItems.length === 0 ||
        getBounds() === 0
      ) {
        return {};
      }

      const tags = window.document.getElementsByClassName("scrubber-tag");
      if (tags.length === 0) {
        return {};
      }

      let tag: Element | null;
      for (let index = 0; index < tags.length; index++) {
        tag = tags.item(index);
        const id = tag?.getAttribute("data-marker-id") ?? null;
        if (id === i.toString()) {
          break;
        }
      }

      const marker = props.scene.scene_markers[i];
      const duration = Number(props.file.duration);
      const percentage = marker.seconds / duration;

      const left =
        scrubberSliderEl.current.scrollWidth * percentage -
        tag!.clientWidth / 2;
      return {
        left: `${left}px`,
        height: 20,
      };
    }

    return props.scene.scene_markers.map((marker, index) => {
      const dataAttrs = {
        "data-marker-id": index,
      };
      return (
        <div
          key={index}
          className="scrubber-tag"
          style={getTagStyle(index)}
          {...dataAttrs}
        >
          {marker.title || marker.primary_tag.name}
        </div>
      );
    });
  }

  function renderSprites() {
    function getStyleForSprite(index: number): CSSProperties {
      if (!props.scene.paths.vtt) {
        return {};
      }
      const sprite = spriteItems[index];
      const left = sprite.w * index;
      const path = props.scene.paths.vtt.replace("_thumbs.vtt", "_sprite.jpg"); // TODO: Gnarly
      return {
        width: `${sprite.w}px`,
        height: `${sprite.h}px`,
        margin: "0px auto",
        backgroundPosition: `${-sprite.x}px ${-sprite.y}px`,
        backgroundImage: `url(${path})`,
        left: `${left}px`,
      };
    }

    return spriteItems.map((spriteItem, index) => {
      const dataAttrs = {
        "data-sprite-item-id": index,
      };
      return (
        <div
          key={index}
          className="scrubber-item"
          style={getStyleForSprite(index)}
          {...dataAttrs}
        >
          <span className="scrubber-item-time">
            {TextUtils.secondsToTimestamp(spriteItem.start)} -{" "}
            {TextUtils.secondsToTimestamp(spriteItem.end)}
          </span>
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
        &lt;
      </Button>
      <div ref={contentEl} className="scrubber-content">
        <div className="scrubber-tags-background" />
        <div ref={positionIndicatorEl} id="scrubber-position-indicator" />
        <div id="scrubber-current-position" />
        <div className="scrubber-viewport">
          <div ref={scrubberSliderEl} className="scrubber-slider">
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
        &gt;
      </Button>
    </div>
  );
};
