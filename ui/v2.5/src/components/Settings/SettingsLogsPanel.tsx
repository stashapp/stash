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

// maximum number of log entries to keep - entries are discarded oldest-first
const MAX_LOG_ENTRIES = 50000;
// maximum number of log entries to display
const MAX_DISPLAY_LOG_ENTRIES = 1000;
const logLevels = ["Trace", "Debug", "Info", "Warning", "Error"];

export const SettingsLogsPanel: React.FC = () => {
  const [entries, setEntries] = useState<LogEntry[]>([]);
  const { data, error } = useLoggingSubscribe();
  const [logLevel, setLogLevel] = useState<string>("Info");
  const intl = useIntl();

  useEffect(() => {
    async function getInitialLogs() {
      const logQuery = await queryLogs();
      if (logQuery.error) return;

      const initEntries = logQuery.data.logs.map((e) => new LogEntry(e));
      if (initEntries.length !== 0) {
        setEntries((prev) => {
          return [...prev, ...initEntries].slice(0, MAX_LOG_ENTRIES);
        });
      }
    }

    getInitialLogs();
  }, []);

  useEffect(() => {
    if (!data) return;

    const newEntries = data.loggingSubscribe.map((e) => new LogEntry(e));
    newEntries.reverse();
    setEntries((prev) => {
      return [...newEntries, ...prev].slice(0, MAX_LOG_ENTRIES);
    });
  }, [data]);

  const displayEntries = entries
    .filter(filterByLogLevel)
    .slice(0, MAX_DISPLAY_LOG_ENTRIES);

  function maybeRenderError() {
    if (error) {
      return (
        <div className="error">
          Error connecting to log server: {error.message}
        </div>
      );
    }
  }

  function filterByLogLevel(logEntry: LogEntry) {
    if (logLevel === "Trace") return true;

    const logLevelIndex = logLevels.indexOf(logLevel);
    const levelIndex = logLevels.indexOf(logEntry.level);

    return levelIndex >= logLevelIndex;
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
          onChange={(v) => setLogLevel(v)}
        >
          {logLevels.map((level) => (
            <option key={level} value={level}>
              {level}
            </option>
          ))}
        </SelectSetting>
      </SettingSection>

      <div className="logs">
        {maybeRenderError()}
        {displayEntries.map((logEntry) => (
          <LogElement logEntry={logEntry} key={logEntry.id} />
        ))}
      </div>
    </>
  );
};
