import { CriterionModifier } from "src/core/generated-graphql";
import { StringCriterion, StringCriterionOption } from "./criterion";

export const PathCriterionOption = new StringCriterionOption({
  messageID: "path",
  type: "path",
  defaultModifier: CriterionModifier.Includes,
  makeCriterion: () => new PathCriterion(),
});

export class PathCriterion extends StringCriterion {
  constructor() {
    super(PathCriterionOption);
  }
}
