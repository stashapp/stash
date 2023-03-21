import * as GQL from "src/core/generated-graphql";

interface IImagePaths {
  image?: GQL.Maybe<string>;
  thumbnail?: GQL.Maybe<string>;
}

interface IImageFiles {
  clip?: boolean;
}

export interface ILightboxImage {
  id?: string;
  title?: GQL.Maybe<string>;
  rating100?: GQL.Maybe<number>;
  o_counter?: GQL.Maybe<number>;
  paths: IImagePaths;
  files?: GQL.Maybe<IImageFiles>[];
}

export interface IChapter {
  id: string;
  title: string;
  image_index: number;
}
