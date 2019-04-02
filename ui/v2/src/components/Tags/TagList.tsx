import { Button, Classes, Dialog, EditableText, FormGroup, HTMLTable, InputGroup, Spinner, Tag } from "@blueprintjs/core";
import _ from "lodash";
import React, { FunctionComponent, useEffect, useState } from "react";
import { QueryHookResult } from "react-apollo-hooks";
import { Link } from "react-router-dom";
import { FindGalleriesQuery, FindGalleriesVariables } from "../../core/generated-graphql";
import * as GQL from "../../core/generated-graphql";
import { StashService } from "../../core/StashService";
import { ListHook } from "../../hooks/ListHook";
import { IBaseProps } from "../../models/base-props";
import { ListFilterModel } from "../../models/list-filter/filter";
import { DisplayMode, FilterMode } from "../../models/list-filter/types";
import { ErrorUtils } from "../../utils/errors";
import { NavigationUtils } from "../../utils/navigation";
import { ToastUtils } from "../../utils/toasts";

interface IProps extends IBaseProps {}

export const TagList: FunctionComponent<IProps> = (props: IProps) => {
  const [tags, setTags] = useState<GQL.AllTagsAllTags[]>([]);
  const [isLoading, setIsLoading] = useState(false);

  // Editing / New state
  const [editingTag, setEditingTag] = useState<Partial<GQL.TagDataFragment> | undefined>(undefined);
  const [name, setName] = useState<string>("");

  const { data, error, loading } = StashService.useAllTags();
  const updateTag = StashService.useTagUpdate(getTagInput() as GQL.TagUpdateInput);
  const createTag = StashService.useTagCreate(getTagInput() as GQL.TagCreateInput);

  useEffect(() => {
    setIsLoading(loading);
    if (!data || !data.allTags || !!error) { return; }
    setTags(data.allTags);
  }, [data]);

  useEffect(() => {
    if (!!editingTag) {
      setName(editingTag.name || "");
    } else {
      setName("");
    }
  }, [editingTag]);

  function getTagInput() {
    const tagInput: Partial<GQL.TagCreateInput | GQL.TagUpdateInput> = { name };
    if (!!editingTag) { (tagInput as Partial<GQL.TagUpdateInput>).id = editingTag.id; }
    return tagInput;
  }

  async function onEdit() {
    try {
      if (!!editingTag && !!editingTag.id) {
        await updateTag();
        ToastUtils.success("Updated tag");
      } else {
        await createTag();
        ToastUtils.success("Created tag");
      }
      setEditingTag(undefined);
    } catch (e) {
      ErrorUtils.handle(e);
    }
  }

  if (!data || !data.allTags || isLoading) { return <Spinner size={Spinner.SIZE_LARGE} />; }
  if (!!error) { return <>{error.message}</>; }

  const tagElements = tags.map((tag) => {
    return (
      <div key={tag.id} className="tag-list-row">
        <span onClick={() => setEditingTag(tag)}>{tag.name}</span>
        <div style={{float: "right"}}>
          <Link className="bp3-button" to={NavigationUtils.makeTagScenesUrl(tag)}>Scenes: {tag.scene_count}</Link>
          <Link className="bp3-button" to={NavigationUtils.makeTagSceneMarkersUrl(tag)}>
            Markers: {tag.scene_marker_count}
          </Link>
          <span>Total: {(tag.scene_count || 0) + (tag.scene_marker_count || 0)}</span>
        </div>
      </div>
    );
  });
  return (
    <div id="tag-list-container">
      <Button intent="primary" style={{marginTop: "20px"}} onClick={() => setEditingTag({})}>New Tag</Button>
      <Dialog
        isOpen={!!editingTag}
        onClose={() => setEditingTag(undefined)}
        title={!!editingTag && !!editingTag.id ? "Edit Tag" : "New Tag"}
      >
        <div className="dialog-content">
          <FormGroup label="Name">
            <InputGroup
              onChange={(newValue: any) => setName(newValue.target.value)}
              value={name}
            />
          </FormGroup>
        </div>
        <div className={Classes.DIALOG_FOOTER}>
          <div className={Classes.DIALOG_FOOTER_ACTIONS}>
            <Button onClick={() => onEdit()}>{!!editingTag && !!editingTag.id ? "Update" : "Create"}</Button>
          </div>
        </div>
      </Dialog>

      {tagElements}
    </div>
  );
};
