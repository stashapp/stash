import React, { useEffect, useState, useCallback } from "react";
import * as GQL from "src/core/generated-graphql";
import Mousetrap from "mousetrap";
import {
  queryScrapeMovieURL,
  useListMovieScrapers,
} from "src/core/StashService";
import {
  LoadingIndicator,
  StudioSelect,
  Icon,
  DetailsEditNavbar,
} from "src/components/Shared";
import { useToast } from "src/hooks";
import { Table, Form, Modal as BSModal, Button } from "react-bootstrap";
import {
  TableUtils,
  EditableTextUtils,
  TextUtils,
  DurationUtils,
  ImageUtils,
} from "src/utils";
import { RatingStars } from "src/components/Scenes/SceneDetails/RatingStars";
import { MovieScrapeDialog } from "./MovieScrapeDialog";

interface IMovieEditPanel {
  movie?: Partial<GQL.MovieDataFragment>;
  onSubmit: (movie: Partial<GQL.MovieCreateInput | GQL.MovieUpdateInput>) => void;
  onCancel: () => void;
  onDelete: () => void;
  setFrontImage: (image?: string | null) => void;
  setBackImage: (image?: string | null) => void;
}

export const MovieEditPanel: React.FC<IMovieEditPanel> = ({movie, onSubmit, onCancel, onDelete, setFrontImage, setBackImage}) => {
  const Toast = useToast();

  const isNew = movie === undefined;

  const [isLoading, setIsLoading] = useState(false);

  // Editing movie state
  const [name, setName] = useState<string | undefined>(undefined);
  const [aliases, setAliases] = useState<string | undefined>(undefined);
  const [duration, setDuration] = useState<number | undefined>(undefined);
  const [date, setDate] = useState<string | undefined>(undefined);
  const [rating, setRating] = useState<number | undefined>(undefined);
  const [studioId, setStudioId] = useState<string>();
  const [director, setDirector] = useState<string | undefined>(undefined);
  const [synopsis, setSynopsis] = useState<string | undefined>(undefined);
  const [url, setUrl] = useState<string | undefined>(undefined);

  const Scrapers = useListMovieScrapers();
  const [scrapedMovie, setScrapedMovie] = useState<
    GQL.ScrapedMovie | undefined
  >();

  // set up hotkeys
  useEffect(() => {
    Mousetrap.bind("r 0", () => setRating(NaN));
    Mousetrap.bind("r 1", () => setRating(1));
    Mousetrap.bind("r 2", () => setRating(2));
    Mousetrap.bind("r 3", () => setRating(3));
    Mousetrap.bind("r 4", () => setRating(4));
    Mousetrap.bind("r 5", () => setRating(5));
    // Mousetrap.bind("u", (e) => {
    //   setStudioFocus()
    //   e.preventDefault();
    // });
    Mousetrap.bind("s s", () => onSubmit(getMovieInput()));

    return () => {
      Mousetrap.unbind("r 0");
      Mousetrap.unbind("r 1");
      Mousetrap.unbind("r 2");
      Mousetrap.unbind("r 3");
      Mousetrap.unbind("r 4");
      Mousetrap.unbind("r 5");
      // Mousetrap.unbind("u");
      Mousetrap.unbind("s s");
    };
  });

  function updateMovieEditState(state: Partial<GQL.MovieDataFragment>) {
    setName(state.name ?? undefined);
    setAliases(state.aliases ?? undefined);
    setDuration(state.duration ?? undefined);
    setDate(state.date ?? undefined);
    setRating(state.rating ?? undefined);
    setStudioId(state?.studio?.id ?? undefined);
    setDirector(state.director ?? undefined);
    setSynopsis(state.synopsis ?? undefined);
    setUrl(state.url ?? undefined);
  }

  const updateMovieData = useCallback(
    (movieData: Partial<GQL.MovieDataFragment>) => {
      setFrontImage(undefined);
      setBackImage(undefined);
      updateMovieEditState(movieData);
    },
    []
  );

  useEffect(() => {
    if (movie) {
      updateMovieData(movie);
    }
  }, [movie]);

  function getMovieInput() {
    const input: Partial<GQL.MovieCreateInput | GQL.MovieUpdateInput> = {
      name,
      aliases,
      duration,
      date,
      rating: rating ?? null,
      studio_id: studioId ?? null,
      director,
      synopsis,
      url,
    };

    if (movie && movie.id) {
      (input as GQL.MovieUpdateInput).id = movie.id;
    }
    return input;
  }

  function updateMovieEditStateFromScraper(
    state: Partial<GQL.ScrapedMovieDataFragment>
  ) {
    if (state.name) {
      setName(state.name);
    }

    if (state.aliases) {
      setAliases(state.aliases ?? undefined);
    }

    if (state.duration) {
      setDuration(DurationUtils.stringToSeconds(state.duration) ?? undefined);
    }

    if (state.date) {
      setDate(state.date ?? undefined);
    }

    if (state.studio && state.studio.id) {
      setStudioId(state.studio.id ?? undefined);
    }

    if (state.director) {
      setDirector(state.director ?? undefined);
    }
    if (state.synopsis) {
      setSynopsis(state.synopsis ?? undefined);
    }
    if (state.url) {
      setUrl(state.url ?? undefined);
    }

    const imageStr = (state as GQL.ScrapedMovieDataFragment).front_image;
    setFrontImage(imageStr ?? undefined);

    const backImageStr = (state as GQL.ScrapedMovieDataFragment).back_image;
    setBackImage(backImageStr ?? undefined);
  }

  async function onScrapeMovieURL() {
    if (!url) return;
    setIsLoading(true);

    try {
      const result = await queryScrapeMovieURL(url);
      if (!result.data || !result.data.scrapeMovieURL) {
        return;
      }

      // if this is a new movie, just dump the data
      if (isNew) {
        updateMovieEditStateFromScraper(result.data.scrapeMovieURL);
      } else {
        setScrapedMovie(result.data.scrapeMovieURL);
      }
    } catch (e) {
      Toast.error(e);
    } finally {
      setIsLoading(false);
    }
  }

  function urlScrapable(scrapedUrl: string) {
    return (
      !!scrapedUrl &&
      (Scrapers?.data?.listMovieScrapers ?? []).some((s) =>
        (s?.movie?.urls ?? []).some((u) => scrapedUrl.includes(u))
      )
    );
  }

  function maybeRenderScrapeButton() {
    if (!url || !urlScrapable(url)) {
      return undefined;
    }
    return (
      <Button
        className="minimal scrape-url-button"
        onClick={() => onScrapeMovieURL()}
      >
        <Icon icon="file-upload" />
      </Button>
    );
  }

  function maybeRenderScrapeDialog() {
    if (!scrapedMovie) {
      return;
    }

    const currentMovie = getMovieInput();

    // Get image paths for scrape gui
    currentMovie.front_image = movie?.front_image_path;
    currentMovie.back_image = movie?.back_image_path;

    return (
      <MovieScrapeDialog
        movie={currentMovie}
        scraped={scrapedMovie}
        onClose={(m) => {
          onScrapeDialogClosed(m);
        }}
      />
    );
  }

  function onScrapeDialogClosed(p?: GQL.ScrapedMovieDataFragment) {
    if (p) {
      updateMovieEditStateFromScraper(p);
    }
    setScrapedMovie(undefined);
  }

  function onFrontImageChange(event: React.FormEvent<HTMLInputElement>) {
    ImageUtils.onImageChange(event, setFrontImage);
  }

  function onBackImageChange(event: React.FormEvent<HTMLInputElement>) {
    ImageUtils.onImageChange(event, setBackImage);
  }

  if (isLoading) return <LoadingIndicator />;

  const isEditing = true;

  // TODO: CSS class
  return (
    <div className="row">
      <div className="movie-details col">
        {isNew && <h2>Add Movie</h2>}

        <Table>
          <tbody>
            {TableUtils.renderInputGroup({
              title: "Name",
              value: name ?? "",
              isEditing,
              onChange: setName,
            })}
            {TableUtils.renderInputGroup({
              title: "Aliases",
              value: aliases,
              isEditing,
              onChange: setAliases,
            })}
            {TableUtils.renderDurationInput({
              title: "Duration",
              value: duration ? duration.toString() : "",
              isEditing,
              onChange: (value: string | undefined) =>
                setDuration(value ? Number.parseInt(value, 10) : undefined),
            })}
            {TableUtils.renderInputGroup({
              title: "Date (YYYY-MM-DD)",
              value: date,
              isEditing,
              onChange: setDate,
            })}
            <tr>
              <td>Studio</td>
              <td>
                <StudioSelect
                  onSelect={(items) =>
                    setStudioId(items.length > 0 ? items[0]?.id : undefined)
                  }
                  ids={studioId ? [studioId] : []}
                />
              </td>
            </tr>
            {TableUtils.renderInputGroup({
              title: "Director",
              value: director,
              isEditing,
              onChange: setDirector,
            })}
            <tr>
              <td>Rating</td>
              <td>
                <RatingStars
                  value={rating}
                  onSetRating={(value) => setRating(value)}
                />
              </td>
            </tr>
          </tbody>
        </Table>

        <Form.Group controlId="url">
          <Form.Label>URL {maybeRenderScrapeButton()}</Form.Label>
          <div>
            {EditableTextUtils.renderInputGroup({
              isEditing,
              onChange: setUrl,
              value: url,
              url: TextUtils.sanitiseURL(url),
            })}
          </div>
        </Form.Group>

        <Form.Group controlId="synopsis">
          <Form.Label>Synopsis</Form.Label>
          <Form.Control
            as="textarea"
            className="movie-synopsis text-input"
            onChange={(newValue: React.ChangeEvent<HTMLTextAreaElement>) =>
              setSynopsis(newValue.currentTarget.value)
            }
            value={synopsis}
          />
        </Form.Group>

        <DetailsEditNavbar
          objectName={movie?.name ?? "movie"}
          isNew={isNew}
          isEditing={isEditing}
          onToggleEdit={onCancel}
          onSave={() => onSubmit(getMovieInput())}
          onImageChange={onFrontImageChange}
          onImageChangeURL={setFrontImage}
          onClearImage={() => { setFrontImage(null)}}
          onBackImageChange={onBackImageChange}
          onBackImageChangeURL={setBackImage}
          onClearBackImage={() => { setBackImage(null)}}
          onDelete={onDelete}
        />
      </div>

      {maybeRenderScrapeDialog()}
    </div>
  );
};
