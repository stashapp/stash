import React, { useState, useMemo, useEffect } from "react";
import { Button, Form } from "react-bootstrap";
import {
  mutateMetadataAutoTag,
  useConfiguration,
  useConfigureDefaults,
} from "src/core/StashService";
import { Icon, Modal, OperationButton } from "src/components/Shared";
import { useToast } from "src/hooks";
import * as GQL from "src/core/generated-graphql";
import { FormattedMessage, useIntl } from "react-intl";
import { DirectorySelectionDialog } from "src/components/Settings/SettingsTasksPanel/DirectorySelectionDialog";
import { Manual } from "src/components/Help/Manual";
import { withoutTypename } from "src/utils";

interface IAutoTagOptions {
  options: GQL.AutoTagMetadataInput;
  setOptions: (s: GQL.AutoTagMetadataInput) => void;
}

const AutoTagOptions: React.FC<IAutoTagOptions> = ({
  options,
  setOptions: setOptionsState,
}) => {
  const intl = useIntl();

  const { performers, studios, tags } = options;
  const wildcard = ["*"];

  function toggle(v?: GQL.Maybe<string[]>) {
    if (!v) {
      return wildcard;
    }
    return [];
  }

  function setOptions(input: Partial<GQL.AutoTagMetadataInput>) {
    setOptionsState({ ...options, ...input });
  }

  return (
    <Form.Group>
      <Form.Check
        id="autotag-performers"
        checked={!!performers?.length}
        label={intl.formatMessage({ id: "performers" })}
        onChange={() => setOptions({ performers: toggle(performers) })}
      />
      <Form.Check
        id="autotag-studios"
        checked={!!studios?.length}
        label={intl.formatMessage({ id: "studios" })}
        onChange={() => setOptions({ studios: toggle(studios) })}
      />
      <Form.Check
        id="autotag-tags"
        checked={!!tags?.length}
        label={intl.formatMessage({ id: "tags" })}
        onChange={() => setOptions({ tags: toggle(tags) })}
      />
    </Form.Group>
  );
};

interface IAutoTagDialogProps {
  onClose: () => void;
}

export const AutoTagDialog: React.FC<IAutoTagDialogProps> = ({ onClose }) => {
  const [configureDefaults] = useConfigureDefaults();

  const [options, setOptions] = useState<GQL.AutoTagMetadataInput>({
    performers: ["*"],
    studios: ["*"],
    tags: ["*"],
  });
  const [paths, setPaths] = useState<string[]>([]);
  const [showManual, setShowManual] = useState(false);
  const [settingPaths, setSettingPaths] = useState(false);
  const [animation, setAnimation] = useState(true);
  const [savingDefaults, setSavingDefaults] = useState(false);

  const intl = useIntl();
  const Toast = useToast();

  const { data: configData, error: configError } = useConfiguration();

  useEffect(() => {
    if (!configData?.configuration.defaults) {
      return;
    }

    const { autoTag } = configData.configuration.defaults;

    if (autoTag) {
      setOptions(withoutTypename(autoTag));
    }
  }, [configData]);

  const selectionStatus = useMemo(() => {
    const message = paths.length ? (
      <div>
        <FormattedMessage id="config.tasks.auto_tag.auto_tagging_paths" />:
        <ul>
          {paths.map((p) => (
            <li key={p}>{p}</li>
          ))}
        </ul>
      </div>
    ) : (
      <span>
        <FormattedMessage id="config.tasks.auto_tag.auto_tagging_all_paths" />.
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
              <Icon icon="folder-open" />
            </Button>
          </div>
        </div>
      </Form.Group>
    );
  }, [intl, paths]);

  if (configError) return <div>{configError}</div>;
  if (!configData) return <div />;

  function makeDefaultAutoTagInput() {
    const ret = options;
    const { paths: _paths, ...withoutSpecifics } = ret;
    return withoutSpecifics;
  }

  async function onAutoTag() {
    try {
      await mutateMetadataAutoTag(options);

      Toast.success({
        content: intl.formatMessage(
          { id: "config.tasks.added_job_to_queue" },
          { operation_name: intl.formatMessage({ id: "actions.auto_tag" }) }
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

  async function setAsDefault() {
    try {
      setSavingDefaults(true);
      await configureDefaults({
        variables: {
          input: {
            autoTag: makeDefaultAutoTagInput(),
          },
        },
      });

      Toast.success({
        content: intl.formatMessage(
          { id: "config.tasks.defaults_set" },
          { action: intl.formatMessage({ id: "actions.auto_tag" }) }
        ),
      });
    } catch (e) {
      Toast.error(e);
    } finally {
      setSavingDefaults(false);
    }
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
        defaultActiveTab="AutoTagging.md"
      />
    );
  }

  return (
    <Modal
      modalProps={{ animation }}
      show
      icon="cogs"
      header={intl.formatMessage({ id: "actions.auto_tag" })}
      accept={{
        onClick: onAutoTag,
        text: intl.formatMessage({ id: "actions.auto_tag" }),
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
        <AutoTagOptions options={options} setOptions={(o) => setOptions(o)} />
      </Form>
    </Modal>
  );
};

export default AutoTagDialog;
