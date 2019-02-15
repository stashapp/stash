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
}

export interface ILabeledId {
  id: string;
  label: string;
}
