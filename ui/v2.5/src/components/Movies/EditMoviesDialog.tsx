import React, { useEffect, useState } from "react";
import { Form, Col, Row } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import { useBulkMovieUpdate } from "src/core/StashService";
import * as GQL from "src/core/generated-graphql";
import { ModalComponent } from "../Shared/Modal";
import { StudioSelect } from "../Shared/Select";
import { useToast } from "src/hooks/Toast";
import * as FormUtils from "src/utils/form";
import { RatingSystem } from "../Shared/Rating/RatingSystem";
import {
  getAggregateInputIDs,
  getAggregateInputValue,
  getAggregateRating,
  getAggregateStudioId,
  getAggregateTagIds,
} from "src/utils/bulkUpdate";
import { faPencilAlt } from "@fortawesome/free-solid-svg-icons";
import { isEqual } from "lodash-es";
import { MultiSet } from "../Shared/MultiSet";

interface IListOperationProps {
  selected: GQL.MovieDataFragment[];
  onClose: (applied: boolean) => void;
}

export const EditMoviesDialog: React.FC<IListOperationProps> = (
  props: IListOperationProps
) => {
  const intl = useIntl();
  const Toast = useToast();
  const [rating100, setRating] = useState<number | undefined>();
  const [studioId, setStudioId] = useState<string | undefined>();
  const [director, setDirector] = useState<string | undefined>();

  const [tagMode, setTagMode] = React.useState<GQL.BulkUpdateIdMode>(
    GQL.BulkUpdateIdMode.Add
  );
  const [tagIds, setTagIds] = useState<string[]>();
  const [existingTagIds, setExistingTagIds] = useState<string[]>();

  const [updateMovies] = useBulkMovieUpdate(getMovieInput());

  const [isUpdating, setIsUpdating] = useState(false);

  function getMovieInput(): GQL.BulkMovieUpdateInput {
    const aggregateRating = getAggregateRating(props.selected);
    const aggregateStudioId = getAggregateStudioId(props.selected);
    const aggregateTagIds = getAggregateTagIds(props.selected);

    const movieInput: GQL.BulkMovieUpdateInput = {
      ids: props.selected.map((movie) => movie.id),
      director,
    };

    // if rating is undefined
    movieInput.rating100 = getAggregateInputValue(rating100, aggregateRating);
    movieInput.studio_id = getAggregateInputValue(studioId, aggregateStudioId);
    movieInput.tag_ids = getAggregateInputIDs(tagMode, tagIds, aggregateTagIds);

    return movieInput;
  }

  async function onSave() {
    setIsUpdating(true);
    try {
      await updateMovies();
      Toast.success(
        intl.formatMessage(
          { id: "toast.updated_entity" },
          {
            entity: intl.formatMessage({ id: "movies" }).toLocaleLowerCase(),
          }
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
    let updateStudioId: string | undefined;
    let updateTagIds: string[] = [];
    let updateDirector: string | undefined;
    let first = true;

    state.forEach((movie: GQL.MovieDataFragment) => {
      const movieTagIDs = (movie.tags ?? []).map((p) => p.id).sort();

      if (first) {
        first = false;
        updateRating = movie.rating100 ?? undefined;
        updateStudioId = movie.studio?.id ?? undefined;
        updateTagIds = movieTagIDs;
        updateDirector = movie.director ?? undefined;
      } else {
        if (movie.rating100 !== updateRating) {
          updateRating = undefined;
        }
        if (movie.studio?.id !== updateStudioId) {
          updateStudioId = undefined;
        }
        if (movie.director !== updateDirector) {
          updateDirector = undefined;
        }
        if (!isEqual(movieTagIDs, updateTagIds)) {
          updateTagIds = [];
        }
      }
    });

    setRating(updateRating);
    setStudioId(updateStudioId);
    setExistingTagIds(updateTagIds);
    setDirector(updateDirector);
  }, [props.selected]);

  function render() {
    return (
      <ModalComponent
        show
        icon={faPencilAlt}
        header={intl.formatMessage(
          { id: "actions.edit_entity" },
          { entityType: intl.formatMessage({ id: "movies" }) }
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
          <Form.Group controlId="director">
            <Form.Label>
              <FormattedMessage id="director" />
            </Form.Label>
            <Form.Control
              className="input-control"
              type="text"
              value={director}
              onChange={(event) => setDirector(event.currentTarget.value)}
              placeholder={intl.formatMessage({ id: "director" })}
            />
          </Form.Group>
          <Form.Group controlId="tags">
            <Form.Label>
              <FormattedMessage id="tags" />
            </Form.Label>
            <MultiSet
              type="tags"
              disabled={isUpdating}
              onUpdate={(itemIDs) => setTagIds(itemIDs)}
              onSetMode={(newMode) => setTagMode(newMode)}
              existingIds={existingTagIds ?? []}
              ids={tagIds ?? []}
              mode={tagMode}
            />
          </Form.Group>
        </Form>
      </ModalComponent>
    );
  }

  return render();
};
