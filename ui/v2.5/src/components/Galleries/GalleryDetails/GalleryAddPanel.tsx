import React from "react";
import * as GQL from "src/core/generated-graphql";
import { GalleriesCriterion } from "src/models/list-filter/criteria/galleries";
import { ListFilterModel } from "src/models/list-filter/filter";
import { ImageList } from "src/components/Images/ImageList";
import { showWhenSelected } from "src/hooks/ListHook";
import { mutateAddGalleryImages } from "src/core/StashService";
import { useToast } from "src/hooks";
import { TextUtils } from "src/utils";
import { useIntl } from "react-intl";
import { IconProp } from "@fortawesome/fontawesome-svg-core";

interface IGalleryAddProps {
  gallery: GQL.GalleryDataFragment;
}

export const GalleryAddPanel: React.FC<IGalleryAddProps> = ({ gallery }) => {
  const Toast = useToast();
  const intl = useIntl();

  function filterHook(filter: ListFilterModel) {
    const galleryValue = {
      id: gallery.id,
      label: gallery.title ?? TextUtils.fileNameFromPath(gallery.path ?? ""),
    };
    // if galleries is already present, then we modify it, otherwise add
    let galleryCriterion = filter.criteria.find((c) => {
      return c.criterionOption.type === "galleries";
    }) as GalleriesCriterion;

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
      Toast.success({
        content: intl.formatMessage(
          { id: "toast.added_entity" },
          {
            count: imageCount,
            singularEntity: intl.formatMessage(
              { id: "countables.images" },
              { count: 1 }
            ),
            pluralEntity: intl.formatMessage(
              { id: "countables.images" },
              { count: imageCount }
            ),
          }
        ),
      });
    } catch (e) {
      Toast.error(e);
    }
  }

  const otherOperations = [
    {
      text: intl.formatMessage(
        { id: "actions.add_to_entity" },
        {
          entityType: intl.formatMessage(
            { id: "countables.galleries" },
            { count: 1 }
          ),
        }
      ),
      onClick: addImages,
      isDisplayed: showWhenSelected,
      postRefetch: true,
      icon: "plus" as IconProp,
    },
  ];

  return (
    <ImageList filterHook={filterHook} extraOperations={otherOperations} />
  );
};
