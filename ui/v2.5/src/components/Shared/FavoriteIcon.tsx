import React from "react";
import { Icon } from "../Shared/Icon";
import { Button } from "react-bootstrap";
import { faHeart } from "@fortawesome/free-solid-svg-icons";
import cx from "classnames";

export const FavoriteIcon: React.FC<{
  favorite: boolean;
  onToggleFavorite: (v: boolean) => void;
}> = ({ favorite, onToggleFavorite }) => {
  return (
    <Button
      className={cx(
        "minimal",
        "mousetrap",
        "favorite-button",
        favorite ? "favorite" : "not-favorite"
      )}
      onClick={() => onToggleFavorite!(!favorite)}
    >
      <Icon icon={faHeart} size="2x" />
    </Button>
  );
};
