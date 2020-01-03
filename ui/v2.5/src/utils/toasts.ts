import { Intent, Position, Toaster } from "@blueprintjs/core";

const toaster = Toaster.create({
  position: Position.TOP,
});

export class ToastUtils {
  public static success(message: string) {
    toaster.show({
      message,
      intent: Intent.SUCCESS,
    });
  }
}
