import {
  MenuItem,
  Menu,
} from "@blueprintjs/core";
import React, { FunctionComponent } from "react";
import { IMenuItem } from "../App";

interface IProps {
  className: string
  menuItems: IMenuItem[]
}

export const Sidebar: FunctionComponent<IProps> = (props) => {
  return (
    <>
      <div className={"sidebar" + props.className}>
        <Menu large={true}>
          {props.menuItems.map((i) => {
            return (
              <MenuItem
                icon={i.icon}
                text={i.text}
                href={i.href}
              />
            )
          })}
        </Menu>
      </div>
    </>
  );
};
