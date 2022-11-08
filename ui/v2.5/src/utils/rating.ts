import * as GQL from "src/core/generated-graphql";

function round(value: number, step: number) {
  let denom = step;
  if (!denom) {
    denom = 1.0;
  }
  const inv = 1.0 / denom;
  return Math.round(value * inv) / inv;
}

export function getMaxStars(ratingSystem: GQL.RatingSystem) {
  switch (ratingSystem) {
    //case GQL.RatingSystem.TenStar:
    //case GQL.RatingSystem.TenPointFiveStar:
    //case GQL.RatingSystem.TenPointTwoFiveStar:
    //  return 10;
    case GQL.RatingSystem.FiveStar:
    case GQL.RatingSystem.FivePointFiveStar:
    case GQL.RatingSystem.FivePointTwoFiveStar:
      return 5;
    case GQL.RatingSystem.TenPointDecimal:
      return 0;
    default:
      return 5;
  }
}

export function getRatingPrecision(ratingSystem: GQL.RatingSystem) {
  switch (ratingSystem) {
    //case GQL.RatingSystem.TenStar:
    case GQL.RatingSystem.FiveStar:
      return 0;
    //case GQL.RatingSystem.TenPointFiveStar:
    case GQL.RatingSystem.FivePointFiveStar:
      return 0.5;
    //case GQL.RatingSystem.TenPointTwoFiveStar:
    case GQL.RatingSystem.FivePointTwoFiveStar:
      return 0.25;
    case GQL.RatingSystem.TenPointDecimal:
      return 0.1;
    default:
      return 0;
  }
}

export function convertToRatingFormat(
  rating: number | undefined,
  ratingSystem: GQL.RatingSystem
) {
  if (!rating) {
    return null;
  }
  let toReturn;
  switch (ratingSystem) {
    //case GQL.RatingSystem.TenStar:
    //  toReturn = round(rating / 10, 1);
    //  break;
    //case GQL.RatingSystem.TenPointFiveStar:
    //  toReturn = round(rating / 10, 0.5);
    //  break;
    //case GQL.RatingSystem.TenPointTwoFiveStar:
    //  toReturn = round(rating / 10, 0.25);
    //  break;
    case GQL.RatingSystem.FiveStar:
      toReturn = round(rating / 20, 1);
      break;
    case GQL.RatingSystem.FivePointFiveStar:
      toReturn = round(rating / 20, 0.5);
      break;
    case GQL.RatingSystem.FivePointTwoFiveStar:
      toReturn = round(rating / 20, 0.25);
      break;
    case GQL.RatingSystem.TenPointDecimal:
      toReturn = round(rating / 10, 0.1);
      break;
    default:
      toReturn = round(rating / 20, 1);
      break;
  }
  return toReturn;
}

export function convertFromRatingFormat(
  rating: number,
  ratingSystem: GQL.Maybe<GQL.RatingSystem> | undefined
) {
  let toReturn;
  switch (ratingSystem) {
    //case GQL.RatingSystem.TenStar:
    //  toReturn = Math.round(rating * 10);
    //  break;
    //case GQL.RatingSystem.TenPointFiveStar:
    //  toReturn = Math.round(rating * 10);
    //  break;
    //case GQL.RatingSystem.TenPointTwoFiveStar:
    //  toReturn = Math.round(rating * 10);
    //  break;
    case GQL.RatingSystem.FiveStar:
      toReturn = Math.round(rating * 20);
      break;
    case GQL.RatingSystem.FivePointFiveStar:
      toReturn = Math.round(rating * 20);
      break;
    case GQL.RatingSystem.FivePointTwoFiveStar:
      toReturn = Math.round(rating * 20);
      break;
    case GQL.RatingSystem.TenPointDecimal:
      toReturn = Math.round(rating * 10);
      break;
    default:
      toReturn = Math.round(rating * 20);
      break;
  }
  return toReturn;
}
