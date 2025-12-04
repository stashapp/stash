import React, { useMemo } from "react";
import Slider from "@ant-design/react-slick";
import { SceneCard } from "./SceneCard";
import { SceneQueue } from "src/models/sceneQueue";
import { getSlickSliderSettings } from "src/core/recommendations";
import { RecommendationRow } from "../FrontPage/RecommendationRow";
import * as GQL from "src/core/generated-graphql";

interface IProps {
  isTouch: boolean;
  limit?: number;
  header: string;
}

export const AISceneRecommendationRow: React.FC<IProps> = ({
  isTouch,
  limit = 100,
  header,
}) => {
  const { data, loading } = GQL.useSceneRecommendationsQuery({
    variables: { limit },
  });

  const scenes = useMemo(() => {
    return data?.sceneRecommendations?.recommendations?.map((r) => r.scene) ?? [];
  }, [data]);

  const queue = useMemo(() => {
    // Create queue from scene ID list
    const sceneIDs = scenes.map((s) => s.id);
    return SceneQueue.fromSceneIDList(sceneIDs);
  }, [scenes]);

  const cardCount = scenes.length;

  if (!loading && !cardCount) {
    return null;
  }

  return (
    <RecommendationRow
      className="scene-recommendations ai-recommendations"
      header={header}
    >
      <Slider
        {...getSlickSliderSettings(cardCount || limit, isTouch)}
      >
        {loading
          ? [...Array(Math.min(limit, 25))].map((_, i) => (
              <div key={`_${i}`} className="scene-skeleton skeleton-card"></div>
            ))
          : scenes.map((scene, index) => (
              <SceneCard
                key={scene.id}
                scene={scene}
                queue={queue}
                index={index}
                zoomIndex={1}
              />
            ))}
      </Slider>
    </RecommendationRow>
  );
};

