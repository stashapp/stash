import { CriterionModifier } from "src/core/generated-graphql";
import {
  ModifierCriterionOption,
  IHierarchicalLabeledIdCriterion,
} from "./criterion";

const modifierOptions = [CriterionModifier.Includes];

const defaultModifier = CriterionModifier.Includes;
const inputType = "folders";

export const FolderCriterionOption = new ModifierCriterionOption({
  messageID: "folder",
  type: "folder",
  modifierOptions,
  defaultModifier,
  inputType,
  makeCriterion: () => new FolderCriterion(),
});

// for galleries, we should use parent folder to distinguish between gallery folder
// and parent folder of the gallery folder
export const ParentFolderCriterionOption = new ModifierCriterionOption({
  messageID: "parent_folder",
  type: "parent_folder",
  modifierOptions,
  defaultModifier,
  inputType,
  makeCriterion: () => new ParentFolderCriterion(),
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

export class ParentFolderCriterion extends IHierarchicalLabeledIdCriterion {
  constructor() {
    super(ParentFolderCriterionOption);
  }

  public applyToCriterionInput(input: Record<string, unknown>) {
    input.parent_folder = this.toCriterionInput();
  }
}
