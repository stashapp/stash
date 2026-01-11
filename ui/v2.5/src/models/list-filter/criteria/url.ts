// This file is separate so the URL duplication criterion can be reused
// for other entity types (performers, galleries, studios, etc.) in the future.
import { UrlDuplicationCriterionInput } from "src/core/generated-graphql";
import { BooleanCriterionOption, StringCriterion } from "./criterion";

export const DuplicatedURLCriterionOption = new BooleanCriterionOption(
  "duplicated_url",
  "duplicated_url",
  () => new DuplicatedURLCriterion()
);

export class DuplicatedURLCriterion extends StringCriterion {
  constructor() {
    super(DuplicatedURLCriterionOption);
  }

  public toCriterionInput(): UrlDuplicationCriterionInput {
    return {
      duplicated: this.value === "true",
    };
  }
}
