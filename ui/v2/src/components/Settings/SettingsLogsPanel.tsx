import {
  H4, FormGroup, HTMLSelect,
} from "@blueprintjs/core";
import React, { FunctionComponent, useState, useEffect } from "react";
import * as GQL from "../../core/generated-graphql";
import { StashService } from "../../core/StashService";

interface IProps {}

function convertTime(logEntry : GQL.LogEntryDataFragment) {
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

export const SettingsLogsPanel: FunctionComponent<IProps> = (props: IProps) => {
  const { data, error } = StashService.useLoggingSubscribe();
  const { data: existingData } = StashService.useLogs();
  
  const [existingLogEntries, setExistingLogEntries] = useState<LogEntry[]>([]);
  const [logEntries, setLogEntries] = useState<LogEntry[]>([]);
  const [logLevel, setLogLevel] = useState<string>("Info");
  const [filteredLogEntries, setFilteredLogEntries] = useState<LogEntry[]>([]);

  useEffect(() => {
    if (!data) { return; }

    // append data to the logEntries
    var convertedData = data.loggingSubscribe.map(convertLogEntry);

    // put newest entries at the top
    convertedData.reverse();
    var newLogEntries = convertedData.concat(logEntries);

    setLogEntries(newLogEntries);
  }, [data]);

  useEffect(() => {
    if (!existingData || !existingData.logs) { return; }

    var convertedData = existingData.logs.map(convertLogEntry);
    setExistingLogEntries(convertedData);
  }, [existingData]);

  useEffect(() => {
    // concatenate and filter the log entries
    
    var filteredEntries : LogEntry[] = [];
    
    if (logEntries) {
      filteredEntries = filteredEntries.concat(logEntries.filter(filterByLogLevel));
    }

    if (existingLogEntries) {
      filteredEntries = filteredEntries.concat(existingLogEntries.filter(filterByLogLevel));
    }
    
    setFilteredLogEntries(filteredEntries);
  }, [logLevel, logEntries, existingLogEntries]);

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
    if (logLevel == "Debug") {
      return true;
    }

    var logLevelIndex = logLevels.indexOf(logLevel);
    var levelIndex = logLevels.indexOf(logEntry.level);

    return levelIndex >= logLevelIndex;
  }

  return (
    <>
      <H4>Logs</H4>
      <div>
      <FormGroup inline={true} label="Log Level">
        <HTMLSelect
          options={logLevels}
          onChange={(event) => setLogLevel(event.target.value)}
          value={logLevel}
        />
        </FormGroup>
      </div>
      <div className="logs">
        {maybeRenderError()}
        {filteredLogEntries.map((logEntry) =>
          <LogElement logEntry={logEntry} key={logEntry.id}/>
        )}
      </div>
    </>
  );
};
