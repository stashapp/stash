import React, { CSSProperties, useState } from "react";
import {
  Badge,
  Button,
  Card,
  Collapse,
  Form as BootstrapForm
} from "react-bootstrap";
import { Field, FieldProps, Form, Formik } from "formik";
import * as GQL from "src/core/generated-graphql";
import { StashService } from "src/core/StashService";
import { TextUtils } from "src/utils";
import { useToast } from "src/hooks";
import {
  DurationInput,
  TagSelect,
  MarkerTitleSuggest
} from "src/components/Shared";
import { WallPanel } from "src/components/Wall/WallPanel";
import { SceneHelpers } from "../helpers";

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

export const SceneMarkersPanel: React.FC<ISceneMarkersPanelProps> = (
  props: ISceneMarkersPanelProps
) => {
  const Toast = useToast();
  const [isEditorOpen, setIsEditorOpen] = useState<boolean>(false);
  const [
    editingMarker,
    setEditingMarker
  ] = useState<GQL.SceneMarkerDataFragment | null>(null);

  const [sceneMarkerCreate] = StashService.useSceneMarkerCreate();
  const [sceneMarkerUpdate] = StashService.useSceneMarkerUpdate();
  const [sceneMarkerDestroy] = StashService.useSceneMarkerDestroy();

  const jwplayer = SceneHelpers.getPlayer();

  function onOpenEditor(marker?: GQL.SceneMarkerDataFragment) {
    setIsEditorOpen(true);
    setEditingMarker(marker ?? null);
  }

  function onClickMarker(marker: GQL.SceneMarkerDataFragment) {
    props.onClickMarker(marker);
  }

  function renderTags() {
    function renderMarkers(primaryTag: GQL.SceneMarkerTag) {
      const markers = primaryTag.scene_markers.map(marker => {
        const markerTags = marker.tags.map(tag => (
          <Badge key={tag.id} variant="secondary" className="tag-item">
            {tag.name}
          </Badge>
        ));

        return (
          <div key={marker.id}>
            <hr />
            <div>
              <Button variant="link" onClick={() => onClickMarker(marker)}>
                {marker.title}
              </Button>
              {!isEditorOpen ? (
                <Button
                  variant="link"
                  style={{ float: "right" }}
                  onClick={() => onOpenEditor(marker)}
                >
                  Edit
                </Button>
              ) : (
                ""
              )}
            </div>
            <div>{TextUtils.secondsToTimestamp(marker.seconds)}</div>
            <div className="card-section centered">{markerTags}</div>
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
      flex: "0 0 auto"
    };
    const tags = (props.scene as any).scene_marker_tags.map(
      (primaryTag: GQL.SceneMarkerTag) => {
        return (
          <div key={primaryTag.tag.id} style={{ padding: "1px" }}>
            <Card style={style}>
              <div className="content" style={{ whiteSpace: "normal" }}>
                <h3>{primaryTag.tag.name}</h3>
                {renderMarkers(primaryTag)}
              </div>
            </Card>
          </div>
        );
      }
    );
    return tags;
  }

  function renderForm() {
    function onSubmit(values: IFormFields) {
      const isEditing = !!editingMarker;
      const variables:
        | GQL.SceneMarkerUpdateInput
        | GQL.SceneMarkerCreateInput = {
        title: values.title,
        seconds: parseFloat(values.seconds),
        scene_id: props.scene.id,
        primary_tag_id: values.primaryTagId,
        tag_ids: values.tagIds
      };
      if (!isEditing) {
        sceneMarkerCreate({ variables })
          .then(() => {
            setIsEditorOpen(false);
            setEditingMarker(null);
          })
          .catch(err => Toast.error(err));
      } else {
        const updateVariables = variables as GQL.SceneMarkerUpdateInput;
        updateVariables.id = editingMarker!.id;
        sceneMarkerUpdate({ variables: updateVariables })
          .then(() => {
            setIsEditorOpen(false);
            setEditingMarker(null);
          })
          .catch(err => Toast.error(err));
      }
    }
    function onDelete() {
      if (!editingMarker) {
        return;
      }
      sceneMarkerDestroy({ variables: { id: editingMarker.id } })
        // eslint-disable-next-line no-console
        .catch(err => console.error(err));
      setIsEditorOpen(false);
      setEditingMarker(null);
    }
    function renderTitleField(fieldProps: FieldProps<IFormFields>) {
      return (
        <MarkerTitleSuggest
          initialMarkerTitle={editingMarker?.title}
          onChange={(query: string) =>
            fieldProps.form.setFieldValue("title", query)
          }
        />
      );
    }
    function renderSecondsField(fieldProps: FieldProps<IFormFields>) {
      return (
        <DurationInput
          onValueChange={s => fieldProps.form.setFieldValue("seconds", s)}
          onReset={() =>
            fieldProps.form.setFieldValue(
              "seconds",
              Math.round(jwplayer.getPosition())
            )
          }
          numericValue={Number.parseInt(fieldProps.field.value.seconds, 10)}
        />
      );
    }
    function renderPrimaryTagField(fieldProps: FieldProps<IFormFields>) {
      return (
        <TagSelect
          onSelect={tags =>
            fieldProps.form.setFieldValue("primaryTagId", tags[0]?.id)
          }
          initialIds={editingMarker ? [editingMarker.primary_tag.id] : []}
        />
      );
    }
    function renderTagsField(fieldProps: FieldProps<IFormFields>) {
      return (
        <TagSelect
          isMulti
          onSelect={tags =>
            fieldProps.form.setFieldValue(
              "tagIds",
              tags.map(tag => tag.id)
            )
          }
          initialIds={editingMarker ? fieldProps.form.values.tagIds : []}
        />
      );
    }
    function renderFormFields() {
      let deleteButton: JSX.Element | undefined;
      if (editingMarker) {
        deleteButton = (
          <Button
            variant="danger"
            style={{ float: "right", marginRight: "10px" }}
            onClick={() => onDelete()}
          >
            Delete
          </Button>
        );
      }
      return (
        <Form style={{ marginTop: "10px" }}>
          <div className="columns is-multiline is-gapless">
            <BootstrapForm.Group>
              <BootstrapForm.Label htmlFor="title">
                Scene Marker Title
              </BootstrapForm.Label>
              <Field name="title" render={renderTitleField} />
            </BootstrapForm.Group>
            <BootstrapForm.Group>
              <BootstrapForm.Label htmlFor="seconds">Time</BootstrapForm.Label>
              <Field name="seconds" render={renderSecondsField} />
            </BootstrapForm.Group>
            <BootstrapForm.Group>
              <BootstrapForm.Label htmlFor="primaryTagId">
                Primary Tag
              </BootstrapForm.Label>
              <Field name="primaryTagId" render={renderPrimaryTagField} />
            </BootstrapForm.Group>
            <BootstrapForm.Group>
              <BootstrapForm.Label htmlFor="tagIds">Tags</BootstrapForm.Label>
              <Field name="tagIds" render={renderTagsField} />
            </BootstrapForm.Group>
          </div>
          <div className="buttons-container">
            <Button variant="primary" type="submit">
              Submit
            </Button>
            <Button type="button" onClick={() => setIsEditorOpen(false)}>
              Cancel
            </Button>
            {deleteButton}
          </div>
        </Form>
      );
    }
    let initialValues: any;
    if (editingMarker) {
      initialValues = {
        title: editingMarker.title,
        seconds: editingMarker.seconds,
        primaryTagId: editingMarker.primary_tag.id,
        tagIds: editingMarker.tags.map(tag => tag.id)
      };
    } else {
      initialValues = {
        title: "",
        seconds: Math.round(jwplayer.getPosition()),
        primaryTagId: "",
        tagIds: []
      };
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
      <div style={{ margin: "5px" }}>
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
      marginBottom: "20px"
    };
    return (
      <>
        {newMarkerForm}
        <div style={containerStyle}>{renderTags()}</div>
        <WallPanel
          sceneMarkers={props.scene.scene_markers}
          clickHandler={marker => {
            window.scrollTo(0, 0);
            onClickMarker(marker as any);
          }}
        />
      </>
    );
  }

  return render();
};
