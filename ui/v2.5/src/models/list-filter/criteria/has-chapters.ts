import { CriterionOption, StringCriterion } from "./criterion";

export const HasChaptersCriterionOption = new CriterionOption({
  messageID: "hasChapters",
  type: "hasChapters",
  parameterName: "has_chapters",
  options: [true.toString(), false.toString()],
});

export class HasChaptersCriterion extends StringCriterion {
  constructor() {
    super(HasChaptersCriterionOption);
  }

  protected toCriterionInput(): string {
    return this.value;
  }
}
