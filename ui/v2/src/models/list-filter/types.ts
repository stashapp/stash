export enum DisplayMode {
  Grid,
  List,
  Wall,
}

export enum FilterMode {
  Scenes,
  Performers,
  Studios,
  Movies,
  Galleries,
  SceneMarkers,
}

export interface ILabeledId {
  id: string;
  label: string;
}

export interface ILabeledValue {
  label: string;
  value: string;
}
