import React, { useState, useEffect, useMemo } from "react";
import { Form, Button } from "react-bootstrap";
import {
  mutateMetadataGenerate,
  useConfigureDefaults,
} from "src/core/StashService";
import { Modal, Icon, OperationButton } from "src/components/Shared";
import { useToast } from "src/hooks";
import * as GQL from "src/core/generated-graphql";
import { FormattedMessage, useIntl } from "react-intl";
import { ConfigurationContext } from "src/hooks/Config";
// import { DirectorySelectionDialog } from "../Settings/SettingsTasksPanel/DirectorySelectionDialog";
import { Manual } from "../Help/Manual";
import { withoutTypename } from "src/utils";
import { GenerateOptions } from "../Tasks/GenerateOptions";

interface ISceneGenerateDialog {
  selectedIds?: string[];
  onClose: () => void;
}

export const GenerateDialog: React.FC<ISceneGenerateDialog> = ({
  selectedIds,
  onClose,
}) => {
  const { configuration } = React.useContext(ConfigurationContext);
  const [configureDefaults] = useConfigureDefaults();

  function getDefaultOptions(): GQL.GenerateMetadataInput {
    return {
      sprites: true,
      phashes: true,
      previews: true,
      markers: true,
      previewOptions: {
        previewSegments: 0,
        previewSegmentDuration: 0,
        previewPreset: GQL.PreviewPreset.Slow,
      },
    };
  }

  const [options, setOptions] = useState<GQL.GenerateMetadataInput>(
    getDefaultOptions()
  );
  const [configRead, setConfigRead] = useState(false);
  const [paths /* , setPaths */] = useState<string[]>([]);
  const [showManual, setShowManual] = useState(false);
  // const [settingPaths, setSettingPaths] = useState(false);
  const [savingDefaults, setSavingDefaults] = useState(false);
  const [animation, setAnimation] = useState(true);

  const intl = useIntl();
  const Toast = useToast();

  useEffect(() => {
    if (configRead) {
      return;
    }

    if (configuration?.defaults.generate) {
      const { generate } = configuration.defaults;
      setOptions(withoutTypename(generate));
      setConfigRead(true);
    } else if (configuration?.general) {
      // backwards compatibility
      const { general } = configuration;
      setOptions((existing) => ({
        ...existing,
        previewOptions: {
          ...existing.previewOptions,
          previewSegments:
            general.previewSegments ?? existing.previewOptions?.previewSegments,
          previewSegmentDuration:
            general.previewSegmentDuration ??
            existing.previewOptions?.previewSegmentDuration,
          previewExcludeStart:
            general.previewExcludeStart ??
            existing.previewOptions?.previewExcludeStart,
          previewExcludeEnd:
            general.previewExcludeEnd ??
            existing.previewOptions?.previewExcludeEnd,
          previewPreset:
            general.previewPreset ?? existing.previewOptions?.previewPreset,
        },
      }));
      setConfigRead(true);
    }
  }, [configuration, configRead]);

  const selectionStatus = useMemo(() => {
    if (selectedIds) {
      return (
        <Form.Group id="selected-generate-ids">
          <FormattedMessage
            id="config.tasks.generate.generating_scenes"
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
        <FormattedMessage id="config.tasks.generate.generating_from_paths" />:
        <ul>
          {paths.map((p) => (
            <li key={p}>{p}</li>
          ))}
        </ul>
      </div>
    ) : (
      <span>
        <FormattedMessage
          id="config.tasks.generate.generating_scenes"
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

    // function onClick() {
    //   setAnimation(false);
    //   setSettingPaths(true);
    // }

    return (
      <Form.Group className="dialog-selected-folders">
        <div>
          {message}
          {/* <div>
            <Button
              title={intl.formatMessage({ id: "actions.select_folders" })}
              onClick={() => onClick()}
            >
              <Icon icon="folder-open" />
            </Button>
          </div> */}
        </div>
      </Form.Group>
    );
  }, [selectedIds, intl, paths]);

  async function onGenerate() {
    try {
      await mutateMetadataGenerate(options);
      Toast.success({
        content: intl.formatMessage(
          { id: "config.tasks.added_job_to_queue" },
          { operation_name: intl.formatMessage({ id: "actions.generate" }) }
        ),
      });
    } catch (e) {
      Toast.error(e);
    } finally {
      onClose();
    }
  }

  function makeDefaultGenerateInput() {
    const ret = options;
    // const { paths: _paths, ...withoutSpecifics } = ret;
    const { overwrite: _overwrite, ...withoutSpecifics } = ret;
    return withoutSpecifics;
  }

  function onShowManual() {
    setAnimation(false);
    setShowManual(true);
  }

  async function setAsDefault() {
    try {
      setSavingDefaults(true);
      await configureDefaults({
        variables: {
          input: {
            generate: makeDefaultGenerateInput(),
          },
        },
      });

      Toast.success({
        content: intl.formatMessage(
          { id: "config.tasks.defaults_set" },
          { action: intl.formatMessage({ id: "actions.generate" }) }
        ),
      });
    } catch (e) {
      Toast.error(e);
    } finally {
      setSavingDefaults(false);
    }
  }

  // if (settingPaths) {
  //   return (
  //     <DirectorySelectionDialog
  //       animation={false}
  //       allowEmpty
  //       initialPaths={paths}
  //       onClose={(p) => {
  //         if (p) {
  //           setPaths(p);
  //         }
  //         setSettingPaths(false);
  //       }}
  //     />
  //   );
  // }

  if (showManual) {
    return (
      <Manual
        animation={false}
        show
        onClose={() => setShowManual(false)}
        defaultActiveTab="Tasks.md"
      />
    );
  }

  return (
    <Modal
      show
      modalProps={{ animation, size: "lg" }}
      icon="cogs"
      header={intl.formatMessage({ id: "actions.generate" })}
      accept={{
        onClick: onGenerate,
        text: intl.formatMessage({ id: "actions.generate" }),
      }}
      cancel={{
        onClick: () => onClose(),
        text: intl.formatMessage({ id: "actions.cancel" }),
        variant: "secondary",
      }}
      disabled={savingDefaults}
      footerButtons={
        <OperationButton variant="secondary" operation={setAsDefault}>
          <FormattedMessage id="actions.set_as_default" />
        </OperationButton>
      }
      leftFooterButtons={
        <Button
          title="Help"
          className="minimal help-button"
          onClick={() => onShowManual()}
        >
          <Icon icon="question-circle" />
        </Button>
      }
    >
      <Form>
        {selectionStatus}
        <GenerateOptions options={options} setOptions={setOptions} />
      </Form>
    </Modal>
  );
};
