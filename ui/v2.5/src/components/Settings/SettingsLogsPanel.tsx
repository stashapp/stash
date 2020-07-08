import React, { useEffect, useReducer, useState } from "react";
import { Form } from "react-bootstrap";
import * as GQL from "src/core/generated-graphql";
import { useLogs, useLoggingSubscribe } from "src/core/StashService";

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

// maximum number of log entries to display. Subsequent entries will truncate
// the list, dropping off the oldest entries first.
const MAX_LOG_ENTRIES = 200;
const logLevels = ["Trace", "Debug", "Info", "Warning", "Error"];

const logReducer = (existingEntries: LogEntry[], newEntries: LogEntry[]) => [
  ...newEntries.reverse(),
  ...existingEntries,
];

export const SettingsLogsPanel: React.FC = () => {
  const { data, error } = useLoggingSubscribe();
  const { data: existingData } = useLogs();
  const [currentData, dispatchLogUpdate] = useReducer(logReducer, []);
  const [logLevel, setLogLevel] = useState<string>("Info");

  useEffect(() => {
    const newData = (data?.loggingSubscribe ?? []).map((e) => new LogEntry(e));
    dispatchLogUpdate(newData);
  }, [data]);

  const oldData = (existingData?.logs ?? []).map((e) => new LogEntry(e));
  const filteredLogEntries = [...currentData, ...oldData]
    .filter(filterByLogLevel)
    .slice(0, MAX_LOG_ENTRIES);

  const maybeRenderError = error ? (
    <div className="error">Error connecting to log server: {error.message}</div>
  ) : (
    ""
  );

  function filterByLogLevel(logEntry: LogEntry) {
    if (logLevel === "Trace") return true;

    const logLevelIndex = logLevels.indexOf(logLevel);
    const levelIndex = logLevels.indexOf(logEntry.level);

    return levelIndex >= logLevelIndex;
  }

  return (
    <>
      <h4>Logs</h4>
      <Form.Row id="log-level">
        <Form.Label className="col-6 col-sm-2">Log Level</Form.Label>
        <Form.Control
          className="col-6 col-sm-2 input-control"
          as="select"
          defaultValue={logLevel}
          onChange={(event) => setLogLevel(event.currentTarget.value)}
        >
          {logLevels.map((level) => (
            <option key={level} value={level}>
              {level}
            </option>
          ))}
        </Form.Control>
      </Form.Row>
      <div className="logs">
        {maybeRenderError}
        {filteredLogEntries.map((logEntry) => (
          <LogElement logEntry={logEntry} key={logEntry.id} />
        ))}
      </div>
    </>
  );
};
