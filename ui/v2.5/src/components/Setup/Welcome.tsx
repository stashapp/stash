import React, { useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { Alert, Button, Card, Container, Form } from "react-bootstrap";
import * as GQL from "src/core/generated-graphql";
import { useConfigureUI, useSystemStatus } from "src/core/StashService";
import { useHistory } from "react-router-dom";
import { useConfigurationContext } from "src/hooks/Config";
import { Icon } from "../Shared/Icon";
import { LoadingIndicator } from "../Shared/LoadingIndicator";
import { faQuestionCircle } from "@fortawesome/free-solid-svg-icons";
import { releaseNotes } from "src/docs/en/ReleaseNotes";
import { ExternalLink } from "../Shared/ExternalLink";

const DiscordLink = (
  <ExternalLink href="https://discord.gg/2TsNFKt">Discord</ExternalLink>
);
const GithubLink = (
  <ExternalLink href="https://github.com/stashapp/stash/issues">
    <FormattedMessage id="setup.github_repository" />
  </ExternalLink>
);

const SuccessStep: React.FC<{
  configuration: GQL.ConfigDataFragment;
  systemStatus: GQL.SystemStatusQuery;
}> = ({ configuration, systemStatus }) => {
  const intl = useIntl();
  const history = useHistory();

  const [mutateDownloadFFMpeg] = GQL.useDownloadFfMpegMutation();
  const [saveUI] = useConfigureUI();

  const [downloadFFmpeg, setDownloadFFmpeg] = useState(true);

  const status = systemStatus?.systemStatus;

  async function markReleaseNotesSeen() {
    try {
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
      // ignore
    }
  }

  async function onFinishClick() {
    await markReleaseNotesSeen();

    if ((!status?.ffmpegPath || !status?.ffprobePath) && downloadFFmpeg) {
      await mutateDownloadFFMpeg();
    }

    history.replace("/settings?tab=library");
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

export const Welcome: React.FC = () => {
  const { configuration } = useConfigurationContext();

  const {
    data: systemStatus,
    loading: statusLoading,
    error: statusError,
  } = useSystemStatus();

  if (statusLoading) {
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
    <Container className="setup-wizard">
      <h1 className="text-center">
        <FormattedMessage id="setup.stash_setup_wizard" />
      </h1>
      <Card>
        <SuccessStep
          systemStatus={systemStatus}
          configuration={configuration}
        />
      </Card>
    </Container>
  );
};

export default Welcome;
