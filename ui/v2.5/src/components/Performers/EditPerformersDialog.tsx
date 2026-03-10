import React, { useEffect, useState } from "react";
import { Form } from "react-bootstrap";
import { useIntl } from "react-intl";
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
import { BulkUpdateFormGroup, BulkUpdateTextInput } from "../Shared/BulkUpdate";
import { faPencilAlt } from "@fortawesome/free-solid-svg-icons";
import { CountrySelect } from "../Shared/CountrySelect";
import { useConfigurationContext } from "src/hooks/Config";
import cx from "classnames";
import { BulkUpdateDateInput } from "../Shared/DateInput";
import { getDateError } from "src/utils/yup";

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
  "career_start",
  "career_end",
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
  const [height, setHeight] = useState<string | undefined | null>();
  const [weight, setWeight] = useState<string | undefined | null>();
  const [penis_length, setPenisLength] = useState<string | undefined | null>();
  const [updateInput, setUpdateInput] = useState<GQL.BulkPerformerUpdateInput>(
    {}
  );
  const genderOptions = [""].concat(genderStrings);
  const circumcisedOptions = [""].concat(circumcisedStrings);

  const unsetDisabled = props.selected.length < 2;

  const [updatePerformers] = useBulkPerformerUpdate(getPerformerInput());

  const [birthdateError, setBirthdateError] = useState<string | undefined>();
  const [deathDateError, setDeathDateError] = useState<string | undefined>();

  useEffect(() => {
    setBirthdateError(getDateError(updateInput.birthdate ?? "", intl));
  }, [updateInput.birthdate, intl]);

  useEffect(() => {
    setDeathDateError(getDateError(updateInput.death_date ?? "", intl));
  }, [updateInput.death_date, intl]);

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
      performerInput.height_cm = height === null ? null : parseFloat(height);
    }
    if (weight !== undefined) {
      performerInput.weight = weight === null ? null : parseFloat(weight);
    }
    if (penis_length !== undefined) {
      performerInput.penis_length =
        penis_length === null ? null : parseFloat(penis_length);
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
          { id: "dialogs.edit_entity_count_title" },
          {
            count: props?.selected?.length ?? 1,
            singularEntity: intl.formatMessage({ id: "performer" }),
            pluralEntity: intl.formatMessage({ id: "performers" }),
          }
        )}
        accept={{
          onClick: onSave,
          text: intl.formatMessage({ id: "actions.apply" }),
        }}
        disabled={isUpdating || !!birthdateError || !!deathDateError}
        cancel={{
          onClick: () => props.onClose(false),
          text: intl.formatMessage({ id: "actions.cancel" }),
          variant: "secondary",
        }}
        isRunning={isUpdating}
      >
        <Form>
          <BulkUpdateFormGroup name="rating">
            <RatingSystem
              value={updateInput.rating100}
              onSetRating={(value) =>
                setUpdateField({ rating100: value ?? undefined })
              }
              disabled={isUpdating}
            />
          </BulkUpdateFormGroup>

          <Form.Group controlId="favorite">
            <IndeterminateCheckbox
              setChecked={(checked) => setUpdateField({ favorite: checked })}
              checked={updateInput.favorite ?? undefined}
              label={intl.formatMessage({ id: "favourite" })}
            />
          </Form.Group>

          <BulkUpdateFormGroup name="gender">
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
          </BulkUpdateFormGroup>

          <BulkUpdateFormGroup name="disambiguation">
            <BulkUpdateTextInput
              value={updateInput.disambiguation}
              valueChanged={(newValue) =>
                setUpdateField({ disambiguation: newValue })
              }
              unsetDisabled={unsetDisabled}
            />
          </BulkUpdateFormGroup>

          <BulkUpdateFormGroup name="birthdate">
            <BulkUpdateDateInput
              value={updateInput.birthdate}
              valueChanged={(newValue) =>
                setUpdateField({ birthdate: newValue })
              }
              unsetDisabled={unsetDisabled}
              error={birthdateError}
            />
          </BulkUpdateFormGroup>
          <BulkUpdateFormGroup name="death_date">
            <BulkUpdateDateInput
              value={updateInput.death_date}
              valueChanged={(newValue) =>
                setUpdateField({ death_date: newValue })
              }
              unsetDisabled={unsetDisabled}
              error={deathDateError}
            />
          </BulkUpdateFormGroup>
          <BulkUpdateFormGroup name="country">
            <CountrySelect
              value={updateInput.country ?? ""}
              onChange={(v) => setUpdateField({ country: v })}
              showFlag
            />
          </BulkUpdateFormGroup>

          <BulkUpdateFormGroup name="ethnicity">
            <BulkUpdateTextInput
              value={updateInput.ethnicity}
              valueChanged={(newValue) =>
                setUpdateField({ ethnicity: newValue })
              }
              unsetDisabled={unsetDisabled}
            />
          </BulkUpdateFormGroup>
          <BulkUpdateFormGroup name="hair_color">
            <BulkUpdateTextInput
              value={updateInput.hair_color}
              valueChanged={(newValue) =>
                setUpdateField({ hair_color: newValue })
              }
              unsetDisabled={unsetDisabled}
            />
          </BulkUpdateFormGroup>
          <BulkUpdateFormGroup name="eye_color">
            <BulkUpdateTextInput
              value={updateInput.eye_color}
              valueChanged={(newValue) =>
                setUpdateField({ eye_color: newValue })
              }
              unsetDisabled={unsetDisabled}
            />
          </BulkUpdateFormGroup>
          <BulkUpdateFormGroup name="height">
            <BulkUpdateTextInput
              value={height}
              valueChanged={(newValue) => setHeight(newValue)}
              unsetDisabled={unsetDisabled}
            />
          </BulkUpdateFormGroup>
          <BulkUpdateFormGroup name="weight">
            <BulkUpdateTextInput
              value={weight}
              valueChanged={(newValue) => setWeight(newValue)}
              unsetDisabled={unsetDisabled}
            />
          </BulkUpdateFormGroup>
          <BulkUpdateFormGroup name="measurements">
            <BulkUpdateTextInput
              value={updateInput.measurements}
              valueChanged={(newValue) =>
                setUpdateField({ measurements: newValue })
              }
              unsetDisabled={unsetDisabled}
            />
          </BulkUpdateFormGroup>
          <BulkUpdateFormGroup name="penis_length">
            <BulkUpdateTextInput
              value={penis_length}
              valueChanged={(newValue) => setPenisLength(newValue)}
              unsetDisabled={unsetDisabled}
            />
          </BulkUpdateFormGroup>

          <BulkUpdateFormGroup name="circumcised">
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
          </BulkUpdateFormGroup>

          <BulkUpdateFormGroup name="fake_tits">
            <BulkUpdateTextInput
              value={updateInput.fake_tits}
              valueChanged={(newValue) =>
                setUpdateField({ fake_tits: newValue })
              }
              unsetDisabled={unsetDisabled}
            />
          </BulkUpdateFormGroup>
          <BulkUpdateFormGroup name="tattoos">
            <BulkUpdateTextInput
              value={updateInput.tattoos}
              valueChanged={(newValue) => setUpdateField({ tattoos: newValue })}
              unsetDisabled={unsetDisabled}
            />
          </BulkUpdateFormGroup>
          <BulkUpdateFormGroup name="piercings">
            <BulkUpdateTextInput
              value={updateInput.piercings}
              valueChanged={(newValue) =>
                setUpdateField({ piercings: newValue })
              }
              unsetDisabled={unsetDisabled}
            />
          </BulkUpdateFormGroup>
          <BulkUpdateFormGroup name="career_start">
            <BulkUpdateTextInput
              value={updateInput.career_start?.toString()}
              valueChanged={(v) =>
                setUpdateField({ career_start: v ? parseInt(v) : undefined })
              }
              unsetDisabled={unsetDisabled}
            />
          </BulkUpdateFormGroup>
          <BulkUpdateFormGroup name="career_end">
            <BulkUpdateTextInput
              value={updateInput.career_end?.toString()}
              valueChanged={(v) =>
                setUpdateField({ career_end: v ? parseInt(v) : undefined })
              }
              unsetDisabled={unsetDisabled}
            />
          </BulkUpdateFormGroup>

          <BulkUpdateFormGroup name="tags" inline={false}>
            <MultiSet
              type={"tags"}
              disabled={isUpdating}
              onUpdate={(itemIDs) => {
                setTagIds((c) => ({ ...c, ids: itemIDs }));
              }}
              onSetMode={(newMode) => {
                setTagIds((c) => ({ ...c, mode: newMode }));
              }}
              ids={tagIds.ids ?? []}
              existingIds={existingTagIds}
              mode={tagIds.mode}
              menuPortalTarget={document.body}
            />
          </BulkUpdateFormGroup>

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
