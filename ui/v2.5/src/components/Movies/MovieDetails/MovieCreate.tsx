import React, { useState } from "react";
import * as GQL from "src/core/generated-graphql";
import { useMovieCreate } from "src/core/StashService";
import { useHistory } from "react-router-dom";
import { LoadingIndicator } from "src/components/Shared";
import { useToast } from "src/hooks";
import { MovieEditPanel } from "./MovieEditPanel";

const MovieCreate: React.FC = () => {
  const history = useHistory();
  const Toast = useToast();

  // Editing movie state
  const [frontImage, setFrontImage] = useState<string | undefined | null>(
    undefined
  );
  const [backImage, setBackImage] = useState<string | undefined | null>(
    undefined
  );
  const [encodingImage, setEncodingImage] = useState<boolean>(false);

  const [createMovie] = useMovieCreate();

  const onImageEncoding = (isEncoding = false) => setEncodingImage(isEncoding);

  function getMovieInput(
    input: Partial<GQL.MovieCreateInput | GQL.MovieUpdateInput>
  ) {
    const ret: Partial<GQL.MovieCreateInput | GQL.MovieUpdateInput> = {
      ...input,
    };

    return ret;
  }

  async function onSave(
    input: Partial<GQL.MovieCreateInput | GQL.MovieUpdateInput>
  ) {
    try {
      const result = await createMovie({
        variables: getMovieInput(input) as GQL.MovieCreateInput,
      });
      if (result.data?.movieCreate?.id) {
        history.push(`/movies/${result.data.movieCreate.id}`);
      }
    } catch (e) {
      Toast.error(e);
    }
  }

  function renderFrontImage() {
    if (frontImage) {
      return (
        <div className="movie-image-container">
          <img alt="Front Cover" src={frontImage} />
        </div>
      );
    }
  }

  function renderBackImage() {
    if (backImage) {
      return (
        <div className="movie-image-container">
          <img alt="Back Cover" src={backImage} />
        </div>
      );
    }
  }

  // TODO: CSS class
  return (
    <div className="row">
      <div className="movie-details mb-3 col">
        <div className="logo w-100">
          {encodingImage ? (
            <LoadingIndicator message="Encoding image..." />
          ) : (
            <div className="movie-images">
              {renderFrontImage()}
              {renderBackImage()}
            </div>
          )}
        </div>

        <MovieEditPanel
          onSubmit={onSave}
          onCancel={() => history.push("/movies")}
          onDelete={() => {}}
          setFrontImage={setFrontImage}
          setBackImage={setBackImage}
          onImageEncoding={onImageEncoding}
        />
      </div>
    </div>
  );
};

export default MovieCreate;
