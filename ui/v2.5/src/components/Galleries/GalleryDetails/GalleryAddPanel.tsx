import React from "react";
import * as GQL from "src/core/generated-graphql";
import { GalleriesCriterion } from "src/models/list-filter/criteria/galleries";
import { ListFilterModel } from "src/models/list-filter/filter";
import { ImageList } from "src/components/Images/ImageList";
import { showWhenSelected } from "src/components/List/ItemList";
import { mutateAddGalleryImages } from "src/core/StashService";
import { useToast } from "src/hooks/Toast";
import { useIntl } from "react-intl";
import { faPlus } from "@fortawesome/free-solid-svg-icons";
import { galleryTitle } from "src/core/galleries";
import { IItemListOperation } from "src/components/List/FilteredListToolbar";
import { PatchComponent } from "src/patch";

interface IGalleryAddProps {
  active: boolean;
  gallery: GQL.GalleryDataFragment;
  extraOperations?: IItemListOperation<GQL.FindImagesQueryResult>[];
}

export const GalleryAddPanel: React.FC<IGalleryAddProps> = PatchComponent(
  "GalleryAddPanel",
  ({ active, gallery, extraOperations = [] }) => {
    const Toast = useToast();
    const intl = useIntl();

    function filterHook(filter: ListFilterModel) {
      const galleryValue = {
        id: gallery.id,
        label: galleryTitle(gallery),
      };
      // if galleries is already present, then we modify it, otherwise add
      let galleryCriterion = filter.criteria.find((c) => {
        return c.criterionOption.type === "galleries";
      }) as GalleriesCriterion | undefined;

      if (
        galleryCriterion &&
        galleryCriterion.modifier === GQL.CriterionModifier.Excludes
      ) {
        // add the gallery if not present
        if (
          !galleryCriterion.value.find((p) => {
            return p.id === gallery.id;
          })
        ) {
          galleryCriterion.value.push(galleryValue);
        }

        galleryCriterion.modifier = GQL.CriterionModifier.Excludes;
      } else {
        // overwrite
        galleryCriterion = new GalleriesCriterion();
        galleryCriterion.modifier = GQL.CriterionModifier.Excludes;
        galleryCriterion.value = [galleryValue];
        filter.criteria.push(galleryCriterion);
      }

      return filter;
    }

    async function addImages(
      result: GQL.FindImagesQueryResult,
      filter: ListFilterModel,
      selectedIds: Set<string>
    ) {
      try {
        await mutateAddGalleryImages({
          gallery_id: gallery.id!,
          image_ids: Array.from(selectedIds.values()),
        });
        const imageCount = selectedIds.size;
        Toast.success(
          intl.formatMessage(
            { id: "toast.added_entity" },
            {
              count: imageCount,
              singularEntity: intl.formatMessage({ id: "image" }),
              pluralEntity: intl.formatMessage({ id: "images" }),
            }
          )
        );
      } catch (e) {
        Toast.error(e);
      }
    }

    const otherOperations = [
      ...extraOperations,
      {
        text: intl.formatMessage(
          { id: "actions.add_to_entity" },
          { entityType: intl.formatMessage({ id: "gallery" }) }
        ),
        onClick: addImages,
        isDisplayed: showWhenSelected,
        postRefetch: true,
        icon: faPlus,
      },
    ];

    return (
      <ImageList
        filterHook={filterHook}
        extraOperations={otherOperations}
        alterQuery={active}
      />
    );
  }
);
