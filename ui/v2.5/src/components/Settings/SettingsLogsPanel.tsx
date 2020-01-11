import React, { useState, useEffect, useRef } from "react";
import { Form, Col } from 'react-bootstrap';
import * as GQL from "src/core/generated-graphql";
import { StashService } from "src/core/StashService";

function convertTime(logEntry: GQL.LogEntryDataFragment) {
  function pad(val : number) {
    var ret = val.toString();
    if (val <= 9) {
      ret = "0" + ret;
    }

    return ret;
  }

  var date = new Date(logEntry.time);
  var month = date.getMonth() + 1;
  var day = date.getDate();
  var dateStr = date.getFullYear() + "-" + pad(month) + "-" + pad(day);
  dateStr += " " + pad(date.getHours()) + ":" + pad(date.getMinutes()) + ":" + pad(date.getSeconds());

  return dateStr;
}

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

    var id = LogEntry.nextId++;
    this.id = id.toString();
  }
}

export const SettingsLogsPanel: React.FC = () => {
  const { data, error } = StashService.useLoggingSubscribe();
  const { data: existingData } = StashService.useLogs();

  const logEntries = useRef<LogEntry[]>([]);
  const [logLevel, setLogLevel] = useState<string>("Info");
  const [filteredLogEntries, setFilteredLogEntries] = useState<LogEntry[]>([]);
  const lastUpdate = useRef<number>(0);
  const updateTimeout = useRef<NodeJS.Timeout>();

  // maximum number of log entries to display. Subsequent entries will truncate
  // the list, dropping off the oldest entries first.
  const MAX_LOG_ENTRIES = 200;

  function truncateLogEntries(entries : LogEntry[]) {
    entries.length = Math.min(entries.length, MAX_LOG_ENTRIES);
  }

  function prependLogEntries(toPrepend : LogEntry[]) {
    var newLogEntries = toPrepend.concat(logEntries.current);
    truncateLogEntries(newLogEntries);
    logEntries.current = newLogEntries;
  }

  function appendLogEntries(toAppend : LogEntry[]) {
    var newLogEntries = logEntries.current.concat(toAppend);
    truncateLogEntries(newLogEntries);
    logEntries.current = newLogEntries;
  }

  useEffect(() => {
    if (!data) { return; }

    // append data to the logEntries
    var convertedData = data.loggingSubscribe.map(convertLogEntry);

    // filter subscribed data as it comes in, otherwise we'll end up
    // truncating stuff that wasn't filtered out
    convertedData = convertedData.filter(filterByLogLevel)

    // put newest entries at the top
    convertedData.reverse();
    prependLogEntries(convertedData);

    updateFilteredEntries();
  }, [data]);

  useEffect(() => {
    if (!existingData || !existingData.logs) { return; }

    var convertedData = existingData.logs.map(convertLogEntry);
    appendLogEntries(convertedData);

    updateFilteredEntries();
  }, [existingData]);

  function updateFilteredEntries() {
    if (!updateTimeout.current) {
      console.log("Updating after timeout");
    }
    updateTimeout.current = undefined;

    var filteredEntries = logEntries.current.filter(filterByLogLevel);
    setFilteredLogEntries(filteredEntries);

    lastUpdate.current = new Date().getTime();
  }

  useEffect(() => {
    updateFilteredEntries();
  }, [logLevel]);

  function convertLogEntry(logEntry : GQL.LogEntryDataFragment) {
    return new LogEntry(logEntry);
  }

  function levelClass(level : string) {
    return level.toLowerCase().trim();
  }

  interface ILogElementProps {
    logEntry : LogEntry
  }

  function LogElement(props : ILogElementProps) {
    // pad to maximum length of level enum
    var level = props.logEntry.level.padEnd(GQL.LogLevel.Progress.length);

    return (
      <>
        <span>{props.logEntry.time}</span>&nbsp;
        <span className={levelClass(props.logEntry.level)}>{level}</span>&nbsp;
        <span>{props.logEntry.message}</span>
        <br/>
      </>
    );
  }

  function maybeRenderError() {
    if (error) {
      return (
        <>
        <span className={"error"}>Error connecting to log server: {error.message}</span><br/>
        </>
      );
    }
  }

  const logLevels = ["Debug", "Info", "Warning", "Error"];

  function filterByLogLevel(logEntry : LogEntry) {
    if (logLevel === "Debug") {
      return true;
    }

    var logLevelIndex = logLevels.indexOf(logLevel);
    var levelIndex = logLevels.indexOf(logEntry.level);

    return levelIndex >= logLevelIndex;
  }

  return (
    <>
      <h4>Logs</h4>
      <Form.Row id="log-level">
        <Col xs={1}>
          <Form.Label>Log Level</Form.Label>
        </Col>
        <Col xs={2}>
          <Form.Control
            as="select"
            defaultValue={logLevel}
            onChange={(event) => setLogLevel(event.currentTarget.value)}
          >
              { logLevels.map(level => (<option key={level} value={level}>{level}</option>)) }
          </Form.Control>
        </Col>
      </Form.Row>
      <div className="logs">
        {maybeRenderError()}
        {filteredLogEntries.map((logEntry) =>
          <LogElement logEntry={logEntry} key={logEntry.id}/>
        )}
      </div>
    </>
  );
};
