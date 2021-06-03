import React, { useEffect, useState } from "react";
import { Button, Form, Col, Row } from "react-bootstrap";
import { FormattedMessage } from "react-intl";
import Mousetrap from "mousetrap";
import * as GQL from "src/core/generated-graphql";
import { useImageUpdate } from "src/core/StashService";
import {
  PerformerSelect,
  TagSelect,
  StudioSelect,
  LoadingIndicator,
} from "src/components/Shared";
import { useToast } from "src/hooks";
import { FormUtils } from "src/utils";
import { RatingStars } from "src/components/Scenes/SceneDetails/RatingStars";

interface IProps {
  image: GQL.ImageDataFragment;
  isVisible: boolean;
  onDelete: () => void;
}

export const ImageEditPanel: React.FC<IProps> = ({
  image,
  isVisible,
  onDelete,
}) => {
  const Toast = useToast();
  const [title, setTitle] = useState<string>(image?.title ?? "");
  const [rating, setRating] = useState<number>(image.rating ?? NaN);
  const [studioId, setStudioId] = useState<string | undefined>(
    image.studio?.id ?? undefined
  );
  const [performerIds, setPerformerIds] = useState<string[]>(
    image.performers.map((p) => p.id)
  );
  const [tagIds, setTagIds] = useState<string[]>(image.tags.map((t) => t.id));

  // Network state
  const [isLoading, setIsLoading] = useState(false);

  const [updateImage] = useImageUpdate();

  useEffect(() => {
    if (isVisible) {
      Mousetrap.bind("s s", () => {
        onSave();
      });
      Mousetrap.bind("d d", () => {
        onDelete();
      });

      // numeric keypresses get caught by jwplayer, so blur the element
      // if the rating sequence is started
      Mousetrap.bind("r", () => {
        if (document.activeElement instanceof HTMLElement) {
          document.activeElement.blur();
        }

        Mousetrap.bind("0", () => setRating(NaN));
        Mousetrap.bind("1", () => setRating(1));
        Mousetrap.bind("2", () => setRating(2));
        Mousetrap.bind("3", () => setRating(3));
        Mousetrap.bind("4", () => setRating(4));
        Mousetrap.bind("5", () => setRating(5));

        setTimeout(() => {
          Mousetrap.unbind("0");
          Mousetrap.unbind("1");
          Mousetrap.unbind("2");
          Mousetrap.unbind("3");
          Mousetrap.unbind("4");
          Mousetrap.unbind("5");
        }, 1000);
      });

      return () => {
        Mousetrap.unbind("s s");
        Mousetrap.unbind("d d");

        Mousetrap.unbind("r");
      };
    }
  });

  function getImageInput(): GQL.ImageUpdateInput {
    return {
      id: image.id,
      title,
      rating: rating ?? null,
      studio_id: studioId ?? null,
      performer_ids: performerIds,
      tag_ids: tagIds,
    };
  }

  async function onSave() {
    setIsLoading(true);
    try {
      const result = await updateImage({
        variables: {
          input: getImageInput(),
        },
      });
      if (result.data?.imageUpdate) {
        Toast.success({ content: "Updated image" });
      }
    } catch (e) {
      Toast.error(e);
    }
    setIsLoading(false);
  }

  if (isLoading) return <LoadingIndicator />;

  return (
    <div id="image-edit-details">
      <div className="form-container row px-3 pt-3">
        <div className="col edit-buttons mb-3 pl-0">
          <Button className="edit-button" variant="primary" onClick={onSave}>
            <FormattedMessage id="actions.save" />
          </Button>
          <Button
            className="edit-button"
            variant="danger"
            onClick={() => onDelete()}
          >
            <FormattedMessage id="actions.delete" />
          </Button>
        </div>
      </div>
      <div className="form-container row px-3">
        <div className="col-12 col-lg-6 col-xl-12">
          {FormUtils.renderInputGroup({
            title: "Title",
            value: title,
            onChange: setTitle,
            isEditing: true,
          })}
          <Form.Group controlId="rating" as={Row}>
            {FormUtils.renderLabel({
              title: "Rating",
            })}
            <Col xs={9}>
              <RatingStars
                value={rating}
                onSetRating={(value) => setRating(value ?? NaN)}
              />
            </Col>
          </Form.Group>

          <Form.Group controlId="studio" as={Row}>
            {FormUtils.renderLabel({
              title: "Studio",
            })}
            <Col xs={9}>
              <StudioSelect
                onSelect={(items) =>
                  setStudioId(items.length > 0 ? items[0]?.id : undefined)
                }
                ids={studioId ? [studioId] : []}
              />
            </Col>
          </Form.Group>

          <Form.Group controlId="performers" as={Row}>
            {FormUtils.renderLabel({
              title: "Performers",
              labelProps: {
                column: true,
                sm: 3,
                xl: 12,
              },
            })}
            <Col sm={9} xl={12}>
              <PerformerSelect
                isMulti
                onSelect={(items) =>
                  setPerformerIds(items.map((item) => item.id))
                }
                ids={performerIds}
              />
            </Col>
          </Form.Group>

          <Form.Group controlId="tags" as={Row}>
            {FormUtils.renderLabel({
              title: "Tags",
              labelProps: {
                column: true,
                sm: 3,
                xl: 12,
              },
            })}
            <Col sm={9} xl={12}>
              <TagSelect
                isMulti
                onSelect={(items) => setTagIds(items.map((item) => item.id))}
                ids={tagIds}
              />
            </Col>
          </Form.Group>
        </div>
      </div>
    </div>
  );
};
