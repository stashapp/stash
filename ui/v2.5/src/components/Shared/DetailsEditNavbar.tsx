import { Button, Form, Modal, Nav, Navbar, OverlayTrigger, Popover } from 'react-bootstrap';
import React, { useState } from "react";
import { Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import { NavUtils } from "src/utils";

interface IProps {
  performer?: Partial<GQL.PerformerDataFragment>;
  studio?: Partial<GQL.StudioDataFragment>;
  isNew: boolean;
  isEditing: boolean;
  onToggleEdit: () => void;
  onSave: () => void;
  onDelete: () => void;
  onAutoTag?: () => void;
  onImageChange: (event: React.FormEvent<HTMLInputElement>) => void;

  // TODO: only for performers.  make generic
  scrapers?: GQL.ListPerformerScrapersListPerformerScrapers[];
  onDisplayScraperDialog?: (scraper: GQL.ListPerformerScrapersListPerformerScrapers) => void;
}

export const DetailsEditNavbar: React.FC<IProps> = (props: IProps) => {
  const [isDeleteAlertOpen, setIsDeleteAlertOpen] = useState<boolean>(false);

  function renderEditButton() {
    if (props.isNew) { return; }
    return (
      <Button
        variant="primary"
        onClick={() => props.onToggleEdit()}
      >
        { props.isEditing ? "Cancel" : "Edit"}
      </Button>
    );
  }

  function renderSaveButton() {
    if (!props.isEditing) { return; }
    return <Button variant="success" onClick={() => props.onSave()}>Save</Button>;
  }

  function renderDeleteButton() {
    if (props.isNew || props.isEditing) { return; }
    return <Button variant="danger" onClick={() => setIsDeleteAlertOpen(true)}>Delete</Button>;
  }

  function renderImageInput() {
    if (!props.isEditing) { return; }
      return (
        <Form.Group controlId="cover-file">
          <Form.Label>Choose image...</Form.Label>
          <Form.Control type="file" accept=".jpg,.jpeg,.png" onChange={props.onImageChange} />
        </Form.Group>
      )
  }

  function renderScraperMenu() {
    if (!props.performer) { return; }
    if (!props.isEditing) { return; }

    const popover = (
      <Popover id="scraper-popover">
        <Popover.Content>
          <div>
            { props.scrapers ? props.scrapers.map((s) => (
              <div onClick={() => props.onDisplayScraperDialog &&  props.onDisplayScraperDialog(s) }>
                {s.name}
              </div>
            )) : ''}
          </div>
        </Popover.Content>
      </Popover>
    );

    return (
      <OverlayTrigger trigger="click" placement="bottom" overlay={popover}>
        <Button>Scrape with...</Button>
      </OverlayTrigger>
    );
  }

  function renderAutoTagButton() {
    if (props.isNew || props.isEditing) { return; }
    if (!!props.onAutoTag) {
      return (<Button onClick={() => {
        if (props.onAutoTag) { props.onAutoTag() }
      }}>Auto Tag</Button>)
    }
  }

  function renderScenesButton() {
    if (props.isEditing) { return; }
    let linkSrc: string = "#";
    if (props.performer) {
      linkSrc = NavUtils.makePerformerScenesUrl(props.performer);
    } else if (props.studio) {
      linkSrc = NavUtils.makeStudioScenesUrl(props.studio);
    }
    return (
      <Link to={linkSrc}>
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

    return (
      <Modal
        show={isDeleteAlertOpen}
      >
        <Modal.Body>
          Are you sure you want to delete {name}?
        </Modal.Body>
        <Modal.Footer>
          <Button variant="danger" onClick={props.onDelete}>Delete</Button>
          <Button variant="secondary" onClick={() => setIsDeleteAlertOpen(false)}>Cancel</Button>
        </Modal.Footer>
      </Modal>
    );
  }


  return (
    <>
    {renderDeleteAlert()}
    <Navbar bg="dark">
      <Nav className="mr-auto ml-auto">
          {renderEditButton()}
          {renderScraperMenu()}
          {renderImageInput()}
          {renderSaveButton()}

          {renderAutoTagButton()}
          {renderScenesButton()}
          {renderDeleteButton()}
      </Nav>
    </Navbar>
    </>
  );
};
