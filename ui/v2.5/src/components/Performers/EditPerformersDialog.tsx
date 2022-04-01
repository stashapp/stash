import React, { useEffect, useState } from "react";
import { Form, Col, Row } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import { useBulkPerformerUpdate } from "src/core/StashService";
import * as GQL from "src/core/generated-graphql";
import { Modal } from "src/components/Shared";
import { useToast } from "src/hooks";
import { FormUtils } from "src/utils";
import MultiSet from "../Shared/MultiSet";
import { RatingStars } from "../Scenes/SceneDetails/RatingStars";
import {
  getAggregateInputIDs,
  getAggregateState,
  getAggregateTagIds,
} from "src/utils/bulkUpdate";
import {
  genderStrings,
  genderToString,
  stringToGender,
} from "src/utils/gender";
import { IndeterminateCheckbox } from "../Shared/IndeterminateCheckbox";

interface IListOperationProps {
  selected: GQL.SlimPerformerDataFragment[];
  onClose: (applied: boolean) => void;
}

export const EditPerformersDialog: React.FC<IListOperationProps> = (
  props: IListOperationProps
) => {
  const intl = useIntl();
  const Toast = useToast();
  const [tagMode, setTagMode] = React.useState<GQL.BulkUpdateIdMode>(
    GQL.BulkUpdateIdMode.Add
  );
  const [tagIds, setTagIds] = useState<string[]>();
  const [existingTagIds, setExistingTagIds] = useState<string[]>();
  const [updateInput, setUpdateInput] = useState<GQL.BulkPerformerUpdateInput>(
    {}
  );
  const genderOptions = [""].concat(genderStrings);

  const [updatePerformers] = useBulkPerformerUpdate(getPerformerInput());

  // Network state
  const [isUpdating, setIsUpdating] = useState(false);

  function setUpdateField(input: Partial<GQL.BulkPerformerUpdateInput>) {
    setUpdateInput({ ...updateInput, ...input });
  }

  function getPerformerInput(): GQL.BulkPerformerUpdateInput {
    // need to determine what we are actually setting on each performer
    const aggregateTagIds = getAggregateTagIds(props.selected);

    const performerInput: GQL.BulkPerformerUpdateInput = {
      ids: props.selected.map((performer) => {
        return performer.id;
      }),
      ...updateInput,
    };

    performerInput.tag_ids = getAggregateInputIDs(
      tagMode,
      tagIds,
      aggregateTagIds
    );

    return performerInput;
  }

  async function onSave() {
    setIsUpdating(true);
    try {
      await updatePerformers();
      Toast.success({
        content: intl.formatMessage(
          { id: "toast.updated_entity" },
          {
            entity: intl
              .formatMessage({ id: "performers" })
              .toLocaleLowerCase(),
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
    const updateState: GQL.BulkPerformerUpdateInput = {};

    const state = props.selected;
    let updateTagIds: string[] = [];
    let first = true;

    state.forEach((performer: GQL.SlimPerformerDataFragment) => {
      const performerTagIDs = (performer.tags ?? []).map((p) => p.id).sort();

      updateState.favorite = getAggregateState(
        updateState.favorite,
        performer.favorite,
        first
      );
      updateState.rating = getAggregateState(
        updateState.rating,
        performer.rating,
        first
      );
      updateState.gender = getAggregateState(
        updateState.gender,
        performer.gender,
        first
      );

      updateTagIds =
        getAggregateState(updateTagIds, performerTagIDs, first) ?? [];

      first = false;
    });

    setExistingTagIds(updateTagIds);
    setUpdateInput(updateState);
  }, [props.selected, tagMode]);

  function renderTextField(
    name: string,
    value: string | undefined | null,
    setter: (newValue: string | undefined) => void
  ) {
    return (
      <Form.Group controlId={name}>
        <Form.Label>
          <FormattedMessage id={name} />
        </Form.Label>
        <Form.Control
          className="input-control"
          type="text"
          value={value ?? ""}
          onChange={(event) => setter(event.currentTarget.value)}
          placeholder={intl.formatMessage({ id: name })}
        />
      </Form.Group>
    );
  }

  function render() {
    return (
      <Modal
        show
        icon="pencil-alt"
        header={intl.formatMessage(
          { id: "actions.edit_entity" },
          { entityType: intl.formatMessage({ id: "performers" }) }
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
        <Form.Group controlId="rating" as={Row}>
          {FormUtils.renderLabel({
            title: intl.formatMessage({ id: "rating" }),
          })}
          <Col xs={9}>
            <RatingStars
              value={updateInput.rating ?? undefined}
              onSetRating={(value) => setUpdateField({ rating: value })}
              disabled={isUpdating}
            />
          </Col>
        </Form.Group>
        <Form>
          <Form.Group controlId="favorite">
            <IndeterminateCheckbox
              setChecked={(checked) => setUpdateField({ favorite: checked })}
              checked={updateInput.favorite ?? undefined}
            />
          </Form.Group>

          <Form.Group>
            <Form.Label>
              <FormattedMessage id="gender" />
            </Form.Label>
            <Form.Control
              as="select"
              className="input-control"
              value={genderToString(updateInput.gender ?? undefined)}
              onChange={(event) =>
                setUpdateField({
                  gender: stringToGender(event.currentTarget.value),
                })
              }
            >
              {genderOptions.map((opt) => (
                <option value={opt} key={opt}>
                  {opt}
                </option>
              ))}
            </Form.Control>
          </Form.Group>

          {renderTextField("country", updateInput.country, (v) =>
            setUpdateField({ country: v })
          )}
          {renderTextField("ethnicity", updateInput.ethnicity, (v) =>
            setUpdateField({ ethnicity: v })
          )}
          {renderTextField("hair_color", updateInput.hair_color, (v) =>
            setUpdateField({ hair_color: v })
          )}
          {renderTextField("eye_color", updateInput.eye_color, (v) =>
            setUpdateField({ eye_color: v })
          )}
          {renderTextField("fake_tits", updateInput.fake_tits, (v) =>
            setUpdateField({ fake_tits: v })
          )}
          {renderTextField("tattoos", updateInput.tattoos, (v) =>
            setUpdateField({ tattoos: v })
          )}
          {renderTextField("piercings", updateInput.piercings, (v) =>
            setUpdateField({ piercings: v })
          )}
          {renderTextField("career_length", updateInput.career_length, (v) =>
            setUpdateField({ career_length: v })
          )}

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
      </Modal>
    );
  }

  return render();
};
