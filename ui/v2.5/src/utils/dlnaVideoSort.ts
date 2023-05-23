export enum VideoSortOrder {
  Created_At = "created_at",
  Date = "date",
  Random = "random",
  Title = "title",
  Updated_At = "updated_at",
}

export const defaultVideoSort = VideoSortOrder.Title;

export const videoSortOrderIntlMap = new Map<VideoSortOrder, string>([
  [VideoSortOrder.Created_At, "created_at"],
  [VideoSortOrder.Date, "date"],
  [VideoSortOrder.Random, "random"],
  [VideoSortOrder.Title, "title"],
  [VideoSortOrder.Updated_At, "updated_at"],
]);
