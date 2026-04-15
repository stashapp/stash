import React, { useState, useCallback, useEffect, useRef } from "react";
import { Overlay, Popover, OverlayProps } from "react-bootstrap";
import { PatchComponent } from "src/patch";
import { Icon } from "./Icon";
import { faExclamationTriangle } from "@fortawesome/free-solid-svg-icons";

interface IHoverPopover {
  enterDelay?: number;
  leaveDelay?: number;
  content: JSX.Element[] | JSX.Element | string;
  className?: string;
  placement?: OverlayProps["placement"];
  onOpen?: () => void;
  onClose?: () => void;
  target?: React.RefObject<HTMLElement>;
}

export const HoverPopover: React.FC<IHoverPopover> = PatchComponent(
  "HoverPopover",
  ({
    enterDelay = 200,
    leaveDelay = 200,
    content,
    children,
    className,
    placement = "top",
    onOpen,
    onClose,
    target,
  }) => {
    const [show, setShow] = useState(false);
    const triggerRef = useRef<HTMLDivElement>(null);
    const enterTimer = useRef<number>();
    const leaveTimer = useRef<number>();

    const handleMouseEnter = useCallback(() => {
      window.clearTimeout(leaveTimer.current);
      enterTimer.current = window.setTimeout(() => {
        setShow(true);
        onOpen?.();
      }, enterDelay);
    }, [enterDelay, onOpen]);

    const handleMouseLeave = useCallback(() => {
      window.clearTimeout(enterTimer.current);
      leaveTimer.current = window.setTimeout(() => {
        setShow(false);
        onClose?.();
      }, leaveDelay);
    }, [leaveDelay, onClose]);

    useEffect(
      () => () => {
        window.clearTimeout(enterTimer.current);
        window.clearTimeout(leaveTimer.current);
      },
      []
    );

    return (
      <>
        <div
          className={className}
          onMouseEnter={handleMouseEnter}
          onMouseLeave={handleMouseLeave}
          ref={triggerRef}
        >
          {children}
        </div>
        {triggerRef.current && (
          <Overlay
            show={show}
            placement={placement}
            target={target?.current ?? triggerRef.current}
          >
            <Popover
              onMouseEnter={handleMouseEnter}
              onMouseLeave={handleMouseLeave}
              id="popover"
              className="hover-popover-content"
            >
              {content}
            </Popover>
          </Overlay>
        )}
      </>
    );
  }
);

// convenience component to set the padding on popover content
export const PopoverCard: React.FC<{ className?: string }> = ({
  className,
  children,
}) => {
  return <div className={`popover-card ${className}`}>{children}</div>;
};

export const WarningHoverPopover: React.FC<IHoverPopover> = PatchComponent(
  "WarningHoverPopover",
  ({ children, ...props }) => (
    <HoverPopover {...props} className="warning-hover-popover">
      <Icon icon={faExclamationTriangle} />
    </HoverPopover>
  )
);
