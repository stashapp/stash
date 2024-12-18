import * as GQL from "src/core/generated-graphql";
import TextUtils from "src/utils/text";

export const scrapedGroupToCreateInput = (toCreate: GQL.ScrapedGroup) => {
  const input: GQL.GroupCreateInput = {
    name: toCreate.name ?? "",
    urls: toCreate.urls,
    aliases: toCreate.aliases,
    front_image: toCreate.front_image,
    back_image: toCreate.back_image,
    synopsis: toCreate.synopsis,
    date: toCreate.date,
    director: toCreate.director,
    // #788 - convert duration and rating to the correct type
    duration: TextUtils.timestampToSeconds(toCreate.duration),
    studio_id: toCreate.studio?.stored_id,
    rating100: parseInt(toCreate.rating ?? "0", 10) * 20,
  };

  if (!input.duration) {
    input.duration = undefined;
  }

  if (!input.rating100 || Number.isNaN(input.rating100)) {
    input.rating100 = undefined;
  }

  return input;
};
