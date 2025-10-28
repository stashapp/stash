import React, { useState, useEffect, useMemo } from "react";
import { Button, Form } from "react-bootstrap";
import {
  mutateMetadataIdentify,
  useConfiguration,
  useConfigureDefaults,
  useListSceneScrapers,
} from "src/core/StashService";
import { Icon } from "src/components/Shared/Icon";
import { ModalComponent } from "src/components/Shared/Modal";
import { OperationButton } from "src/components/Shared/OperationButton";
import { useToast } from "src/hooks/Toast";
import * as GQL from "src/core/generated-graphql";
import { FormattedMessage, useIntl } from "react-intl";
import { withoutTypename } from "src/utils/data";
import {
  SCRAPER_PREFIX,
  STASH_BOX_PREFIX,
} from "src/components/Tagger/constants";
import { DirectorySelectionDialog } from "src/components/Settings/Tasks/DirectorySelectionDialog";
import { Manual } from "src/components/Help/Manual";
import { IScraperSource } from "./constants";
import { OptionsEditor } from "./Options";
import { SourcesEditor, SourcesList } from "./Sources";
import {
  faCogs,
  faFolderOpen,
  faQuestionCircle,
} from "@fortawesome/free-solid-svg-icons";

const autoTagScraperID = "builtin_autotag";

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
      fieldOptions: [
        {
          field: "title",
          strategy: GQL.IdentifyFieldStrategy.Overwrite,
        },
        {
          field: "studio",
          strategy: GQL.IdentifyFieldStrategy.Merge,
          createMissing: true,
        },
        {
          field: "performers",
          strategy: GQL.IdentifyFieldStrategy.Merge,
          createMissing: true,
        },
        {
          field: "tags",
          strategy: GQL.IdentifyFieldStrategy.Merge,
          createMissing: true,
        },
      ],
      includeMalePerformers: true,
      setCoverImage: true,
      setOrganized: false,
      skipMultipleMatches: true,
      skipMultipleMatchTag: undefined,
      skipSingleNamePerformers: true,
      skipSingleNamePerformerTag: undefined,
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
  const [paths, setPaths] = useState<string[]>([]);
  const [showManual, setShowManual] = useState(false);
  const [settingPaths, setSettingPaths] = useState(false);
  const [animation, setAnimation] = useState(true);
  const [editingField, setEditingField] = useState(false);
  const [savingDefaults, setSavingDefaults] = useState(false);

  const intl = useIntl();
  const Toast = useToast();

  const { data: configData, error: configError } = useConfiguration();
  const { data: scraperData, error: scraperError } = useListSceneScrapers();

  const allSources = useMemo(() => {
    if (!configData || !scraperData) return;

    const ret: IScraperSource[] = [];

    ret.push(
      ...configData.configuration.general.stashBoxes.map((b, i) => {
        return {
          id: `${STASH_BOX_PREFIX}${i}`,
          displayName: `stash-box: ${b.name}`,
          stash_box_endpoint: b.endpoint,
        };
      })
    );

    const scrapers = scraperData.listScrapers;

    const fragmentScrapers = scrapers.filter((s) =>
      s.scene?.supported_scrapes.includes(GQL.ScrapeType.Fragment)
    );

    ret.push(
      ...fragmentScrapers.map((s) => {
        return {
          id: `${SCRAPER_PREFIX}${s.id}`,
          displayName: s.name,
          scraper_id: s.id,
        };
      })
    );

    return ret;
  }, [configData, scraperData]);

  const selectionStatus = useMemo(() => {
    if (selectedIds) {
      return (
        <Form.Group id="selected-identify-ids">
          <FormattedMessage
            id="config.tasks.identify.identifying_scenes"
            values={{
              num: selectedIds.length,
              scene: intl.formatMessage(
                {
                  id: "countables.scenes",
                },
                {
                  count: selectedIds.length,
                }
              ),
            }}
          />
          .
        </Form.Group>
      );
    }
    const message = paths.length ? (
      <div>
        <FormattedMessage id="config.tasks.identify.identifying_from_paths" />:
        <ul>
          {paths.map((p) => (
            <li key={p}>{p}</li>
          ))}
        </ul>
      </div>
    ) : (
      <span>
        <FormattedMessage
          id="config.tasks.identify.identifying_scenes"
          values={{
            num: intl.formatMessage({ id: "all" }),
            scene: intl.formatMessage(
              {
                id: "countables.scenes",
              },
              {
                count: 0,
              }
            ),
          }}
        />
        .
      </span>
    );

    function onClick() {
      setAnimation(false);
      setSettingPaths(true);
    }

    return (
      <Form.Group className="dialog-selected-folders">
        <div>
          {message}
          <div>
            <Button
              title={intl.formatMessage({ id: "actions.select_folders" })}
              onClick={() => onClick()}
            >
              <Icon icon={faFolderOpen} />
            </Button>
          </div>
        </div>
      </Form.Group>
    );
  }, [selectedIds, intl, paths]);

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
            sourceOptions.fieldOptions =
              sourceOptions.fieldOptions?.map(withoutTypename);
            ret.options = sourceOptions;
          }

          return ret;
        })
        .filter((s) => s) as IScraperSource[];

      setSources(mappedSources);
      if (identifyDefaults.options) {
        const defaultOptions = withoutTypename(identifyDefaults.options);
        defaultOptions.fieldOptions =
          defaultOptions.fieldOptions?.map(withoutTypename);
        setOptions(defaultOptions);
      }
    } else {
      // default to first stash-box instance only
      const stashBox = allSources.find((s) => s.stash_box_endpoint);

      // add auto-tag as well
      const autoTag = allSources.find(
        (s) => s.id === `${SCRAPER_PREFIX}${autoTagScraperID}`
      );

      const newSources: IScraperSource[] = [];
      if (stashBox) {
        newSources.push(stashBox);
      }

      // sanity check - this should always be true
      if (autoTag) {
        // don't set organised by default
        const autoTagCopy = { ...autoTag };
        autoTagCopy.options = {
          setOrganized: false,
          skipMultipleMatches: true,
          skipSingleNamePerformers: true,
        };
        newSources.push(autoTagCopy);
      }

      setSources(newSources);
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
      sceneIDs: selectedIds,
      paths,
    };
  }

  function makeDefaultIdentifyInput() {
    const ret = makeIdentifyInput();
    const { sceneIDs, paths: _paths, ...withoutSpecifics } = ret;
    return withoutSpecifics;
  }

  async function onIdentify() {
    try {
      await mutateMetadataIdentify(makeIdentifyInput());

      Toast.success(
        intl.formatMessage(
          { id: "config.tasks.added_job_to_queue" },
          { operation_name: intl.formatMessage({ id: "actions.identify" }) }
        )
      );
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

  function onShowManual() {
    setAnimation(false);
    setShowManual(true);
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
            identify: makeDefaultIdentifyInput(),
          },
        },
      });

      Toast.success(
        intl.formatMessage(
          { id: "config.tasks.defaults_set" },
          { action: intl.formatMessage({ id: "actions.identify" }) }
        )
      );
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

  if (settingPaths) {
    return (
      <DirectorySelectionDialog
        animation={false}
        allowEmpty
        initialPaths={paths}
        onClose={(p) => {
          if (p) {
            setPaths(p);
          }
          setSettingPaths(false);
        }}
      />
    );
  }

  if (showManual) {
    return (
      <Manual
        animation={false}
        show
        onClose={() => setShowManual(false)}
        defaultActiveTab="Identify.md"
      />
    );
  }

  return (
    <ModalComponent
      modalProps={{ animation, size: "lg" }}
      show
      icon={faCogs}
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
      disabled={editingField || savingDefaults || sources.length === 0}
      footerButtons={
        <OperationButton
          variant="secondary"
          disabled={editingField || savingDefaults}
          operation={setAsDefault}
        >
          <FormattedMessage id="actions.set_as_default" />
        </OperationButton>
      }
      leftFooterButtons={
        <Button
          title="Help"
          className="minimal help-button"
          onClick={() => onShowManual()}
        >
          <Icon icon={faQuestionCircle} />
        </Button>
      }
    >
      <Form>
        {selectionStatus}
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
    </ModalComponent>
  );
};

export default IdentifyDialog;
