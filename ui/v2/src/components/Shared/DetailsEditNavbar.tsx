import {
  Button,
  FileInput,
  Menu,
  MenuItem,
  Navbar,
  NavbarDivider,
  Popover,
} from "@blueprintjs/core";
import _ from "lodash";
import React, { FunctionComponent } from "react";
import { Link } from "react-router-dom";
import * as GQL from "../../core/generated-graphql";
import { NavigationUtils } from "../../utils/navigation";

interface IProps {
  performer?: Partial<GQL.PerformerDataFragment>;
  studio?: Partial<GQL.StudioDataFragment>;
  isNew: boolean;
  isEditing: boolean;
  onToggleEdit: () => void;
  onSave: () => void;
  onDelete: () => void;
  onImageChange: (event: React.FormEvent<HTMLInputElement>) => void;

  // TODO: only for performers.  make generic
  onDisplayFreeOnesDialog?: () => void;
}

export const DetailsEditNavbar: FunctionComponent<IProps> = (props: IProps) => {
  function renderEditButton() {
    if (props.isNew) { return; }
    return (
      <Button
        intent="primary"
        text={props.isEditing ? "Cancel" : "Edit"}
        onClick={() => props.onToggleEdit()}
      />
    );
  }

  function renderSaveButton() {
    if (!props.isEditing) { return; }
    return <Button intent="success" text="Save" onClick={() => props.onSave()} />;
  }

  function renderDeleteButton() {
    if (props.isNew || props.isEditing) { return; }
    return <Button intent="danger" text="Delete" onClick={() => props.onDelete()} />;
  }

  function renderImageInput() {
    if (!props.isEditing) { return; }
    return <FileInput text="Choose image..." onInputChange={props.onImageChange} inputProps={{accept: ".jpg,.jpeg"}} />;
  }

  function renderScraperMenu() {
    if (!props.performer) { return; }
    if (!props.isEditing) { return; }
    const scraperMenu = (
      <Menu>
        <MenuItem
          text="FreeOnes"
          onClick={() => { if (props.onDisplayFreeOnesDialog) { props.onDisplayFreeOnesDialog(); }}}
        />
      </Menu>
    );
    return (
      <Popover content={scraperMenu} position="bottom">
        <Button text="Scrape with..."/>
      </Popover>
    );
  }

  function renderScenesButton() {
    if (props.isEditing) { return; }
    let linkSrc: string = "#";
    if (!!props.performer) {
      linkSrc = NavigationUtils.makePerformerScenesUrl(props.performer);
    } else if (!!props.studio) {
      linkSrc = NavigationUtils.makeStudioScenesUrl(props.studio);
    }
    return (
      <Link className="bp3-button" to={linkSrc}>
        Scenes
      </Link>
    );
  }

  return (
    <Navbar>
      <Navbar.Group>
        {renderEditButton()}
        {props.isEditing && !props.isNew ? <NavbarDivider /> : undefined}
        {renderScraperMenu()}
        {renderImageInput()}
        {renderSaveButton()}

        {renderScenesButton()}
        {renderDeleteButton()}
      </Navbar.Group>
    </Navbar>
  );
};
