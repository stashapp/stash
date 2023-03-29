import * as GQL from "src/core/generated-graphql";

interface IImagePaths {
  image?: GQL.Maybe<string>;
  thumbnail?: GQL.Maybe<string>;
}

interface IFiles {
  __typename?: string;
  width: number;
  height: number;
}

export interface ILightboxImage {
  id?: string;
  title?: GQL.Maybe<string>;
  rating100?: GQL.Maybe<number>;
  o_counter?: GQL.Maybe<number>;
  paths: IImagePaths;
  visual_files?: GQL.Maybe<IFiles>[];
}

export interface IChapter {
  id: string;
  title: string;
  image_index: number;
}
