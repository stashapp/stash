import React, { useEffect, useState } from "react";
import { Button, Form, InputGroup } from "react-bootstrap";
import * as GQL from "src/core/generated-graphql";
import { useConfiguration, useConfigureGeneral } from "src/core/StashService";
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
  const [cachePath, setCachePath] = useState<string | undefined>(undefined);
  const [maxTranscodeSize, setMaxTranscodeSize] = useState<
    GQL.StreamingResolutionEnum | undefined
  >(undefined);
  const [maxStreamingTranscodeSize, setMaxStreamingTranscodeSize] = useState<
    GQL.StreamingResolutionEnum | undefined
  >(undefined);
  const [forceMkv, setForceMkv] = useState<boolean>(false);
  const [forceHevc, setForceHevc] = useState<boolean>(false);
  const [username, setUsername] = useState<string | undefined>(undefined);
  const [password, setPassword] = useState<string | undefined>(undefined);
  const [maxSessionAge, setMaxSessionAge] = useState<number>(0);
  const [logFile, setLogFile] = useState<string | undefined>();
  const [logOut, setLogOut] = useState<boolean>(true);
  const [logLevel, setLogLevel] = useState<string>("Info");
  const [logAccess, setLogAccess] = useState<boolean>(true);
  const [excludes, setExcludes] = useState<string[]>([]);
  const [scraperUserAgent, setScraperUserAgent] = useState<string | undefined>(
    undefined
  );

  const { data, error, loading } = useConfiguration();

  const [updateGeneralConfig] = useConfigureGeneral({
    stashes,
    databasePath,
    generatedPath,
    cachePath,
    maxTranscodeSize,
    maxStreamingTranscodeSize,
    forceMkv,
    forceHevc,
    username,
    password,
    maxSessionAge,
    logFile,
    logOut,
    logLevel,
    logAccess,
    excludes,
    scraperUserAgent,
  });

  useEffect(() => {
    if (!data?.configuration || error) return;

    const conf = data.configuration;
    if (conf.general) {
      setStashes(conf.general.stashes ?? []);
      setDatabasePath(conf.general.databasePath);
      setGeneratedPath(conf.general.generatedPath);
      setCachePath(conf.general.cachePath);
      setMaxTranscodeSize(conf.general.maxTranscodeSize ?? undefined);
      setMaxStreamingTranscodeSize(
        conf.general.maxStreamingTranscodeSize ?? undefined
      );
      setForceMkv(conf.general.forceMkv);
      setForceHevc(conf.general.forceHevc);
      setUsername(conf.general.username);
      setPassword(conf.general.password);
      setMaxSessionAge(conf.general.maxSessionAge);
      setLogFile(conf.general.logFile ?? undefined);
      setLogOut(conf.general.logOut);
      setLogLevel(conf.general.logLevel);
      setLogAccess(conf.general.logAccess);
      setExcludes(conf.general.excludes);
      setScraperUserAgent(conf.general.scraperUserAgent ?? undefined);
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
    GQL.StreamingResolutionEnum.Original,
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
  if (!data?.configuration || loading) return <LoadingIndicator />;

  return (
    <>
      <h4>Library</h4>
      <Form.Group>
        <Form.Group id="stashes">
          <h6>Stashes</h6>
          <FolderSelect
            directories={stashes}
            onDirectoriesChanged={onStashesChanged}
          />
          <Form.Text className="text-muted">
            Directory locations to your content
          </Form.Text>
        </Form.Group>

        <Form.Group id="database-path">
          <h6>Database Path</h6>
          <Form.Control
            className="col col-sm-6 text-input"
            defaultValue={databasePath}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
              setDatabasePath(e.currentTarget.value)
            }
          />
          <Form.Text className="text-muted">
            File location for the SQLite database (requires restart)
          </Form.Text>
        </Form.Group>

        <Form.Group id="generated-path">
          <h6>Generated Path</h6>
          <Form.Control
            className="col col-sm-6 text-input"
            defaultValue={generatedPath}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
              setGeneratedPath(e.currentTarget.value)
            }
          />
          <Form.Text className="text-muted">
            Directory location for the generated files (scene markers, scene
            previews, sprites, etc)
          </Form.Text>
        </Form.Group>

        <Form.Group id="cache-path">
          <h6>Cache Path</h6>
          <Form.Control
            className="col col-sm-6 text-input"
            defaultValue={cachePath}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
              setCachePath(e.currentTarget.value)
            }
          />
          <Form.Text className="text-muted">
            Directory location of the cache
          </Form.Text>
        </Form.Group>

        <Form.Group>
          <h6>Excluded Patterns</h6>
          <Form.Group>
            {excludes &&
              excludes.map((regexp, i) => (
                <InputGroup>
                  <Form.Control
                    className="col col-sm-6 text-input"
                    value={regexp}
                    onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                      excludeRegexChanged(i, e.currentTarget.value)
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
              ))}
          </Form.Group>
          <Button className="minimal" onClick={() => excludeAddRegex()}>
            <Icon icon="plus" />
          </Button>
          <Form.Text className="text-muted">
            Regexps of files/paths to exclude from Scan and add to Clean
            <a
              href="https://github.com/stashapp/stash/wiki/Exclude-file-configuration"
              rel="noopener noreferrer"
              target="_blank"
            >
              <Icon icon="question-circle" />
            </a>
          </Form.Text>
        </Form.Group>
      </Form.Group>

      <hr />

      <Form.Group>
        <h4>Video</h4>
        <Form.Group id="transcode-size">
          <h6>Maximum transcode size</h6>
          <Form.Control
            className="col col-sm-6 input-control"
            as="select"
            onChange={(event: React.ChangeEvent<HTMLSelectElement>) =>
              setMaxTranscodeSize(translateQuality(event.currentTarget.value))
            }
            value={resolutionToString(maxTranscodeSize)}
          >
            {transcodeQualities.map((q) => (
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
          <h6>Maximum streaming transcode size</h6>
          <Form.Control
            className="col col-sm-6 input-control"
            as="select"
            onChange={(event: React.ChangeEvent<HTMLSelectElement>) =>
              setMaxStreamingTranscodeSize(
                translateQuality(event.currentTarget.value)
              )
            }
            value={resolutionToString(maxStreamingTranscodeSize)}
          >
            {transcodeQualities.map((q) => (
              <option key={q} value={q}>
                {q}
              </option>
            ))}
          </Form.Control>
          <Form.Text className="text-muted">
            Maximum size for transcoded streams
          </Form.Text>
        </Form.Group>
        <Form.Group id="force-options-mkv">
          <Form.Check
            id="force-mkv"
            checked={forceMkv}
            label="Force Matroska as supported"
            onChange={() => setForceMkv(!forceMkv)}
          />
          <Form.Text className="text-muted">
            Treat Matroska (MKV) as a supported container. Recommended for
            Chromium based browsers
          </Form.Text>
        </Form.Group>
        <Form.Group id="force-options-hevc">
          <Form.Check
            id="force-hevc"
            checked={forceHevc}
            label="Force HEVC as supported"
            onChange={() => setForceHevc(!forceHevc)}
          />
          <Form.Text className="text-muted">
            Treat HEVC as a supported codec. Recommended for Safari or some
            Android based browsers
          </Form.Text>
        </Form.Group>
      </Form.Group>

      <hr />

      <Form.Group id="generated-path">
        <h6>Scraping</h6>
        <Form.Control
          className="col col-sm-6 text-input"
          defaultValue={scraperUserAgent}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            setScraperUserAgent(e.currentTarget.value)
          }
        />
        <Form.Text className="text-muted">
          User-Agent string used during scrape http requests
        </Form.Text>
      </Form.Group>

      <hr />

      <Form.Group>
        <h4>Authentication</h4>
        <Form.Group id="username">
          <h6>Username</h6>
          <Form.Control
            className="col col-sm-6 text-input"
            defaultValue={username}
            onInput={(e: React.FormEvent<HTMLInputElement>) =>
              setUsername(e.currentTarget.value)
            }
          />
          <Form.Text className="text-muted">
            Username to access Stash. Leave blank to disable user authentication
          </Form.Text>
        </Form.Group>
        <Form.Group id="password">
          <h6>Password</h6>
          <Form.Control
            className="col col-sm-6 text-input"
            type="password"
            defaultValue={password}
            onInput={(e: React.FormEvent<HTMLInputElement>) =>
              setPassword(e.currentTarget.value)
            }
          />
          <Form.Text className="text-muted">
            Password to access Stash. Leave blank to disable user authentication
          </Form.Text>
        </Form.Group>

        <Form.Group id="maxSessionAge">
          <h6>Maximum Session Age</h6>
          <Form.Control
            className="col col-sm-6 text-input"
            type="number"
            value={maxSessionAge.toString()}
            onInput={(e: React.FormEvent<HTMLInputElement>) =>
              setMaxSessionAge(Number.parseInt(e.currentTarget.value, 10))
            }
          />
          <Form.Text className="text-muted">
            Maximum idle time before a login session is expired, in seconds.
          </Form.Text>
        </Form.Group>
      </Form.Group>

      <hr />

      <h4>Logging</h4>
      <Form.Group id="log-file">
        <h6>Log file</h6>
        <Form.Control
          className="col col-sm-6 text-input"
          defaultValue={logFile}
          onInput={(e: React.FormEvent<HTMLInputElement>) =>
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
          id="log-terminal"
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
        <h6>Log Level</h6>
        <Form.Control
          className="col col-sm-6 input-control"
          as="select"
          onChange={(event: React.ChangeEvent<HTMLSelectElement>) =>
            setLogLevel(event.currentTarget.value)
          }
          value={logLevel}
        >
          {["Trace", "Debug", "Info", "Warning", "Error"].map((o) => (
            <option key={o} value={o}>
              {o}
            </option>
          ))}
        </Form.Control>
      </Form.Group>

      <Form.Group>
        <Form.Check
          id="log-http"
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
