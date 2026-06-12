import React, { useCallback, useEffect, useRef, useState } from "react";
import { useIntl } from "react-intl";
import { Overlay, Popover } from "react-bootstrap";
import { faLink, faTriangleExclamation } from "@fortawesome/free-solid-svg-icons";

import { OperationButton } from "../Shared/OperationButton";
import { Icon } from "../Shared/Icon";
import { PerformerCollisionPopoverContent } from "./scenes/TaggerPerformerPopover";

const ENTER_DELAY = 200;
const LEAVE_DELAY = 200;

export const LinkButton: React.FC<{
  disabled: boolean;
  onLink: () => Promise<void>;
  collisionMessageIds?: string[];
}> = ({ disabled, onLink, collisionMessageIds = [] }) => {
  const intl = useIntl();
  const buttonRef = useRef<HTMLButtonElement>(null);
  const [showPopover, setShowPopover] = useState(false);
  const enterTimer = useRef<number>();
  const leaveTimer = useRef<number>();
  const popoverId = useRef(
    `link-button-collision-${Math.random().toString(36).slice(2)}`
  );

  const hasCollision = collisionMessageIds.length > 0;

  const handleMouseEnter = useCallback(() => {
    if (!hasCollision) return;

    window.clearTimeout(leaveTimer.current);
    enterTimer.current = window.setTimeout(
      () => setShowPopover(true),
      ENTER_DELAY
    );
  }, [hasCollision]);

  const handleMouseLeave = useCallback(() => {
    if (!hasCollision) return;

    window.clearTimeout(enterTimer.current);
    leaveTimer.current = window.setTimeout(
      () => setShowPopover(false),
      LEAVE_DELAY
    );
  }, [hasCollision]);

  useEffect(
    () => () => {
      window.clearTimeout(enterTimer.current);
      window.clearTimeout(leaveTimer.current);
    },
    []
  );

  return (
    <>
      <OperationButton
        ref={buttonRef}
        variant="secondary"
        disabled={disabled}
        operation={onLink}
        hideChildrenWhenLoading
        title={intl.formatMessage({ id: "component_tagger.verb_link_existing" })}
        onMouseEnter={handleMouseEnter}
        onMouseLeave={handleMouseLeave}
      >
        <Icon
          icon={hasCollision ? faTriangleExclamation : faLink}
          className={hasCollision ? "text-warning" : undefined}
        />
      </OperationButton>
      {hasCollision && (
        <Overlay
          show={showPopover}
          placement="bottom-end"
          target={buttonRef}
          container={document.body}
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
            <PerformerCollisionPopoverContent
              messageIds={collisionMessageIds}
            />
          </Popover>
        </Overlay>
      )}
    </>
  );
};
