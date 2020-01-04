import React from "react";
import { Table, Spinner } from 'react-bootstrap';
import { StashService } from "../../core/StashService";

export const SettingsAboutPanel: React.FC = () => {
  const { data, error, loading } = StashService.useVersion();

  function maybeRenderTag() {
    if (!data || !data.version || !data.version.version) { return; }
    return (
      <tr>
        <td>Version:</td>
        <td>{data.version.version}</td>
      </tr>
    );
  }

  function renderVersion() {
    if (!data || !data.version) { return; }
    return (
      <>
      <Table>
        <tbody>
          {maybeRenderTag()}
          <tr>
            <td>Build hash:</td>
            <td>{data.version.hash}</td>
          </tr>
          <tr>
            <td>Build time:</td>
            <td>{data.version.build_time}</td>
          </tr>
        </tbody>  
      </Table>
      </>
    );
  }
  return (
    <>
      <h4>About</h4>
      {!data || loading ? <Spinner animation="border" variant="light" /> : undefined}
      {!!error ? <span>error.message</span> : undefined}
      {renderVersion()}
    </>
  );
};
