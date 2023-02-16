import React, { useState, useEffect } from "react";
import { Form, Button, ListGroup } from "react-bootstrap";
import { ModalComponent } from "src/components/Shared/Modal";
import { Icon } from "src/components/Shared/Icon";
import { FormattedMessage, useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import { IScraperSource } from "./constants";
import { OptionsEditor } from "./Options";
import {
  faCog,
  faGripVertical,
  faMinus,
  faPencilAlt,
  faPlus,
} from "@fortawesome/free-solid-svg-icons";

interface ISourceEditor {
  isNew: boolean;
  availableSources: IScraperSource[];
  source: IScraperSource;
  saveSource: (s?: IScraperSource) => void;
  defaultOptions: GQL.IdentifyMetadataOptionsInput;
}

export const SourcesEditor: React.FC<ISourceEditor> = ({
  isNew,
  availableSources,
  source: initialSource,
  saveSource,
  defaultOptions,
}) => {
  const [source, setSource] = useState<IScraperSource>(initialSource);
  const [editingField, setEditingField] = useState(false);

  const intl = useIntl();

  // if id is empty, then we are adding a new source
  const headerMsgId = isNew ? "actions.add" : "dialogs.edit_entity_title";
  const acceptMsgId = isNew ? "actions.add" : "actions.confirm";

  function handleSourceSelect(e: React.ChangeEvent<HTMLSelectElement>) {
    const selectedSource = availableSources.find(
      (s) => s.id === e.currentTarget.value
    );
    if (!selectedSource) return;

    setSource({
      ...source,
      id: selectedSource.id,
      displayName: selectedSource.displayName,
      scraper_id: selectedSource.scraper_id,
      stash_box_endpoint: selectedSource.stash_box_endpoint,
    });
  }

  return (
    <ModalComponent
      dialogClassName="identify-source-editor"
      modalProps={{ animation: false, size: "lg" }}
      show
      icon={isNew ? faPlus : faPencilAlt}
      header={intl.formatMessage(
        { id: headerMsgId },
        {
          count: 1,
          singularEntity: source?.displayName,
          pluralEntity: source?.displayName,
        }
      )}
      accept={{
        onClick: () => saveSource(source),
        text: intl.formatMessage({ id: acceptMsgId }),
      }}
      cancel={{
        onClick: () => saveSource(),
        text: intl.formatMessage({ id: "actions.cancel" }),
        variant: "secondary",
      }}
      disabled={
        (!source.scraper_id && !source.stash_box_endpoint) || editingField
      }
    >
      <Form>
        {isNew && (
          <Form.Group>
            <h5>
              <FormattedMessage id="config.tasks.identify.source" />
            </h5>
            <Form.Control
              as="select"
              value={source.id}
              className="input-control"
              onChange={handleSourceSelect}
            >
              {availableSources.map((i) => (
                <option value={i.id} key={i.id}>
                  {i.displayName}
                </option>
              ))}
            </Form.Control>
          </Form.Group>
        )}
        <OptionsEditor
          options={source.options ?? {}}
          setOptions={(o) => setSource({ ...source, options: o })}
          source={source}
          setEditingField={(v) => setEditingField(v)}
          defaultOptions={defaultOptions}
        />
      </Form>
    </ModalComponent>
  );
};

interface ISourcesList {
  sources: IScraperSource[];
  setSources: (s: IScraperSource[]) => void;
  editSource: (s?: IScraperSource) => void;
  canAdd: boolean;
}

export const SourcesList: React.FC<ISourcesList> = ({
  sources,
  setSources,
  editSource,
  canAdd,
}) => {
  const [tempSources, setTempSources] = useState(sources);
  const [dragIndex, setDragIndex] = useState<number | undefined>();
  const [mouseOverIndex, setMouseOverIndex] = useState<number | undefined>();

  useEffect(() => {
    setTempSources([...sources]);
  }, [sources]);

  function removeSource(index: number) {
    const newSources = [...sources];
    newSources.splice(index, 1);
    setSources(newSources);
  }

  function onDragStart(event: React.DragEvent<HTMLElement>, index: number) {
    event.dataTransfer.effectAllowed = "move";
    setDragIndex(index);
  }

  function onDragOver(event: React.DragEvent<HTMLElement>, index?: number) {
    if (dragIndex !== undefined && index !== undefined && index !== dragIndex) {
      const newSources = [...tempSources];
      const moved = newSources.splice(dragIndex, 1);
      newSources.splice(index, 0, moved[0]);
      setTempSources(newSources);
      setDragIndex(index);
    }

    event.dataTransfer.dropEffect = "move";
    event.preventDefault();
  }

  function onDragOverDefault(event: React.DragEvent<HTMLDivElement>) {
    event.dataTransfer.dropEffect = "move";
    event.preventDefault();
  }

  function onDrop() {
    // assume we've already set the temp source list
    // feed it up
    setSources(tempSources);
    setDragIndex(undefined);
    setMouseOverIndex(undefined);
  }

  return (
    <Form.Group className="scraper-sources" onDragOver={onDragOverDefault}>
      <h5>
        <FormattedMessage id="config.tasks.identify.sources" />
      </h5>
      <ListGroup as="ul" className="scraper-source-list">
        {tempSources.map((s, index) => (
          <ListGroup.Item
            as="li"
            key={s.id}
            className="d-flex justify-content-between align-items-center"
            draggable={mouseOverIndex === index}
            onDragStart={(e) => onDragStart(e, index)}
            onDragEnter={(e) => onDragOver(e, index)}
            onDrop={() => onDrop()}
          >
            <div>
              <div
                className="minimal text-muted drag-handle"
                onMouseEnter={() => setMouseOverIndex(index)}
                onMouseLeave={() => setMouseOverIndex(undefined)}
              >
                <Icon icon={faGripVertical} />
              </div>
              {s.displayName}
            </div>
            <div>
              <Button className="minimal" onClick={() => editSource(s)}>
                <Icon icon={faCog} />
              </Button>
              <Button
                className="minimal text-danger"
                onClick={() => removeSource(index)}
              >
                <Icon icon={faMinus} />
              </Button>
            </div>
          </ListGroup.Item>
        ))}
      </ListGroup>
      {canAdd && (
        <div className="text-right">
          <Button
            className="minimal add-scraper-source-button"
            onClick={() => editSource()}
          >
            <Icon icon={faPlus} />
          </Button>
        </div>
      )}
    </Form.Group>
  );
};
