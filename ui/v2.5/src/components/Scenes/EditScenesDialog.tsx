import React, { useEffect, useState } from "react";
import { Form, Col, Row } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import isEqual from "lodash-es/isEqual";
import { useBulkSceneUpdate } from "src/core/StashService";
import * as GQL from "src/core/generated-graphql";
import { StudioSelect } from "../Shared/Select";
import { ModalComponent } from "../Shared/Modal";
import { MultiSelect, MultiString } from "../Shared/MultiSet";
import { useToast } from "src/hooks/Toast";
import * as FormUtils from "src/utils/form";
import { RatingSystem } from "../Shared/Rating/RatingSystem";
import {
  getAggregateInputIDs,
  getAggregateInputStrings,
  getAggregateInputValue,
  getAggregateStateObject,
  getAggregateMovieIds,
  getAggregatePerformerIds,
  getAggregateRating,
  getAggregateStudioId,
  getAggregateTagIds,
  getAggregateGalleryIds,
  getAggregateUrls,
} from "src/utils/bulkUpdate";
import { BulkUpdateTextInput } from "../Shared/BulkUpdateTextInput";
import { faPencilAlt } from "@fortawesome/free-solid-svg-icons";

interface IListOperationProps {
  selected: GQL.SlimSceneDataFragment[];
  onClose: (applied: boolean) => void;
  showAllFields?: boolean;
}

const sceneFields = [
  "title",
  "scene_code",
  "details",
  "director",
  "date",
  "urls",
];

export const EditScenesDialog: React.FC<IListOperationProps> = (
  props: IListOperationProps
) => {
  const intl = useIntl();
  const Toast = useToast();
  const [rating100, setRating] = useState<number>();
  const [studioId, setStudioId] = useState<string>();
  const [urlsMode, setUrlsMode] =
    React.useState<GQL.BulkUpdateIdMode>(GQL.BulkUpdateIdMode.Add);
  const [urls, setUrls] = useState<string[]>();
  const [existingUrls, setExistingUrls] = useState<string[]>();
  const selectedUrls = props.selected.map((scene) => ({
    urls: scene.urls.map((url) => ({ value: url }))
  }));
  const [performerMode, setPerformerMode] =
    React.useState<GQL.BulkUpdateIdMode>(GQL.BulkUpdateIdMode.Add);
  const [performerIds, setPerformerIds] = useState<string[]>();
  const [existingPerformerIds, setExistingPerformerIds] = useState<string[]>();
  const [tagMode, setTagMode] = React.useState<GQL.BulkUpdateIdMode>(
    GQL.BulkUpdateIdMode.Add
  );
  const [tagIds, setTagIds] = useState<string[]>();
  const [existingTagIds, setExistingTagIds] = useState<string[]>();
  const [movieMode, setMovieMode] = React.useState<GQL.BulkUpdateIdMode>(
    GQL.BulkUpdateIdMode.Add
  );
  const [movieIds, setMovieIds] = useState<string[]>();
  const [existingMovieIds, setExistingMovieIds] = useState<string[]>();
  const [galleryMode, setGalleryMode] = React.useState<GQL.BulkUpdateIdMode>(
    GQL.BulkUpdateIdMode.Add
  );
  const [galleryIds, setGalleryIds] = useState<string[]>();
  const [existingGalleryIds, setExistingGalleryIds] = useState<string[]>();
  const [organized, setOrganized] = useState<boolean | undefined>();
  const [updateInput, setUpdateInput] = useState<GQL.BulkSceneUpdateInput>(
    {}
  );

  const [showAllFields, setShowAllFields] = useState(props.showAllFields ?? false);

  const [updateScenes] = useBulkSceneUpdate(getSceneInput());

  // Network state
  const [isUpdating, setIsUpdating] = useState(false);

  const checkboxRef = React.createRef<HTMLInputElement>();

  function setUpdateField(input: Partial<GQL.BulkSceneUpdateInput>) {
    setUpdateInput({ ...updateInput, ...input });
  }

  function getSceneInput(): GQL.BulkSceneUpdateInput {
    // need to determine what we are actually setting on each scene
    const aggregateRating = getAggregateRating(props.selected);
    const aggregateStudioId = getAggregateStudioId(props.selected);
    const aggregatePerformerIds = getAggregatePerformerIds(props.selected);
    const aggregateTagIds = getAggregateTagIds(props.selected);
    const aggregateMovieIds = getAggregateMovieIds(props.selected);
    const aggregateGalleryIds = getAggregateGalleryIds(props.selected);
    const aggregateUrls = getAggregateUrls(selectedUrls);

    const sceneInput: GQL.BulkSceneUpdateInput = {
      ids: props.selected.map((scene) => {
        return scene.id;
      }),
      ...updateInput,
    };

    sceneInput.rating100 = getAggregateInputValue(rating100, aggregateRating);
    sceneInput.studio_id = getAggregateInputValue(studioId, aggregateStudioId);

    sceneInput.urls = getAggregateInputStrings(
      urlsMode,
      urls,
      aggregateUrls
    );

    sceneInput.performer_ids = getAggregateInputIDs(
      performerMode,
      performerIds,
      aggregatePerformerIds
    );
    sceneInput.tag_ids = getAggregateInputIDs(tagMode, tagIds, aggregateTagIds);
    sceneInput.movie_ids = getAggregateInputIDs(
      movieMode,
      movieIds,
      aggregateMovieIds
    );
    sceneInput.gallery_ids = getAggregateInputIDs(
      galleryMode, 
      galleryIds, 
      aggregateGalleryIds
    );

    if (organized !== undefined) {
      sceneInput.organized = organized;
    }

    return sceneInput;
  }

  async function onSave() {
    setIsUpdating(true);
    try {
      await updateScenes();
      Toast.success(
        intl.formatMessage(
          { id: "toast.updated_entity" },
          { entity: intl.formatMessage({ id: "scenes" }).toLocaleLowerCase() }
        )
      );
      props.onClose(true);
    } catch (e) {
      Toast.error(e);
    }
    setIsUpdating(false);
  }

  useEffect(() => {
    const state = props.selected;
    let updateRating: number | undefined;
    let updateStudioID: string | undefined;
    let updatePerformerIds: string[] = [];
    let updateTagIds: string[] = [];
    let updateMovieIds: string[] = [];
    let updateGalleryIds: string[] = [];
    let updateUrls: string[] = [];
    let updateOrganized: boolean | undefined;
    let first = true;

    state.forEach((scene: GQL.SlimSceneDataFragment) => {
      getAggregateStateObject(state, scene, sceneFields, first);
      const sceneRating = scene.rating100;
      const sceneStudioID = scene?.studio?.id;
      const scenePerformerIDs = (scene.performers ?? [])
        .map((p) => p.id)
        .sort();
      const sceneTagIDs = (scene.tags ?? []).map((p) => p.id).sort();
      const sceneMovieIDs = (scene.movies ?? []).map((m) => m.movie.id).sort();
      const sceneGalleryIDs = (scene.galleries ?? []).map((g) => g.id).sort();
      const sceneUrls = (scene.urls ?? []);

      if (first) {
        updateRating = sceneRating ?? undefined;
        updateStudioID = sceneStudioID;
        updatePerformerIds = scenePerformerIDs;
        updateTagIds = sceneTagIDs;
        updateMovieIds = sceneMovieIDs;
        updateGalleryIds = sceneGalleryIDs;
        updateUrls = sceneUrls;
        first = false;
        updateOrganized = scene.organized;
      } else {
        if (sceneRating !== updateRating) {
          updateRating = undefined;
        }
        if (sceneStudioID !== updateStudioID) {
          updateStudioID = undefined;
        }
        if (!isEqual(scenePerformerIDs, updatePerformerIds)) {
          updatePerformerIds = [];
        }
        if (!isEqual(sceneTagIDs, updateTagIds)) {
          updateTagIds = [];
        }
        if (!isEqual(sceneMovieIDs, updateMovieIds)) {
          updateMovieIds = [];
        }
        if (!isEqual(sceneGalleryIDs, updateGalleryIds)) {
          updateGalleryIds = [];
        }
        if (!isEqual(sceneUrls, updateUrls)) {
          updateUrls = [];
        }
        if (scene.organized !== updateOrganized) {
          updateOrganized = undefined;
        }
      }
    });

    setRating(updateRating);
    setStudioId(updateStudioID);
    setExistingPerformerIds(updatePerformerIds);
    setExistingTagIds(updateTagIds);
    setExistingMovieIds(updateMovieIds);
    setExistingGalleryIds(updateGalleryIds);
    setExistingUrls(updateUrls)
    setOrganized(updateOrganized);
  }, [props.selected, performerMode, tagMode, movieMode, urlsMode]);

  useEffect(() => {
    if (checkboxRef.current) {
      checkboxRef.current.indeterminate = organized === undefined;
    }
  }, [organized, checkboxRef]);

  function renderURLMultiSelect(
    urls: string[] | undefined
  ) {
    return (
      <MultiString
        disabled={isUpdating}
        onUpdate={(itemIDs) => {setUrls(itemIDs)}}
        onSetMode={(newMode) => {setUrlsMode(newMode)}}
        strings={urls ?? []}
        existing={existingUrls ?? []}
        mode={urlsMode}
      />
    );
  }

  function renderMultiSelect(
    type: "performers" | "tags" | "movies" | "galleries",
    ids: string[] | undefined
  ) {
    let mode = GQL.BulkUpdateIdMode.Add;
    let existingIds: string[] | undefined = [];
    switch (type) {
      case "performers":
        mode = performerMode;
        existingIds = existingPerformerIds;
        break;
      case "tags":
        mode = tagMode;
        existingIds = existingTagIds;
        break;
      case "movies":
        mode = movieMode;
        existingIds = existingMovieIds;
        break;
      case "galleries":
        mode = galleryMode;
        existingIds = existingGalleryIds;
        break;
    }

    return (
      <MultiSelect
        type={type}
        disabled={isUpdating}
        onUpdate={(itemIDs) => {
          switch (type) {
            case "performers":
              setPerformerIds(itemIDs);
              break;
            case "tags":
              setTagIds(itemIDs);
              break;
            case "movies":
              setMovieIds(itemIDs);
              break;
            case "galleries":
              setGalleryIds(itemIDs);
              break;
          }
        }}
        onSetMode={(newMode) => {
          switch (type) {
            case "performers":
              setPerformerMode(newMode);
              break;
            case "tags":
              setTagMode(newMode);
              break;
            case "movies":
              setMovieMode(newMode);
              break;
            case "galleries":
              setGalleryMode(newMode);
              break;
          }
        }}
        ids={ids ?? []}
        existing={existingIds ?? []}
        mode={mode}
      />
    );
  }

  function cycleOrganized() {
    if (organized) {
      setOrganized(undefined);
    } else if (organized === undefined) {
      setOrganized(false);
    } else {
      setOrganized(true);
    }
  }

  function renderTextField(
    value: string | undefined | null,
    setter: (newValue: string | undefined) => void,
    isDetails: Boolean = false
  ) {
    return (
      <Form.Group>
        <BulkUpdateTextInput
          as={isDetails ? 'textarea' : undefined}
          value={value === null ? "" : value ?? undefined}
          valueChanged={(newValue) => setter(newValue)}
          unsetDisabled={props.selected.length < 2}
        />
      </Form.Group>
    );
  }

  function render() {
    return (
      <ModalComponent
        show
        icon={faPencilAlt}
        header={intl.formatMessage(
          { id: "dialogs.edit_entity_title" },
          {
            count: props?.selected?.length ?? 1,
            singularEntity: intl.formatMessage({ id: "scene" }),
            pluralEntity: intl.formatMessage({ id: "scenes" }),
          }
        )}
        accept={{
          onClick: onSave,
          text: intl.formatMessage({ id: "actions.apply" }),
        }}
        cancel={{
          onClick: () => props.onClose(false),
          text: intl.formatMessage({ id: "actions.cancel" }),
          variant: "secondary",
        }}
        leftFooterButtons={
          <Form.Group controlId="toggle-all">
            <Form.Switch
              label={intl.formatMessage({ id: "actions.all_fields" })}
              checked={showAllFields}
              onChange={() => setShowAllFields(!showAllFields)}
            />
          </Form.Group>
        }
        isRunning={isUpdating}
      >
        <Form>
          <Form.Group controlId="rating" as={Row}>
            {FormUtils.renderLabel({
              title: intl.formatMessage({ id: "rating" }),
            })}
            <Col xs={9}>
              <RatingSystem
                value={rating100}
                onSetRating={(value) => setRating(value ?? undefined)}
                disabled={isUpdating}
              />
            </Col>
          </Form.Group>
          
          {showAllFields && 
          <Form.Group controlId="text-input" as={Row}>
            {FormUtils.renderLabel({title: intl.formatMessage({ id: "title" })})}
            <Col xs={9}>
              {renderTextField(updateInput.title, (v) => setUpdateField({ title: v }))}
            </Col>
            {FormUtils.renderLabel({title: intl.formatMessage({ id: "scene_code" })})}
            <Col xs={9}>
              {renderTextField(updateInput.code, (v) => setUpdateField({ code: v }))}
            </Col>
            {FormUtils.renderLabel({title: intl.formatMessage({ id: "urls" })})}
            <Col xs={9}>
              {renderURLMultiSelect(urls)}
            </Col>
            {FormUtils.renderLabel({title: intl.formatMessage({ id: "date" })})}
            <Col xs={9}>
              {renderTextField(updateInput.date, (v) => setUpdateField({ date: v }))}
            </Col>
            {FormUtils.renderLabel({title: intl.formatMessage({ id: "director" })})}
            <Col xs={9}>
              {renderTextField(updateInput.director, (v) => setUpdateField({ director: v }))}
            </Col>
          </Form.Group>
          }

          <Form.Group controlId="studio" as={Row}>
            {FormUtils.renderLabel({
              title: intl.formatMessage({ id: "studio" }),
            })}
            <Col xs={9}>
              <StudioSelect
                onSelect={(items) =>
                  setStudioId(items.length > 0 ? items[0]?.id : undefined)
                }
                ids={studioId ? [studioId] : []}
                isDisabled={isUpdating}
              />
            </Col>
          </Form.Group>

          <Form.Group controlId="performers">
            <Form.Label>
              <FormattedMessage id="performers" />
            </Form.Label>
            {renderMultiSelect("performers", performerIds)}
          </Form.Group>

          <Form.Group controlId="movies">
            <Form.Label>
              <FormattedMessage id="movies" />
            </Form.Label>
            {renderMultiSelect("movies", movieIds)}
          </Form.Group>

          {showAllFields && <Form.Group controlId="galleries">
            <Form.Label>
              <FormattedMessage id="galleries" />
            </Form.Label>
            {renderMultiSelect("galleries", galleryIds)}
          </Form.Group>}

          <Form.Group controlId="tags">
            <Form.Label>
              <FormattedMessage id="tags" />
            </Form.Label>
            {renderMultiSelect("tags", tagIds)}
          </Form.Group>

          {showAllFields && 
          <Form.Group controlId="details">
            <Form.Label>
              <FormattedMessage id="details" />
            </Form.Label>
            {renderTextField(updateInput.details, (v) => setUpdateField({ details: v }), true)}
          </Form.Group>
          }

          <Form.Group controlId="organized">
            <Form.Check
              type="checkbox"
              label={intl.formatMessage({ id: "organized" })}
              checked={organized}
              ref={checkboxRef}
              onChange={() => cycleOrganized()}
            />
          </Form.Group>
        </Form>
      </ModalComponent>
    );
  }

  return render();
};