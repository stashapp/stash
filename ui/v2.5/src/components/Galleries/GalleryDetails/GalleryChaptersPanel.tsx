import React, { useState, useEffect } from "react";
import { Button } from "react-bootstrap";
import { FormattedMessage } from "react-intl";
import Mousetrap from "mousetrap";
import * as GQL from "src/core/generated-graphql";
import { ChapterEntries } from "./ChapterEntry";
import { GalleryChapterForm } from "./GalleryChapterForm";

interface IGalleryChapterPanelProps {
  gallery: GQL.GalleryDataFragment;
  isVisible: boolean;
  onClickChapter: (index: number) => void;
}

export const GalleryChapterPanel: React.FC<IGalleryChapterPanelProps> = (
  props: IGalleryChapterPanelProps
) => {
  const [isEditorOpen, setIsEditorOpen] = useState<boolean>(false);
  const [editingChapter, setEditingChapter] =
    useState<GQL.GalleryChapterDataFragment>();

  // set up hotkeys
  useEffect(() => {
    if (props.isVisible) {
      Mousetrap.bind("n", () => onOpenEditor());

      return () => {
        Mousetrap.unbind("n");
      };
    }
  });

  function onOpenEditor(chapter?: GQL.GalleryChapterDataFragment) {
    setIsEditorOpen(true);
    setEditingChapter(chapter ?? undefined);
  }

  function onClickChapter(image_index: number) {
    props.onClickChapter(image_index);
  }

  const closeEditor = () => {
    setEditingChapter(undefined);
    setIsEditorOpen(false);
  };

  if (isEditorOpen)
    return (
      <GalleryChapterForm
        galleryID={props.gallery.id}
        editingChapter={editingChapter}
        onClose={closeEditor}
      />
    );

  return (
    <div>
      <Button onClick={() => onOpenEditor()}>
        <FormattedMessage id="actions.create_chapters" />
      </Button>
      <div className="container">
        <ChapterEntries
          galleryChapters={props.gallery.chapters}
          onClickChapter={onClickChapter}
          onEdit={onOpenEditor}
        />
      </div>
    </div>
  );
};

export default GalleryChapterPanel;
