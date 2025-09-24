import * as GQL from "src/core/generated-graphql";

interface IImagePaths {
  image?: GQL.Maybe<string>;
  thumbnail?: GQL.Maybe<string>;
  preview?: GQL.Maybe<string>;
}

interface IFiles {
  __typename?: string;
  path: string;
  width: number;
  height: number;
  video_codec?: GQL.Maybe<string>;
}

interface IWithPath {
  path: string;
}

export interface IGallery {
  id: string;
  title?: GQL.Maybe<string>;
  files?: GQL.Maybe<IWithPath[]>;
  folder?: GQL.Maybe<IWithPath>;
}

export interface ILightboxImage {
  id?: string;
  title?: GQL.Maybe<string>;
  rating100?: GQL.Maybe<number>;
  o_counter?: GQL.Maybe<number>;
  paths: IImagePaths;
  visual_files?: IFiles[];
  galleries?: GQL.Maybe<IGallery[]>;
}

export interface IChapter {
  id: string;
  title: string;
  image_index: number;
}
