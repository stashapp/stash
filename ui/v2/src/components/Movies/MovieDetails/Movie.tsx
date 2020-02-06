import {
  EditableText,
  HTMLTable,
  Spinner,
} from "@blueprintjs/core";
import React, { FunctionComponent, useEffect, useState } from "react";
import * as GQL from "../../../core/generated-graphql";
import { StashService } from "../../../core/StashService";
import { IBaseProps } from "../../../models";
import { ErrorUtils } from "../../../utils/errors";
import { TableUtils } from "../../../utils/table";
import { DetailsEditNavbar } from "../../Shared/DetailsEditNavbar";
import { ImageUtils } from "../../../utils/image";

interface IProps extends IBaseProps {}

export const Movie: FunctionComponent<IProps> = (props: IProps) => {
  const isNew = props.match.params.id === "new";

  // Editing state
  const [isEditing, setIsEditing] = useState<boolean>(isNew);

  // Editing movie state
  const [front_image, setFront_Image] = useState<string | undefined>(undefined);
  const [back_image, setBack_Image] = useState<string | undefined>(undefined);
  const [name, setName] = useState<string | undefined>(undefined);
  const [aliases, setAliases] = useState<string | undefined>(undefined);
  const [duration_movie, setDuration_movie] = useState<string | undefined>(undefined);
  const [date_movie, setDate_movie] = useState<string | undefined>(undefined);
  const [rating_movie, setRating] = useState<string | undefined>(undefined);
  const [director, setDirector] = useState<string | undefined>(undefined);
  const [synopsis, setSynopsis] = useState<string | undefined>(undefined);
  const [url, setUrl] = useState<string | undefined>(undefined);

  // Movie state
  const [movie, setMovie] = useState<Partial<GQL.MovieDataFragment>>({});
  const [imagePreview, setImagePreview] = useState<string | undefined>(undefined);
  const [backimagePreview, setBackImagePreview] = useState<string | undefined>(undefined);

  // Network state
  const [isLoading, setIsLoading] = useState(false);

  const { data, error, loading } = StashService.useFindMovie(props.match.params.id);
  const updateMovie = StashService.useMovieUpdate(getMovieInput() as GQL.MovieUpdateInput);
  const createMovie = StashService.useMovieCreate(getMovieInput() as GQL.MovieCreateInput);
  const deleteMovie = StashService.useMovieDestroy(getMovieInput() as GQL.MovieDestroyInput);

  function updateMovieEditState(state: Partial<GQL.MovieDataFragment>) {
    setName(state.name);
    setAliases(state.aliases);
    setDuration_movie(state.duration_movie);
    setDate_movie(state.date_movie);
    setRating(state.rating_movie);
    setDirector(state.director);
    setSynopsis(state.synopsis);
    setUrl(state.url);
  }

  useEffect(() => {
    setIsLoading(loading);
    if (!data || !data.findMovie || !!error) { return; }
    setMovie(data.findMovie);
  }, [data, loading, error]);

  useEffect(() => {
    setImagePreview(movie.front_image_path);
    setBackImagePreview(movie.back_image_path);
    setFront_Image(undefined);
    setBack_Image(undefined);
    updateMovieEditState(movie);
    if (!isNew) {
      setIsEditing(false);
    }
  }, [movie, isNew]);

  function onImageLoad(this: FileReader) {
    setImagePreview(this.result as string);
    setFront_Image(this.result as string);
    
  }

  function onBackImageLoad(this: FileReader) {
    setBackImagePreview(this.result as string);
    setBack_Image(this.result as string);
  }


  ImageUtils.addPasteImageHook(onImageLoad);
  ImageUtils.addPasteImageHook(onBackImageLoad);

  if (!isNew && !isEditing) {
    if (!data || !data.findMovie || isLoading) { return <Spinner size={Spinner.SIZE_LARGE} />; }
    if (!!error) { return <>error...</>; }
  }

  function getMovieInput() {
    const input: Partial<GQL.MovieCreateInput | GQL.MovieUpdateInput> = {
      name,
      aliases,
      duration_movie,
      date_movie,
	    rating_movie,
      director,
      synopsis,
      url,
      front_image,
      back_image
      
    };

    if (!isNew) {
      (input as GQL.MovieUpdateInput).id = props.match.params.id;
    }
    return input;
  }

  async function onSave() {
    setIsLoading(true);
    try {
      if (!isNew) {
        const result = await updateMovie();
        setMovie(result.data.movieUpdate);
      } else {
        const result = await createMovie();
        setMovie(result.data.movieCreate);
        props.history.push(`/movies/${result.data.movieCreate.id}`);
      }
    } catch (e) {
      ErrorUtils.handle(e);
    }
    setIsLoading(false);
  }

  async function onDelete() {
    setIsLoading(true);
    try {
      await deleteMovie();
    } catch (e) {
      ErrorUtils.handle(e);
    }
    setIsLoading(false);
    
    // redirect to movies page
    props.history.push(`/movies`);
  }

  function onImageChange(event: React.FormEvent<HTMLInputElement>) {
    ImageUtils.onImageChange(event, onImageLoad);
  }

  function onBackImageChange(event: React.FormEvent<HTMLInputElement>) {
    ImageUtils.onImageChange(event, onBackImageLoad);
  }

  // TODO: CSS class
  return (
    <>
      <div className="columns is-multiline no-spacing">
        <div className="column is-half details-image-container">
          <img alt={name} className="movie" src={imagePreview} />
          <img alt={name} className="movie" src={backimagePreview} />
       </div>
        <div className="column is-half details-detail-container">
          <DetailsEditNavbar
            movie={movie}
            isNew={isNew}
            isEditing={isEditing}
            onToggleEdit={() => { setIsEditing(!isEditing); updateMovieEditState(movie); }}
            onSave={onSave}
            onDelete={onDelete}
            onImageChange={onImageChange}
            onBackImageChange={onBackImageChange}
            />
          <h1 className="bp3-heading">
            <EditableText
              disabled={!isEditing}
              value={name}
              placeholder="Name"
              onChange={(value) => setName(value)}
            />
          </h1>

          <HTMLTable style={{width: "100%"}}>
            <tbody>
              {TableUtils.renderInputGroup({title: "Aliases", value: aliases, isEditing, onChange: setAliases})}
              {TableUtils.renderInputGroup({title: "Duration", value: duration_movie, isEditing, onChange: setDuration_movie})}
              {TableUtils.renderInputGroup({title: "Date (YYYY-MM-DD)", value: date_movie, isEditing, onChange: setDate_movie})}
              {TableUtils.renderInputGroup({title: "Director", value: director, isEditing, onChange: setDirector})}
              {TableUtils.renderHtmlSelect({
                title: "Rating", 
                value: rating_movie, 
                isEditing, 
                onChange: (value: string) => setRating(value), 
                selectOptions: ["","1","2","3","4","5"]
                })}
              {TableUtils.renderInputGroup({title: "URL", value: url, isEditing, onChange: setUrl})}
              {TableUtils.renderTextArea({title: "Synopsis", value: synopsis, isEditing, onChange: setSynopsis})}            
            
            </tbody>
          </HTMLTable>
        </div>
      </div>
    </>
  );
};
