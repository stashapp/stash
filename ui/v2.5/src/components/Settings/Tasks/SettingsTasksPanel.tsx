import React, { useState } from "react";
import { useIntl } from "react-intl";
import { LoadingIndicator } from "src/components/Shared";
import { LibraryTasks } from "./LibraryTasks";
import { DataManagementTasks } from "./DataManagementTasks";
import { PluginTasks } from "./PluginTasks";
import { JobTable } from "./JobTable";

export const SettingsTasksPanel: React.FC = () => {
  const intl = useIntl();
  const [isBackupRunning, setIsBackupRunning] = useState<boolean>(false);

  if (isBackupRunning) {
    return (
      <LoadingIndicator
        message={intl.formatMessage({ id: "config.tasks.backing_up_database" })}
      />
    );
  }

  return (
    <>
      <h4>{intl.formatMessage({ id: "config.tasks.job_queue" })}</h4>

      <JobTable />

      <hr />

      <LibraryTasks />
      <hr />
      <DataManagementTasks setIsBackupRunning={setIsBackupRunning} />
      <hr />
      <PluginTasks />
    </>
  );
};
