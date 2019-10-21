import {
  H4,
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

  public constructor(logEntry: GQL.LogEntryDataFragment) {
    this.time = convertTime(logEntry);
    
    // pad to maximum length of level enum
    this.level = logEntry.level.padEnd(GQL.LogLevel.Progress.length);
    
    this.message = logEntry.message;
  }
}

export const SettingsLogsPanel: FunctionComponent<IProps> = (props: IProps) => {
  const { data } = StashService.useLoggingSubscribe();
  const { data: existingData } = StashService.useLogs();
  
  const [existingLogEntries, setExistingLogEntries] = useState<LogEntry[]>([]);
  const [logEntries, setLogEntries] = useState<LogEntry[]>([]);

  useEffect(() => {
    if (!data) { return; }

    // append data to the logEntries
    var convertedData = data.loggingSubscribe.map(convertLogEntry);
    var newLogEntries = logEntries.concat(convertedData);

    setLogEntries(newLogEntries);
  }, [data]);

  useEffect(() => {
    if (!existingData || !existingData.logs) { return; }

    var convertedData = existingData.logs.map(convertLogEntry);
    setExistingLogEntries(convertedData);
  }, [existingData])

  function convertLogEntry(logEntry : GQL.LogEntryDataFragment) {
    return new LogEntry(logEntry);
  }

  function levelClass(level : string) {
    return level.toLowerCase().trim();
  }

  function renderLogEntry(logEntry : LogEntry) {
    return (
      <>
        <span>{logEntry.time}</span>&nbsp;
        <span className={levelClass(logEntry.level)}>{logEntry.level}</span>&nbsp;
        <span>{logEntry.message}</span>
        <br/>
      </>
    );
  }

  return (
    <>
      <H4>Logs</H4>
      <div className="logs">
        {existingLogEntries.map(renderLogEntry)}
        {logEntries.map(renderLogEntry)}
      </div>
    </>
  );
};
