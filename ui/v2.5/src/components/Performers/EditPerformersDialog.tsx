import React, { useEffect, useState } from "react";
import { Col, Form, Row } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import { useBulkPerformerUpdate } from "src/core/StashService";
import * as GQL from "src/core/generated-graphql";
import { ModalComponent } from "../Shared/Modal";
import { useToast } from "src/hooks/Toast";
import { MultiSet } from "../Shared/MultiSet";
import { RatingSystem } from "../Shared/Rating/RatingSystem";
import {
  getAggregateInputValue,
  getAggregateState,
  getAggregateStateObject,
} from "src/utils/bulkUpdate";
import {
  genderStrings,
  genderToString,
  stringToGender,
} from "src/utils/gender";
import {
  circumcisedStrings,
  circumcisedToString,
  stringToCircumcised,
} from "src/utils/circumcised";
import { IndeterminateCheckbox } from "../Shared/IndeterminateCheckbox";
import { BulkUpdateTextInput } from "../Shared/BulkUpdateTextInput";
import { faPencilAlt } from "@fortawesome/free-solid-svg-icons";
import * as FormUtils from "src/utils/form";
import { CountrySelect } from "../Shared/CountrySelect";
import { useConfigurationContext } from "src/hooks/Config";
import cx from "classnames";

interface IListOperationProps {
  selected: GQL.SlimPerformerDataFragment[];
  onClose: (applied: boolean) => void;
}

const performerFields = [
  "favorite",
  "disambiguation",
  "rating100",
  "gender",
  "birthdate",
  "death_date",
  "career_length",
  "country",
  "ethnicity",
  "eye_color",
  // "height",
  // "weight",
  "measurements",
  "fake_tits",
  "penis_length",
  "circumcised",
  "hair_color",
  "tattoos",
  "piercings",
  "ignore_auto_tag",
];

export const EditPerformersDialog: React.FC<IListOperationProps> = (
  props: IListOperationProps
) => {
  const intl = useIntl();
  const Toast = useToast();

  const { configuration } = useConfigurationContext();
  const { sfwContentMode } = configuration.interface;

  const [tagIds, setTagIds] = useState<GQL.BulkUpdateIds>({
    mode: GQL.BulkUpdateIdMode.Add,
  });
  const [existingTagIds, setExistingTagIds] = useState<string[]>();
  const [aggregateState, setAggregateState] =
    useState<GQL.BulkPerformerUpdateInput>({});
  // height and weight needs conversion to/from number
  const [height, setHeight] = useState<string | undefined>();
  const [weight, setWeight] = useState<string | undefined>();
  const [penis_length, setPenisLength] = useState<string | undefined>();
  const [updateInput, setUpdateInput] = useState<GQL.BulkPerformerUpdateInput>(
    {}
  );
  const genderOptions = [""].concat(genderStrings);
  const circumcisedOptions = [""].concat(circumcisedStrings);

  const [updatePerformers] = useBulkPerformerUpdate(getPerformerInput());

  // Network state
  const [isUpdating, setIsUpdating] = useState(false);

  function setUpdateField(input: Partial<GQL.BulkPerformerUpdateInput>) {
    setUpdateInput({ ...updateInput, ...input });
  }

  function getPerformerInput(): GQL.BulkPerformerUpdateInput {
    const performerInput: GQL.BulkPerformerUpdateInput = {
      ids: props.selected.map((performer) => {
        return performer.id;
      }),
      ...updateInput,
      tag_ids: tagIds,
    };

    // we don't have unset functionality for the rating star control
    // so need to determine if we are setting a rating or not
    performerInput.rating100 = getAggregateInputValue(
      updateInput.rating100,
      aggregateState.rating100
    );

    // gender dropdown doesn't have unset functionality
    // so need to determine what we are setting
    performerInput.gender = getAggregateInputValue(
      updateInput.gender,
      aggregateState.gender
    );
    performerInput.circumcised = getAggregateInputValue(
      updateInput.circumcised,
      aggregateState.circumcised
    );

    if (height !== undefined) {
      performerInput.height_cm = parseFloat(height);
    }
    if (weight !== undefined) {
      performerInput.weight = parseFloat(weight);
    }

    if (penis_length !== undefined) {
      performerInput.penis_length = parseFloat(penis_length);
    }

    return performerInput;
  }

  async function onSave() {
    setIsUpdating(true);
    try {
      await updatePerformers();
      Toast.success(
        intl.formatMessage(
          { id: "toast.updated_entity" },
          {
            entity: intl
              .formatMessage({ id: "performers" })
              .toLocaleLowerCase(),
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
    const updateState: GQL.BulkPerformerUpdateInput = {};

    const state = props.selected;
    let updateTagIds: string[] = [];
    let updateHeight: string | undefined | null = undefined;
    let updateWeight: string | undefined | null = undefined;
    let updatePenisLength: string | undefined | null = undefined;
    let first = true;

    state.forEach((performer: GQL.SlimPerformerDataFragment) => {
      getAggregateStateObject(updateState, performer, performerFields, first);

      const performerTagIDs = (performer.tags ?? []).map((p) => p.id).sort();

      updateTagIds =
        getAggregateState(updateTagIds, performerTagIDs, first) ?? [];

      const thisHeight =
        performer.height_cm !== undefined && performer.height_cm !== null
          ? performer.height_cm.toString()
          : performer.height_cm;
      updateHeight = getAggregateState(updateHeight, thisHeight, first);

      const thisWeight =
        performer.weight !== undefined && performer.weight !== null
          ? performer.weight.toString()
          : performer.weight;
      updateWeight = getAggregateState(updateWeight, thisWeight, first);

      const thisPenisLength =
        performer.penis_length !== undefined && performer.penis_length !== null
          ? performer.penis_length.toString()
          : performer.penis_length;
      updatePenisLength = getAggregateState(
        updatePenisLength,
        thisPenisLength,
        first
      );

      first = false;
    });

    setExistingTagIds(updateTagIds);
    setHeight(updateHeight);
    setWeight(updateWeight);
    setAggregateState(updateState);
    setUpdateInput(updateState);
  }, [props.selected]);

  function renderTextField(
    name: string,
    value: string | undefined | null,
    setter: (newValue: string | undefined) => void
  ) {
    return (
      <Form.Group controlId={name} data-field={name}>
        <Form.Label>
          <FormattedMessage id={name} />
        </Form.Label>
        <BulkUpdateTextInput
          value={value === null ? "" : value ?? undefined}
          valueChanged={(newValue) => setter(newValue)}
          unsetDisabled={props.selected.length < 2}
        />
      </Form.Group>
    );
  }

  function render() {
    // sfw class needs to be set because it is outside body

    return (
      <ModalComponent
        dialogClassName={cx("edit-performers-dialog", {
          "sfw-content-mode": sfwContentMode,
        })}
        show
        icon={faPencilAlt}
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
        <Form.Group controlId="rating" as={Row} data-field={name}>
          {FormUtils.renderLabel({
            title: intl.formatMessage({ id: "rating" }),
          })}
          <Col xs={9}>
            <RatingSystem
              value={updateInput.rating100}
              onSetRating={(value) =>
                setUpdateField({ rating100: value ?? undefined })
              }
              disabled={isUpdating}
            />
          </Col>
        </Form.Group>
        <Form>
          <Form.Group controlId="favorite">
            <IndeterminateCheckbox
              setChecked={(checked) => setUpdateField({ favorite: checked })}
              checked={updateInput.favorite ?? undefined}
              label={intl.formatMessage({ id: "favourite" })}
            />
          </Form.Group>

          <Form.Group>
            <Form.Label>
              <FormattedMessage id="gender" />
            </Form.Label>
            <Form.Control
              as="select"
              className="input-control"
              value={genderToString(updateInput.gender)}
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

          {renderTextField("disambiguation", updateInput.disambiguation, (v) =>
            setUpdateField({ disambiguation: v })
          )}
          {renderTextField("birthdate", updateInput.birthdate, (v) =>
            setUpdateField({ birthdate: v })
          )}
          {renderTextField("death_date", updateInput.death_date, (v) =>
            setUpdateField({ death_date: v })
          )}

          <Form.Group>
            <Form.Label>
              <FormattedMessage id="country" />
            </Form.Label>
            <CountrySelect
              value={updateInput.country ?? ""}
              onChange={(v) => setUpdateField({ country: v })}
              showFlag
            />
          </Form.Group>

          {renderTextField("ethnicity", updateInput.ethnicity, (v) =>
            setUpdateField({ ethnicity: v })
          )}
          {renderTextField("hair_color", updateInput.hair_color, (v) =>
            setUpdateField({ hair_color: v })
          )}
          {renderTextField("eye_color", updateInput.eye_color, (v) =>
            setUpdateField({ eye_color: v })
          )}
          {renderTextField("height", height, (v) => setHeight(v))}
          {renderTextField("weight", weight, (v) => setWeight(v))}
          {renderTextField("measurements", updateInput.measurements, (v) =>
            setUpdateField({ measurements: v })
          )}
          {renderTextField("penis_length", penis_length, (v) =>
            setPenisLength(v)
          )}

          <Form.Group data-field="circumcised">
            <Form.Label>
              <FormattedMessage id="circumcised" />
            </Form.Label>
            <Form.Control
              as="select"
              className="input-control"
              value={circumcisedToString(updateInput.circumcised)}
              onChange={(event) =>
                setUpdateField({
                  circumcised: stringToCircumcised(event.currentTarget.value),
                })
              }
            >
              {circumcisedOptions.map((opt) => (
                <option value={opt} key={opt}>
                  {opt}
                </option>
              ))}
            </Form.Control>
          </Form.Group>

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
              onUpdate={(itemIDs) => setTagIds({ ...tagIds, ids: itemIDs })}
              onSetMode={(newMode) => setTagIds({ ...tagIds, mode: newMode })}
              existingIds={existingTagIds ?? []}
              ids={tagIds.ids ?? []}
              mode={tagIds.mode}
              menuPortalTarget={document.body}
            />
          </Form.Group>

          <Form.Group controlId="ignore-auto-tags">
            <IndeterminateCheckbox
              label={intl.formatMessage({ id: "ignore_auto_tag" })}
              setChecked={(checked) =>
                setUpdateField({ ignore_auto_tag: checked })
              }
              checked={updateInput.ignore_auto_tag ?? undefined}
            />
          </Form.Group>
        </Form>
      </ModalComponent>
    );
  }

  return render();
};
