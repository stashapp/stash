import { CriterionModifier } from "src/core/generated-graphql";
import {
  ModifierCriterionOption,
  IHierarchicalLabeledIdCriterion,
} from "./criterion";

const modifierOptions = [CriterionModifier.Includes];

const defaultModifier = CriterionModifier.Includes;
const inputType = "text"; // FIXME

export const FolderCriterionOption = new ModifierCriterionOption({
  messageID: "folder",
  type: "folder",
  modifierOptions,
  defaultModifier,
  inputType,
  makeCriterion: () => new FolderCriterion(),
});

export class FolderCriterion extends IHierarchicalLabeledIdCriterion {
  constructor() {
    super(FolderCriterionOption);
  }

  public applyToCriterionInput(input: Record<string, unknown>) {
    input.files_filter = {
      parent_folder: this.toCriterionInput(),
    };
  }
}
