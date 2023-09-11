import { StringCriterion, StringCriterionOption } from "./criterion";

export const PathCriterionOption = new StringCriterionOption(
  "path",
  "path",
  () => new PathCriterion()
);

export class PathCriterion extends StringCriterion {
  constructor() {
    super(PathCriterionOption);
  }
}
