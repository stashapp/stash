import React, { useCallback, useEffect, useRef, useState } from "react";
import { Button, Overlay, Popover } from "react-bootstrap";
import { faTriangleExclamation } from "@fortawesome/free-solid-svg-icons";

import { Icon } from "../Shared/Icon";
import { PerformerCollisionPopoverContent } from "./scenes/TaggerPerformerPopover";

const ENTER_DELAY = 200;
const LEAVE_DELAY = 200;

export const PerformerCollisionWarningButton: React.FC<{
  messageIds: string[];
}> = ({ messageIds }) => {
  const buttonRef = useRef<HTMLButtonElement>(null);
  const [showPopover, setShowPopover] = useState(false);
  const [pinned, setPinned] = useState(false);
  const enterTimer = useRef<number>();
  const leaveTimer = useRef<number>();
  const popoverId = useRef(
    `performer-collision-warning-${Math.random().toString(36).slice(2)}`
  );

  const dismissPopover = useCallback(() => {
    window.clearTimeout(enterTimer.current);
    window.clearTimeout(leaveTimer.current);
    setPinned(false);
    setShowPopover(false);
  }, []);

  const handleMouseEnter = useCallback(() => {
    window.clearTimeout(leaveTimer.current);
    enterTimer.current = window.setTimeout(
      () => setShowPopover(true),
      ENTER_DELAY
    );
  }, []);

  const handleMouseLeave = useCallback(() => {
    if (pinned) return;

    window.clearTimeout(enterTimer.current);
    leaveTimer.current = window.setTimeout(
      () => setShowPopover(false),
      LEAVE_DELAY
    );
  }, [pinned]);

  const handleClick = useCallback(() => {
    window.clearTimeout(enterTimer.current);
    window.clearTimeout(leaveTimer.current);
    setPinned(true);
    setShowPopover(true);
  }, []);

  useEffect(
    () => () => {
      window.clearTimeout(enterTimer.current);
      window.clearTimeout(leaveTimer.current);
    },
    []
  );

  if (messageIds.length === 0) {
    return null;
  }

  return (
    <>
      <Button
        ref={buttonRef}
        variant="warning"
        className="performer-collision-warning-button"
        onClick={handleClick}
        onMouseEnter={handleMouseEnter}
        onMouseLeave={handleMouseLeave}
      >
        <Icon icon={faTriangleExclamation} />
      </Button>
      <Overlay
        show={showPopover}
        placement="bottom-end"
        target={buttonRef}
        container={document.body}
        rootClose={pinned}
        onHide={dismissPopover}
        popperConfig={{
          modifiers: [
            {
              name: "offset",
              options: {
                offset: [0, 6],
              },
            },
          ],
        }}
      >
        <Popover
          id={popoverId.current}
          className="hover-popover-content"
          onMouseEnter={handleMouseEnter}
          onMouseLeave={handleMouseLeave}
        >
          <PerformerCollisionPopoverContent messageIds={messageIds} />
        </Popover>
      </Overlay>
    </>
  );
};
