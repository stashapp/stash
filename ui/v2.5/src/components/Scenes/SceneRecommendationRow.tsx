import React, { FunctionComponent } from "react";
import { FindScenesQueryResult } from "src/core/generated-graphql";
import Slider from "react-slick";
import { SceneCard } from "./SceneCard";
import { SceneQueue } from "src/models/sceneQueue";
import { ListFilterModel } from "src/models/list-filter/filter";

interface IProps {
  isTouch: boolean;
  filter: ListFilterModel;
  result: FindScenesQueryResult;
  queue: SceneQueue;
  header: String;
  linkText: String;
}

export const SceneRecommendationRow: FunctionComponent<IProps> = (
  props: IProps
) => {
  function determineSlidesToScroll(prefered: number, cardCount: number) {
    if (props.isTouch) {
      return 1;
    } else if (cardCount! > prefered) {
      return prefered;
    } else {
      return cardCount;
    }
  }

  const cardCount = props.result.data?.findScenes.count;
  var settings = {
    dots: !props.isTouch,
    arrows: !props.isTouch,
    infinite: !props.isTouch,
    speed: 300,
    variableWidth: true,
    swipeToSlide: true,
    slidesToShow: cardCount! > 5 ? 5 : cardCount,
    slidesToScroll: determineSlidesToScroll(5, cardCount!),
    responsive: [
      {
        breakpoint: 1909,
        settings: {
          slidesToShow: cardCount! > 4 ? 4 : cardCount,
          slidesToScroll: determineSlidesToScroll(4, cardCount!),
        },
      },
      {
        breakpoint: 1542,
        settings: {
          slidesToShow: cardCount! > 3 ? 3 : cardCount,
          slidesToScroll: determineSlidesToScroll(3, cardCount!),
        },
      },
      {
        breakpoint: 1170,
        settings: {
          slidesToShow: cardCount! > 2 ? 2 : cardCount,
          slidesToScroll: determineSlidesToScroll(2, cardCount!),
        },
      },
      {
        breakpoint: 801,
        settings: {
          slidesToShow: 1,
          slidesToScroll: 1,
          dots: false,
        },
      },
    ],
  };

  return (
    <div className="recommendation-row scene-recommendations">
      <div className="recommendation-row-head">
        <div>
          <h2>{props.header}</h2>
        </div>
        <a href={`/scenes?${props.filter.makeQueryParameters()}`}>
          {props.linkText}
        </a>
      </div>
      <Slider {...settings}>
        {props.result.data?.findScenes.scenes.map((scene, index) => (
          <SceneCard
            key={scene.id}
            scene={scene}
            queue={props.queue}
            index={index}
            zoomIndex={1}
          />
        ))}
      </Slider>
    </div>
  );
};
