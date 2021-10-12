import React, { useState, useEffect, useMemo } from "react";
import { Form, Button, ListGroup, Row, Col } from "react-bootstrap";
import {
  mutateMetadataIdentify,
  useConfiguration,
  useListSceneScrapers,
} from "src/core/StashService";
import { Modal, Icon } from "src/components/Shared";
import { useToast } from "src/hooks";
import * as GQL from "src/core/generated-graphql";
import { FormattedMessage, useIntl } from "react-intl";
import { ThreeStateCheckbox } from "../Shared/ThreeStateCheckbox";

interface IScraperSource {
  id: string;
  displayName: string;
  stash_box_endpoint?: string;
  scraper_id?: string;
  options?: GQL.IdentifyMetadataOptionsInput;
}

interface IFieldOptionsEditor {
  availableFields: string[];
  options: GQL.IdentifyFieldOptions;
  editOptions: (o?: GQL.IdentifyFieldOptions) => void;
  removeField: () => void;
  editing: boolean;
}

const FieldOptionsEditor: React.FC<IFieldOptionsEditor> = ({availableFields, options, removeField, editOptions, editing}) => {
  const intl = useIntl();

  const [localOptions, setLocalOptions] = useState(options);

  useEffect(() => {
    setLocalOptions(options);
  }, [options]);

  function renderFieldSelect() {
    return (
      <Form.Group>
        <Form.Label>
          Field
        </Form.Label>
        <Form.Control
          disabled={!editing}
          className="w-auto input-control"
          as="select"
          value={localOptions.field}
          onChange={(e: React.ChangeEvent<HTMLSelectElement>) =>
            setLocalOptions({...localOptions, field: e.currentTarget.value})
          }
        >
          {availableFields.map((f) => (
            <option key={f} value={f}>
              {f}
            </option>
          ))}
        </Form.Control>
      </Form.Group>
    );
  }

  function renderStrategySelect() {
    const strategyStrings = Object.keys(GQL.IdentifyFieldStrategy);

    return (
      <Form.Group>
        <Form.Label>
          Strategy
        </Form.Label>
        <Form.Control
          disabled={!editing}
          className="w-auto input-control"
          as="select"
          value={localOptions.strategy}
          onChange={(e: React.ChangeEvent<HTMLSelectElement>) =>
            setLocalOptions({...localOptions, strategy: e.currentTarget.value as GQL.IdentifyFieldStrategy})
          }
        >
          {strategyStrings.map((f) => (
            <option key={f} value={f}>
              {f}
            </option>
          ))}
        </Form.Control>
      </Form.Group>
    )
  }

  function maybeRenderCreateMissing() {
    const createMissingFields = ["studio", "performers", "tags"];

    if (createMissingFields.includes(localOptions.field)) {
      return (
        <ThreeStateCheckbox
          value={localOptions.createMissing === null ? undefined : localOptions.createMissing}
          setValue={(v) => setLocalOptions({...localOptions, createMissing: v})}
          label={intl.formatMessage({id: "config.tasks.identify.create_missing"})}
        />
      )
    }
  }

  function render() {
    return (
      <Row className="mx-2 align-items-center">
        <Col sm={3}>
          {renderFieldSelect()}
        </Col>
        <Col sm={3}>
          {renderStrategySelect()}
        </Col>
        <Col sm={3}>
          {maybeRenderCreateMissing()}
        </Col>
        
        <div className="col-3 d-flex justify-content-end">
          {editing ? (
            <>
              <Button
                className="minimal text-success"
                onClick={() => editOptions(localOptions)}
              >
                <Icon icon="check" />
              </Button>
              <Button
                className="minimal text-danger"
                onClick={() => editOptions()}
              >
                <Icon icon="times" />
              </Button>
            </>
          ) : (
            <Button
              className="minimal text-danger"
              onClick={() => removeField()}
            >
              <Icon icon="minus" />
            </Button>
          )}
        </div>
      </Row>
    )
  }
  
  return (
    <ListGroup.Item
      as="li"
      key={options.field}
    >
      {render()}
    </ListGroup.Item>
  )
};

interface IFieldOptionsList {
  fieldOptions?: GQL.IdentifyFieldOptions[];
  setFieldOptions: (o: GQL.IdentifyFieldOptions[]) => void;
  setEditingField: (v: boolean) => void;
}

const FieldOptionsList: React.FC<IFieldOptionsList> = ({fieldOptions, setFieldOptions, setEditingField}) => {
  const [localFieldOptions, setLocalFieldOptions] = useState<GQL.IdentifyFieldOptions[]>([]);
  const [editField, setEditField] = useState<GQL.IdentifyFieldOptions | undefined>();

  useEffect(() => {
    if (fieldOptions) {
      setLocalFieldOptions([...fieldOptions]);
    }
  }, [fieldOptions]);

  const allFields = [
    "title",
    "date",
    "details",
    "url",
    "studio",
    "performers",
    "tags",
    "stash_ids",
  ];

  const availableFields = useMemo(() => {
    return allFields.filter(f => !localFieldOptions?.some(o => o !== editField && o.field === f));
  }, [localFieldOptions]);

  function onAdd() {
    const newOptions = [ ...localFieldOptions ];
    const newOption = {
      field: availableFields[0],
      strategy: GQL.IdentifyFieldStrategy.Ignore,
    };
    newOptions.push(newOption);
    setLocalFieldOptions(newOptions);
    setEditField(newOption);
    setEditingField(true);
  }

  function handleEditOptions(o?: GQL.IdentifyFieldOptions) {
    if (!o) {
      if (localFieldOptions.length > (fieldOptions?.length ?? 0)) {
        // must be new field option. remove it
        const newOptions = [ ...localFieldOptions ];
        newOptions.pop();
        setLocalFieldOptions(newOptions);
      }
    } else {
      const newOptions = [ ...localFieldOptions ];
      newOptions.splice(newOptions.indexOf(editField!), 1, o);
      setFieldOptions(newOptions);
    }

    setEditField(undefined);
    setEditingField(false);
  }

  function removeField(index: number) {
    const newOptions = [ ...localFieldOptions ];
    newOptions.splice(index, 1);
    setFieldOptions(newOptions);
  }

  return (
    <Form.Group className="scraper-sources">
      <h5>
        <FormattedMessage id="config.tasks.identify.field_options" />
      </h5>
      <ListGroup as="ul" className="scraper-source-list">
        {localFieldOptions?.map((s, index) => (
          <FieldOptionsEditor 
            availableFields={availableFields}
            options={s} 
            removeField={() => removeField(index)} 
            editOptions={handleEditOptions} 
            editing={s === editField}
          />
        ))}
      </ListGroup>
      {!editField && availableFields.length > 0 ? (
        <div className="text-right">
          <Button
            className="minimal add-scraper-source-button"
            onClick={() => onAdd()}
          >
            <Icon icon="plus" />
          </Button>
        </div>
      ): undefined}
    </Form.Group>
  )
};

interface IOptionsEditor {
  options: GQL.IdentifyMetadataOptionsInput;
  setOptions: (s: GQL.IdentifyMetadataOptionsInput) => void;
  source?: IScraperSource;
  setEditingField: (v: boolean) => void;
}

const OptionsEditor: React.FC<IOptionsEditor> = ({
  options,
  setOptions: setOptionsState,
  source,
  setEditingField,
}) => {
  const intl = useIntl();

  function setOptions(v: Partial<GQL.IdentifyMetadataOptionsInput>) {
    setOptionsState({ ...options, ...v });
  }

  const headingID = !source
    ? "config.tasks.identify.default_options"
    : "config.tasks.identify.source_options";
  const checkboxProps = {
    allowUndefined: !!source,
    indeterminateClassname: "text-muted",
  };

  return (
    <Form.Group>
      <h5>
        <FormattedMessage
          id={headingID}
          values={{ source: source?.displayName }}
        />
      </h5>
      <Form.Group>
        <ThreeStateCheckbox
          value={
            options.includeMalePerformers === null
              ? undefined
              : options.includeMalePerformers
          }
          setValue={(v) =>
            setOptions({
              includeMalePerformers: v,
            })
          }
          label={intl.formatMessage({
            id: "config.tasks.identify.include_male_performers",
          })}
          {...checkboxProps}
        />
        <ThreeStateCheckbox
          value={
            options.setCoverImage === null ? undefined : options.setCoverImage
          }
          setValue={(v) =>
            setOptions({
              setCoverImage: v,
            })
          }
          label={intl.formatMessage({
            id: "config.tasks.identify.set_cover_images",
          })}
          {...checkboxProps}
        />
        <ThreeStateCheckbox
          value={options.setOrganized === null ? undefined : options.setOrganized}
          setValue={(v) =>
            setOptions({
              setOrganized: v,
            })
          }
          label={intl.formatMessage({
            id: "config.tasks.identify.set_organized",
          })}
          {...checkboxProps}
        />
      </Form.Group>

      <FieldOptionsList
        fieldOptions={options.fieldOptions ?? undefined}
        setFieldOptions={(o) => setOptions({ fieldOptions: o })}
        setEditingField={setEditingField}
      />

      {!source && (
        <Form.Text className="text-muted">
          {intl.formatMessage({
            id: "config.tasks.identify.explicit_set_description",
          })}
        </Form.Text>
      )}
    </Form.Group>
  );
};

interface ISourceEditor {
  isNew: boolean;
  availableSources: IScraperSource[];
  source: IScraperSource;
  saveSource: (s?: IScraperSource) => void;
}

const SourcesEditor: React.FC<ISourceEditor> = ({
  isNew,
  availableSources,
  source: initialSource,
  saveSource,
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
    <Modal
      modalProps={{ animation: false }}
      show
      icon={isNew ? "plus" : "pencil-alt"}
      header={intl.formatMessage(
        { id: headerMsgId },
        {
          count: 1,
          singularEntity: source?.displayName,
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
      disabled={(!source.scraper_id && !source.stash_box_endpoint) || editingField}
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
        />
      </Form>
    </Modal>
  );
};

interface ISourcesList {
  sources: IScraperSource[];
  setSources: (s: IScraperSource[]) => void;
  editSource: (s?: IScraperSource) => void;
  canAdd: boolean;
}

const SourcesList: React.FC<ISourcesList> = ({
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
                <Icon icon="grip-vertical" />
              </div>
              {s.displayName}
            </div>
            <div>
              <Button className="minimal" onClick={() => editSource(s)}>
                <Icon icon="cog" />
              </Button>
              <Button
                className="minimal text-danger"
                onClick={() => removeSource(index)}
              >
                <Icon icon="minus" />
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
            <Icon icon="plus" />
          </Button>
        </div>
      )}
    </Form.Group>
  );
};

interface IIdentifyDialogProps {
  selectedIds?: string[];
  onClose: () => void;
}

export const IdentifyDialog: React.FC<IIdentifyDialogProps> = ({
  selectedIds,
  onClose,
}) => {
  function getDefaultOptions(): GQL.IdentifyMetadataOptionsInput {
    return {
      fieldOptions: [],
      includeMalePerformers: true,
      setCoverImage: true,
      setOrganized: false,
    };
  }

  const [options, setOptions] = useState<GQL.IdentifyMetadataOptionsInput>(
    getDefaultOptions()
  );
  const [sources, setSources] = useState<IScraperSource[]>([]);
  const [editingSource, setEditingSource] = useState<
    IScraperSource | undefined
  >();
  const [animation, setAnimation] = useState(true);
  const [editingField, setEditingField] = useState(false);

  const intl = useIntl();
  const Toast = useToast();

  const { data: configData, error: configError } = useConfiguration();
  const { data: scraperData, error: scraperError } = useListSceneScrapers();

  const allSources = useMemo(() => {
    if (!configData || !scraperData) return;

    const defaultSources: IScraperSource[] = [];

    // TODO - use tagger constants

    defaultSources.push(
      ...configData.configuration.general.stashBoxes.map((b, i) => {
        return {
          id: `stash-box: ${i}`,
          displayName: `stash-box: ${b.name}`,
          stash_box_endpoint: b.endpoint,
        };
      })
    );

    const scrapers = scraperData.listSceneScrapers;

    const fragmentScrapers = scrapers.filter((s) =>
      s.scene?.supported_scrapes.includes(GQL.ScrapeType.Fragment)
    );

    // TODO - ensure auto-tag is last when we add auto-tag PR

    defaultSources.push(
      ...fragmentScrapers.map((s) => {
        return {
          id: `scraper: ${s.id}`,
          displayName: s.name,
          scraper_id: s.id,
        };
      })
    );

    return defaultSources;
  }, [configData, scraperData]);

  useEffect(() => {
    if (!allSources) return;

    // set default sources
    setSources(allSources);
  }, [allSources]);

  if (configError || scraperError)
    return <div>{configError ?? scraperError}</div>;
  if (!allSources) return <div></div>;

  async function onIdentify() {
    try {
      await mutateMetadataIdentify({
        sources: sources.map((s) => {
          return {
            source: {
              scraper_id: s.scraper_id,
              stash_box_endpoint: s.stash_box_endpoint,
            },
            options: s.options,
          };
        }),
        options,
      });

      Toast.success({
        content: intl.formatMessage(
          { id: "config.tasks.added_job_to_queue" },
          { operation_name: intl.formatMessage({ id: "actions.identify" }) }
        ),
      });
    } catch (e) {
      Toast.error(e);
    } finally {
      onClose();
    }
  }

  function getAvailableSources() {
    // only include scrapers not already present
    return !editingSource?.id === undefined
      ? []
      : allSources?.filter((s) => {
          return !sources.some((ss) => ss.id === s.id);
        }) ?? [];
  }

  function onEditSource(s?: IScraperSource) {
    setAnimation(false);

    // if undefined, then set a dummy source to create a new one
    if (!s) {
      setEditingSource(getAvailableSources()[0]);
    } else {
      setEditingSource(s);
    }
  }

  function isNewSource() {
    return !!editingSource && !sources.includes(editingSource);
  }

  function onSaveSource(s?: IScraperSource) {
    if (s) {
      let found = false;
      const newSources = sources.map((ss) => {
        if (ss.id === s.id) {
          found = true;
          return s;
        }
        return ss;
      });

      if (!found) {
        newSources.push(s);
      }

      setSources(newSources);
    }
    setEditingSource(undefined);
  }

  if (editingSource) {
    return (
      <SourcesEditor
        availableSources={getAvailableSources()}
        source={editingSource}
        saveSource={onSaveSource}
        isNew={isNewSource()}
      />
    );
  }

  return (
    <Modal
      modalProps={{ animation, size: "lg" }}
      show
      icon="cogs"
      header={intl.formatMessage({ id: "actions.identify" })}
      accept={{
        onClick: onIdentify,
        text: intl.formatMessage({ id: "actions.identify" }),
      }}
      cancel={{
        onClick: () => onClose(),
        text: intl.formatMessage({ id: "actions.cancel" }),
        variant: "secondary",
      }}
      disabled={editingField}
    >
      <Form>
        <SourcesList
          sources={sources}
          setSources={(s) => setSources(s)}
          editSource={onEditSource}
          canAdd={sources.length < allSources.length}
        />
        <OptionsEditor 
          options={options} 
          setOptions={(o) => setOptions(o)} 
          setEditingField={(v) => setEditingField(v)} 
        />
      </Form>
    </Modal>
  );
};

export default IdentifyDialog;
