import { Maybe } from "src/core/generated-graphql";

// returns true if the file should be treated as a video in the UI
export function isVideo(o: {
  __typename?: string;
  video_codec?: Maybe<string>;
}) {
  return o.__typename == "VideoFile" && o.video_codec != "gif";
}
