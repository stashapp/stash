import {
  StringBooleanCriterion,
  StringBooleanCriterionOption,
} from "./criterion";

export const HasChaptersCriterionOption = new StringBooleanCriterionOption(
  "hasChapters",
  "has_chapters",
  () => new HasChaptersCriterion()
);

export class HasChaptersCriterion extends StringBooleanCriterion {
  constructor() {
    super(HasChaptersCriterionOption);
  }
}
