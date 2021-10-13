import React, { useState, useEffect, useMemo } from "react";
import { Button, Form, Spinner } from "react-bootstrap";
import {
  mutateMetadataIdentify,
  useConfiguration,
  useConfigureDefaults,
  useListSceneScrapers,
} from "src/core/StashService";
import { Modal } from "src/components/Shared";
import { useToast } from "src/hooks";
import * as GQL from "src/core/generated-graphql";
import { FormattedMessage, useIntl } from "react-intl";
import { withoutTypename } from "src/utils";
import { IScraperSource } from "./constants";
import { OptionsEditor } from "./Options";
import { SourcesEditor, SourcesList } from "./Sources";

interface IIdentifyDialogProps {
  selectedIds?: string[];
  onClose: () => void;
}

export const IdentifyDialog: React.FC<IIdentifyDialogProps> = ({
  // selectedIds,
  onClose,
}) => {
  function getDefaultOptions(): GQL.IdentifyMetadataOptionsInput {
    return {
      fieldOptions: [
        {
          field: "title",
          strategy: GQL.IdentifyFieldStrategy.Overwrite,
        },
      ],
      includeMalePerformers: true,
      setCoverImage: true,
      setOrganized: false,
    };
  }

  const [configureDefaults] = useConfigureDefaults();

  const [options, setOptions] = useState<GQL.IdentifyMetadataOptionsInput>(
    getDefaultOptions()
  );
  const [sources, setSources] = useState<IScraperSource[]>([]);
  const [editingSource, setEditingSource] = useState<
    IScraperSource | undefined
  >();
  const [animation, setAnimation] = useState(true);
  const [editingField, setEditingField] = useState(false);
  const [savingDefaults, setSavingDefaults] = useState(false);

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
    if (!configData || !allSources) return;

    const { identify: identifyDefaults } = configData.configuration.defaults;

    if (identifyDefaults) {
      const mappedSources = identifyDefaults.sources
        .map((s) => {
          const found = allSources.find(
            (ss) =>
              ss.scraper_id === s.source.scraper_id ||
              ss.stash_box_endpoint === s.source.stash_box_endpoint
          );

          if (!found) return;

          const ret: IScraperSource = {
            ...found,
          };

          if (s.options) {
            const sourceOptions = withoutTypename(s.options);
            sourceOptions.fieldOptions = sourceOptions.fieldOptions?.map(
              withoutTypename
            );
            ret.options = sourceOptions;
          }

          return ret;
        })
        .filter((s) => s) as IScraperSource[];

      setSources(mappedSources);
      if (identifyDefaults.options) {
        const defaultOptions = withoutTypename(identifyDefaults.options);
        defaultOptions.fieldOptions = defaultOptions.fieldOptions?.map(
          withoutTypename
        );
        setOptions(defaultOptions);
      }
    } else {
      // default to first stash-box instance only
      const stashBox = allSources.find((s) => s.stash_box_endpoint);
      if (stashBox) {
        setSources([stashBox]);
      } else {
        setSources([]);
      }
    }
  }, [allSources, configData]);

  if (configError || scraperError)
    return <div>{configError ?? scraperError}</div>;
  if (!allSources || !configData) return <div />;

  function makeIdentifyInput(): GQL.IdentifyMetadataInput {
    return {
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
    };
  }

  async function onIdentify() {
    try {
      await mutateMetadataIdentify(makeIdentifyInput());

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

  async function setAsDefault() {
    try {
      setSavingDefaults(true);
      await configureDefaults({
        variables: {
          input: {
            identify: makeIdentifyInput(),
          },
        },
      });
    } catch (e) {
      Toast.error(e);
    } finally {
      setSavingDefaults(false);
    }
  }

  if (editingSource) {
    return (
      <SourcesEditor
        availableSources={getAvailableSources()}
        source={editingSource}
        saveSource={onSaveSource}
        isNew={isNewSource()}
        defaultOptions={options}
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
      disabled={editingField || savingDefaults}
      footerButtons={
        <Button
          variant="secondary"
          disabled={editingField || savingDefaults}
          onClick={() => setAsDefault()}
        >
          {savingDefaults && (
            <Spinner animation="border" role="status" size="sm" />
          )}
          <FormattedMessage id="actions.set_as_default" />
        </Button>
      }
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
