import {
  Button,
  Classes,
  Checkbox,
  Dialog,
  FormGroup,
  HTMLSelect,
  InputGroup,
  Spinner,
  TextArea,
  Collapse,
  Icon,
  FileInput,
  Menu,
  Popover,
  MenuItem,
} from "@blueprintjs/core";
import React, { FunctionComponent, useEffect, useState } from "react";
import * as GQL from "../../../core/generated-graphql";
import { StashService } from "../../../core/StashService";
import { ErrorUtils } from "../../../utils/errors";
import { ToastUtils } from "../../../utils/toasts";
import { FilterMultiSelect } from "../../select/FilterMultiSelect";
import { FilterSelect } from "../../select/FilterSelect";
import { ValidGalleriesSelect } from "../../select/ValidGalleriesSelect";
import { ImageUtils } from "../../../utils/image";
import { SceneMovieTable }  from "./SceneMovieTable";

interface IProps {
  scene: GQL.SceneDataFragment;
  onUpdate: (scene: GQL.SceneDataFragment) => void;
  onDelete: () => void;
}

export const SceneEditPanel: FunctionComponent<IProps> = (props: IProps) => {
  // Editing scene state
  const [title, setTitle] = useState<string | undefined>(undefined);
  const [details, setDetails] = useState<string | undefined>(undefined);
  const [url, setUrl] = useState<string | undefined>(undefined);
  const [date, setDate] = useState<string | undefined>(undefined);
  const [rating, setRating] = useState<number | undefined>(undefined);
  const [galleryId, setGalleryId] = useState<string | undefined>(undefined);
  const [studioId, setStudioId] = useState<string | undefined>(undefined);
  const [performerIds, setPerformerIds] = useState<string[] | undefined>(undefined);
  const [movieIds, setMovieIds] = useState<string[] | undefined>(undefined); 
  const [sceneIdx, setSceneIdx] = useState<string[] | undefined>(undefined); 
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
    const moviIds = !!state.movies ? state.movies.map((sceneMovie) => sceneMovie.movie.id) : undefined;
    const scenIdx = !!state.movies ? state.movies.map((movie) => movie.scene_index!) : undefined; 
  
    const tIds = !!state.tags ? state.tags.map((tag) => tag.id) : undefined;
       
    setTitle(state.title);
    setDetails(state.details); 
    setUrl(state.url);
    setDate(state.date);
    setRating(state.rating == null ? NaN : state.rating);
    setGalleryId(state.gallery ? state.gallery.id : undefined);
    setStudioId(state.studio ? state.studio.id : undefined);
    setMovieIds(moviIds);
    setPerformerIds(perfIds);
    setSceneIdx(scenIdx);
    setTagIds(tIds);
  }

  useEffect(() => {
    updateSceneEditState(props.scene);
    setCoverImagePreview(props.scene.paths.screenshot);
  }, [props.scene]);

  ImageUtils.addPasteImageHook(onImageLoad);

  // if (!isNew && !isEditing) {
  //   if (!data || !data.findPerformer || isLoading) { return <Spinner size={Spinner.SIZE_LARGE} />; }
  //   if (!!error) { return <>error...</>; }
  // }

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
      movies: makeMovieInputs(),
      tag_ids: tagIds,
      cover_image: coverImage,
    };
  }

  function makeMovieInputs(): GQL.SceneMovieInput[] | undefined {
    if (!movieIds) {
      return undefined;
    }
    
    let ret = movieIds.map((id) => {
      let r : GQL.SceneMovieInput = {
        movie_id: id
      };
      return r;
    });

    if (sceneIdx) {
      sceneIdx.forEach((idx, i) => {
        if (!!idx && ret.length > i) {
          ret[i].scene_index = idx;
        }
      });
    }

    return ret;
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

  function renderMultiSelect(type: "performers" | "movies" | "tags", initialIds: string[] | undefined) {
    return (
      <FilterMultiSelect
        type={type}
        onUpdate={(items) => {
          const ids = items.map((i) => i.id);
          switch (type) {
            case "performers": setPerformerIds(ids); break;
            case "movies": setMovieIds(ids); break;
            case "tags": setTagIds(ids); break;
          }
        }}
        initialIds={initialIds}
      />
    );
  }

  function renderTableMovies( initialIds: string[] | undefined, initialIdx: string[] | undefined ) {
    return (
      <SceneMovieTable
          initialIds={initialIds}
          initialIdx={initialIdx}
          onUpdate={(items) => {
            const idx = items.map((i) => i);
            setSceneIdx(idx);
          }}
      />
    );
  }

  function renderDeleteAlert() {
    return (
      <>
      <Dialog
        canOutsideClickClose={false}
        canEscapeKeyClose={false}
        icon="trash"
        isCloseButtonShown={false}
        isOpen={isDeleteAlertOpen}
        title="Delete Scene?"
      >
        <div className={Classes.DIALOG_BODY}>
          <p>
            Are you sure you want to delete this scene? Unless the file is also deleted, this scene will be re-added when scan is performed.
          </p>
          <Checkbox checked={deleteFile} label="Delete file" onChange={() => setDeleteFile(!deleteFile)} />
          <Checkbox checked={deleteGenerated} label="Delete generated supporting files" onChange={() => setDeleteGenerated(!deleteGenerated)} />
        </div>

        <div className={Classes.DIALOG_FOOTER}>
          <div className={Classes.DIALOG_FOOTER_ACTIONS}>
            <Button intent="danger" onClick={() => onDelete()}>Delete</Button>
            <Button onClick={() => setIsDeleteAlertOpen(false)}>Cancel</Button>
          </div>
        </div>
      </Dialog>
      </>
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

  function renderScraperMenuItem(scraper : GQL.ListSceneScrapersListSceneScrapers) {
    return (
      <MenuItem
        text={scraper.name}
        onClick={() => { onScrapeClicked(scraper); }}
      />
    );
  }

  function renderScraperMenu() {
    if (!queryableScrapers || queryableScrapers.length === 0) {
      return;
    }

    const scraperMenu = (
      <Menu>
        {queryableScrapers ? queryableScrapers.map((s) => renderScraperMenuItem(s)) : undefined}
      </Menu>
    );
    return (
      <Popover content={scraperMenu} position="bottom">
        <Button text="Scrape with..."/>
      </Popover>
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
    
    if ((!movieIds || movieIds.length === 0) && scene.movies && scene.movies.length > 0) {
      let idMovis = scene.movies.filter((p) => {
        return p.id !== undefined && p.id !== null;
      });

      if (idMovis.length > 0) {
        let newIds = idMovis.map((p) => p.id);
        setMovieIds(newIds as string[]);
      }
    }

    if ((!sceneIdx || sceneIdx.length === 0) && scene.movies && scene.movies.length > 0) {
      let idxScen= scene.movies.filter((p) => {
        return p.id !== undefined && p.id !== null;
      });

      if (idxScen.length > 0) {
        let newIds = idxScen.map((p) => p.id);
        setSceneIdx(newIds as string[]);
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
        minimal={true} 
        icon="import" 
        id="scrape-url-button"
        onClick={() => onScrapeSceneURL()}/>
    )
  }

  return (
    <>
      {renderDeleteAlert()}
      {isLoading ? <Spinner size={Spinner.SIZE_LARGE} /> : undefined}
      <div className="form-container " style={{width: "50%"}}>
        <FormGroup label="Title">
          <InputGroup
            onChange={(newValue: any) => setTitle(newValue.target.value)}
            value={title}
          />
        </FormGroup>

        <FormGroup label="Details">
          <TextArea
            fill={true}
            onChange={(newValue) => setDetails(newValue.target.value)}
            value={details}
          />
        </FormGroup>

        <FormGroup label="URL">
          <InputGroup
            onChange={(newValue: any) => setUrl(newValue.target.value)}
            value={url}
          />
          {maybeRenderScrapeButton()}
        </FormGroup>

        <FormGroup label="Date" helperText="YYYY-MM-DD">
          <InputGroup
            onChange={(newValue: any) => setDate(newValue.target.value)}
            value={date}
          />
        </FormGroup>

        <FormGroup label="Rating">
          <HTMLSelect
            options={["", 1, 2, 3, 4, 5]}
            onChange={(event) => setRating(parseInt(event.target.value, 10))}
            value={rating}
          />
        </FormGroup>

        <FormGroup label="Gallery">
          <ValidGalleriesSelect
            sceneId={props.scene.id}
            initialId={galleryId}
            onSelectItem={(item) => setGalleryId(item ? item.id : undefined)}
          />
        </FormGroup>

        <FormGroup label="Studio">
          <FilterSelect
            type="studios"
            onSelectItem={(item) => setStudioId(item ? item.id : undefined)}
            initialId={studioId}
          />
        </FormGroup>

        <FormGroup label="Performers">
          {renderMultiSelect("performers", performerIds)}
        </FormGroup>
        
        <FormGroup label="Movies/Scenes">
          {renderMultiSelect("movies", movieIds)}
          {renderTableMovies(movieIds, sceneIdx)}
        </FormGroup>
 
        <FormGroup label="Tags">
          {renderMultiSelect("tags", tagIds)}
        </FormGroup>

        <div className="bp3-form-group">
          <label className="bp3-label collapsible-label" onClick={() => setIsCoverImageOpen(!isCoverImageOpen)}>
            <Icon className="label-icon" icon={isCoverImageOpen ? "chevron-down" : "chevron-right"}/>
            <span>Cover Image</span>
          </label>
          <Collapse isOpen={isCoverImageOpen}>
            <img alt="Scene cover" className="scene-cover" src={coverImagePreview} />
            <FileInput text="Choose image..." onInputChange={onCoverImageChange} inputProps={{accept: ".jpg,.jpeg,.png"}} />
          </Collapse>
        </div>
        
      </div>
      <Button className="edit-button" text="Save" intent="primary" onClick={() => onSave()}/>
      <Button className="edit-button" text="Delete" intent="danger" onClick={() => setIsDeleteAlertOpen(true)}/>
      {renderScraperMenu()}
    </>
  );
};
