import * as GQL from "src/core/generated-graphql";
import DurationUtils from "src/utils/duration";

export const scrapedMovieToCreateInput = (toCreate: GQL.ScrapedMovie) => {
  const input: GQL.MovieCreateInput = {
    name: toCreate.name ?? "",
    url: toCreate.url,
    aliases: toCreate.aliases,
    front_image: toCreate.front_image,
    back_image: toCreate.back_image,
    synopsis: toCreate.synopsis,
    date: toCreate.date,
    // #788 - convert duration and rating to the correct type
    duration: toCreate.duration
      ? DurationUtils.stringToSeconds(toCreate.duration)
      : undefined,
    studio_id: toCreate.studio?.stored_id,
    rating: parseInt(toCreate.rating ?? "0", 10),
  };

  if (!input.duration) {
    input.duration = undefined;
  }

  if (!input.rating || Number.isNaN(input.rating)) {
    input.rating = undefined;
  }

  return input;
};
