import React, { useEffect, useState } from "react";
import { Button, Form, InputGroup } from "react-bootstrap";
import * as GQL from "src/core/generated-graphql";
import { StashService } from "src/core/StashService";
import { useToast } from "src/hooks";
import { Icon, LoadingIndicator } from "src/components/Shared";
import { FolderSelect } from "src/components/Shared/FolderSelect/FolderSelect";

export const SettingsConfigurationPanel: React.FC = () => {
  const Toast = useToast();
  // Editing config state
  const [stashes, setStashes] = useState<string[]>([]);
  const [databasePath, setDatabasePath] = useState<string | undefined>(
    undefined
  );
  const [generatedPath, setGeneratedPath] = useState<string | undefined>(
    undefined
  );
  const [maxTranscodeSize, setMaxTranscodeSize] = useState<
    GQL.StreamingResolutionEnum | undefined
  >(undefined);
  const [maxStreamingTranscodeSize, setMaxStreamingTranscodeSize] = useState<
    GQL.StreamingResolutionEnum | undefined
  >(undefined);
  const [username, setUsername] = useState<string | undefined>(undefined);
  const [password, setPassword] = useState<string | undefined>(undefined);
  const [logFile, setLogFile] = useState<string | undefined>();
  const [logOut, setLogOut] = useState<boolean>(true);
  const [logLevel, setLogLevel] = useState<string>("Info");
  const [logAccess, setLogAccess] = useState<boolean>(true);
  const [excludes, setExcludes] = useState<string[]>([]);

  const { data, error, loading } = StashService.useConfiguration();

  const [updateGeneralConfig] = StashService.useConfigureGeneral({
    stashes,
    databasePath,
    generatedPath,
    maxTranscodeSize,
    maxStreamingTranscodeSize,
    username,
    password,
    logFile,
    logOut,
    logLevel,
    logAccess,
    excludes
  });

  useEffect(() => {
    if (!data?.configuration || error) return;

    const conf = data.configuration;
    if (conf.general) {
      setStashes(conf.general.stashes ?? []);
      setDatabasePath(conf.general.databasePath);
      setGeneratedPath(conf.general.generatedPath);
      setMaxTranscodeSize(conf.general.maxTranscodeSize ?? undefined);
      setMaxStreamingTranscodeSize(
        conf.general.maxStreamingTranscodeSize ?? undefined
      );
      setUsername(conf.general.username);
      setPassword(conf.general.password);
      setLogFile(conf.general.logFile ?? undefined);
      setLogOut(conf.general.logOut);
      setLogLevel(conf.general.logLevel);
      setLogAccess(conf.general.logAccess);
      setExcludes(conf.general.excludes);
    }
  }, [data, error]);

  function onStashesChanged(directories: string[]) {
    setStashes(directories);
  }

  function excludeRegexChanged(idx: number, value: string) {
    const newExcludes = excludes.map((regex, i) => {
      const ret = idx !== i ? regex : value;
      return ret;
    });
    setExcludes(newExcludes);
  }

  function excludeRemoveRegex(idx: number) {
    const newExcludes = excludes.filter((_regex, i) => i !== idx);

    setExcludes(newExcludes);
  }

  function excludeAddRegex() {
    const demo = "sample\\.mp4$";
    const newExcludes = excludes.concat(demo);

    setExcludes(newExcludes);
  }

  async function onSave() {
    try {
      const result = await updateGeneralConfig();
      // eslint-disable-next-line no-console
      console.log(result);
      Toast.success({ content: "Updated config" });
    } catch (e) {
      Toast.error(e);
    }
  }

  const transcodeQualities = [
    GQL.StreamingResolutionEnum.Low,
    GQL.StreamingResolutionEnum.Standard,
    GQL.StreamingResolutionEnum.StandardHd,
    GQL.StreamingResolutionEnum.FullHd,
    GQL.StreamingResolutionEnum.FourK,
    GQL.StreamingResolutionEnum.Original
  ].map(resolutionToString);

  function resolutionToString(r: GQL.StreamingResolutionEnum | undefined) {
    switch (r) {
      case GQL.StreamingResolutionEnum.Low:
        return "240p";
      case GQL.StreamingResolutionEnum.Standard:
        return "480p";
      case GQL.StreamingResolutionEnum.StandardHd:
        return "720p";
      case GQL.StreamingResolutionEnum.FullHd:
        return "1080p";
      case GQL.StreamingResolutionEnum.FourK:
        return "4k";
      case GQL.StreamingResolutionEnum.Original:
        return "Original";
    }

    return "Original";
  }

  function translateQuality(quality: string) {
    switch (quality) {
      case "240p":
        return GQL.StreamingResolutionEnum.Low;
      case "480p":
        return GQL.StreamingResolutionEnum.Standard;
      case "720p":
        return GQL.StreamingResolutionEnum.StandardHd;
      case "1080p":
        return GQL.StreamingResolutionEnum.FullHd;
      case "4k":
        return GQL.StreamingResolutionEnum.FourK;
      case "Original":
        return GQL.StreamingResolutionEnum.Original;
    }

    return GQL.StreamingResolutionEnum.Original;
  }

  if (error) return <h1>{error.message}</h1>;
  if (!data?.configuration || loading)
    return <LoadingIndicator />;

  return (
    <>
      <h4>Library</h4>
      <Form.Group>
        <Form.Group id="stashes">
          <Form.Label>Stashes</Form.Label>
          <FolderSelect
            directories={stashes}
            onDirectoriesChanged={onStashesChanged}
          />
          <Form.Text className="text-muted">
            Directory locations to your content
          </Form.Text>
        </Form.Group>

        <Form.Group id="database-path">
          <Form.Label>Database Path</Form.Label>
          <Form.Control
            defaultValue={databasePath}
            onChange={(e: any) => setDatabasePath(e.target.value)}
          />
          <Form.Text className="text-muted">
            File location for the SQLite database (requires restart)
          </Form.Text>
        </Form.Group>

        <Form.Group id="generated-path">
          <Form.Label>Generated Path</Form.Label>
          <Form.Control
            defaultValue={generatedPath}
            onChange={(e: any) => setGeneratedPath(e.target.value)}
          />
          <Form.Text className="text-muted">
            Directory location for the generated files (scene markers, scene
            previews, sprites, etc)
          </Form.Text>
        </Form.Group>

        <Form.Group>
          <Form.Label>Excluded Patterns</Form.Label>
          {excludes
            ? excludes.map((regexp, i) => (
                <InputGroup>
                  <Form.Control
                    value={regexp}
                    onChange={(e: any) =>
                      excludeRegexChanged(i, e.target.value)
                    }
                  />
                  <InputGroup.Append>
                    <Button
                      variant="danger"
                      onClick={() => excludeRemoveRegex(i)}
                    >
                      <Icon icon="minus" />
                    </Button>
                  </InputGroup.Append>
                </InputGroup>
              ))
            : ""}

          <Button variant="danger" onClick={() => excludeAddRegex()}>
            <Icon icon="plus" />
          </Button>
          <div>
            <p>
              <a
                href="https://github.com/stashapp/stash/wiki/Exclude-file-configuration"
                rel="noopener noreferrer"
                target="_blank"
              >
                <span>
                  Regexps of files/paths to exclude from Scan and add to Clean
                </span>
                <Icon icon="question-circle" />
              </a>
            </p>
          </div>
        </Form.Group>
      </Form.Group>

      <hr />

      <Form.Group>
        <h4>Video</h4>
        <Form.Group id="transcode-size">
          <Form.Label>Maximum transcode size</Form.Label>
          <Form.Control
            as="select"
            onChange={(event: React.FormEvent<HTMLSelectElement>) =>
              setMaxTranscodeSize(translateQuality(event.currentTarget.value))
            }
            value={resolutionToString(maxTranscodeSize)}
          >
            {transcodeQualities.map(q => (
              <option key={q} value={q}>
                {q}
              </option>
            ))}
          </Form.Control>
          <Form.Text className="text-muted">
            Maximum size for generated transcodes
          </Form.Text>
        </Form.Group>
        <Form.Group id="streaming-transcode-size">
          <Form.Label>Maximum streaming transcode size</Form.Label>
          <Form.Control
            as="select"
            onChange={(event: React.FormEvent<HTMLSelectElement>) =>
              setMaxStreamingTranscodeSize(
                translateQuality(event.currentTarget.value)
              )
            }
            value={resolutionToString(maxStreamingTranscodeSize)}
          >
            {transcodeQualities.map(q => (
              <option key={q} value={q}>
                {q}
              </option>
            ))}
          </Form.Control>
          <Form.Text className="text-muted">
            Maximum size for transcoded streams
          </Form.Text>
        </Form.Group>
      </Form.Group>

      <hr />

      <Form.Group>
        <h4>Authentication</h4>
        <Form.Group id="username">
          <Form.Label>Username</Form.Label>
          <Form.Control
            defaultValue={username}
            onChange={(e: React.FormEvent<HTMLInputElement>) =>
              setUsername(e.currentTarget.value)
            }
          />
          <Form.Text className="text-muted">
            Username to access Stash. Leave blank to disable user authentication
          </Form.Text>
        </Form.Group>
        <Form.Group id="password">
          <Form.Label>Password</Form.Label>
          <Form.Control
            type="password"
            defaultValue={password}
            onChange={(e: React.FormEvent<HTMLInputElement>) =>
              setPassword(e.currentTarget.value)
            }
          />
          <Form.Text className="text-muted">
            Password to access Stash. Leave blank to disable user authentication
          </Form.Text>
        </Form.Group>
      </Form.Group>

      <hr />

      <h4>Logging</h4>
      <Form.Group id="log-file">
        <Form.Label>Log file</Form.Label>
        <Form.Control
          defaultValue={logFile}
          onChange={(e: React.FormEvent<HTMLInputElement>) =>
            setLogFile(e.currentTarget.value)
          }
        />
        <Form.Text className="text-muted">
          Path to the file to output logging to. Blank to disable file logging.
          Requires restart.
        </Form.Text>
      </Form.Group>

      <Form.Group>
        <Form.Check
          checked={logOut}
          label="Log to terminal"
          onChange={() => setLogOut(!logOut)}
        />
        <Form.Text className="text-muted">
          Logs to the terminal in addition to a file. Always true if file
          logging is disabled. Requires restart.
        </Form.Text>
      </Form.Group>

      <Form.Group id="log-level">
        <Form.Label>Log Level</Form.Label>
        <Form.Control
          as="select"
          onChange={(event: React.FormEvent<HTMLSelectElement>) =>
            setLogLevel(event.currentTarget.value)
          }
          value={logLevel}
        >
          {["Debug", "Info", "Warning", "Error"].map(o => (
            <option key={o} value={o}>
              {o}
            </option>
          ))}
        </Form.Control>
      </Form.Group>

      <Form.Group>
        <Form.Check
          checked={logAccess}
          label="Log http access"
          onChange={() => setLogAccess(!logAccess)}
        />
        <Form.Text className="text-muted">
          Logs http access to the terminal. Requires restart.
        </Form.Text>
      </Form.Group>

      <hr />

      <Button variant="primary" onClick={() => onSave()}>
        Save
      </Button>
    </>
  );
};
