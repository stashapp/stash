import React, { useState, useEffect, useMemo } from "react";
import { Form, Button } from "react-bootstrap";
import { mutateMetadataGenerate } from "src/core/StashService";
import { ModalComponent } from "../Shared/Modal";
import { Icon } from "src/components/Shared/Icon";
import { useToast } from "src/hooks/Toast";
import * as GQL from "src/core/generated-graphql";
import { FormattedMessage, useIntl } from "react-intl";
import { useConfigurationContext } from "src/hooks/Config";
import { Manual } from "../Help/Manual";
import { withoutTypename } from "src/utils/data";
import { GenerateOptions } from "../Settings/Tasks/GenerateOptions";
import { SettingSection } from "../Settings/SettingSection";
import { faCogs, faQuestionCircle } from "@fortawesome/free-solid-svg-icons";
import { SettingsContext } from "../Settings/context";

interface ISceneGenerateDialog {
  selectedIds?: string[];
  onClose: () => void;
  type: "scene"; // TODO - add image generate
}

export const GenerateDialog: React.FC<ISceneGenerateDialog> = ({
  selectedIds,
  onClose,
  type,
}) => {
  const { configuration } = useConfigurationContext();

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
  const [showManual, setShowManual] = useState(false);
  const [animation, setAnimation] = useState(true);

  const intl = useIntl();
  const Toast = useToast();

  useEffect(() => {
    if (configRead) {
      return;
    }

    // combine the defaults with the system preview generation settings
    if (configuration?.defaults.generate) {
      const { generate } = configuration.defaults;
      setOptions(withoutTypename(generate));
      setConfigRead(true);
    }

    if (configuration?.general) {
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
    const message = (
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

    return (
      <Form.Group className="dialog-selected-folders">
        <div>{message}</div>
      </Form.Group>
    );
  }, [selectedIds, intl]);

  async function onGenerate() {
    try {
      await mutateMetadataGenerate({
        ...options,
        sceneIDs: selectedIds,
      });
      Toast.success(
        intl.formatMessage(
          { id: "config.tasks.added_job_to_queue" },
          { operation_name: intl.formatMessage({ id: "actions.generate" }) }
        )
      );
    } catch (e) {
      Toast.error(e);
    } finally {
      onClose();
    }
  }

  function onShowManual() {
    setAnimation(false);
    setShowManual(true);
  }

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
    <ModalComponent
      show
      modalProps={{ animation, size: "lg" }}
      icon={faCogs}
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
        <SettingsContext>
          <SettingSection>
            <GenerateOptions
              type={type}
              options={options}
              setOptions={setOptions}
              selection
            />
          </SettingSection>
        </SettingsContext>
      </Form>
    </ModalComponent>
  );
};

export default GenerateDialog;
