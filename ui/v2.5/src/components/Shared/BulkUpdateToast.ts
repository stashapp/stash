import { IntlShape } from "react-intl";
import { useToast } from "src/hooks/Toast";

// Bulk edits skip organized items server-side, so the result can hold fewer
// items than were selected; warn instead of reporting success.
// The server returns null rather than [] when every item was skipped.
export function toastBulkUpdated(
  Toast: ReturnType<typeof useToast>,
  intl: IntlShape,
  entityId: string,
  selected: number,
  updated: readonly unknown[] | null | undefined
) {
  const skipped = selected - (updated?.length ?? 0);
  if (skipped > 0) {
    Toast.toast({
      variant: "warning",
      content: intl.formatMessage(
        { id: "organized_items_skipped" },
        { count: skipped }
      ),
    });
    return;
  }

  Toast.success(
    intl.formatMessage(
      { id: "toast.updated_entity" },
      { entity: intl.formatMessage({ id: entityId }).toLocaleLowerCase() }
    )
  );
}
