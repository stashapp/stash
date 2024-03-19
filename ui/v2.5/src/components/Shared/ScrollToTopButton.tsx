import React, { useEffect, useState } from "react";
import { Button } from "react-bootstrap";
import { Icon } from "./Icon";
import { faChevronUp } from "@fortawesome/free-solid-svg-icons";
import cx from "classnames";

export function useScrollTop() {
  const [scrollTop, setScrollTop] = useState(0);

  useEffect(() => {
    const onScroll = () => {
      setScrollTop(window.document.documentElement.scrollTop);
    };

    window.addEventListener("scroll", onScroll);

    return () => {
      window.removeEventListener("scroll", onScroll);
    };
  }, []);

  return scrollTop;
}

// minimum scroll to show the button
const defaultMinScroll = 300;

export const ScrollToTopButton: React.FC<{
  minScroll?: number;
  scrollTop: number;
  onClick: () => void;
}> = ({ minScroll = defaultMinScroll, scrollTop, onClick }) => {
  const [visible, setVisible] = useState(false);

  useEffect(() => {
    if (scrollTop > minScroll) {
      setVisible(true);
    } else {
      setVisible(false);
    }
  }, [scrollTop, minScroll]);

  return (
    <Button
      className={cx("scroll-to-top-button", { show: visible })}
      onClick={onClick}
      size="lg"
      variant="secondary"
    >
      <Icon icon={faChevronUp} />
    </Button>
  );
};
