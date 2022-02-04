import React, { useEffect, useState } from "react";
import { Form, Col, Row } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import _ from "lodash";
import { useBulkMovieUpdate } from "src/core/StashService";
import * as GQL from "src/core/generated-graphql";
import { Modal, StudioSelect } from "src/components/Shared";
import { useToast } from "src/hooks";
import { FormUtils } from "src/utils";
import MultiSet from "../Shared/MultiSet";
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
    const movieInput: GQL.BulkMovieUpdateInput = {
      ids: props.selected.map((movie) => movie.id),
      rating,
      studio_id: studioId,
      director,
    };

    return movieInput;
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
    setRating(undefined);
    setStudioId(undefined);
    setDirector(undefined);
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
