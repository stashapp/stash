export enum ImageWallDirection {
  Column = "column",
  Row = "row",
}

export type ImageWallOptions = {
  margin: number;
  direction: ImageWallDirection;
};

export const defaultImageWallDirection: ImageWallDirection =
  ImageWallDirection.Row;
export const defaultImageWallMargin = 3;

export const imageWallDirectionIntlMap = new Map<ImageWallDirection, string>([
  [ImageWallDirection.Column, "dialogs.imagewall.direction.column"],
  [ImageWallDirection.Row, "dialogs.imagewall.direction.row"],
]);

export const defaultImageWallOptions = {
  margin: defaultImageWallMargin,
  direction: defaultImageWallDirection,
};
