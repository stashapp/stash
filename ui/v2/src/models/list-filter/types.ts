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
}

export interface ILabeledId {
  id: string;
  label: string;
}

export function encodeILabeledId(o: ILabeledId) {
  let ret = Object.assign({}, o);
  ret.label = encodeURIComponent(o.label);
  return ret;
}

export interface ILabeledValue {
  label: string;
  value: string;
}
