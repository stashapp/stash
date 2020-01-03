import { Badge, Button, Card, Collapse, Form as BootstrapForm } from 'react-bootstrap';
import { Field, FieldProps, Form, Formik, FormikActions, FormikProps } from "formik";
import React, { CSSProperties, FunctionComponent, useState } from "react";
import * as GQL from "../../../core/generated-graphql";
import { StashService } from "../../../core/StashService";
import { TextUtils } from "../../../utils/text";
import { FilterMultiSelect } from "../../select/FilterMultiSelect";
import { FilterSelect } from "../../select/FilterSelect";
import { MarkerTitleSuggest } from "../../select/MarkerTitleSuggest";
import { WallPanel } from "../../Wall/WallPanel";
import { SceneHelpers } from "../helpers";
import { ErrorUtils } from "../../../utils/errors";
import { DurationInput } from "../../Shared/DurationInput";

interface ISceneMarkersPanelProps {
  scene: GQL.SceneDataFragment;
  onClickMarker: (marker: GQL.SceneMarkerDataFragment) => void;
}

interface IFormFields {
  title: string;
  seconds: string;
  primaryTagId: string;
  tagIds: string[];
}

export const SceneMarkersPanel: FunctionComponent<ISceneMarkersPanelProps> = (props: ISceneMarkersPanelProps) => {
  const [isEditorOpen, setIsEditorOpen] = useState<boolean>(false);
  const [editingMarker, setEditingMarker] = useState<GQL.SceneMarkerDataFragment | null>(null);

  const sceneMarkerCreate = StashService.useSceneMarkerCreate();
  const sceneMarkerUpdate = StashService.useSceneMarkerUpdate();
  const sceneMarkerDestroy = StashService.useSceneMarkerDestroy();

  const jwplayer = SceneHelpers.getPlayer();

  function onOpenEditor(marker: GQL.SceneMarkerDataFragment | null = null) {
    setIsEditorOpen(true);
    setEditingMarker(marker);
  }

  function onClickMarker(marker: GQL.SceneMarkerDataFragment) {
    props.onClickMarker(marker);
  }

  function renderTags() {
    function renderMarkers(primaryTag: GQL.FindSceneSceneMarkerTags) {
      const markers = primaryTag.scene_markers.map((marker) => {
        const markerTags = marker.tags.map((tag) => (
          <Badge key={tag.id} variant="secondary" className="tag-item">{tag.name}</Badge>
        ));

        return (
          <div key={marker.id}>
            <hr />
            <div>
              <a onClick={() => onClickMarker(marker)}>{marker.title}</a>
              {!isEditorOpen ? <a style={{float: "right"}} onClick={() => onOpenEditor(marker)}>Edit</a> : undefined}
            </div>
            <div>
              {TextUtils.secondsToTimestamp(marker.seconds)}
            </div>
            <div className="card-section centered">
              {markerTags}
            </div>
          </div>
        );
      });
      return markers;
    }

    const style: CSSProperties = {
      height: "300px",
      overflowY: "auto",
      overflowX: "hidden",
      display: "inline-block",
      margin: "5px",
      width: "300px",
      flex: "0 0 auto",
    };
    const tags = (props.scene as any).scene_marker_tags.map((primaryTag: GQL.FindSceneSceneMarkerTags) => {

      return (
        <div key={primaryTag.tag.id} style={{padding: "1px"}}>
          <Card style={style}>
            <div className="content" style={{whiteSpace: "normal"}}>
              <h3>{primaryTag.tag.name}</h3>
              {renderMarkers(primaryTag)}
            </div>
          </Card>
        </div>
      );
    });
    return tags;
  }

  function renderForm() {
    function onSubmit(values: IFormFields, _: FormikActions<IFormFields>) {
      const isEditing = !!editingMarker;
      const variables: GQL.SceneMarkerCreateVariables | GQL.SceneMarkerUpdateVariables = {
        title: values.title,
        seconds: parseFloat(values.seconds),
        scene_id: props.scene.id,
        primary_tag_id: values.primaryTagId,
        tag_ids: values.tagIds,
      };
      if (!isEditing) {
        sceneMarkerCreate({ variables }).then((response) => {
          setIsEditorOpen(false);
          setEditingMarker(null);
        }).catch((err) => ErrorUtils.handleApolloError(err));
      } else {
        const updateVariables = variables as GQL.SceneMarkerUpdateVariables;
        updateVariables.id = editingMarker!.id;
        sceneMarkerUpdate({ variables: updateVariables }).then((response) => {
          setIsEditorOpen(false);
          setEditingMarker(null);
        }).catch((err) => ErrorUtils.handleApolloError(err));
      }
    }
    function onDelete() {
      if (!editingMarker) { return; }
      sceneMarkerDestroy({variables: {id: editingMarker.id}}).then((response) => {
        console.log(response);
      }).catch((err) => console.error(err));
      setIsEditorOpen(false);
      setEditingMarker(null);
    }
    function renderTitleField(fieldProps: FieldProps<IFormFields>) {
      return (
        <MarkerTitleSuggest
          initialMarkerString={!!editingMarker ? editingMarker.title : undefined}
          placeholder="Title"
          name={fieldProps.field.name}
          onBlur={fieldProps.field.onBlur}
          value={fieldProps.field.value}
          onQueryChange={(query) => fieldProps.form.setFieldValue("title", query)}
        />
      );
    }
    function renderSecondsField(fieldProps: FieldProps<IFormFields>) {
      return (
        <DurationInput
          onValueChange={(s) => fieldProps.form.setFieldValue("seconds", s)}
          onReset={() => fieldProps.form.setFieldValue("seconds", Math.round(jwplayer.getPosition()))}
          numericValue={fieldProps.field.value}
        />
      );
    }
    function renderPrimaryTagField(fieldProps: FieldProps<IFormFields>) {
      return (
        <FilterSelect
          type="tags"
          onSelectItem={(tag) => fieldProps.form.setFieldValue("primaryTagId", tag ? tag.id : undefined)}
          initialId={!!editingMarker ? editingMarker.primary_tag.id : undefined}
        />
      );
    }
    function renderTagsField(fieldProps: FieldProps<IFormFields>) {
      return (
        <FilterMultiSelect
          type="tags"
          onUpdate={(tags) => fieldProps.form.setFieldValue("tagIds", tags.map((tag) => tag.id))}
          initialIds={!!editingMarker ? fieldProps.form.values.tagIds : undefined}
        />
      );
    }
    function renderFormFields(formikProps: FormikProps<IFormFields>) {
      let deleteButton: JSX.Element | undefined;
      if (!!editingMarker) {
        deleteButton = (
          <Button
            variant="danger"
            style={{float: "right", marginRight: "10px"}}
            onClick={() => onDelete()}
          >
            Delete
          </Button>
        );
      }
      return (
        <Form style={{marginTop: "10px"}}>
          <div className="columns is-multiline is-gapless">
            <BootstrapForm.Group>
              <BootstrapForm.Label htmlFor="title">Scene Marker Title</BootstrapForm.Label>
              <Field name="title" render={renderTitleField} />
            </BootstrapForm.Group>
            <BootstrapForm.Group>
              <BootstrapForm.Label htmlFor="seconds">Time</BootstrapForm.Label>
              <Field name="seconds" render={renderSecondsField} />
            </BootstrapForm.Group>
            <BootstrapForm.Group>
              <BootstrapForm.Label htmlFor="primaryTagId">Primary Tag</BootstrapForm.Label>
              <Field name="primaryTagId" render={renderPrimaryTagField} />
            </BootstrapForm.Group>
            <BootstrapForm.Group>
              <BootstrapForm.Label htmlFor="tagIds">Tags</BootstrapForm.Label>
              <Field name="tagIds" render={renderTagsField} />
            </BootstrapForm.Group>
          </div>
          <div className="buttons-container">
            <Button variant="primary" type="submit">Submit</Button>
            <Button type="button" onClick={() => setIsEditorOpen(false)}>Cancel</Button>
            {deleteButton}
          </div>
        </Form>
      );
    }
    let initialValues: any;
    if (!!editingMarker) {
      initialValues = {
        title: editingMarker.title,
        seconds: editingMarker.seconds,
        primaryTagId: editingMarker.primary_tag.id,
        tagIds: editingMarker.tags.map((tag) => tag.id),
      };
    } else {
      initialValues = {title: "", seconds: Math.round(jwplayer.getPosition()), primaryTagId: "", tagIds: []};
    }
    return (
      <Collapse in={isEditorOpen}>
        <div className="">
          <Formik
            initialValues={initialValues}
            onSubmit={onSubmit}
            render={renderFormFields}
          />
        </div>
      </Collapse>
    );
  }

  function render() {
    const newMarkerForm = (
      <div style={{margin: "5px"}}>
        <Button onClick={() => onOpenEditor()}>Create</Button>
        {renderForm()}
      </div>
    );
    if (props.scene.scene_markers.length === 0) {
      return newMarkerForm;
    }

    const containerStyle: CSSProperties = {
      overflowY: "hidden",
      overflowX: "scroll",
      whiteSpace: "nowrap",
      display: "flex",
      flexWrap: "nowrap",
      marginBottom: "20px",
    };
    return (
      <>
        {newMarkerForm}
        <div style={containerStyle}>
          {renderTags()}
        </div>
        <WallPanel
          sceneMarkers={props.scene.scene_markers}
          clickHandler={(marker) => { window.scrollTo(0, 0); onClickMarker(marker as any); }}
        />
      </>
    );
  }

  return render();
};
