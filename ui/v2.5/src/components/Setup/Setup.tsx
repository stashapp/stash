import React, { useState, useCallback } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import {
  Alert,
  Button,
  Card,
  Container,
  Form,
  InputGroup,
} from "react-bootstrap";
import * as GQL from "src/core/generated-graphql";
import {
  mutateSetup,
  useConfigureUI,
  useSystemStatus,
} from "src/core/StashService";
import { useHistory } from "react-router-dom";
import { useConfigurationContext } from "src/hooks/Config";
import StashConfiguration from "../Settings/StashConfiguration";
import { Icon } from "../Shared/Icon";
import { LoadingIndicator } from "../Shared/LoadingIndicator";
import { ModalComponent } from "../Shared/Modal";
import { FolderSelectDialog } from "../Shared/FolderSelect/FolderSelectDialog";
import {
  faEllipsisH,
  faExclamationTriangle,
  faQuestionCircle,
} from "@fortawesome/free-solid-svg-icons";
import { releaseNotes } from "src/docs/en/ReleaseNotes";
import { ExternalLink } from "../Shared/ExternalLink";

interface ISetupContextState {
  configuration: GQL.ConfigDataFragment;
  systemStatus: GQL.SystemStatusQuery;

  setupState: Partial<GQL.SetupInput>;
  setupError: string | undefined;

  pathJoin: (...paths: string[]) => string;
  pathDir(path: string): string;

  homeDir: string;
  windows: boolean;
  macApp: boolean;
  homeDirPath: string;
  pwd: string;
  workingDir: string;
}

const SetupStateContext = React.createContext<ISetupContextState | null>(null);

const useSetupContext = () => {
  const context = React.useContext(SetupStateContext);

  if (context === null) {
    throw new Error("useSettings must be used within a SettingsContext");
  }

  return context;
};

const SetupContext: React.FC<{
  setupState: Partial<GQL.SetupInput>;
  setupError: string | undefined;
  systemStatus: GQL.SystemStatusQuery;
  configuration: GQL.ConfigDataFragment;
}> = ({ setupState, setupError, systemStatus, configuration, children }) => {
  const status = systemStatus?.systemStatus;

  const windows = status?.os === "windows";
  const pathSep = windows ? "\\" : "/";
  const homeDir = windows ? "%USERPROFILE%" : "$HOME";
  const pwd = windows ? "%CD%" : "$PWD";

  const pathJoin = useCallback(
    (...paths: string[]) => {
      return paths.join(pathSep);
    },
    [pathSep]
  );

  // simply returns everything preceding the last path separator
  function pathDir(path: string) {
    const lastSep = path.lastIndexOf(pathSep);
    if (lastSep === -1) return "";
    return path.slice(0, lastSep);
  }

  const workingDir = status?.workingDir ?? ".";

  // When running Stash.app, the working directory is (usually) set to /.
  // Assume that the user doesn't want to set up in / (it's usually mounted read-only anyway),
  // so in this situation disallow setting up in the working directory.
  const macApp = status?.os === "darwin" && workingDir === "/";

  const homeDirPath = pathJoin(status?.homeDir ?? homeDir, ".stash");

  const state: ISetupContextState = {
    systemStatus,
    configuration,
    windows,
    macApp,
    pathJoin,
    pathDir,
    homeDir,
    homeDirPath,
    pwd,
    workingDir,
    setupState,
    setupError,
  };

  return (
    <SetupStateContext.Provider value={state}>
      {children}
    </SetupStateContext.Provider>
  );
};

interface IWizardStep {
  next: (input?: Partial<GQL.SetupInput>) => void;
  goBack: () => void;
}

const WelcomeSpecificConfig: React.FC<IWizardStep> = ({ next }) => {
  const { systemStatus } = useSetupContext();
  const status = systemStatus?.systemStatus;
  const overrideConfig = status?.configPath;

  function onNext() {
    next({ configLocation: overrideConfig! });
  }

  return (
    <>
      <section>
        <h2 className="mb-5">
          <FormattedMessage id="setup.welcome_to_stash" />
        </h2>
        <p className="lead text-center">
          <FormattedMessage id="setup.welcome_specific_config.unable_to_locate_specified_config" />
        </p>
        <p>
          <FormattedMessage
            id="setup.welcome_specific_config.config_path"
            values={{
              path: overrideConfig,
              code: (chunks: string) => <code>{chunks}</code>,
            }}
          />
        </p>
        <p>
          <FormattedMessage id="setup.welcome_specific_config.next_step" />
        </p>
      </section>

      <section className="mt-5">
        <div className="d-flex justify-content-center">
          <Button variant="primary mx-2 p-5" onClick={() => onNext()}>
            <FormattedMessage id="actions.next_action" />
          </Button>
        </div>
      </section>
    </>
  );
};

const DefaultWelcomeStep: React.FC<IWizardStep> = ({ next }) => {
  const { pathJoin, homeDir, macApp, homeDirPath, pwd, workingDir } =
    useSetupContext();

  const fallbackStashDir = pathJoin(homeDir, ".stash");
  const fallbackConfigPath = pathJoin(fallbackStashDir, "config.yml");

  function onConfigLocationChosen(inWorkingDir: boolean) {
    const configLocation = inWorkingDir ? "config.yml" : "";
    next({ configLocation });
  }

  return (
    <>
      <section>
        <h2 className="mb-5">
          <FormattedMessage id="setup.welcome_to_stash" />
        </h2>
        <p className="lead text-center">
          <FormattedMessage id="setup.welcome.unable_to_locate_config" />
        </p>
        <p>
          <FormattedMessage
            id="setup.welcome.config_path_logic_explained"
            values={{
              code: (chunks: string) => <code>{chunks}</code>,
              fallback_path: fallbackConfigPath,
            }}
          />
        </p>
        <Alert variant="info text-center">
          <FormattedMessage
            id="setup.welcome.unexpected_explained"
            values={{
              code: (chunks: string) => <code>{chunks}</code>,
            }}
          />
        </Alert>
        <p>
          <FormattedMessage id="setup.welcome.next_step" />
        </p>
      </section>

      <section className="mt-5">
        <h3 className="text-center mb-5">
          <FormattedMessage id="setup.welcome.store_stash_config" />
        </h3>

        <div className="d-flex justify-content-center">
          <Button
            variant="secondary mx-2 p-5"
            onClick={() => onConfigLocationChosen(false)}
          >
            <FormattedMessage
              id="setup.welcome.in_current_stash_directory"
              values={{
                code: (chunks: string) => <code>{chunks}</code>,
                path: fallbackStashDir,
              }}
            />
            <br />
            <code>{homeDirPath}</code>
          </Button>
          <Button
            variant="secondary mx-2 p-5"
            onClick={() => onConfigLocationChosen(true)}
            disabled={macApp}
          >
            {macApp ? (
              <>
                <FormattedMessage
                  id="setup.welcome.in_the_current_working_directory_disabled"
                  values={{
                    code: (chunks: string) => <code>{chunks}</code>,
                    path: pwd,
                  }}
                />
                <br />
                <b>
                  <FormattedMessage
                    id="setup.welcome.in_the_current_working_directory_disabled_macos"
                    values={{
                      code: (chunks: string) => <code>{chunks}</code>,
                      br: () => <br />,
                    }}
                  />
                </b>
              </>
            ) : (
              <>
                <FormattedMessage
                  id="setup.welcome.in_the_current_working_directory"
                  values={{
                    code: (chunks: string) => <code>{chunks}</code>,
                    path: pwd,
                  }}
                />
                <br />
                <code>{workingDir}</code>
              </>
            )}
          </Button>
        </div>
      </section>
    </>
  );
};

const WelcomeStep: React.FC<IWizardStep> = (props) => {
  const { systemStatus } = useSetupContext();
  const status = systemStatus?.systemStatus;
  const overrideConfig = status?.configPath;

  return overrideConfig ? (
    <WelcomeSpecificConfig {...props} />
  ) : (
    <DefaultWelcomeStep {...props} />
  );
};

const StashAlert: React.FC<{ close: (confirm: boolean) => void }> = ({
  close,
}) => {
  const intl = useIntl();

  return (
    <ModalComponent
      show
      icon={faExclamationTriangle}
      accept={{
        text: intl.formatMessage({ id: "actions.confirm" }),
        variant: "danger",
        onClick: () => close(true),
      }}
      cancel={{ onClick: () => close(false) }}
    >
      <p>
        <FormattedMessage id="setup.paths.stash_alert" />
      </p>
    </ModalComponent>
  );
};

const DatabaseSection: React.FC<{
  databaseFile: string;
  setDatabaseFile: React.Dispatch<React.SetStateAction<string>>;
}> = ({ databaseFile, setDatabaseFile }) => {
  const intl = useIntl();

  return (
    <Form.Group id="database">
      <h3>
        <FormattedMessage id="setup.paths.where_can_stash_store_its_database" />
      </h3>
      <p>
        <FormattedMessage
          id="setup.paths.where_can_stash_store_its_database_description"
          values={{
            code: (chunks: string) => <code>{chunks}</code>,
          }}
        />
        <br />
        <FormattedMessage
          id="setup.paths.where_can_stash_store_its_database_warning"
          values={{
            strong: (chunks: string) => <strong>{chunks}</strong>,
          }}
        />
      </p>
      <Form.Control
        className="text-input"
        defaultValue={databaseFile}
        placeholder={intl.formatMessage({
          id: "setup.paths.database_filename_empty_for_default",
        })}
        onChange={(e) => setDatabaseFile(e.currentTarget.value)}
      />
    </Form.Group>
  );
};

const DirectorySelector: React.FC<{
  value: string;
  setValue: React.Dispatch<React.SetStateAction<string>>;
  placeholder: string;
  disabled?: boolean;
}> = ({ value, setValue, placeholder, disabled = false }) => {
  const [showSelectDialog, setShowSelectDialog] = useState(false);

  function onSelectClosed(dir?: string) {
    if (dir) {
      setValue(dir);
    }
    setShowSelectDialog(false);
  }

  return (
    <>
      {showSelectDialog ? (
        <FolderSelectDialog onClose={onSelectClosed} />
      ) : null}
      <InputGroup>
        <Form.Control
          className="text-input"
          value={disabled ? "" : value}
          placeholder={placeholder}
          onChange={(e) => setValue(e.currentTarget.value)}
          disabled={disabled}
        />
        <InputGroup.Append>
          <Button
            variant="secondary"
            className="text-input"
            onClick={() => setShowSelectDialog(true)}
            disabled={disabled}
          >
            <Icon icon={faEllipsisH} />
          </Button>
        </InputGroup.Append>
      </InputGroup>
    </>
  );
};

const GeneratedSection: React.FC<{
  generatedLocation: string;
  setGeneratedLocation: React.Dispatch<React.SetStateAction<string>>;
}> = ({ generatedLocation, setGeneratedLocation }) => {
  const intl = useIntl();

  return (
    <Form.Group id="generated">
      <h3>
        <FormattedMessage id="setup.paths.where_can_stash_store_its_generated_content" />
      </h3>
      <p>
        <FormattedMessage
          id="setup.paths.where_can_stash_store_its_generated_content_description"
          values={{
            code: (chunks: string) => <code>{chunks}</code>,
          }}
        />
      </p>
      <DirectorySelector
        value={generatedLocation}
        setValue={setGeneratedLocation}
        placeholder={intl.formatMessage({
          id: "setup.paths.path_to_generated_directory_empty_for_default",
        })}
      />
    </Form.Group>
  );
};

const CacheSection: React.FC<{
  cacheLocation: string;
  setCacheLocation: React.Dispatch<React.SetStateAction<string>>;
}> = ({ cacheLocation, setCacheLocation }) => {
  const intl = useIntl();

  return (
    <Form.Group id="cache">
      <h3>
        <FormattedMessage id="setup.paths.where_can_stash_store_cache_files" />
      </h3>
      <p>
        <FormattedMessage
          id="setup.paths.where_can_stash_store_cache_files_description"
          values={{
            code: (chunks: string) => <code>{chunks}</code>,
          }}
        />
      </p>
      <DirectorySelector
        value={cacheLocation}
        setValue={setCacheLocation}
        placeholder={intl.formatMessage({
          id: "setup.paths.path_to_cache_directory_empty_for_default",
        })}
      />
    </Form.Group>
  );
};

const BlobsSection: React.FC<{
  blobsLocation: string;
  setBlobsLocation: React.Dispatch<React.SetStateAction<string>>;
  storeBlobsInDatabase: boolean;
  setStoreBlobsInDatabase: React.Dispatch<React.SetStateAction<boolean>>;
}> = ({
  blobsLocation,
  setBlobsLocation,
  storeBlobsInDatabase,
  setStoreBlobsInDatabase,
}) => {
  const intl = useIntl();

  return (
    <Form.Group id="blobs">
      <h3>
        <FormattedMessage id="setup.paths.where_can_stash_store_blobs" />
      </h3>
      <p>
        <FormattedMessage
          id="setup.paths.where_can_stash_store_blobs_description"
          values={{
            code: (chunks: string) => <code>{chunks}</code>,
          }}
        />
      </p>
      <p>
        <FormattedMessage
          id="setup.paths.where_can_stash_store_blobs_description_addendum"
          values={{
            code: (chunks: string) => <code>{chunks}</code>,
            strong: (chunks: string) => <strong>{chunks}</strong>,
          }}
        />
      </p>

      <div>
        <Form.Check
          id="store-blobs-in-database"
          checked={storeBlobsInDatabase}
          label={intl.formatMessage({
            id: "setup.paths.store_blobs_in_database",
          })}
          onChange={() => setStoreBlobsInDatabase(!storeBlobsInDatabase)}
        />
      </div>

      <div>
        <DirectorySelector
          value={blobsLocation}
          setValue={setBlobsLocation}
          placeholder={intl.formatMessage({
            id: "setup.paths.path_to_blobs_directory_empty_for_default",
          })}
          disabled={storeBlobsInDatabase}
        />
      </div>
    </Form.Group>
  );
};

const SetPathsStep: React.FC<IWizardStep> = ({ goBack, next }) => {
  const { configuration, setupState } = useSetupContext();

  const [showStashAlert, setShowStashAlert] = useState(false);

  const [stashes, setStashes] = useState<GQL.StashConfig[]>(
    setupState.stashes ?? []
  );
  const [sfwContentMode, setSfwContentMode] = useState(
    setupState.sfwContentMode ?? false
  );

  const [databaseFile, setDatabaseFile] = useState(
    setupState.databaseFile ?? ""
  );
  const [generatedLocation, setGeneratedLocation] = useState(
    setupState.generatedLocation ?? ""
  );
  const [cacheLocation, setCacheLocation] = useState(
    setupState.cacheLocation ?? ""
  );
  const [storeBlobsInDatabase, setStoreBlobsInDatabase] = useState(
    setupState.storeBlobsInDatabase ?? false
  );
  const [blobsLocation, setBlobsLocation] = useState(
    setupState.blobsLocation ?? ""
  );

  const overrideDatabase = configuration?.general.databasePath;
  const overrideGenerated = configuration?.general.generatedPath;
  const overrideCache = configuration?.general.cachePath;
  const overrideBlobs = configuration?.general.blobsPath;

  function preNext() {
    if (stashes.length === 0) {
      setShowStashAlert(true);
    } else {
      onNext();
    }
  }

  function onNext() {
    const input: Partial<GQL.SetupInput> = {
      stashes,
      databaseFile,
      generatedLocation,
      cacheLocation,
      blobsLocation: storeBlobsInDatabase ? "" : blobsLocation,
      storeBlobsInDatabase,
      sfwContentMode,
    };
    next(input);
  }

  return (
    <>
      {showStashAlert ? (
        <StashAlert
          close={(confirm) => {
            setShowStashAlert(false);
            if (confirm) {
              onNext();
            }
          }}
        />
      ) : null}
      <section>
        <h2 className="mb-3">
          <FormattedMessage id="setup.paths.set_up_your_paths" />
        </h2>
        <p>
          <FormattedMessage id="setup.paths.description" />
        </p>
      </section>
      <section>
        <Form.Group id="stashes">
          <h3>
            <FormattedMessage id="setup.paths.where_is_your_porn_located" />
          </h3>
          <p>
            <FormattedMessage id="setup.paths.where_is_your_porn_located_description" />
          </p>
          <Card>
            <StashConfiguration
              stashes={stashes}
              setStashes={(s) => setStashes(s)}
            />
          </Card>
        </Form.Group>
        <Form.Group id="sfw_content">
          <h3>
            <FormattedMessage id="setup.paths.sfw_content_settings" />
          </h3>
          <p>
            <FormattedMessage id="setup.paths.sfw_content_settings_description" />
          </p>
          <Card>
            <Form.Check
              id="use-sfw-content-mode"
              checked={sfwContentMode}
              label={<FormattedMessage id="setup.paths.use_sfw_content_mode" />}
              onChange={() => setSfwContentMode(!sfwContentMode)}
            />
          </Card>
        </Form.Group>
        {overrideDatabase ? null : (
          <DatabaseSection
            databaseFile={databaseFile}
            setDatabaseFile={setDatabaseFile}
          />
        )}
        {overrideGenerated ? null : (
          <GeneratedSection
            generatedLocation={generatedLocation}
            setGeneratedLocation={setGeneratedLocation}
          />
        )}
        {overrideCache ? null : (
          <CacheSection
            cacheLocation={cacheLocation}
            setCacheLocation={setCacheLocation}
          />
        )}
        {overrideBlobs ? null : (
          <BlobsSection
            blobsLocation={blobsLocation}
            setBlobsLocation={setBlobsLocation}
            storeBlobsInDatabase={storeBlobsInDatabase}
            setStoreBlobsInDatabase={setStoreBlobsInDatabase}
          />
        )}
      </section>
      <section className="mt-5">
        <div className="d-flex justify-content-center">
          <Button variant="secondary mx-2 p-5" onClick={() => goBack()}>
            <FormattedMessage id="actions.previous_action" />
          </Button>
          <Button variant="primary mx-2 p-5" onClick={() => preNext()}>
            <FormattedMessage id="actions.next_action" />
          </Button>
        </div>
      </section>
    </>
  );
};

const StashExclusions: React.FC<{ stash: GQL.StashConfig }> = ({ stash }) => {
  if (!stash.excludeImage && !stash.excludeVideo) {
    return null;
  }

  const excludes = [];
  if (stash.excludeVideo) {
    excludes.push("videos");
  }
  if (stash.excludeImage) {
    excludes.push("images");
  }

  return <span>{`(excludes ${excludes.join(" and ")})`}</span>;
};

const ConfirmStep: React.FC<IWizardStep> = ({ goBack, next }) => {
  const {
    configuration,
    pathDir,
    pathJoin,
    setupState,
    homeDirPath,
    workingDir,
  } = useSetupContext();

  // if unset, means use homeDirPath
  const cfgFile = setupState.configLocation
    ? pathJoin(workingDir, setupState.configLocation)
    : pathJoin(homeDirPath, "config.yml");
  const cfgDir = pathDir(cfgFile);
  const stashes = setupState.stashes ?? [];
  const {
    databaseFile,
    generatedLocation,
    cacheLocation,
    blobsLocation,
    storeBlobsInDatabase,
  } = setupState;

  const overrideDatabase = configuration?.general.databasePath;
  const overrideGenerated = configuration?.general.generatedPath;
  const overrideCache = configuration?.general.cachePath;
  const overrideBlobs = configuration?.general.blobsPath;

  function joinCfgDir(path: string) {
    if (cfgDir) {
      return pathJoin(cfgDir, path);
    } else {
      return path;
    }
  }

  return (
    <>
      <section>
        <h2 className="mb-3">
          <FormattedMessage id="setup.confirm.nearly_there" />
        </h2>
        <p>
          <FormattedMessage id="setup.confirm.almost_ready" />
        </p>
        <dl>
          <dt>
            <FormattedMessage id="setup.confirm.configuration_file_location" />
          </dt>
          <dd>
            <code>{cfgFile}</code>
          </dd>
        </dl>
        <dl>
          <dt>
            <FormattedMessage id="setup.confirm.stash_library_directories" />
          </dt>
          <dd>
            <ul>
              {stashes.map((s) => (
                <li key={s.path}>
                  <code>{s.path} </code>
                  <StashExclusions stash={s} />
                </li>
              ))}
            </ul>
          </dd>
        </dl>
        {!overrideDatabase && (
          <dl>
            <dt>
              <FormattedMessage id="setup.confirm.database_file_path" />
            </dt>
            <dd>
              <code>{databaseFile || joinCfgDir("stash-go.sqlite")}</code>
            </dd>
          </dl>
        )}
        {!overrideGenerated && (
          <dl>
            <dt>
              <FormattedMessage id="setup.confirm.generated_directory" />
            </dt>
            <dd>
              <code>{generatedLocation || joinCfgDir("generated")}</code>
            </dd>
          </dl>
        )}
        {!overrideCache && (
          <dl>
            <dt>
              <FormattedMessage id="setup.confirm.cache_directory" />
            </dt>
            <dd>
              <code>{cacheLocation || joinCfgDir("cache")}</code>
            </dd>
          </dl>
        )}
        {!overrideBlobs && (
          <dl>
            <dt>
              <FormattedMessage id="setup.confirm.blobs_directory" />
            </dt>
            <dd>
              <code>
                {storeBlobsInDatabase ? (
                  <FormattedMessage id="setup.confirm.blobs_use_database" />
                ) : (
                  blobsLocation || joinCfgDir("blobs")
                )}
              </code>
            </dd>
          </dl>
        )}
      </section>
      <section className="mt-5">
        <div className="d-flex justify-content-center">
          <Button variant="secondary mx-2 p-5" onClick={() => goBack()}>
            <FormattedMessage id="actions.previous_action" />
          </Button>
          <Button variant="success mx-2 p-5" onClick={() => next()}>
            <FormattedMessage id="actions.confirm" />
          </Button>
        </div>
      </section>
    </>
  );
};

const DiscordLink = (
  <ExternalLink href="https://discord.gg/2TsNFKt">Discord</ExternalLink>
);
const GithubLink = (
  <ExternalLink href="https://github.com/stashapp/stash/issues">
    <FormattedMessage id="setup.github_repository" />
  </ExternalLink>
);

const ErrorStep: React.FC<{ error: string; goBack: () => void }> = ({
  error,
  goBack,
}) => {
  return (
    <>
      <section>
        <h2>
          <FormattedMessage id="setup.errors.something_went_wrong" />
        </h2>
        <p>
          <FormattedMessage
            id="setup.errors.something_went_wrong_while_setting_up_your_system"
            values={{ error: <pre>{error}</pre> }}
          />
        </p>
        <p>
          <FormattedMessage
            id="setup.errors.something_went_wrong_description"
            values={{ githubLink: GithubLink, discordLink: DiscordLink }}
          />
        </p>
      </section>
      <section className="mt-5">
        <div className="d-flex justify-content-center">
          <Button variant="secondary mx-2 p-5" onClick={goBack}>
            <FormattedMessage id="actions.previous_action" />
          </Button>
        </div>
      </section>
    </>
  );
};

const SuccessStep: React.FC<{}> = () => {
  const intl = useIntl();
  const history = useHistory();

  const [mutateDownloadFFMpeg] = GQL.useDownloadFfMpegMutation();

  const [downloadFFmpeg, setDownloadFFmpeg] = useState(true);

  const { systemStatus } = useSetupContext();
  const status = systemStatus?.systemStatus;

  function onFinishClick() {
    if ((!status?.ffmpegPath || !status?.ffprobePath) && downloadFFmpeg) {
      mutateDownloadFFMpeg();
    }

    history.push("/settings?tab=library");
  }

  return (
    <>
      <section>
        <h2>
          <FormattedMessage id="setup.success.your_system_has_been_created" />
        </h2>
        <p>
          <FormattedMessage id="setup.success.next_config_step_one" />
        </p>
        <p>
          <FormattedMessage
            id="setup.success.next_config_step_two"
            values={{
              code: (chunks: string) => <code>{chunks}</code>,
              localized_task: intl.formatMessage({
                id: "config.categories.tasks",
              }),
              localized_scan: intl.formatMessage({ id: "actions.scan" }),
            }}
          />
        </p>
        {!status?.ffmpegPath || !status?.ffprobePath ? (
          <>
            <Alert variant="warning text-center">
              <FormattedMessage
                id="setup.success.missing_ffmpeg"
                values={{
                  code: (chunks: string) => <code>{chunks}</code>,
                }}
              />
            </Alert>
            <p>
              <Form.Check
                id="download-ffmpeg"
                checked={downloadFFmpeg}
                label={intl.formatMessage({
                  id: "setup.success.download_ffmpeg",
                })}
                onChange={() => setDownloadFFmpeg(!downloadFFmpeg)}
              />
            </p>
          </>
        ) : null}
      </section>
      <section>
        <h3>
          <FormattedMessage id="setup.success.getting_help" />
        </h3>
        <p>
          <FormattedMessage
            id="setup.success.in_app_manual_explained"
            values={{ icon: <Icon icon={faQuestionCircle} /> }}
          />
        </p>
        <p>
          <FormattedMessage
            id="setup.success.help_links"
            values={{ discordLink: DiscordLink, githubLink: GithubLink }}
          />
        </p>
      </section>
      <section>
        <h3>
          <FormattedMessage id="setup.success.support_us" />
        </h3>
        <p>
          <FormattedMessage
            id="setup.success.open_collective"
            values={{
              open_collective_link: (
                <ExternalLink href="https://opencollective.com/stashapp">
                  Open Collective
                </ExternalLink>
              ),
            }}
          />
        </p>
        <p>
          <FormattedMessage id="setup.success.welcome_contrib" />
        </p>
      </section>
      <section>
        <p className="lead text-center">
          <FormattedMessage id="setup.success.thanks_for_trying_stash" />
        </p>
      </section>
      <section className="mt-5">
        <div className="d-flex justify-content-center">
          <Button variant="success mx-2 p-5" onClick={() => onFinishClick()}>
            <FormattedMessage id="actions.finish" />
          </Button>
        </div>
      </section>
    </>
  );
};

const FinishStep: React.FC<IWizardStep> = ({ goBack }) => {
  const { setupError } = useSetupContext();

  if (setupError !== undefined) {
    return <ErrorStep error={setupError} goBack={goBack} />;
  }

  return <SuccessStep />;
};

export const Setup: React.FC = () => {
  const intl = useIntl();
  const { configuration } = useConfigurationContext();

  const [saveUI] = useConfigureUI();

  const {
    data: systemStatus,
    loading: statusLoading,
    error: statusError,
  } = useSystemStatus();

  const [step, setStep] = useState(0);
  const [setupInput, setSetupInput] = useState<Partial<GQL.SetupInput>>({});
  const [creating, setCreating] = useState(false);
  const [setupError, setSetupError] = useState<string | undefined>(undefined);

  const history = useHistory();

  const steps: React.FC<IWizardStep>[] = [
    WelcomeStep,
    SetPathsStep,
    ConfirmStep,
    FinishStep,
  ];
  const Step = steps[step];

  async function createSystem() {
    try {
      setCreating(true);
      setSetupError(undefined);
      await mutateSetup(setupInput as GQL.SetupInput);
      // Set lastNoteSeen to hide release notes dialog
      await saveUI({
        variables: {
          input: {
            ...configuration?.ui,
            lastNoteSeen: releaseNotes[0].date,
          },
        },
      });
    } catch (e) {
      if (e instanceof Error && e.message) {
        setSetupError(e.message);
      } else {
        setSetupError(String(e));
      }
    } finally {
      setCreating(false);
      setStep(step + 1);
    }
  }

  function next(input?: Partial<GQL.SetupInput>) {
    setSetupInput({ ...setupInput, ...input });

    if (Step === ConfirmStep) {
      // create the system
      createSystem();
    } else {
      setStep(step + 1);
    }
  }

  function goBack() {
    if (Step === FinishStep) {
      // go back to the step before ConfirmStep
      setStep(step - 2);
    } else {
      setStep(step - 1);
    }
  }

  if (statusLoading) {
    return <LoadingIndicator />;
  }

  if (
    step === 0 &&
    systemStatus &&
    systemStatus.systemStatus.status !== GQL.SystemStatusEnum.Setup
  ) {
    // redirect to main page
    history.push("/");
    return <LoadingIndicator />;
  }

  if (statusError) {
    return (
      <Container>
        <Alert variant="danger">
          <FormattedMessage
            id="setup.errors.unable_to_retrieve_system_status"
            values={{ error: statusError.message }}
          />
        </Alert>
      </Container>
    );
  }

  if (!configuration || !systemStatus) {
    return (
      <Container>
        <Alert variant="danger">
          <FormattedMessage
            id="setup.errors.unable_to_retrieve_configuration"
            values={{ error: "configuration or systemStatus === undefined" }}
          />
        </Alert>
      </Container>
    );
  }

  return (
    <SetupContext
      setupState={setupInput}
      setupError={setupError}
      configuration={configuration}
      systemStatus={systemStatus}
    >
      <Container className="setup-wizard">
        <h1 className="text-center">
          <FormattedMessage id="setup.stash_setup_wizard" />
        </h1>
        <Card>
          {creating ? (
            <LoadingIndicator
              message={intl.formatMessage({
                id: "setup.creating.creating_your_system",
              })}
            />
          ) : (
            <Step next={next} goBack={goBack} />
          )}
        </Card>
      </Container>
    </SetupContext>
  );
};

export default Setup;
