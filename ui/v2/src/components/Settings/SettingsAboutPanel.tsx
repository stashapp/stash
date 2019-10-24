import {
  H1,
  H4,
  H6,
  HTMLTable,
  Spinner,
  Tag,
} from "@blueprintjs/core";
import React, { FunctionComponent } from "react";
import * as GQL from "../../core/generated-graphql";
import { TextUtils } from "../../utils/text";
import { StashService } from "../../core/StashService";

interface IProps {}

export const SettingsAboutPanel: FunctionComponent<IProps> = (props: IProps) => {
  const { data, error, loading } = StashService.useVersion();

  function renderVersion() {
    if (!data || !data.version) { return; }
    return (
      <>
      <HTMLTable>
        <tbody>
          <tr>
            <td>Build hash:</td>
            <td>{data.version.hash}</td>
          </tr>
          <tr>
            <td>Build time:</td>
            <td>{data.version.build_time}</td>
          </tr>
        </tbody>  
      </HTMLTable>
      </>
    );
  }
  return (
    <>
      <H4>About</H4>
      {!data || loading ? <Spinner size={Spinner.SIZE_LARGE} /> : undefined}
      {!!error ? <span>error.message</span> : undefined}
      {renderVersion()}
    </>
  );
};
