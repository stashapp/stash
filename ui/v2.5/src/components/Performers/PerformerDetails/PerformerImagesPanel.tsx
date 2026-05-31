import React, { useMemo } from "react";
import { useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import { FilteredImageList } from "src/components/Images/ImageList";
import { IItemListOperation } from "src/components/List/FilteredListToolbar";
import { usePerformerUpdate } from "src/core/StashService";
import { usePerformerFilterHook } from "src/core/performers";
import { useToast } from "src/hooks/Toast";
import { View } from "src/components/List/views";
import { PatchComponent } from "src/patch";
import ImageUtils from "src/utils/image";

interface IPerformerImagesPanel {
  active: boolean;
  performer: GQL.PerformerDataFragment;
}

export const PerformerImagesPanel: React.FC<IPerformerImagesPanel> =
  PatchComponent("PerformerImagesPanel", ({ active, performer }) => {
    const intl = useIntl();
    const Toast = useToast();
    const [updatePerformer] = usePerformerUpdate();
    const filterHook = usePerformerFilterHook(performer);

    const extraOperations = useMemo<
      IItemListOperation<GQL.FindImagesQueryResult>[]
    >(
      () => [
        {
          text: intl.formatMessage({
            id: "actions.set_as_performer_image",
          }),
          isDisplayed: (_result, _filter, selectedIds) =>
            selectedIds.size === 1,
          onClick: async (result, _filter, selectedIds) => {
            try {
              const [selectedId] = Array.from(selectedIds);
              const selectedImage = result.data?.findImages.images.find(
                (image) => image.id === selectedId
              );
              const imagePath = selectedImage?.paths.image;

              if (!imagePath) {
                throw new Error("Selected image does not have an image path");
              }

              const imageData = await ImageUtils.imageToDataURL(imagePath);
              await updatePerformer({
                variables: {
                  input: {
                    id: performer.id,
                    image: imageData,
                  },
                },
              });

              Toast.success(
                intl.formatMessage(
                  { id: "toast.updated_entity" },
                  { entity: intl.formatMessage({ id: "performer" }) }
                )
              );
            } catch (e) {
              Toast.error(e);
            }
          },
        },
      ],
      [Toast, intl, performer.id, updatePerformer]
    );

    return (
      <FilteredImageList
        filterHook={filterHook}
        alterQuery={active}
        extraOperations={extraOperations}
        view={View.PerformerImages}
      />
    );
  });
