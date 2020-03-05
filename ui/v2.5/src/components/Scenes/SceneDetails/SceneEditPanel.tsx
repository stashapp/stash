/* eslint-disable react/no-this-in-sfc */

import React, { useEffect, useState } from "react";
import { Button, Dropdown, DropdownButton, Form, Table } from "react-bootstrap";
import * as GQL from "src/core/generated-graphql";
import { StashService } from "src/core/StashService";
import {
  PerformerSelect,
  TagSelect,
  StudioSelect,
  SceneGallerySelect,
  Modal,
  Icon,
  LoadingIndicator,
  ImageInput
} from "src/components/Shared";
import { useToast } from "src/hooks";
import { ImageUtils, TableUtils } from "src/utils";

interface IProps {
  scene: GQL.SceneDataFragment;
  onUpdate: (scene: GQL.SceneDataFragment) => void;
  onDelete: () => void;
}

export const SceneEditPanel: React.FC<IProps> = (props: IProps) => {
  const Toast = useToast();
  const [title, setTitle] = useState<string>();
  const [details, setDetails] = useState<string>();
  const [url, setUrl] = useState<string>();
  const [date, setDate] = useState<string>();
  const [rating, setRating] = useState<number>();
  const [galleryId, setGalleryId] = useState<string>();
  const [studioId, setStudioId] = useState<string>();
  const [performerIds, setPerformerIds] = useState<string[]>();
  const [tagIds, setTagIds] = useState<string[]>();
  const [coverImage, setCoverImage] = useState<string>();

  const Scrapers = StashService.useListSceneScrapers();
  const [queryableScrapers, setQueryableScrapers] = useState<GQL.Scraper[]>([]);

  const [isDeleteAlertOpen, setIsDeleteAlertOpen] = useState<boolean>(false);
  const [deleteFile, setDeleteFile] = useState<boolean>(false);
  const [deleteGenerated, setDeleteGenerated] = useState<boolean>(true);

  const [coverImagePreview, setCoverImagePreview] = useState<string>();

  // Network state
  const [isLoading, setIsLoading] = useState(true);

  const [updateScene] = StashService.useSceneUpdate(getSceneInput());
  const [deleteScene] = StashService.useSceneDestroy(getSceneDeleteInput());

  useEffect(() => {
    const newQueryableScrapers = (
      Scrapers?.data?.listSceneScrapers ?? []
    ).filter(s => s.scene?.supported_scrapes.includes(GQL.ScrapeType.Fragment));

    setQueryableScrapers(newQueryableScrapers);
  }, [Scrapers]);

  function updateSceneEditState(state: Partial<GQL.SceneDataFragment>) {
    const perfIds = state.performers?.map(performer => performer.id);
    const tIds = state.tags ? state.tags.map(tag => tag.id) : undefined;

    setTitle(state.title ?? undefined);
    setDetails(state.details ?? undefined);
    setUrl(state.url ?? undefined);
    setDate(state.date ?? undefined);
    setRating(state.rating === null ? NaN : state.rating);
    setGalleryId(state?.gallery?.id ?? undefined);
    setStudioId(state?.studio?.id ?? undefined);
    setPerformerIds(perfIds);
    setTagIds(tIds);
  }

  useEffect(() => {
    updateSceneEditState(props.scene);
    setCoverImagePreview(props.scene?.paths?.screenshot ?? undefined);
    setIsLoading(false);
  }, [props.scene]);

  ImageUtils.usePasteImage(onImageLoad);

  function getSceneInput(): GQL.SceneUpdateInput {
    return {
      id: props.scene.id,
      title,
      details,
      url,
      date,
      rating,
      gallery_id: galleryId,
      studio_id: studioId,
      performer_ids: performerIds,
      tag_ids: tagIds,
      cover_image: coverImage
    };
  }

  async function onSave() {
    setIsLoading(true);
    try {
      const result = await updateScene();
      if (result.data?.sceneUpdate) {
        props.onUpdate(result.data.sceneUpdate);
        Toast.success({ content: "Updated scene" });
      }
    } catch (e) {
      Toast.error(e);
    }
    setIsLoading(false);
  }

  function getSceneDeleteInput(): GQL.SceneDestroyInput {
    return {
      id: props.scene.id,
      delete_file: deleteFile,
      delete_generated: deleteGenerated
    };
  }

  async function onDelete() {
    setIsDeleteAlertOpen(false);
    setIsLoading(true);
    try {
      await deleteScene();
      Toast.success({ content: "Deleted scene" });
    } catch (e) {
      Toast.error(e);
    }
    setIsLoading(false);
    props.onDelete();
  }

  function renderDeleteAlert() {
    return (
      <Modal
        show={isDeleteAlertOpen}
        icon="trash-alt"
        header="Delete Scene?"
        accept={{ variant: "danger", onClick: onDelete, text: "Delete" }}
        cancel={{ onClick: () => setIsDeleteAlertOpen(false), text: "Cancel" }}
      >
        <p>
          Are you sure you want to delete this scene? Unless the file is also
          deleted, this scene will be re-added when scan is performed.
        </p>
        <Form>
          <Form.Check
            checked={deleteFile}
            label="Delete file"
            onChange={() => setDeleteFile(!deleteFile)}
          />
          <Form.Check
            checked={deleteGenerated}
            label="Delete generated supporting files"
            onChange={() => setDeleteGenerated(!deleteGenerated)}
          />
        </Form>
      </Modal>
    );
  }

  function onImageLoad(this: FileReader) {
    setCoverImagePreview(this.result as string);
    setCoverImage(this.result as string);
  }

  function onCoverImageChange(event: React.FormEvent<HTMLInputElement>) {
    ImageUtils.onImageChange(event, onImageLoad);
  }

  async function onScrapeClicked(scraper: GQL.Scraper) {
    setIsLoading(true);
    try {
      const result = await StashService.queryScrapeScene(
        scraper.id,
        getSceneInput()
      );
      if (!result.data || !result.data.scrapeScene) {
        return;
      }
      updateSceneFromScrapedScene(result.data.scrapeScene);
    } catch (e) {
      Toast.error(e);
    } finally {
      setIsLoading(false);
    }
  }

  function renderScraperMenu() {
    if (!queryableScrapers || queryableScrapers.length === 0) {
      return;
    }

    return (
      <DropdownButton id="scene-scrape" title="Scrape with...">
        {queryableScrapers.map(s => (
          <Dropdown.Item onClick={() => onScrapeClicked(s)}>
            {s.name}
          </Dropdown.Item>
        ))}
      </DropdownButton>
    );
  }

  function urlScrapable(scrapedUrl: string): boolean {
    return (Scrapers?.data?.listSceneScrapers ?? []).some(s =>
      (s?.scene?.urls ?? []).some(u => scrapedUrl.includes(u))
    );
  }

  function updateSceneFromScrapedScene(scene: GQL.ScrapedSceneDataFragment) {
    if (!title && scene.title) {
      setTitle(scene.title);
    }

    if (!details && scene.details) {
      setDetails(scene.details);
    }

    if (!date && scene.date) {
      setDate(scene.date);
    }

    if (!url && scene.url) {
      setUrl(scene.url);
    }

    if (!studioId && scene.studio && scene.studio.id) {
      setStudioId(scene.studio.id);
    }

    if (
      (!performerIds || performerIds.length === 0) &&
      scene.performers &&
      scene.performers.length > 0
    ) {
      const idPerfs = scene.performers.filter(p => {
        return p.id !== undefined && p.id !== null;
      });

      if (idPerfs.length > 0) {
        const newIds = idPerfs.map(p => p.id);
        setPerformerIds(newIds as string[]);
      }
    }

    if (!tagIds?.length && scene?.tags?.length) {
      const idTags = scene.tags.filter(p => {
        return p.id !== undefined && p.id !== null;
      });

      if (idTags.length > 0) {
        const newIds = idTags.map(p => p.id);
        setTagIds(newIds as string[]);
      }
    }
  }

  async function onScrapeSceneURL() {
    if (!url) {
      return;
    }
    setIsLoading(true);
    try {
      const result = await StashService.queryScrapeSceneURL(url);
      if (!result.data || !result.data.scrapeSceneURL) {
        return;
      }
      updateSceneFromScrapedScene(result.data.scrapeSceneURL);
    } catch (e) {
      Toast.error(e);
    } finally {
      setIsLoading(false);
    }
  }

  function maybeRenderScrapeButton() {
    if (!url || !urlScrapable(url)) {
      return undefined;
    }
    return (
      <Button id="scrape-url-button" onClick={onScrapeSceneURL}>
        <Icon icon="file-download" />
      </Button>
    );
  }

  if (isLoading) return <LoadingIndicator />;

  return (
    <div className="form-container row">
      <div className="col-12 col-lg-6">
        <Table id="scene-edit-details">
          <tbody>
            {TableUtils.renderInputGroup({
              title: "Title",
              value: title,
              onChange: setTitle,
              isEditing: true
            })}
            <tr>
              <td>URL</td>
              <td>
                <Form.Control
                  onChange={(newValue: React.FormEvent<HTMLInputElement>) =>
                    setUrl(newValue.currentTarget.value)
                  }
                  value={url}
                  placeholder="URL"
                  className="text-input"
                />
                {maybeRenderScrapeButton()}
              </td>
            </tr>
            {TableUtils.renderInputGroup({
              title: "Date",
              value: date,
              isEditing: true,
              onChange: setDate,
              placeholder: "YYYY-MM-DD"
            })}
            {TableUtils.renderHtmlSelect({
              title: "Rating",
              value: rating,
              isEditing: true,
              onChange: (value: string) =>
                setRating(Number.parseInt(value, 10)),
              selectOptions: ["", 1, 2, 3, 4, 5]
            })}
            <tr>
              <td>Gallery</td>
              <td>
                <SceneGallerySelect
                  sceneId={props.scene.id}
                  initialId={galleryId}
                  onSelect={item => setGalleryId(item ? item.id : undefined)}
                />
              </td>
            </tr>
            <tr>
              <td>Studio</td>
              <td>
                <StudioSelect
                  onSelect={items => items.length && setStudioId(items[0]?.id)}
                  ids={studioId ? [studioId] : []}
                />
              </td>
            </tr>
            <tr>
              <td>Performers</td>
              <td>
                <PerformerSelect
                  isMulti
                  onSelect={items =>
                    setPerformerIds(items.map(item => item.id))
                  }
                  ids={performerIds}
                />
              </td>
            </tr>
            <tr>
              <td>Tags</td>
              <td>
                <TagSelect
                  isMulti
                  onSelect={items => setTagIds(items.map(item => item.id))}
                  ids={tagIds}
                />
              </td>
            </tr>
          </tbody>
        </Table>
      </div>
      <div className="col-12 col-lg-6">
        <Form.Group controlId="details">
          <Form.Label>Details</Form.Label>
          <Form.Control
            as="textarea"
            className="scene-description text-input"
            onChange={(newValue: React.FormEvent<HTMLTextAreaElement>) =>
              setDetails(newValue.currentTarget.value)
            }
            value={details}
          />
        </Form.Group>

        <div>
          <Form.Group className="test" controlId="cover">
            <Form.Label>Cover Image</Form.Label>
            <img
              className="scene-cover"
              src={coverImagePreview}
              alt="Scene cover"
            />
            <ImageInput isEditing onImageChange={onCoverImageChange} />
          </Form.Group>
        </div>
      </div>
      <div className="col edit-buttons">
        <Button className="edit-button" variant="primary" onClick={onSave}>
          Save
        </Button>
        <Button
          className="edit-button"
          variant="danger"
          onClick={() => setIsDeleteAlertOpen(true)}
        >
          Delete
        </Button>
      </div>
      {renderScraperMenu()}
      {renderDeleteAlert()}
    </div>
  );
};
