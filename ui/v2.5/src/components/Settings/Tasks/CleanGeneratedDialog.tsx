import React, { useEffect, useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { ModalComponent } from "src/components/Shared/Modal";
import * as GQL from "src/core/generated-graphql";
import { BooleanSetting } from "../Inputs";
import { faTrashAlt } from "@fortawesome/free-solid-svg-icons";
import { SettingSection } from "../SettingSection";
import { useSettings } from "../context";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";

const CleanGeneratedOptions: React.FC<{
  options: GQL.CleanGeneratedInput;
  setOptions: (s: GQL.CleanGeneratedInput) => void;
}> = ({ options, setOptions: setOptionsState }) => {
  function setOptions(input: Partial<GQL.CleanGeneratedInput>) {
    setOptionsState({ ...options, ...input });
  }

  return (
    <>
      <BooleanSetting
        id="clean-generated-blob-files"
        checked={options.blobFiles ?? false}
        headingID="config.tasks.clean_generated.blob_files"
        onChange={(v) => setOptions({ blobFiles: v })}
      />
      <BooleanSetting
        id="clean-generated-screenshots"
        checked={options.screenshots ?? false}
        headingID="config.tasks.clean_generated.previews"
        subHeadingID="config.tasks.clean_generated.previews_desc"
        onChange={(v) => setOptions({ screenshots: v })}
      />
      <BooleanSetting
        id="clean-generated-sprites"
        checked={options.sprites ?? false}
        headingID="config.tasks.clean_generated.sprites"
        onChange={(v) => setOptions({ sprites: v })}
      />
      <BooleanSetting
        id="clean-generated-transcodes"
        checked={options.transcodes ?? false}
        headingID="config.tasks.clean_generated.transcodes"
        onChange={(v) => setOptions({ transcodes: v })}
      />
      <BooleanSetting
        id="clean-generated-markers"
        checked={options.markers ?? false}
        headingID="config.tasks.clean_generated.markers"
        onChange={(v) => setOptions({ markers: v })}
      />
      <BooleanSetting
        id="clean-generated-image-thumbnails"
        checked={options.imageThumbnails ?? false}
        headingID="config.tasks.clean_generated.image_thumbnails"
        subHeadingID="config.tasks.clean_generated.image_thumbnails_desc"
        onChange={(v) => setOptions({ imageThumbnails: v })}
      />
      <BooleanSetting
        id="clean-generated-dryrun"
        checked={options.dryRun ?? false}
        headingID="config.tasks.only_dry_run"
        onChange={(v) => setOptions({ dryRun: v })}
      />
    </>
  );
};

export const CleanGeneratedDialog: React.FC<{
  onClose: (input?: GQL.CleanGeneratedInput) => void;
}> = ({ onClose }) => {
  const intl = useIntl();

  const { ui, saveUI, loading } = useSettings();

  const [options, setOptions] = useState<GQL.CleanGeneratedInput>({
    blobFiles: true,
    imageThumbnails: true,
    markers: true,
    screenshots: true,
    sprites: true,
    transcodes: true,
    dryRun: false,
  });

  useEffect(() => {
    const defaults = ui.taskDefaults?.cleanGenerated;
    if (defaults) {
      setOptions(defaults);
    }
  }, [ui?.taskDefaults?.cleanGenerated]);

  function confirm() {
    saveUI({
      taskDefaults: {
        ...ui.taskDefaults,
        cleanGenerated: options,
      },
    });
    onClose(options);
  }

  if (loading) return <LoadingIndicator />;

  return (
    <ModalComponent
      show
      header={<FormattedMessage id="actions.clean_generated" />}
      icon={faTrashAlt}
      accept={{
        text: intl.formatMessage({ id: "actions.clean_generated" }),
        variant: "danger",
        onClick: () => confirm(),
      }}
      cancel={{ onClick: () => onClose() }}
    >
      <div className="dialog-container">
        <p>
          <FormattedMessage id="config.tasks.clean_generated.description" />
        </p>
        <SettingSection>
          <CleanGeneratedOptions options={options} setOptions={setOptions} />
        </SettingSection>
        {options.dryRun && (
          <p>
            <FormattedMessage id="actions.tasks.dry_mode_selected" />
          </p>
        )}
      </div>
    </ModalComponent>
  );
};
