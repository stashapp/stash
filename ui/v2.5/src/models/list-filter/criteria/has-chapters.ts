import { CriterionOption, StringCriterion } from "./criterion";

export const HasChaptersCriterionOption = new CriterionOption({
  messageID: "hasChapters",
  type: "has_chapters",
  options: [true.toString(), false.toString()],
  makeCriterion: () => new HasChaptersCriterion(),
});

export class HasChaptersCriterion extends StringCriterion {
  constructor() {
    super(HasChaptersCriterionOption);
  }

  protected toCriterionInput(): string {
    return this.value;
  }
}
