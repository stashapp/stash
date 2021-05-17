import React, { useEffect, useState } from "react";
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
import { Link } from "react-router-dom";
import StashConfiguration from "../Settings/StashConfiguration";
import { Icon, LoadingIndicator } from "../Shared";
import { FolderSelectDialog } from "../Shared/FolderSelect/FolderSelectDialog";

export const Setup: React.FC = () => {
  const [step, setStep] = useState(0);
  const [configLocation, setConfigLocation] = useState("");
  const [stashes, setStashes] = useState<GQL.StashConfig[]>([]);
  const [generatedLocation, setGeneratedLocation] = useState("");
  const [databaseFile, setDatabaseFile] = useState("");
  const [loading, setLoading] = useState(false);
  const [setupError, setSetupError] = useState("");

  const [showGeneratedDialog, setShowGeneratedDialog] = useState(false);

  const { data: systemStatus, loading: statusLoading } = useSystemStatus();

  useEffect(() => {
    if (systemStatus?.systemStatus.configPath) {
      setConfigLocation(systemStatus.systemStatus.configPath);
    }
  }, [systemStatus]);

  const discordLink = (
    <a href="https://discord.gg/2TsNFKt" target="_blank" rel="noreferrer">
      Discord
    </a>
  );
  const githubLink = (
    <a
      href="https://github.com/stashapp/stash/issues"
      target="_blank"
      rel="noreferrer"
    >
      Github repository
    </a>
  );

  function onConfigLocationChosen(loc: string) {
    setConfigLocation(loc);
    next();
  }

  function goBack(n?: number) {
    let dec = n;
    if (!dec) {
      dec = 1;
    }
    setStep(Math.max(0, step - dec));
  }

  function next() {
    setStep(step + 1);
  }

  function renderWelcomeSpecificConfig() {
    return (
      <>
        <section>
          <h2 className="mb-5">Welcome to Stash</h2>
          <p className="lead text-center">
            If you&apos;re reading this, then Stash couldn&apos;t find the
            configuration file specified at the command line or the environment.
            This wizard will guide you through the process of setting up a new
            configuration.
          </p>
          <p>
            Stash will use the following configuration file path:{" "}
            <code>{configLocation}</code>
          </p>
          <p>
            When you&apos;re ready to proceed with setting up a new system,
            click Next.
          </p>
        </section>

        <section className="mt-5">
          <div className="d-flex justify-content-center">
            <Button variant="primary mx-2 p-5" onClick={() => next()}>
              Next
            </Button>
          </div>
        </section>
      </>
    );
  }

  function renderWelcome() {
    return (
      <>
        <section>
          <h2 className="mb-5">Welcome to Stash</h2>
          <p className="lead text-center">
            If you&apos;re reading this, then Stash couldn&apos;t find an
            existing configuration. This wizard will guide you through the
            process of setting up a new configuration.
          </p>
          <p>
            Stash tries to find its configuration file (<code>config.yml</code>)
            from the current working directory first, and if it does not find it
            there, it falls back to <code>$HOME/.stash/config.yml</code> (on
            Windows, this will be <code>%USERPROFILE%\.stash\config.yml</code>).
            You can also make Stash read from a specific configuration file by
            running it with the <code>-c &lt;path to config file&gt;</code> or{" "}
            <code>--config &lt;path to config file&gt;</code> options.
          </p>
          <Alert variant="info text-center">
            If you&apos;re getting this screen unexpectedly, please try
            restarting Stash in the correct working directory or with the{" "}
            <code>-c</code> flag.
          </Alert>
          <p>
            With all of that out of the way, if you&apos;re ready to proceed
            with setting up a new system, choose where you&apos;d like to store
            your configuration file and click Next.
          </p>
        </section>

        <section className="mt-5">
          <h3 className="text-center mb-5">
            Where do you want to store your Stash configuration?
          </h3>

          <div className="d-flex justify-content-center">
            <Button
              variant="secondary mx-2 p-5"
              onClick={() => onConfigLocationChosen("")}
            >
              In the <code>$HOME/.stash</code> directory
            </Button>
            <Button
              variant="secondary mx-2 p-5"
              onClick={() => onConfigLocationChosen("config.yml")}
            >
              In the current working directory
            </Button>
          </div>
        </section>
      </>
    );
  }

  function onGeneratedClosed(d?: string) {
    if (d) {
      setGeneratedLocation(d);
    }

    setShowGeneratedDialog(false);
  }

  function maybeRenderGeneratedSelectDialog() {
    if (!showGeneratedDialog) {
      return;
    }

    return <FolderSelectDialog onClose={onGeneratedClosed} />;
  }

  function renderSetPaths() {
    return (
      <>
        <section>
          <h2 className="mb-3">Set up your paths</h2>
          <p>
            Next up, we need to determine where to find your porn collection,
            where to store the stash database and generated files. These
            settings can be changed later if needed.
          </p>
        </section>
        <section>
          <Form.Group id="stashes">
            <h3>Where is your porn located?</h3>
            <p>
              Add directories containing your porn videos and images. Stash will
              use these directories to find videos and images during scanning.
            </p>
            <Card>
              <StashConfiguration
                stashes={stashes}
                setStashes={(s) => setStashes(s)}
              />
            </Card>
          </Form.Group>
          <Form.Group id="database">
            <h3>Where can Stash store its database?</h3>
            <p>
              Stash uses an sqlite database to store your porn metadata. By
              default, this will be created as <code>stash-go.sqlite</code> in
              the directory containing your config file. If you want to change
              this, please enter an absolute or relative (to the current working
              directory) filename.
            </p>
            <Form.Control
              className="text-input"
              defaultValue={databaseFile}
              placeholder="database filename (empty for default)"
              onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                setDatabaseFile(e.currentTarget.value)
              }
            />
          </Form.Group>
          <Form.Group id="generated">
            <h3>Where can Stash store its generated content?</h3>
            <p>
              In order to provide thumbnails, previews and sprites, Stash
              generates images and videos. This also includes transcodes for
              unsupported file formats. By default, Stash will create a{" "}
              <code>generated</code> directory within the directory containing
              your config file. If you want to change where this generated media
              will be stored, please enter an absolute or relative (to the
              current working directory) path. Stash will create this directory
              if it does not already exist.
            </p>
            <InputGroup>
              <Form.Control
                className="text-input"
                value={generatedLocation}
                placeholder="path to generated directory (empty for default)"
                onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                  setGeneratedLocation(e.currentTarget.value)
                }
              />
              <InputGroup.Append>
                <Button
                  variant="secondary"
                  className="text-input"
                  onClick={() => setShowGeneratedDialog(true)}
                >
                  <Icon icon="ellipsis-h" />
                </Button>
              </InputGroup.Append>
            </InputGroup>
          </Form.Group>
        </section>
        <section className="mt-5">
          <div className="d-flex justify-content-center">
            <Button variant="secondary mx-2 p-5" onClick={() => goBack()}>
              Back
            </Button>
            <Button variant="primary mx-2 p-5" onClick={() => next()}>
              Next
            </Button>
          </div>
        </section>
      </>
    );
  }

  function renderConfigLocation() {
    if (configLocation === "config.yml") {
      return <code>&lt;current working directory&gt;/config.yml</code>;
    }

    if (configLocation === "") {
      return <code>$HOME/.stash/config.yml</code>;
    }

    return <code>{configLocation}</code>;
  }

  function maybeRenderExclusions(s: GQL.StashConfig) {
    if (!s.excludeImage && !s.excludeVideo) {
      return;
    }

    const excludes = [];
    if (s.excludeVideo) {
      excludes.push("videos");
    }
    if (s.excludeImage) {
      excludes.push("images");
    }

    return `(excludes ${excludes.join(" and ")})`;
  }

  function renderStashLibraries() {
    return (
      <ul>
        {stashes.map((s) => (
          <li>
            <code>{s.path} </code>
            {maybeRenderExclusions(s)}
          </li>
        ))}
      </ul>
    );
  }

  async function onSave() {
    try {
      setLoading(true);
      await mutateSetup({
        configLocation,
        databaseFile,
        generatedLocation,
        stashes,
      });
    } catch (e) {
      setSetupError(e.message ?? e.toString());
    } finally {
      setLoading(false);
      next();
    }
  }

  function renderConfirm() {
    return (
      <>
        <section>
          <h2 className="mb-3">Nearly there!</h2>
          <p>
            We&apos;re almost ready to complete the configuration. Please
            confirm the following settings. You can click back to change
            anything incorrect. If everything looks good, click Confirm to
            create your system.
          </p>
          <dl>
            <dt>Configuration file location:</dt>
            <dd>{renderConfigLocation()}</dd>
          </dl>
          <dl>
            <dt>Stash library directories</dt>
            <dd>{renderStashLibraries()}</dd>
          </dl>
          <dl>
            <dt>Database file path</dt>
            <dd>
              <code>
                {databaseFile !== ""
                  ? databaseFile
                  : `<path containing configuration file>/stash-go.sqlite`}
              </code>
            </dd>
          </dl>
          <dl>
            <dt>Generated directory</dt>
            <dd>
              <code>
                {generatedLocation !== ""
                  ? generatedLocation
                  : `<path containing configuration file>/generated`}
              </code>
            </dd>
          </dl>
        </section>
        <section className="mt-5">
          <div className="d-flex justify-content-center">
            <Button variant="secondary mx-2 p-5" onClick={() => goBack()}>
              Back
            </Button>
            <Button variant="success mx-2 p-5" onClick={() => onSave()}>
              Confirm
            </Button>
          </div>
        </section>
      </>
    );
  }

  function renderError() {
    return (
      <>
        <section>
          <h2>Oh no! Something went wrong!</h2>
          <p>
            Something went wrong while setting up your system. Here is the error
            we received:
            <pre>{setupError}</pre>
          </p>
          <p>
            If this looks like a problem with your inputs, go ahead and click
            back to fix them up. Otherwise, raise a bug on the {githubLink}
            or seek help in the {discordLink}.
          </p>
        </section>
        <section className="mt-5">
          <div className="d-flex justify-content-center">
            <Button variant="secondary mx-2 p-5" onClick={() => goBack(2)}>
              Back
            </Button>
          </div>
        </section>
      </>
    );
  }

  function renderSuccess() {
    return (
      <>
        <section>
          <h2>Success! Your system has been created!</h2>
          <p>
            You will be taken to the Configuration page next. This page will
            allow you to customize what files to include and exclude, set a
            username and password to protect your system, and a whole bunch of
            other options.
          </p>
          <p>
            When you are satisfied with these settings, you can begin scanning
            your content into Stash by clicking on <code>Tasks</code>, then{" "}
            <code>Scan</code>.
          </p>
        </section>
        <section>
          <h3>Getting help</h3>
          <p>
            You are encouraged to check out the in-app manual which can be
            accessed from the icon in the top-right corner of the screen that
            looks like this: <Icon icon="question-circle" />
          </p>
          <p>
            If you run into issues or have any questions or suggestions, feel
            free to open an issue in the {githubLink}, or ask the community in
            the {discordLink}.
          </p>
        </section>
        <section>
          <h3>Support us</h3>
          <p>
            Check out our{" "}
            <a
              href="https://opencollective.com/stashapp"
              target="_blank"
              rel="noreferrer"
            >
              OpenCollective
            </a>{" "}
            to see how you can contribute to the continued development of Stash.
          </p>
          <p>
            We also welcome contributions in the form of code (bug fixes,
            improvements and new features), testing, bug reports, improvement
            and feature requests, and user support. Details can be found in the
            Contribution section of the in-app manual.
          </p>
        </section>
        <section>
          <p className="lead text-center">Thanks for trying Stash!</p>
        </section>
        <section className="mt-5">
          <div className="d-flex justify-content-center">
            <Link to="/settings?tab=configuration">
              <Button variant="success mx-2 p-5" onClick={() => goBack(2)}>
                Finish
              </Button>
            </Link>
          </div>
        </section>
      </>
    );
  }

  function renderFinish() {
    if (setupError) {
      return renderError();
    }

    return renderSuccess();
  }

  // only display setup wizard if system is not setup
  if (statusLoading) {
    return <LoadingIndicator />;
  }

  if (
    systemStatus &&
    systemStatus.systemStatus.status !== GQL.SystemStatusEnum.Setup
  ) {
    // redirect to main page
    const newURL = new URL("/", window.location.toString());
    window.location.href = newURL.toString();
    return <LoadingIndicator />;
  }

  const welcomeStep =
    systemStatus && systemStatus.systemStatus.configPath !== ""
      ? renderWelcomeSpecificConfig
      : renderWelcome;
  const steps = [welcomeStep, renderSetPaths, renderConfirm, renderFinish];

  return (
    <Container>
      {maybeRenderGeneratedSelectDialog()}
      <h1 className="text-center">Stash Setup Wizard</h1>
      {loading ? (
        <LoadingIndicator message="Creating your system" />
      ) : (
        <Card>{steps[step]()}</Card>
      )}
    </Container>
  );
};
