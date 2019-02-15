import { Intent, Position, Toaster } from "@blueprintjs/core";
import { ApolloError } from "apollo-boost";

const toaster = Toaster.create({
  position: Position.TOP,
});

export class ErrorUtils {
  public static handle(error: any) {
    console.error(error);
    toaster.show({
      message: error.toString(),
      intent: Intent.DANGER,
    });
  }

  public static handleApolloError(error: ApolloError) {
    console.error(error);
    toaster.show({
      message: error.message,
      intent: Intent.DANGER,
    });
  }
}
