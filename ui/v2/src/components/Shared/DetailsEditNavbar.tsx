import {
  Alert,
  Button,
  FileInput,
  Menu,
  MenuItem,
  Navbar,
  NavbarDivider,
  Popover,
} from "@blueprintjs/core";
import React, { FunctionComponent, useState } from "react";
import { Link } from "react-router-dom";
import * as GQL from "../../core/generated-graphql";
import { NavigationUtils } from "../../utils/navigation";

interface IProps {
  performer?: Partial<GQL.PerformerDataFragment>;
  studio?: Partial<GQL.StudioDataFragment>;
  movie?: Partial<GQL.MovieDataFragment>;
  isNew: boolean;
  isEditing: boolean;
  onToggleEdit: () => void;
  onSave: () => void;
  onDelete: () => void;
  onAutoTag?: () => void;
  onImageChange: (event: React.FormEvent<HTMLInputElement>) => void;
  onBackImageChange?: (event: React.FormEvent<HTMLInputElement>) => void;

  // TODO: only for performers.  make generic
  scrapers?: GQL.ListPerformerScrapersListPerformerScrapers[];
  onDisplayScraperDialog?: (scraper: GQL.ListPerformerScrapersListPerformerScrapers) => void;
}

export const DetailsEditNavbar: FunctionComponent<IProps> = (props: IProps) => {
  const [isDeleteAlertOpen, setIsDeleteAlertOpen] = useState<boolean>(false);

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
    return <Button intent="danger" text="Delete" onClick={() => setIsDeleteAlertOpen(true)} />;
  }

  function renderImageInput() {
    if (!props.isEditing) { return; }
    return <FileInput text="Choose image..." onInputChange={props.onImageChange} inputProps={{accept: ".jpg,.jpeg"}} />;
  }
   
  function renderBackImageInput() {
    if (!props.movie) { return; }
    if (!props.isEditing) { return; }
    return <FileInput text="Choose back image..." onInputChange={props.onBackImageChange} inputProps={{accept: ".jpg,.jpeg"}} />;
  }

  function renderScraperMenuItem(scraper : GQL.ListPerformerScrapersListPerformerScrapers) {
    return (
      <MenuItem
        text={scraper.name}
        onClick={() => { if (props.onDisplayScraperDialog) { props.onDisplayScraperDialog(scraper); }}}
      />
    );
  }

  function renderScraperMenu() {
    if (!props.performer) { return; }
    if (!props.isEditing) { return; }
    const scraperMenu = (
      <Menu>
        {props.scrapers ? props.scrapers.map((s) => renderScraperMenuItem(s)) : undefined}
      </Menu>
    );
    return (
      <Popover content={scraperMenu} position="bottom">
        <Button text="Scrape with..."/>
      </Popover>
    );
  }

  function renderAutoTagButton() {
    if (props.isNew || props.isEditing) { return; }
    if (!!props.onAutoTag) {
      return (<Button text="Auto Tag" onClick={() => {
        if (props.onAutoTag) { props.onAutoTag() }
      }}></Button>)
    }
  }

  function renderScenesButton() {
    if (props.isEditing) { return; }
    let linkSrc: string = "#";
    if (!!props.performer) {
      linkSrc = NavigationUtils.makePerformerScenesUrl(props.performer);
    } else if (!!props.studio) {
      linkSrc = NavigationUtils.makeStudioScenesUrl(props.studio);
    } else if (!!props.movie) {
      linkSrc = NavigationUtils.makeMovieScenesUrl(props.movie);
    }
    return (
      <Link className="bp3-button" to={linkSrc}>
        Scenes
      </Link>
    );
  }

  function renderDeleteAlert() {
    var name;

    if (props.performer) {
      name = props.performer.name;
    }
    if (props.studio) {
      name = props.studio.name;
    }
    if (props.movie) {
      name = props.movie.name;
    }

    return (
      <Alert
        cancelButtonText="Cancel"
        confirmButtonText="Delete"
        icon="trash"
        intent="danger"
        isOpen={isDeleteAlertOpen}
        onCancel={() => setIsDeleteAlertOpen(false)}
        onConfirm={() => props.onDelete()}
      >
        <p>
          Are you sure you want to delete {name}?
        </p>
      </Alert>
    );
  }


  return (
    <>
    {renderDeleteAlert()}
    <Navbar>
      <Navbar.Group>
        {renderEditButton()}
        {props.isEditing && !props.isNew ? <NavbarDivider /> : undefined}
        {renderScraperMenu()}
        {renderImageInput()}
        {renderBackImageInput()} 
        {renderSaveButton()}

        {renderAutoTagButton()}
        {renderScenesButton()}
        {renderDeleteButton()}
      </Navbar.Group>
    </Navbar>
    </>
  );
};
