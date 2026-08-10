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
import { mutateSetup, useSystemStatus } from "src/core/StashService";
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
  faEye,
  faEyeSlash,
} from "@fortawesome/free-solid-svg-icons";
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
    throw new Error("useSetupContext must be used within a SetupContext");
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
  if (!stash.excludeImage && !stash.excludeVideo && !stash.excludeAudio) {
    return null;
  }

  const excludes = [];
  if (stash.excludeVideo) {
    excludes.push("videos");
  }
  if (stash.excludeAudio) {
    excludes.push("audio");
  }
  if (stash.excludeImage) {
    excludes.push("images");
  }

  return <span>{`(excludes ${excludes.join(" and ")})`}</span>;
};

function validateUsername(username: string) {
  if (username.trim() !== username) {
    return false;
  }

  return true;
}

function validatePassword(username: string, password: string) {
  if (!username) return true;

  if (!password.length) return false;
  return true;
}

const PasswordField: React.FC<{
  password: string;
  setPassword: React.Dispatch<React.SetStateAction<string>>;
  isInvalid?: boolean;
}> = ({ password, setPassword, isInvalid }) => {
  const [showPassword, setShowPassword] = useState(false);
  const intl = useIntl();

  const type = showPassword ? "text" : "password";
  const hideShowTextID = showPassword ? "actions.hide" : "actions.show";
  const icon = showPassword ? faEyeSlash : faEye;

  return (
    <div className="password-field-group">
      <Form.Control
        className="text-input"
        type={type}
        placeholder={intl.formatMessage({
          id: "login.password",
        })}
        value={password}
        onChange={(e) => setPassword(e.currentTarget.value)}
        isInvalid={isInvalid}
      />
      <Button
        variant="secondary"
        onClick={() => setShowPassword(!showPassword)}
        title={intl.formatMessage({ id: hideShowTextID })}
        className="show-password-button"
      >
        <Icon icon={icon} />
      </Button>
      <Form.Control.Feedback type="invalid">
        {isInvalid
          ? intl.formatMessage({
              id: "setup.credentials.password_invalid",
            })
          : null}
      </Form.Control.Feedback>
    </div>
  );
};

const CredentialsStep: React.FC<IWizardStep> = ({ goBack, next }) => {
  const { setupState } = useSetupContext();
  const intl = useIntl();

  const [username, setUsername] = useState(setupState.initialUsername || "");
  const [password, setPassword] = useState(setupState.initialPassword ?? "");
  const usernameValid = validateUsername(username);
  const passwordValid = validatePassword(username, password);
  const valid = usernameValid && passwordValid;

  function onNext() {
    const input: Partial<GQL.SetupInput> = {
      initialUsername: username.trim() === "" ? undefined : username,
      initialPassword: password === "" ? undefined : password,
    };
    next(input);
  }

  return (
    <>
      <section>
        <h2 className="mb-3">
          <FormattedMessage id="setup.credentials.heading" />
        </h2>
        <p>
          <FormattedMessage id="setup.credentials.description" />
        </p>

        <Form.Group controlId="username">
          <Form.Label>
            <FormattedMessage id="login.username" />:
          </Form.Label>
          <Form.Control
            className="text-input"
            placeholder={intl.formatMessage({
              id: "login.username",
            })}
            value={username}
            onChange={(e) => setUsername(e.currentTarget.value)}
            isInvalid={!usernameValid}
          />
          <Form.Control.Feedback type="invalid">
            {!usernameValid
              ? intl.formatMessage({ id: "setup.credentials.username_invalid" })
              : null}
          </Form.Control.Feedback>
        </Form.Group>
        <Form.Group controlId="password">
          <Form.Label>
            <FormattedMessage id="login.password" />:
          </Form.Label>
          <Form.Group id="password">
            <PasswordField
              password={password}
              setPassword={setPassword}
              isInvalid={!passwordValid}
            />
          </Form.Group>
        </Form.Group>
      </section>

      <section className="mt-5">
        <div className="d-flex justify-content-center">
          <Button variant="secondary mx-2 p-5" onClick={() => goBack()}>
            <FormattedMessage id="actions.previous_action" />
          </Button>
          <Button
            variant="primary mx-2 p-5"
            onClick={() => onNext()}
            disabled={!valid}
          >
            <FormattedMessage id="actions.next_action" />
          </Button>
        </div>
      </section>
    </>
  );
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
    initialUsername,
    initialPassword,
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
        {initialUsername?.length && initialPassword?.length && (
          <>
            <dl>
              <dt>
                <FormattedMessage id="login.username" />
              </dt>
              <dd>
                <code>{initialUsername}</code>
              </dd>
            </dl>
            <dl>
              <dt>
                <FormattedMessage id="login.password" />
              </dt>
              <dd>
                <code>********</code>
              </dd>
            </dl>
          </>
        )}
      </section>
      {initialUsername?.length && initialPassword?.length && (
        <p className="lead">
          <Icon icon={faExclamationTriangle} className="text-warning" />
          <FormattedMessage id="setup.confirm.password_set_warning" />
        </p>
      )}
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

const FinishStep: React.FC<IWizardStep> = ({ goBack }) => {
  const { setupError } = useSetupContext();

  if (!setupError) {
    throw new Error(
      "FinishStep should only be shown when there is a setupError"
    );
  }

  return <ErrorStep error={setupError} goBack={goBack} />;
};

export const Setup: React.FC = () => {
  const intl = useIntl();
  const { configuration } = useConfigurationContext();

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
    CredentialsStep,
    ConfirmStep,
    FinishStep,
  ];
  const Step = steps[step];

  async function createSystem() {
    try {
      setCreating(true);
      setSetupError(undefined);
      await mutateSetup(setupInput as GQL.SetupInput);
      history.replace("/welcome");
    } catch (e) {
      if (e instanceof Error && e.message) {
        setSetupError(e.message);
      } else {
        setSetupError(String(e));
      }
      setStep(step + 1);
    } finally {
      setCreating(false);
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
    // redirect to welcome page
    history.push("/welcome");
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
