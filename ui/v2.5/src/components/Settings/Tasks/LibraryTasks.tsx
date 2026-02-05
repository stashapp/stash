import React, { useState, useEffect } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { Button, Form } from "react-bootstrap";
import {
  mutateMetadataScan,
  mutateMetadataAutoTag,
  mutateMetadataGenerate,
} from "src/core/StashService";
import { withoutTypename } from "src/utils/data";
import { useConfigurationContext } from "src/hooks/Config";
import { IdentifyDialog } from "../../Dialogs/IdentifyDialog/IdentifyDialog";
import * as GQL from "src/core/generated-graphql";
import { DirectorySelectionDialog } from "./DirectorySelectionDialog";
import { ScanOptions } from "./ScanOptions";
import { useToast } from "src/hooks/Toast";
import { GenerateOptions } from "./GenerateOptions";
import { SettingSection } from "../SettingSection";
import { BooleanSetting, Setting, SettingGroup } from "../Inputs";
import { ManualLink } from "src/components/Help/context";
import { Icon } from "src/components/Shared/Icon";
import { faQuestionCircle } from "@fortawesome/free-solid-svg-icons";
import { useSettings } from "../context";

interface IAutoTagOptions {
  options: GQL.AutoTagMetadataInput;
  setOptions: (s: GQL.AutoTagMetadataInput) => void;
}

const AutoTagOptions: React.FC<IAutoTagOptions> = ({
  options,
  setOptions: setOptionsState,
}) => {
  const { performers, studios, tags } = options;
  const wildcard = ["*"];

  function set(v?: boolean) {
    if (v) {
      return wildcard;
    }
    return [];
  }

  function setOptions(input: Partial<GQL.AutoTagMetadataInput>) {
    setOptionsState({ ...options, ...input });
  }

  return (
    <>
      <BooleanSetting
        id="autotag-performers"
        checked={!!performers?.length}
        headingID="performers"
        onChange={(v) => setOptions({ performers: set(v) })}
      />
      <BooleanSetting
        id="autotag-studios"
        checked={!!studios?.length}
        headingID="studios"
        onChange={(v) => setOptions({ studios: set(v) })}
      />
      <BooleanSetting
        id="autotag-tags"
        checked={!!tags?.length}
        headingID="tags"
        onChange={(v) => setOptions({ tags: set(v) })}
      />
    </>
  );
};

export const LibraryTasks: React.FC = () => {
  const intl = useIntl();
  const Toast = useToast();
  const { ui, saveUI, loading } = useSettings();

  const { taskDefaults } = ui;

  const [dialogOpen, setDialogOpenState] = useState({
    scan: false,
    autoTag: false,
    identify: false,
  });

  function getDefaultScanOptions(): GQL.ScanMetadataInput {
    return {
      scanGenerateCovers: true,
      scanGeneratePreviews: false,
      scanGenerateImagePreviews: false,
      scanGenerateSprites: false,
      scanGeneratePhashes: false,
      scanGenerateThumbnails: false,
      scanGenerateClipPreviews: false,
    };
  }

  const [scanOptions, setScanOptions] = useState<GQL.ScanMetadataInput>(
    getDefaultScanOptions()
  );
  const [autoTagOptions, setAutoTagOptions] =
    useState<GQL.AutoTagMetadataInput>({
      performers: ["*"],
      studios: ["*"],
      tags: ["*"],
    });

  function getDefaultGenerateOptions(): GQL.GenerateMetadataInput {
    return {
      covers: true,
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

  const [generateOptions, setGenerateOptions] =
    useState<GQL.GenerateMetadataInput>(getDefaultGenerateOptions());

  type DialogOpenState = typeof dialogOpen;

  const { configuration } = useConfigurationContext();
  const [configRead, setConfigRead] = useState(false);

  useEffect(() => {
    if (!configuration?.defaults || loading) {
      return;
    }

    const { scan, autoTag } = configuration.defaults;

    // prefer UI defaults over system defaults
    // other defaults should be deprecated
    if (taskDefaults?.scan) {
      setScanOptions(taskDefaults.scan);
    } else if (scan) {
      setScanOptions(withoutTypename(scan));
    }

    if (taskDefaults?.autoTag) {
      setAutoTagOptions(taskDefaults.autoTag);
    } else if (autoTag) {
      setAutoTagOptions(withoutTypename(autoTag));
    }

    if (taskDefaults?.generate) {
      setGenerateOptions(taskDefaults.generate);
    }

    // combine the defaults with the system preview generation settings
    // only do this once
    // don't do this if UI had a default
    if (!configRead && !taskDefaults?.generate) {
      if (configuration?.defaults.generate) {
        const { generate } = configuration.defaults;
        setGenerateOptions(withoutTypename(generate));
      }

      setConfigRead(true);
    }
  }, [configuration, configRead, taskDefaults, loading]);

  function configureDefaults(partial: Record<string, {}>) {
    saveUI({ taskDefaults: { ...partial } });
  }

  function onSetScanOptions(s: GQL.ScanMetadataInput) {
    configureDefaults({ scan: s });
    setScanOptions(s);
  }

  function onSetGenerateOptions(s: GQL.GenerateMetadataInput) {
    configureDefaults({ generate: s });
    setGenerateOptions(s);
  }

  function onSetAutoTagOptions(s: GQL.AutoTagMetadataInput) {
    configureDefaults({ autoTag: s });
    setAutoTagOptions(s);
  }

  function setDialogOpen(s: Partial<DialogOpenState>) {
    setDialogOpenState((v) => {
      return { ...v, ...s };
    });
  }

  function renderScanDialog() {
    if (!dialogOpen.scan) {
      return;
    }

    return <DirectorySelectionDialog onClose={onScanDialogClosed} />;
  }

  function onScanDialogClosed(paths?: string[]) {
    if (paths) {
      runScan(paths);
    }

    setDialogOpen({ scan: false });
  }

  async function runScan(paths?: string[]) {
    try {
      await mutateMetadataScan({
        ...scanOptions,
        paths,
      });

      Toast.success(
        intl.formatMessage(
          { id: "config.tasks.added_job_to_queue" },
          { operation_name: intl.formatMessage({ id: "actions.scan" }) }
        )
      );
    } catch (e) {
      Toast.error(e);
    }
  }

  function renderAutoTagDialog() {
    if (!dialogOpen.autoTag) {
      return;
    }

    return <DirectorySelectionDialog onClose={onAutoTagDialogClosed} />;
  }

  function onAutoTagDialogClosed(paths?: string[]) {
    if (paths) {
      runAutoTag(paths);
    }

    setDialogOpen({ autoTag: false });
  }

  async function runAutoTag(paths?: string[]) {
    try {
      await mutateMetadataAutoTag({
        ...autoTagOptions,
        paths,
      });

      Toast.success(
        intl.formatMessage(
          { id: "config.tasks.added_job_to_queue" },
          { operation_name: intl.formatMessage({ id: "actions.auto_tag" }) }
        )
      );
    } catch (e) {
      Toast.error(e);
    }
  }

  function maybeRenderIdentifyDialog() {
    if (!dialogOpen.identify) return;

    return (
      <IdentifyDialog onClose={() => setDialogOpen({ identify: false })} />
    );
  }

  async function onGenerateClicked() {
    try {
      // insert preview options here instead of loading them
      const general = configuration?.general;

      await mutateMetadataGenerate({
        ...generateOptions,
        previewOptions: {
          ...generateOptions.previewOptions,
          previewSegments:
            general?.previewSegments ??
            generateOptions.previewOptions?.previewSegments,
          previewSegmentDuration:
            general?.previewSegmentDuration ??
            generateOptions.previewOptions?.previewSegmentDuration,
          previewExcludeStart:
            general?.previewExcludeStart ??
            generateOptions.previewOptions?.previewExcludeStart,
          previewExcludeEnd:
            general?.previewExcludeEnd ??
            generateOptions.previewOptions?.previewExcludeEnd,
          previewPreset:
            general?.previewPreset ??
            generateOptions.previewOptions?.previewPreset,
        },
      });
      Toast.success(
        intl.formatMessage(
          { id: "config.tasks.added_job_to_queue" },
          { operation_name: intl.formatMessage({ id: "actions.generate" }) }
        )
      );
    } catch (e) {
      Toast.error(e);
    }
  }

  return (
    <Form.Group>
      {renderScanDialog()}
      {renderAutoTagDialog()}
      {maybeRenderIdentifyDialog()}

      <SettingSection headingID="library">
        <SettingGroup
          settingProps={{
            heading: (
              <>
                <FormattedMessage id="actions.scan" />
                <ManualLink tab="Tasks">
                  <Icon icon={faQuestionCircle} />
                </ManualLink>
              </>
            ),
            subHeadingID: "config.tasks.scan_for_content_desc",
          }}
          topLevel={
            <>
              <Button
                variant="secondary"
                type="submit"
                className="mr-2"
                onClick={() => runScan()}
              >
                <FormattedMessage id="actions.scan" />
              </Button>

              <Button
                variant="secondary"
                type="submit"
                className="mr-2"
                onClick={() => setDialogOpen({ scan: true })}
              >
                <FormattedMessage id="actions.selective_scan" />…
              </Button>
            </>
          }
          collapsible
        >
          <ScanOptions options={scanOptions} setOptions={onSetScanOptions} />
        </SettingGroup>
      </SettingSection>

      <SettingSection advanced>
        <Setting
          heading={
            <>
              <FormattedMessage id="config.tasks.identify.heading" />
              <ManualLink tab="Identify">
                <Icon icon={faQuestionCircle} />
              </ManualLink>
            </>
          }
          subHeadingID="config.tasks.identify.description"
        >
          <Button
            variant="secondary"
            type="submit"
            onClick={() => setDialogOpen({ identify: true })}
          >
            <FormattedMessage id="actions.identify" />…
          </Button>
        </Setting>
      </SettingSection>

      <SettingSection advanced>
        <SettingGroup
          settingProps={{
            heading: (
              <>
                <FormattedMessage id="actions.auto_tag" />
                <ManualLink tab="AutoTagging">
                  <Icon icon={faQuestionCircle} />
                </ManualLink>
              </>
            ),
            subHeadingID: "config.tasks.auto_tag_based_on_filenames",
          }}
          topLevel={
            <>
              <Button
                variant="secondary"
                type="submit"
                className="mr-2"
                onClick={() => runAutoTag()}
              >
                <FormattedMessage id="actions.auto_tag" />
              </Button>
              <Button
                variant="secondary"
                type="submit"
                onClick={() => setDialogOpen({ autoTag: true })}
              >
                <FormattedMessage id="actions.selective_auto_tag" />…
              </Button>
            </>
          }
          collapsible
        >
          <AutoTagOptions
            options={autoTagOptions}
            setOptions={onSetAutoTagOptions}
          />
        </SettingGroup>
      </SettingSection>

      <SettingSection headingID="config.tasks.generated_content">
        <SettingGroup
          settingProps={{
            heading: (
              <>
                <FormattedMessage id="actions.generate" />
                <ManualLink tab="Tasks">
                  <Icon icon={faQuestionCircle} />
                </ManualLink>
              </>
            ),
            subHeadingID: "config.tasks.generate_desc",
          }}
          topLevel={
            <Button
              variant="secondary"
              type="submit"
              onClick={() => onGenerateClicked()}
            >
              <FormattedMessage id="actions.generate" />
            </Button>
          }
          collapsible
        >
          <GenerateOptions
            options={generateOptions}
            setOptions={onSetGenerateOptions}
          />
        </SettingGroup>
      </SettingSection>
    </Form.Group>
  );
};
