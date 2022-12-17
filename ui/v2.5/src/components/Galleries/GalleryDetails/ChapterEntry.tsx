import React from "react";
import { FormattedMessage } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import { Button, Badge, Card } from "react-bootstrap";
import TextUtils from "src/utils/text";

interface IChapterEntries {
  galleryChapters: GQL.GalleryChapterDataFragment[];
  onClickChapter: (page_number: int) => void;
  onEdit: (chapter: GQL.GalleryChapterDataFragment) => void;
}

export const ChapterEntries: React.FC<IChapterEntries> = ({
  galleryChapters,
  onClickChapter,
  onEdit,
}) => {
  if (!galleryChapters?.length) return <div />;

  const chapterCards = galleryChapters.map(chapter => {
      return (
        <div key={chapter.id}>
          <hr />
          <div className="row">
            <Button variant="link" onClick={() => onClickChapter(chapter.page_number)}>
              <div className="row">{chapter.title} - <FormattedMessage id="page" /> {chapter.page_number}</div>
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
