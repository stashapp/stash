import React, { useEffect, useState } from "react";
import * as GQL from "../../../core/generated-graphql";
import { StashService } from "../../../core/StashService";
import { ErrorUtils } from "../../../utils/errors";
import { ToastUtils } from "../../../utils/toasts";
import { FilterMultiSelect } from "../../select/FilterMultiSelect";
import { FilterSelect } from "../../select/FilterSelect";
import { ValidGalleriesSelect } from "../../select/ValidGalleriesSelect";
import { ImageUtils } from "../../../utils/image";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { Collapse, Dropdown, DropdownButton, Form, Button, Modal, Spinner } from 'react-bootstrap';

interface IProps {
  scene: GQL.SceneDataFragment;
  onUpdate: (scene: GQL.SceneDataFragment) => void;
  onDelete: () => void;
}

export const SceneEditPanel: React.FC<IProps> = (props: IProps) => {
  const [title, setTitle] = useState<string | undefined>(undefined);
  const [details, setDetails] = useState<string | undefined>(undefined);
  const [url, setUrl] = useState<string | undefined>(undefined);
  const [date, setDate] = useState<string | undefined>(undefined);
  const [rating, setRating] = useState<number | undefined>(undefined);
  const [galleryId, setGalleryId] = useState<string | undefined>(undefined);
  const [studioId, setStudioId] = useState<string | undefined>(undefined);
  const [performerIds, setPerformerIds] = useState<string[] | undefined>(undefined);
  const [tagIds, setTagIds] = useState<string[] | undefined>(undefined);
  const [coverImage, setCoverImage] = useState<string | undefined>(undefined);

  const Scrapers = StashService.useListSceneScrapers();
  const [queryableScrapers, setQueryableScrapers] = useState<GQL.ListSceneScrapersListSceneScrapers[]>([]);

  const [isDeleteAlertOpen, setIsDeleteAlertOpen] = useState<boolean>(false);
  const [deleteFile, setDeleteFile] = useState<boolean>(false);
  const [deleteGenerated, setDeleteGenerated] = useState<boolean>(true);

  const [isCoverImageOpen, setIsCoverImageOpen] = useState<boolean>(false);
  const [coverImagePreview, setCoverImagePreview] = useState<string | undefined>(undefined);

  // Network state
  const [isLoading, setIsLoading] = useState(false);

  const updateScene = StashService.useSceneUpdate(getSceneInput());
  const deleteScene = StashService.useSceneDestroy(getSceneDeleteInput());

  useEffect(() => {
    var newQueryableScrapers : GQL.ListSceneScrapersListSceneScrapers[] = [];

    if (!!Scrapers.data && Scrapers.data.listSceneScrapers) {
      newQueryableScrapers = Scrapers.data.listSceneScrapers.filter((s) => {
        return s.scene && s.scene.supported_scrapes.includes(GQL.ScrapeType.Fragment);
      });
    }

    setQueryableScrapers(newQueryableScrapers);

  }, [Scrapers.data])

  function updateSceneEditState(state: Partial<GQL.SceneDataFragment>) {
    const perfIds = !!state.performers ? state.performers.map((performer) => performer.id) : undefined;
    const tIds = !!state.tags ? state.tags.map((tag) => tag.id) : undefined;

    setTitle(state.title);
    setDetails(state.details);
    setUrl(state.url);
    setDate(state.date);
    setRating(state.rating === null ? NaN : state.rating);
    setGalleryId(state.gallery ? state.gallery.id : undefined);
    setStudioId(state.studio ? state.studio.id : undefined);
    setPerformerIds(perfIds);
    setTagIds(tIds);
  }

  useEffect(() => {
    updateSceneEditState(props.scene);
    setCoverImagePreview(props.scene.paths.screenshot);
  }, [props.scene]);

  ImageUtils.addPasteImageHook(onImageLoad);

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
      cover_image: coverImage,
    };
  }

  async function onSave() {
    setIsLoading(true);
    try {
      const result = await updateScene();
      props.onUpdate(result.data.sceneUpdate);
      ToastUtils.success("Updated scene");
    } catch (e) {
      ErrorUtils.handle(e);
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
      ToastUtils.success("Deleted scene");
    } catch (e) {
      ErrorUtils.handle(e);
    }
    setIsLoading(false);

    props.onDelete();
  }

  function renderMultiSelect(type: "performers" | "tags", initialIds: string[] | undefined) {
    return (
      <FilterMultiSelect
        type={type}
        onUpdate={(items) => {
          const ids = items.map((i) => i.id);
          switch (type) {
            case "performers": setPerformerIds(ids); break;
            case "tags": setTagIds(ids); break;
          }
        }}
        initialIds={initialIds}
      />
    );
  }

  function renderDeleteAlert() {
    return (
      <Modal
        keyboard={false}
        onHide={() => {}}
        show={isDeleteAlertOpen}
      >
        <Modal.Header>
            <FontAwesomeIcon icon="trash-alt" />
            <span>Delete Scene?</span>
        </Modal.Header>
        <Modal.Body>
          <p>
            Are you sure you want to delete this scene? Unless the file is also deleted, this scene will be re-added when scan is performed.
          </p>
          <Form>
            <Form.Check checked={deleteFile} label="Delete file" onChange={() => setDeleteFile(!deleteFile)} />
            <Form.Check checked={deleteGenerated} label="Delete generated supporting files" onChange={() => setDeleteGenerated(!deleteGenerated)} />
          </Form>
        </Modal.Body>

        <Modal.Footer>
          <div>
            <Button variant="danger" onClick={onDelete}>Delete</Button>
            <Button onClick={() => setIsDeleteAlertOpen(false)}>Cancel</Button>
          </div>
        </Modal.Footer>
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
  
  async function onScrapeClicked(scraper : GQL.ListSceneScrapersListSceneScrapers) {
    setIsLoading(true);
    try {
      const result = await StashService.queryScrapeScene(scraper.id, getSceneInput());
      if (!result.data || !result.data.scrapeScene) { return; }
      updateSceneFromScrapedScene(result.data.scrapeScene);
    } catch (e) {
      ErrorUtils.handle(e);
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
        { queryableScrapers.map(s => (
            <Dropdown.Item onClick={() => onScrapeClicked(s)}>{s.name}</Dropdown.Item>
          ))
        }
      </DropdownButton>
    );
  }

  function urlScrapable(url: string) : boolean {
    return !!url && !!Scrapers.data && Scrapers.data.listSceneScrapers && Scrapers.data.listSceneScrapers.some((s) => {
      return !!s.scene && !!s.scene.urls && s.scene.urls.some((u) => { return url.includes(u); });
    });
  }

  function updateSceneFromScrapedScene(scene : GQL.ScrapedSceneDataFragment) {
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

    if ((!performerIds || performerIds.length === 0) && scene.performers && scene.performers.length > 0) {
      let idPerfs = scene.performers.filter((p) => {
        return p.id !== undefined && p.id !== null;
      });

      if (idPerfs.length > 0) {
        let newIds = idPerfs.map((p) => p.id);
        setPerformerIds(newIds as string[]);
      }
    }

    if ((!tagIds || tagIds.length === 0) && scene.tags && scene.tags.length > 0) {
      let idTags = scene.tags.filter((p) => {
        return p.id !== undefined && p.id !== null;
      });

      if (idTags.length > 0) {
        let newIds = idTags.map((p) => p.id);
        setTagIds(newIds as string[]);
      }
    }
  }

  async function onScrapeSceneURL() {
    if (!url) { return; }
    setIsLoading(true);
    try {
      const result = await StashService.queryScrapeSceneURL(url);
      if (!result.data || !result.data.scrapeSceneURL) { return; }
      updateSceneFromScrapedScene(result.data.scrapeSceneURL);
    } catch (e) {
      ErrorUtils.handle(e);
    } finally {
      setIsLoading(false);
    }
  }

  function maybeRenderScrapeButton() {
    if (!url || !urlScrapable(url)) {
      return undefined;
    }
    return (
      <Button 
        id="scrape-url-button"
        onClick={onScrapeSceneURL}>
        <FontAwesomeIcon icon="file-download" />
      </Button>
    )
  }

  return (
    <>
      {renderDeleteAlert()}
      {isLoading ? <Spinner animation="border" variant="light" /> : undefined}
      <div className="form-container " style={{width: "50%"}}>
        <Form.Group controlId="title">
          <Form.Label>Title</Form.Label>
          <Form.Control
            onChange={(newValue: any) => setTitle(newValue.target.value)}
            value={title}
          />
        </Form.Group>

        <Form.Group controlId="details">
          <Form.Label>Details</Form.Label>
            <Form.Control
            as="textarea"
            onChange={(newValue: any) => setDetails(newValue.target.value)}
            value={details}
          />
        </Form.Group>

        <Form.Group controlId="url">
          <Form.Label>URL</Form.Label>
          <Form.Control
            onChange={(newValue: any) => setUrl(newValue.target.value)}
            value={url}
          />
          {maybeRenderScrapeButton()}
        </Form.Group>

        <Form.Group controlId="date">
          <Form.Label>Date</Form.Label>
          <Form.Control
            onChange={(newValue: any) => setDate(newValue.target.value)}
            value={date}
          />
          <div>YYYY-MM-DD</div>
        </Form.Group>

        <Form.Group controlId="rating">
          <Form.Label>Rating</Form.Label>
          <Form.Control
            as="select"
            onChange={(event: any) => setRating(parseInt(event.target.value, 10))}>
              { ["", 1, 2, 3, 4, 5].map(opt => (
                  <option selected={opt === rating} value={opt}>{opt}</option>
              )) }
          </Form.Control>
        </Form.Group>

        <Form.Group controlId="gallery">
          <Form.Label>Gallery</Form.Label>
          <ValidGalleriesSelect
            sceneId={props.scene.id}
            initialId={galleryId}
            onSelectItem={(item) => setGalleryId(item ? item.id : undefined)}
          />
        </Form.Group>

        <Form.Group controlId="studio">
          <Form.Label>Studio</Form.Label>
          <FilterSelect
            type="studios"
            onSelectItem={(item) => setStudioId(item ? item.id : undefined)}
            initialId={studioId}
          />
        </Form.Group>

        <Form.Group controlId="performers">
          <Form.Label>Performers</Form.Label>
          {renderMultiSelect("performers", performerIds)}
        </Form.Group>

        <Form.Group controlId="tags">
          <Form.Label>Tags</Form.Label>
          {renderMultiSelect("tags", tagIds)}
        </Form.Group>

        <div>
          <label onClick={() => setIsCoverImageOpen(!isCoverImageOpen)}>
            <FontAwesomeIcon icon={isCoverImageOpen ? "chevron-down" : "chevron-right"} />
            <span>Cover Image</span>
          </label>
          <Collapse in={isCoverImageOpen}>
            <div>
              <img className="scene-cover" src={coverImagePreview} alt="" />
              <Form.Group className="test" controlId="cover">
                <Form.Control type="file" onChange={onCoverImageChange} accept=".jpg,.jpeg,.png" />
              </Form.Group>
            </div>
          </Collapse>
        </div>
        
      </div>
      <Button className="edit-button" variant="primary" onClick={onSave}>Save</Button>
      <Button className="edit-button" variant="danger" onClick={() => setIsDeleteAlertOpen(true)}>Delete</Button>
      {renderScraperMenu()}
    </>
  );
};
