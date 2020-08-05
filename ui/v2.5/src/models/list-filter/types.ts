// NOTE: add new enum values to the end, to ensure existing data
// is not impacted
export enum DisplayMode {
  Grid,
  List,
  Wall,
}

export enum FilterMode {
  Scenes,
  Performers,
  Studios,
  Galleries,
  SceneMarkers,
  Movies,
  Tags,
}

export interface ILabeledId {
  id: string;
  label: string;
}

export interface ILabeledValue {
  label: string;
  value: string;
}

export function encodeILabeledId(o: ILabeledId) {
  return { ...o, label: encodeURIComponent(o.label) };
}

export interface IOptionType {
  id: string;
  name?: string;
  image_path?: string;
}
