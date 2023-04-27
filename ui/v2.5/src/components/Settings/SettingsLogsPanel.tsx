import React, { useEffect, useState } from "react";
import { useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import { useLoggingSubscribe, queryLogs } from "src/core/StashService";
import { SelectSetting } from "./Inputs";
import { SettingSection } from "./SettingSection";
import { JobTable } from "./Tasks/JobTable";

function convertTime(logEntry: GQL.LogEntryDataFragment) {
  function pad(val: number) {
    let ret = val.toString();
    if (val <= 9) {
      ret = `0${ret}`;
    }

    return ret;
  }

  const date = new Date(logEntry.time);
  const month = date.getMonth() + 1;
  const day = date.getDate();
  let dateStr = `${date.getFullYear()}-${pad(month)}-${pad(day)}`;
  dateStr += ` ${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(
    date.getSeconds()
  )}`;

  return dateStr;
}

function levelClass(level: string) {
  return level.toLowerCase().trim();
}

interface ILogElementProps {
  logEntry: LogEntry;
}

const LogElement: React.FC<ILogElementProps> = ({ logEntry }) => {
  // pad to maximum length of level enum
  const level = logEntry.level.padEnd(GQL.LogLevel.Progress.length);

  return (
    <div className="row">
      <span className="log-time">{logEntry.time}</span>
      <span className={`${levelClass(logEntry.level)}`}>{level}</span>
      <span className="col col-sm-9">{logEntry.message}</span>
    </div>
  );
};

class LogEntry {
  public time: string;
  public level: string;
  public message: string;
  public id: string;

  private static nextId: number = 0;

  public constructor(logEntry: GQL.LogEntryDataFragment) {
    this.time = convertTime(logEntry);
    this.level = logEntry.level;
    this.message = logEntry.message;

    const id = LogEntry.nextId++;
    this.id = id.toString();
  }
}

// maximum number of log entries to display - entries are discarded oldest-first
const MAX_LOG_ENTRIES = 1000;

const logLevels = {
  Trace: GQL.LogLevel.Trace,
  Debug: GQL.LogLevel.Debug,
  Info: GQL.LogLevel.Info,
  Warning: GQL.LogLevel.Warning,
  Error: GQL.LogLevel.Error,
};

export const SettingsLogsPanel: React.FC = () => {
  const intl = useIntl();
  const [entries, setEntries] = useState<LogEntry[]>([]);
  const [logLevel, setLogLevel] = useState(GQL.LogLevel.Info);

  const [subscribe, setSubscribe] = useState(false);
  const { data, error: subscriptionError } = useLoggingSubscribe(
    logLevel,
    !subscribe
  );

  const [error, setError] = useState<Error>();

  function onChangeLogLevel(v: string) {
    const level = logLevels[v as keyof typeof logLevels];
    setLogLevel(level);
  }

  useEffect(() => {
    async function setInitialLogs() {
      let logQuery;
      try {
        logQuery = await queryLogs(logLevel);
      } catch (e) {
        setError(e as Error);
        return;
      }
      if (logQuery.error) {
        setError(logQuery.error);
        return;
      }
      setError(undefined);

      const initEntries = logQuery.data.logs.map((e) => new LogEntry(e));
      setEntries(initEntries.slice(0, MAX_LOG_ENTRIES));
      setSubscribe(true);
    }

    setSubscribe(false);
    setInitialLogs();
  }, [logLevel]);

  useEffect(() => {
    if (subscriptionError) {
      setError(subscriptionError);
    }
  }, [subscriptionError]);

  useEffect(() => {
    if (!data) return;

    const newEntries = data.loggingSubscribe.map((e) => new LogEntry(e));
    if (newEntries.length === 0) return;

    newEntries.reverse();
    setEntries((prev) => {
      return [...newEntries, ...prev].slice(0, MAX_LOG_ENTRIES);
    });
  }, [data]);

  function maybeRenderError() {
    if (error) {
      return (
        <div className="error">
          Error connecting to log server: {error.message}
        </div>
      );
    }
  }

  return (
    <>
      <h2>{intl.formatMessage({ id: "config.tasks.job_queue" })}</h2>
      <JobTable />
      <SettingSection headingID="config.categories.logs">
        <SelectSetting
          id="log-level"
          headingID="config.logs.log_level"
          value={logLevel}
          onChange={onChangeLogLevel}
        >
          {Object.keys(logLevels).map((level) => (
            <option key={level} value={level}>
              {level}
            </option>
          ))}
        </SelectSetting>
      </SettingSection>

      <div className="logs">
        {maybeRenderError()}
        {entries.map((logEntry) => (
          <LogElement logEntry={logEntry} key={logEntry.id} />
        ))}
      </div>
    </>
  );
};
