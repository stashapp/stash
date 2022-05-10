import * as GQL from "src/core/generated-graphql";
import { defineMessages, useIntl } from "react-intl";
import React from "react";
import {
  useFindScenes,
  useFindMovies,
  useFindStudios,
  useFindGalleries,
  useFindPerformers,
} from "src/core/StashService";
import { SceneCard } from "src/components/Scenes/SceneCard";
import { StudioCard } from "src/components/Studios/StudioCard";
import { MovieCard } from "src/components/Movies/MovieCard";
import { PerformerCard } from "src/components/Performers/PerformerCard";
import { GalleryCard } from "src/components/Galleries/GalleryCard";
import { SceneQueue } from "src/models/sceneQueue";
import { ListFilterModel } from "src/models/list-filter/filter";
import Slider from "react-slick";

const Recommendations: React.FC = () => {
  function isTouchEnabled() {
    return "ontouchstart" in window || navigator.maxTouchPoints > 0;
  }

  const isTouch = isTouchEnabled();

  const intl = useIntl();
  const itemsPerPage = 25;
  const scenefilter = new ListFilterModel(GQL.FilterMode.Scenes);
  scenefilter.sortBy = "date";
  scenefilter.sortDirection = GQL.SortDirectionEnum.Desc;
  scenefilter.itemsPerPage = itemsPerPage;
  const sceneResult = useFindScenes(scenefilter);
  const hasScenes =
    sceneResult.data &&
    sceneResult.data.findScenes &&
    sceneResult.data.findScenes.count > 0;

  const studiofilter = new ListFilterModel(GQL.FilterMode.Studios);
  studiofilter.sortBy = "scenes_count";
  studiofilter.sortDirection = GQL.SortDirectionEnum.Desc;
  studiofilter.itemsPerPage = itemsPerPage;
  const studioResult = useFindStudios(studiofilter);
  const hasStudios =
    studioResult.data &&
    studioResult.data.findStudios &&
    studioResult.data.findStudios.count > 0;

  const moviefilter = new ListFilterModel(GQL.FilterMode.Movies);
  moviefilter.sortBy = "date";
  moviefilter.sortDirection = GQL.SortDirectionEnum.Desc;
  moviefilter.itemsPerPage = itemsPerPage;
  const movieResult = useFindMovies(moviefilter);
  const hasMovies =
    movieResult.data &&
    movieResult.data.findMovies &&
    movieResult.data.findMovies.count > 0;

  const performerfilter = new ListFilterModel(GQL.FilterMode.Performers);
  performerfilter.sortBy = "created_at";
  performerfilter.sortDirection = GQL.SortDirectionEnum.Desc;
  performerfilter.itemsPerPage = itemsPerPage;
  const performerResult = useFindPerformers(performerfilter);
  const hasPerformers =
    performerResult.data &&
    performerResult.data.findPerformers &&
    performerResult.data.findPerformers.count > 0;

  const galleryfilter = new ListFilterModel(GQL.FilterMode.Galleries);
  galleryfilter.sortBy = "date";
  galleryfilter.sortDirection = GQL.SortDirectionEnum.Desc;
  galleryfilter.itemsPerPage = itemsPerPage;
  const galleryResult = useFindGalleries(galleryfilter);
  const hasGalleries =
    galleryResult.data &&
    galleryResult.data.findGalleries &&
    galleryResult.data.findGalleries.count > 0;

  const messages = defineMessages({
    emptyServer: {
      id: "emptyServer",
      defaultMessage:
        "Add some scenes to your server to view recommendations on this page.",
    },
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
    infinite: !isTouch,
    speed: 300,
    variableWidth: true,
    swipeToSlide: true,
    slidesToShow: 5,
    slidesToScroll: 5,
    responsive: [
      {
      breakpoint: 1909,
      settings: {
        slidesToShow: 4,
        slidesToScroll: 4
      }
      },
      {
      breakpoint: 1542,
      settings: {
        slidesToShow: 3,
        slidesToScroll: 3
      }
      },
      {
        breakpoint: 1170,
        settings: {
          slidesToShow: 2,
          slidesToScroll: 2
        }
      },
      {
        breakpoint: 801,
        settings: {
          slidesToShow: 1,
          slidesToScroll: 1,
          dots: false,
        }
      }
    ]
  };
  const queue = SceneQueue.fromListFilterModel(scenefilter);

  return (
    <div className="recommendations-container">
      {!hasScenes &&
      !hasStudios &&
      !hasMovies &&
      !hasPerformers &&
      !hasGalleries ? (
        <div className="no-recommendations">
          {intl.formatMessage(messages.emptyServer)}
        </div>
      ) : (
        <div>
          {hasScenes && (
            <div className="recommendation-row">
              <div className="recommendation-row-head">
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
                  />
                ))}
              </Slider>
            </div>
          )}

          {hasStudios && (
            <div className="recommendation-row">
              <div className="recommendation-row-head">
                <div>
                  <h2>{intl.formatMessage(messages.mostActiveStudios)}</h2>
                </div>
                <a href="/studios?sortby=scenes_count&sortdir=desc">View all</a>
              </div>
              <Slider {...settings}>
                {studioResult.data?.findStudios.studios.map((studio) => (
                  <StudioCard
                    key={studio.id}
                    studio={studio}
                    hideParent={true}
                  />
                ))}
              </Slider>
            </div>
          )}

          {hasMovies && (
            <div className="recommendation-row">
              <div className="recommendation-row-head">
                <div>
                  <h2>{intl.formatMessage(messages.latestMovies)}</h2>
                </div>
                <a href="/movies?sortby=date&sortdir=desc">View all</a>
              </div>
              <Slider {...settings}>
                {movieResult.data?.findMovies.movies.map((p) => (
                  <MovieCard key={p.id} movie={p} />
                ))}
              </Slider>
            </div>
          )}

          {hasPerformers && (
            <div className="recommendation-row">
              <div className="recommendation-row-head">
                <div>
                  <h2>{intl.formatMessage(messages.latestPerformers)}</h2>
                </div>
                <a href="/performers?sortby=created_at&sortdir=desc">
                  View all
                </a>
              </div>
              <Slider {...settings}>
                {performerResult.data?.findPerformers.performers.map((p) => (
                  <PerformerCard key={p.id} performer={p} />
                ))}
              </Slider>
            </div>
          )}

          {hasGalleries && (
            <div className="recommendation-row">
              <div className="recommendation-row-head">
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
                  />
                ))}
              </Slider>
            </div>
          )}
        </div>
      )}
    </div>
  );
};

export default Recommendations;
