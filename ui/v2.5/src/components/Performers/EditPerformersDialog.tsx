import React, { useEffect, useState } from "react";
import { Form } from "react-bootstrap";
import _ from "lodash";
import { useBulkPerformerUpdate } from "src/core/StashService";
import * as GQL from "src/core/generated-graphql";
import { Modal } from "src/components/Shared";
import { useToast } from "src/hooks";
import MultiSet from "../Shared/MultiSet";

interface IListOperationProps {
  selected: GQL.SlimPerformerDataFragment[];
  onClose: (applied: boolean) => void;
}

export const EditPerformersDialog: React.FC<IListOperationProps> = (
  props: IListOperationProps
) => {
  const Toast = useToast();

  const [tagMode, setTagMode] = React.useState<GQL.BulkUpdateIdMode>(
    GQL.BulkUpdateIdMode.Add
  );
  const [tagIds, setTagIds] = useState<string[]>();
  const [favorite, setFavorite] = useState<boolean | undefined>();

  const [updatePerformers] = useBulkPerformerUpdate(getPerformerInput());

  // Network state
  const [isUpdating, setIsUpdating] = useState(false);

  const checkboxRef = React.createRef<HTMLInputElement>();

  function makeBulkUpdateIds(
    ids: string[],
    mode: GQL.BulkUpdateIdMode
  ): GQL.BulkUpdateIds {
    return {
      mode,
      ids,
    };
  }

  function getPerformerInput(): GQL.BulkPerformerUpdateInput {
    // need to determine what we are actually setting on each performer
    const aggregateTagIds = getTagIds(props.selected);

    const performerInput: GQL.BulkPerformerUpdateInput = {
      ids: props.selected.map((performer) => {
        return performer.id;
      }),
    };

    // if tagIds non-empty, then we are setting them
    if (
      tagMode === GQL.BulkUpdateIdMode.Set &&
      (!tagIds || tagIds.length === 0)
    ) {
      // and all performers have the same ids,
      if (aggregateTagIds.length > 0) {
        // then unset the tagIds, otherwise ignore
        performerInput.tag_ids = makeBulkUpdateIds(tagIds || [], tagMode);
      }
    } else {
      // if tagIds non-empty, then we are setting them
      performerInput.tag_ids = makeBulkUpdateIds(tagIds || [], tagMode);
    }

    if (favorite !== undefined) {
      performerInput.favorite = favorite;
    }

    return performerInput;
  }

  async function onSave() {
    setIsUpdating(true);
    try {
      await updatePerformers();
      Toast.success({ content: "Updated performers" });
      props.onClose(true);
    } catch (e) {
      Toast.error(e);
    }
    setIsUpdating(false);
  }

  function getTagIds(state: GQL.SlimPerformerDataFragment[]) {
    let ret: string[] = [];
    let first = true;

    state.forEach((performer: GQL.SlimPerformerDataFragment) => {
      if (first) {
        ret = performer.tags ? performer.tags.map((t) => t.id).sort() : [];
        first = false;
      } else {
        const tIds = performer.tags
          ? performer.tags.map((t) => t.id).sort()
          : [];

        if (!_.isEqual(ret, tIds)) {
          ret = [];
        }
      }
    });

    return ret;
  }

  useEffect(() => {
    const state = props.selected;
    let updateTagIds: string[] = [];
    let updateFavorite: boolean | undefined;
    let first = true;

    state.forEach((performer: GQL.SlimPerformerDataFragment) => {
      const performerTagIDs = (performer.tags ?? []).map((p) => p.id).sort();

      if (first) {
        updateTagIds = performerTagIDs;
        first = false;
        updateFavorite = performer.favorite;
      } else {
        if (!_.isEqual(performerTagIDs, updateTagIds)) {
          updateTagIds = [];
        }
        if (performer.favorite !== updateFavorite) {
          updateFavorite = undefined;
        }
      }
    });

    if (tagMode === GQL.BulkUpdateIdMode.Set) {
      setTagIds(updateTagIds);
    }
    setFavorite(updateFavorite);
  }, [props.selected, tagMode]);

  useEffect(() => {
    if (checkboxRef.current) {
      checkboxRef.current.indeterminate = favorite === undefined;
    }
  }, [favorite, checkboxRef]);

  function renderMultiSelect(
    type: "performers" | "tags",
    ids: string[] | undefined
  ) {
    let mode = GQL.BulkUpdateIdMode.Add;
    switch (type) {
      case "tags":
        mode = tagMode;
        break;
    }

    return (
      <MultiSet
        type={type}
        disabled={isUpdating}
        onUpdate={(items) => {
          const itemIDs = items.map((i) => i.id);
          switch (type) {
            case "tags":
              setTagIds(itemIDs);
              break;
          }
        }}
        onSetMode={(newMode) => {
          switch (type) {
            case "tags":
              setTagMode(newMode);
              break;
          }
        }}
        ids={ids ?? []}
        mode={mode}
      />
    );
  }

  function cycleFavorite() {
    if (favorite) {
      setFavorite(undefined);
    } else if (favorite === undefined) {
      setFavorite(false);
    } else {
      setFavorite(true);
    }
  }

  function render() {
    return (
      <Modal
        show
        icon="pencil-alt"
        header="Edit Performers"
        accept={{ onClick: onSave, text: "Apply" }}
        cancel={{
          onClick: () => props.onClose(false),
          text: "Cancel",
          variant: "secondary",
        }}
        isRunning={isUpdating}
      >
        <Form>
          <Form.Group controlId="tags">
            <Form.Label>Tags</Form.Label>
            {renderMultiSelect("tags", tagIds)}
          </Form.Group>

          <Form.Group controlId="favorite">
            <Form.Check
              type="checkbox"
              label="Favorite"
              checked={favorite}
              ref={checkboxRef}
              onChange={() => cycleFavorite()}
            />
          </Form.Group>
        </Form>
      </Modal>
    );
  }

  return render();
};
