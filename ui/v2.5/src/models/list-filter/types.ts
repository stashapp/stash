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
  SceneMarkers
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
  const ret = { ...o };
  ret.label = encodeURIComponent(o.label);
  return ret;
}

export interface IOptionType {
  id: string;
  name?: string;
  image_path?: string;
}
