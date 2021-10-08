import React, { useState, useEffect } from "react";
import { Form, Button, ListGroup } from "react-bootstrap";
import {
  mutateMetadataIdentify,
  useConfiguration,
  useListSceneScrapers,
} from "src/core/StashService";
import { Modal, Icon } from "src/components/Shared";
import { useToast } from "src/hooks";
import * as GQL from "src/core/generated-graphql";
import { FormattedMessage, useIntl } from "react-intl";

interface IScraperSource {
  id: string;
  displayName: string;
  stash_box_endpoint?: string;
  scraper_id?: string;
  options?: GQL.IdentifyMetadataOptionsInput;
}

interface ISourcesList {
  sources: IScraperSource[];
  setSources: (s: IScraperSource[]) => void;
}

const SourcesList: React.FC<ISourcesList> = ({ sources, setSources }) => {
  function removeSource(index: number) {
    const newSources = [...sources];
    newSources.splice(index, 1);
    setSources(newSources);
  }

  return (
    <Form.Group className="scraper-sources">
      <h5>
        <FormattedMessage id="config.tasks.identify.sources" />
      </h5>
      <ListGroup as="ul" className="scraper-source-list">
        {sources.map((s, index) => (
          <ListGroup.Item
            as="li"
            key={s.id}
            className="d-flex justify-content-between"
          >
            <div>
              <Button className="minimal text-muted drag-handle">
                <Icon icon="grip-lines" />
              </Button>
              {s.displayName}
            </div>
            <div>
              {/* <Button className="minimal"><Icon icon="cog" /></Button> */}
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
      {/* <div className="text-right">
    <Button className="minimal add-scraper-source-button" onClick={() => {}}>
      <Icon icon="plus" />
    </Button>
    </div> */}
    </Form.Group>
  );
};

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

  function setOptions(v: Partial<GQL.IdentifyMetadataOptionsInput>) {
    setOptionsState({ ...options, ...v });
  }

  return (
    <Form.Group>
      <h5>
        <FormattedMessage id="config.tasks.identify.default_options" />
      </h5>
      <Form.Check
        checked={options.includeMalePerformers ?? false}
        label={intl.formatMessage({
          id: "config.tasks.identify.include_male_performers",
        })}
        onChange={() =>
          setOptions({ includeMalePerformers: !options.includeMalePerformers })
        }
      />
      <Form.Check
        checked={options.setCoverImage ?? false}
        label={intl.formatMessage({
          id: "config.tasks.identify.set_cover_images",
        })}
        onChange={() => setOptions({ setCoverImage: !options.setCoverImage })}
      />
      <Form.Check
        checked={options.setOrganized ?? false}
        label={intl.formatMessage({
          id: "config.tasks.identify.set_organized",
        })}
        onChange={() => setOptions({ setOrganized: !options.setOrganized })}
      />

      <FieldOptionsList
        fieldOptions={options.fieldOptions ?? undefined}
        setFieldOptions={(o) => setOptions({ fieldOptions: o })}
      />
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
  const intl = useIntl();
  const Toast = useToast();

  const { data: configData, error: configError } = useConfiguration();
  const { data: scraperData, error: scraperError } = useListSceneScrapers();

  useEffect(() => {
    if (!configData || !scraperData) return;

    // set default sources
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

    setSources(defaultSources);
  }, [configData, scraperData]);

  if (configError || scraperError)
    return <div>{configError ?? scraperError}</div>;
  if (!configData || !scraperData) return <div></div>;

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

  return (
    <Modal
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
        <SourcesList sources={sources} setSources={(s) => setSources(s)} />
        <OptionsEditor options={options} setOptions={(o) => setOptions(o)} />
      </Form>
    </Modal>
  );
};

export default IdentifyDialog;
