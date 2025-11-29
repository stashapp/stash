import { TitleDuplicationCriterionInput } from "src/core/generated-graphql";
import { BooleanCriterionOption, StringCriterion } from "./criterion";

export const DuplicatedTitleCriterionOption = new BooleanCriterionOption(
  "duplicated_title",
  "duplicated_title",
  () => new DuplicatedTitleCriterion()
);

export class DuplicatedTitleCriterion extends StringCriterion {
  constructor() {
    super(DuplicatedTitleCriterionOption);
  }

  public toCriterionInput(): TitleDuplicationCriterionInput {
    return {
      duplicated: this.value === "true",
    };
  }
}