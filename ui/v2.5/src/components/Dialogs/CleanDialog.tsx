import React, { useState, useMemo } from "react";
import { Button, Form } from "react-bootstrap";
import {
  mutateMetadataClean,
  useConfiguration,
  // useConfigureDefaults,
} from "src/core/StashService";
import { Icon, Modal } from "src/components/Shared";
import { useToast } from "src/hooks";
import * as GQL from "src/core/generated-graphql";
import { useIntl } from "react-intl";
import { Manual } from "src/components/Help/Manual";

interface ICleanOptions {
  options: GQL.CleanMetadataInput;
  setOptions: (s: GQL.CleanMetadataInput) => void;
}

const CleanOptions: React.FC<ICleanOptions> = ({
  options,
  setOptions: setOptionsState,
}) => {
  const intl = useIntl();

  function setOptions(input: Partial<GQL.CleanMetadataInput>) {
    setOptionsState({ ...options, ...input });
  }

  return (
    <Form.Group>
      <Form.Check
        id="clean-dryrun"
        checked={options.dryRun}
        label={intl.formatMessage({ id: "config.tasks.only_dry_run" })}
        onChange={() => setOptions({ dryRun: !options.dryRun })}
      />
    </Form.Group>
  );
};

interface ICleanDialog {
  onClose: () => void;
}

export const CleanDialog: React.FC<ICleanDialog> = ({ onClose }) => {
  const [options, setOptions] = useState<GQL.CleanMetadataInput>({
    dryRun: false,
  });
  // TODO - selective clean
  // const [paths, setPaths] = useState<string[]>([]);
  // const [settingPaths, setSettingPaths] = useState(false);
  const [showManual, setShowManual] = useState(false);
  const [animation, setAnimation] = useState(true);

  const intl = useIntl();
  const Toast = useToast();

  const { data: configData, error: configError } = useConfiguration();

  const message = useMemo(() => {
    if (options.dryRun) {
      return (
        <p>{intl.formatMessage({ id: "actions.tasks.dry_mode_selected" })}</p>
      );
    } else {
      return (
        <p>
          {intl.formatMessage({ id: "actions.tasks.clean_confirm_message" })}
        </p>
      );
    }
  }, [options.dryRun, intl]);

  // const selectionStatus = useMemo(() => {
  //   const message = paths.length ? (
  //     <div>
  //       <FormattedMessage id="config.tasks.auto_tag.auto_tagging_paths" />:
  //       <ul>
  //         {paths.map((p) => (
  //           <li key={p}>{p}</li>
  //         ))}
  //       </ul>
  //     </div>
  //   ) : (
  //     <span>
  //       <FormattedMessage id="config.tasks.auto_tag.auto_tagging_all_paths" />.
  //     </span>
  //   );

  //   function onClick() {
  //     setAnimation(false);
  //     setSettingPaths(true);
  //   }

  //   return (
  //     <Form.Group className="dialog-selected-folders">
  //       <div>
  //         {message}
  //         <div>
  //           <Button
  //             title={intl.formatMessage({ id: "actions.select_folders" })}
  //             onClick={() => onClick()}
  //           >
  //             <Icon icon="folder-open" />
  //           </Button>
  //         </div>
  //       </div>
  //     </Form.Group>
  //   );
  // }, [intl, paths]);

  if (configError) return <div>{configError}</div>;
  if (!configData) return <div />;

  async function onClean() {
    try {
      await mutateMetadataClean(options);

      Toast.success({
        content: intl.formatMessage(
          { id: "config.tasks.added_job_to_queue" },
          { operation_name: intl.formatMessage({ id: "actions.clean" }) }
        ),
      });
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
      modalProps={{ animation }}
      show
      icon="cogs"
      header={intl.formatMessage({ id: "actions.clean" })}
      accept={{
        onClick: onClean,
        variant: "danger",
        text: intl.formatMessage({ id: "actions.clean" }),
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
          <Icon icon="question-circle" />
        </Button>
      }
    >
      <Form>
        <CleanOptions options={options} setOptions={(o) => setOptions(o)} />
        {message}
      </Form>
    </Modal>
  );
};

export default CleanDialog;
