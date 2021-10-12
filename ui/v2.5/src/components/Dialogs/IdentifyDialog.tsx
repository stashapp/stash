import React, { useState, useEffect, useMemo } from "react";
import { Form, Button, ListGroup } from "react-bootstrap";
import {
  mutateMetadataIdentify,
  useConfiguration,
  useListSceneScrapers,
} from "src/core/StashService";
import { Modal, Icon } from "src/components/Shared";
import { useIndeterminate, useToast } from "src/hooks";
import * as GQL from "src/core/generated-graphql";
import { FormattedMessage, useIntl } from "react-intl";

interface IScraperSource {
  id: string;
  displayName: string;
  stash_box_endpoint?: string;
  scraper_id?: string;
  options?: GQL.IdentifyMetadataOptionsInput;
}

interface IFieldOptionsEditor {
  options: GQL.IdentifyFieldOptions;
  setOptions: (o: GQL.IdentifyFieldOptions) => void;
}

const FieldOptionsEditor: React.FC<IFieldOptionsEditor> = ({}) => <div></div>;

interface IFieldOptionsList {
  fieldOptions?: GQL.IdentifyFieldOptions[];
  setFieldOptions: (o: GQL.IdentifyFieldOptions[]) => void;
}

const FieldOptionsList: React.FC<IFieldOptionsList> = ({}) => <div></div>;

interface IOptionsEditor {
  options: GQL.IdentifyMetadataOptionsInput;
  setOptions: (s: GQL.IdentifyMetadataOptionsInput) => void;
}

const OptionsEditor: React.FC<IOptionsEditor> = ({
  options,
  setOptions: setOptionsState,
}) => {
  const intl = useIntl();
  const malePerformerRef = React.createRef<HTMLInputElement>();
  const coverImageRef = React.createRef<HTMLInputElement>();
  const organizedRef = React.createRef<HTMLInputElement>();

  function setOptions(v: Partial<GQL.IdentifyMetadataOptionsInput>) {
    setOptionsState({ ...options, ...v });
  }

  useIndeterminate(
    malePerformerRef,
    options.includeMalePerformers ?? undefined
  );
  useIndeterminate(coverImageRef, options.setCoverImage ?? undefined);
  useIndeterminate(organizedRef, options.setOrganized ?? undefined);

  function cycleState(existingState: boolean | undefined) {
    if (existingState) {
      return undefined;
    }
    if (existingState === undefined) {
      return false;
    }
    return true;
  }

  return (
    <Form.Group>
      <h5>
        <FormattedMessage id="config.tasks.identify.default_options" />
      </h5>
      <Form.Check
        ref={malePerformerRef}
        checked={options.includeMalePerformers ?? false}
        label={intl.formatMessage({
          id: "config.tasks.identify.include_male_performers",
        })}
        onChange={() =>
          setOptions({
            includeMalePerformers: cycleState(
              options.includeMalePerformers ?? undefined
            ),
          })
        }
      />
      <Form.Check
        ref={coverImageRef}
        checked={options.setCoverImage ?? false}
        label={intl.formatMessage({
          id: "config.tasks.identify.set_cover_images",
        })}
        onChange={() =>
          setOptions({
            setCoverImage: cycleState(options.setCoverImage ?? undefined),
          })
        }
      />
      <Form.Check
        ref={organizedRef}
        checked={options.setOrganized ?? false}
        label={intl.formatMessage({
          id: "config.tasks.identify.set_organized",
        })}
        onChange={() =>
          setOptions({
            setOrganized: cycleState(options.setOrganized ?? undefined),
          })
        }
      />

      <FieldOptionsList
        fieldOptions={options.fieldOptions ?? undefined}
        setFieldOptions={(o) => setOptions({ fieldOptions: o })}
      />
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

  const intl = useIntl();

  // if id is empty, then we are adding a new source
  const headerMsgId = isNew ? "actions.add" : "actions.edit";
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
      icon="plus"
      header={intl.formatMessage({ id: headerMsgId })}
      accept={{
        onClick: () => saveSource(source),
        text: intl.formatMessage({ id: acceptMsgId }),
      }}
      cancel={{
        onClick: () => saveSource(),
        text: intl.formatMessage({ id: "actions.cancel" }),
        variant: "secondary",
      }}
      disabled={!source.scraper_id && !source.stash_box_endpoint}
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
    event.dataTransfer.setData("text/plain", "");
    event.dataTransfer.setDragImage(new Image(), 0, 0);
    event.dataTransfer.effectAllowed = "move";
    // event.dataTransfer.dropEffect = "move";
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
            className="d-flex justify-content-between"
            draggable={mouseOverIndex === index}
            onDragStart={(e) => onDragStart(e, index)}
            onDragEnter={(e) => onDragOver(e, index)}
            onDrop={() => onDrop()}
          >
            <div>
              <Button
                className="minimal text-muted drag-handle"
                onMouseEnter={() => setMouseOverIndex(index)}
                onMouseLeave={() => setMouseOverIndex(undefined)}
              >
                <Icon icon="grip-lines" />
              </Button>
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
      modalProps={{ animation }}
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
    >
      <Form>
        <SourcesList
          sources={sources}
          setSources={(s) => setSources(s)}
          editSource={onEditSource}
          canAdd={sources.length < allSources.length}
        />
        <OptionsEditor options={options} setOptions={(o) => setOptions(o)} />
      </Form>
    </Modal>
  );
};

export default IdentifyDialog;
