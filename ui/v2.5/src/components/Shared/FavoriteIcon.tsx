import React from "react";
import { Icon } from "../Shared/Icon";
import { Button } from "react-bootstrap";
import { faHeart } from "@fortawesome/free-solid-svg-icons";
import cx from "classnames";
import { SizeProp } from "@fortawesome/fontawesome-svg-core";

export const FavoriteIcon: React.FC<{
  favorite: boolean;
  onToggleFavorite: (v: boolean) => void;
  size?: SizeProp;
  className?: string;
}> = ({ favorite, onToggleFavorite, size, className }) => {
  return (
    <Button
      className={cx(
        "minimal",
        "mousetrap",
        "favorite-button",
        className,
        favorite ? "favorite" : "not-favorite"
      )}
      onClick={() => onToggleFavorite!(!favorite)}
    >
      <Icon icon={faHeart} size={size} />
    </Button>
  );
};
