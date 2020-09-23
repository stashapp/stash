import React from "react";
import * as GQL from "src/core/generated-graphql";
import { GalleriesCriterion } from "src/models/list-filter/criteria/galleries";
import { ListFilterModel } from "src/models/list-filter/filter";
import { ImageList } from "src/components/Images/ImageList";
import { showWhenSelected } from "src/hooks/ListHook";
import { mutateAddGalleryImages } from "src/core/StashService";

interface IGalleryAddProps {
  gallery: Partial<GQL.GalleryDataFragment>;
}

export const GalleryAddPanel: React.FC<IGalleryAddProps> = ({
  gallery,
}) => {
  function filterHook(filter: ListFilterModel) {
    const galleryValue = { id: gallery.id!, label: gallery.title ?? gallery.path ?? "" };
    // if galleries is already present, then we modify it, otherwise add
    let galleryCriterion = filter.criteria.find((c) => {
      return c.type === "galleries";
    }) as GalleriesCriterion;

    if (
      galleryCriterion &&
      (galleryCriterion.modifier === GQL.CriterionModifier.Excludes)
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
    selectedIds: Set<string>) {
    await mutateAddGalleryImages({
      gallery_id: gallery.id!,
      image_ids: Array.from(selectedIds.values()),
    });
  }

  const otherOperations = [
    {
      text: "Add to Gallery",
      onClick: addImages,
      isDisplayed: showWhenSelected,
      postRefetch: true,
    }
  ];

  return <ImageList filterHook={filterHook} extraOperations={otherOperations} persistState={false}/>;
};
