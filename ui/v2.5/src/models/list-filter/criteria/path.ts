import { StringCriterion, StringCriterionOption } from "./criterion";

export const PathCriterionOption = new CriterionOption({
  messageID: "path",
  type: "path",
  modifierOptions: StringCriterionOption.modifierOptions,
  defaultModifier: StringCriterionOption.defaultModifier,
  makeCriterion: () => new PathCriterion(),
  inputType: StringCriterionOption.inputType,
});

export class PathCriterion extends StringCriterion {
  constructor() {
    super(PathCriterionOption);
  }
}
