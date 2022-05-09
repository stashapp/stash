import * as GQL from "src/core/generated-graphql";
import { defineMessages } from "react-intl";
import React from "react";
import { useIntl } from "react-intl";
import { TITLE_SUFFIX } from "src/components/Shared";
import {
  useFindDefaultFilter,
  useFindScenes,
  useFindSceneMarkers,
  useFindImages,
  useFindMovies,
  useFindStudios,
  useFindGalleries,
  useFindPerformers,
  useFindTags,
} from "src/core/StashService";
import { SceneCard } from "src/components/Scenes/SceneCard";
import { StudioCard } from "src/components/Studios/StudioCard";
import { MovieCard } from "src/components/Movies/MovieCard";
import { PerformerCard } from "src/components/Performers/PerformerCard";
import { GalleryCard } from "src/components/Galleries/GalleryCard";
import { SceneQueue } from "src/models/sceneQueue";
import { ListFilterModel } from "src/models/list-filter/filter";
import Slider from "react-slick";
// import { RecommendationList } from "./RecommendationList";

const Recommendations: React.FC = () => {
  function isTouchEnabled() {
    return "ontouchstart" in window || navigator.maxTouchPoints > 0;
  }

  function onSelectChange(id: string, selected: boolean, shiftKey: boolean) {}

  const isTouch = isTouchEnabled();

  const intl = useIntl();
  const itemsPerPage = 20;
  const scenefilter = new ListFilterModel(GQL.FilterMode.Scenes);
  scenefilter.sortBy = "date";
  scenefilter.sortDirection = GQL.SortDirectionEnum.Desc;
  scenefilter.itemsPerPage = itemsPerPage;
  const sceneResult = useFindScenes(scenefilter);
  const hasScenes = !sceneResult.data || !sceneResult.data.findScenes;

  const studiofilter = new ListFilterModel(GQL.FilterMode.Studios);
  studiofilter.sortBy = "scenes_count";
  studiofilter.sortDirection = GQL.SortDirectionEnum.Desc;
  studiofilter.itemsPerPage = itemsPerPage;
  const studioResult = useFindStudios(studiofilter);
  const hasStudios = !studioResult.data || !studioResult.data.findStudios;

  const moviefilter = new ListFilterModel(GQL.FilterMode.Movies);
  moviefilter.sortBy = "date";
  moviefilter.sortDirection = GQL.SortDirectionEnum.Desc;
  moviefilter.itemsPerPage = itemsPerPage;
  const movieResult = useFindMovies(moviefilter);
  const hasMovies = !movieResult.data || !movieResult.data.findMovies;

  const performerfilter = new ListFilterModel(GQL.FilterMode.Performers);
  performerfilter.sortBy = "created_at";
  performerfilter.sortDirection = GQL.SortDirectionEnum.Desc;
  performerfilter.itemsPerPage = itemsPerPage;
  const performerResult = useFindPerformers(performerfilter);

  const galleryfilter = new ListFilterModel(GQL.FilterMode.Galleries);
  galleryfilter.sortBy = "created_at";
  galleryfilter.sortDirection = GQL.SortDirectionEnum.Desc;
  galleryfilter.itemsPerPage = itemsPerPage;
  const galleryResult = useFindGalleries(galleryfilter);

  const messages = defineMessages({
    latestScenes: {
      id: "latestScenes",
      defaultMessage: "Latest Scenes",
    },
    mostActiveStudios: {
      id: "mostActiveStudios",
      defaultMessage: "Most Active Studios",
    },
    latestMovies: {
      id: "latestMovies",
      defaultMessage: "Latest Movies",
    },
    latestPerformers: {
      id: "latestPerformers",
      defaultMessage: "Latest Performers",
    },
    latestGalleries: {
      id: "latestGalleries",
      defaultMessage: "Latest Galleries",
    },
  });

  var settings = {
    dots: !isTouch,
    arrows: !isTouch,
    infinite: true,
    speed: 500,
    variableWidth: true,
    swipeToSlide: true,
    slidesToShow: 1,
  };
  const queue = SceneQueue.fromListFilterModel(scenefilter);
  const title_template = `${intl.formatMessage({
    id: "recommendations",
  })} ${TITLE_SUFFIX}`;
  return (
    <div className="recommendations-container">
      <div className="recommendation-head">
        <div>
          <h2>{intl.formatMessage(messages.latestScenes)}</h2>
        </div>
        <a href="/scenes?sortby=date&sortdir=desc">View all</a>
      </div>
      <Slider {...settings}>
        {sceneResult.data?.findScenes.scenes.map((scene, index) => (
          <SceneCard
            key={scene.id}
            scene={scene}
            queue={queue}
            index={index}
            zoomIndex={1}
            selecting={false}
            selected={false}
            onSelectedChanged={(selected: boolean, shiftKey: boolean) =>
              onSelectChange(scene.id, selected, shiftKey)
            }
          />
        ))}
      </Slider>

      <div className="recommendation-head">
        <div>
          <h2>{intl.formatMessage(messages.mostActiveStudios)}</h2>
        </div>
        <a href="/studios?sortby=scenes_count">View all</a>
      </div>
      <Slider {...settings}>
        {studioResult.data?.findStudios.studios.map((studio) => (
          <StudioCard
            key={studio.id}
            studio={studio}
            hideParent={true}
            selecting={false}
            selected={false}
            onSelectedChanged={(selected: boolean, shiftKey: boolean) =>
              onSelectChange(studio.id, selected, shiftKey)
            }
          />
        ))}
      </Slider>

      <div className="recommendation-head">
        <div>
          <h2>{intl.formatMessage(messages.latestMovies)}</h2>
        </div>
        <a href="/movies?sortby=date&sortdir=desc">View all</a>
      </div>
      <Slider {...settings}>
        {movieResult.data?.findMovies.movies.map((p) => (
          <MovieCard
            key={p.id}
            movie={p}
            selecting={false}
            selected={false}
            onSelectedChanged={(selected: boolean, shiftKey: boolean) =>
              onSelectChange(p.id, selected, shiftKey)
            }
          />
        ))}
      </Slider>

      <div className="recommendation-head">
        <div>
          <h2>{intl.formatMessage(messages.latestPerformers)}</h2>
        </div>
        <a href="/performers?sortby=created_at&sortdir=desc">View all</a>
      </div>
      <Slider {...settings}>
        {performerResult.data?.findPerformers.performers.map((p) => (
          <PerformerCard
            key={p.id}
            performer={p}
            selecting={false}
            selected={false}
            onSelectedChanged={(selected: boolean, shiftKey: boolean) =>
              onSelectChange(p.id, selected, shiftKey)
            }
            extraCriteria={undefined}
          />
        ))}
      </Slider>

      <div className="recommendation-head">
        <div>
          <h2>{intl.formatMessage(messages.latestGalleries)}</h2>
        </div>
        <a href="/galleries?sortby=date&sortdir=desc">View all</a>
      </div>
      <Slider {...settings}>
        {galleryResult.data?.findGalleries.galleries.map((gallery) => (
          <GalleryCard
            key={gallery.id}
            gallery={gallery}
            zoomIndex={1}
            selecting={false}
            selected={false}
            onSelectedChanged={(selected: boolean, shiftKey: boolean) =>
              onSelectChange(gallery.id, selected, shiftKey)
            }
          />
        ))}
      </Slider>
    </div>
  );
};

export default Recommendations;
