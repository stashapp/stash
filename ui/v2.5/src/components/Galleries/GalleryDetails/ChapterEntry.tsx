import React from "react";
import { FormattedMessage } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import { Button } from "react-bootstrap";

interface IChapterEntries {
  galleryChapters: GQL.GalleryChapterDataFragment[];
  onClickChapter: (image_index: number) => void;
  onEdit: (chapter: GQL.GalleryChapterDataFragment) => void;
}

export const ChapterEntries: React.FC<IChapterEntries> = ({
  galleryChapters,
  onClickChapter,
  onEdit,
}) => {
  if (!galleryChapters?.length) return <div />;

  const chapterCards = galleryChapters.map((chapter) => {
    return (
      <div key={chapter.id}>
        <hr />
        <div className="row">
          <Button
            variant="link"
            onClick={() => onClickChapter(chapter.image_index)}
          >
            <div className="row">
              {chapter.title}
              {chapter.title.length > 0 ? " - #" : "#"}
              {chapter.image_index}
            </div>
          </Button>
          <Button
            variant="link"
            className="ml-auto"
            onClick={() => onEdit(chapter)}
          >
            <FormattedMessage id="actions.edit" />
          </Button>
        </div>
      </div>
    );
  });

  return <div>{chapterCards}</div>;
};
