import React, { useEffect, useState } from "react";
import { Form, Col, Row } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import { useBulkMovieUpdate } from "src/core/StashService";
import * as GQL from "src/core/generated-graphql";
import { Modal, StudioSelect } from "src/components/Shared";
import { useToast } from "src/hooks";
import { FormUtils } from "src/utils";
import { RatingStars } from "../Scenes/SceneDetails/RatingStars";

interface IListOperationProps {
  selected: GQL.MovieDataFragment[];
  onClose: (applied: boolean) => void;
}

export const EditMoviesDialog: React.FC<IListOperationProps> = (
  props: IListOperationProps
) => {
  const intl = useIntl();
  const Toast = useToast();
  const [rating, setRating] = useState<number | undefined>();
  const [studioId, setStudioId] = useState<string | undefined>();
  const [director, setDirector] = useState<string | undefined>();

  const [updateMovies] = useBulkMovieUpdate(getMovieInput());

  const [isUpdating, setIsUpdating] = useState(false);

  function getMovieInput(): GQL.BulkMovieUpdateInput {
    const aggregateRating = getRating(props.selected);

    const movieInput: GQL.BulkMovieUpdateInput = {
      ids: props.selected.map((movie) => movie.id),
      studio_id: studioId,
      director,
    };

    // if rating is undefined
    if (rating === undefined) {
      // and all galleries have the same rating, then we are unsetting the rating.
      if (aggregateRating) {
        // null to unset rating
        movieInput.rating = null;
      }
      // otherwise not setting the rating
    } else {
      // if rating is set, then we are setting the rating for all
      movieInput.rating = rating;
    }

    return movieInput;
  }

  function getRating(state: GQL.MovieDataFragment[]) {
    let ret: number | undefined;
    let first = true;

    state.forEach((movie) => {
      if (first) {
        ret = movie.rating ?? undefined;
        first = false;
      } else if (ret !== movie.rating) {
        ret = undefined;
      }
    });

    return ret;
  }

  async function onSave() {
    setIsUpdating(true);
    try {
      await updateMovies();
      Toast.success({
        content: intl.formatMessage(
          { id: "toast.updated_entity" },
          {
            entity: intl.formatMessage({ id: "movies" }).toLocaleLowerCase(),
          }
        ),
      });
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
    let updateDirector: string | undefined;
    let first = true;

    state.forEach((movie: GQL.MovieDataFragment) => {
      if (first) {
        first = false;
        updateRating = movie.rating ?? undefined;
        updateStudioId = movie.studio?.id ?? undefined;
        updateDirector = movie.director ?? undefined;
      } else {
        if (movie.rating !== updateRating) {
          updateRating = undefined;
        }
        if (movie.studio?.id !== updateStudioId) {
          updateStudioId = undefined;
        }
        if (movie.director !== updateDirector) {
          updateDirector = undefined;
        }
      }
    });

    setRating(updateRating);
    setStudioId(updateStudioId);
    setDirector(updateDirector);
  }, [props.selected]);

  function render() {
    return (
      <Modal
        show
        icon="pencil-alt"
        header="Edit Movies"
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
              <RatingStars
                value={rating}
                onSetRating={(value) => setRating(value)}
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
        </Form>
      </Modal>
    );
  }

  return render();
};
