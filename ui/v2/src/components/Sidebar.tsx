import {
  MenuItem,
  Menu,
} from "@blueprintjs/core";
import React, { FunctionComponent } from "react";

interface IProps {
  className: string
}

export const Sidebar: FunctionComponent<IProps> = (props) => {
  return (
    <>
      <div className={"sidebar" + props.className}>
        <Menu large={true}>
          <MenuItem
            icon="video"
            text="Scenes"
            href="/scenes"
          />
          <MenuItem
            href="/scenes/markers"
            icon="map-marker"
            text="Markers"
          />
          <MenuItem
            href="/galleries"
            icon="media"
            text="Galleries"
          />
          <MenuItem
            href="/performers"
            icon="person"
            text="Performers"
          />
          <MenuItem
            href="/studios"
            icon="mobile-video"
            text="Studios"
          />
          <MenuItem
            href="/tags"
            icon="tag"
            text="Tags"
          />
        </Menu>
      </div>
    </>
  );
};
