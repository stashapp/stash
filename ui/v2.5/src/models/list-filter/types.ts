// NOTE: add new enum values to the end, to ensure existing data
// is not impacted
export enum DisplayMode {
  Grid,
  List,
  Wall
}

export enum FilterMode {
  Scenes,
  Performers,
  Studios,
  Galleries,
  SceneMarkers,
  Movies,
}

export interface ILabeledId {
  id: string;
  label: string;
}

export interface ILabeledValue {
  label: string;
  value: string;
}

export interface IOptionType {
  id: string;
  name?: string;
  image_path?: string;
}
