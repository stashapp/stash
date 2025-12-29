import { useMemo } from "react";
import { Button, ButtonGroup } from "react-bootstrap";
import * as GQL from "src/core/generated-graphql";
import { Icon } from "../Shared/Icon";
import { TagLink } from "../Shared/TagLink";
import { HoverPopover } from "../Shared/HoverPopover";
import NavUtils from "src/utils/navigation";
import TextUtils from "src/utils/text";
import { useConfigurationContext } from "src/hooks/Config";
import { GridCard } from "../Shared/GridCard/GridCard";
import { faTag } from "@fortawesome/free-solid-svg-icons";
import { markerTitle } from "src/core/markers";
import { Link } from "react-router-dom";
import { objectTitle } from "src/core/files";
import { PatchComponent } from "src/patch";
import { PerformerPopoverButton } from "../Shared/PerformerPopoverButton";
import { ScenePreview } from "./SceneCard";
import { TruncatedText } from "../Shared/TruncatedText";

interface ISceneMarkerCardProps {
  marker: GQL.SceneMarkerDataFragment;
  cardWidth?: number;
  previewHeight?: number;
  index?: number;
  compact?: boolean;
  selecting?: boolean;
  selected?: boolean | undefined;
  zoomIndex?: number;
  onSelectedChanged?: (selected: boolean, shiftKey: boolean) => void;
}

const SceneMarkerCardPopovers = PatchComponent(
  "SceneMarkerCard.Popovers",
  (props: ISceneMarkerCardProps) => {
    function maybeRenderPerformerPopoverButton() {
      if (props.marker.scene.performers.length <= 0) return;

      return (
        <PerformerPopoverButton
          performers={props.marker.scene.performers}
          linkType="scene_marker"
        />
      );
    }

    function renderTagPopoverButton() {
      const popoverContent = [
        <TagLink
          key={props.marker.primary_tag.id}
          tag={props.marker.primary_tag}
          linkType="scene_marker"
        />,
      ];

      props.marker.tags.map((tag) =>
        popoverContent.push(
          <TagLink key={tag.id} tag={tag} linkType="scene_marker" />
        )
      );

      return (
        <HoverPopover
          className="tag-count"
          placement="bottom"
          content={popoverContent}
        >
          <Button className="minimal">
            <Icon icon={faTag} />
            <span>{popoverContent.length}</span>
          </Button>
        </HoverPopover>
      );
    }

    function renderPopoverButtonGroup() {
      if (!props.compact) {
        return (
          <>
            <hr />
            <ButtonGroup className="card-popovers">
              {maybeRenderPerformerPopoverButton()}
              {renderTagPopoverButton()}
            </ButtonGroup>
          </>
        );
      }
    }

    return <>{renderPopoverButtonGroup()}</>;
  }
);

const SceneMarkerCardDetails = PatchComponent(
  "SceneMarkerCard.Details",
  (props: ISceneMarkerCardProps) => {
    return (
      <div className="scene-marker-card__details">
        <span className="scene-marker-card__time">
          {TextUtils.formatTimestampRange(
            props.marker.seconds,
            props.marker.end_seconds ?? undefined
          )}
        </span>
        <TruncatedText
          className="scene-marker-card__scene"
          lineCount={3}
          text={
            <Link to={NavUtils.makeSceneMarkersSceneUrl(props.marker.scene)}>
              {objectTitle(props.marker.scene)}
            </Link>
          }
        />
      </div>
    );
  }
);

const SceneMarkerCardImage = PatchComponent(
  "SceneMarkerCard.Image",
  (props: ISceneMarkerCardProps) => {
    const { configuration } = useConfigurationContext();

    const file = useMemo(
      () =>
        props.marker.scene.files.length > 0
          ? props.marker.scene.files[0]
          : undefined,
      [props.marker.scene]
    );

    function isPortrait() {
      const width = file?.width ? file.width : 0;
      const height = file?.height ? file.height : 0;
      return height > width;
    }

    function maybeRenderSceneSpecsOverlay() {
      return (
        <div className="scene-specs-overlay">
          {props.marker.end_seconds && (
            <span className="overlay-duration">
              {TextUtils.secondsToTimestamp(
                props.marker.end_seconds - props.marker.seconds
              )}
            </span>
          )}
        </div>
      );
    }

    return (
      <>
        <ScenePreview
          image={props.marker.screenshot ?? undefined}
          video={props.marker.stream ?? undefined}
          soundActive={configuration?.interface?.soundOnPreview ?? false}
          isPortrait={isPortrait()}
        />
        {maybeRenderSceneSpecsOverlay()}
      </>
    );
  }
);

export const SceneMarkerCard = PatchComponent(
  "SceneMarkerCard",
  (props: ISceneMarkerCardProps) => {
    function zoomIndex() {
      if (!props.compact && props.zoomIndex !== undefined) {
        return `zoom-${props.zoomIndex}`;
      }

      return "";
    }

    return (
      <GridCard
        className={`scene-marker-card ${zoomIndex()}`}
        url={NavUtils.makeSceneMarkerUrl(props.marker)}
        title={markerTitle(props.marker)}
        width={props.cardWidth}
        linkClassName="scene-marker-card-link"
        thumbnailSectionClassName="video-section"
        resumeTime={props.marker.seconds}
        image={<SceneMarkerCardImage {...props} />}
        details={<SceneMarkerCardDetails {...props} />}
        popovers={<SceneMarkerCardPopovers {...props} />}
        selected={props.selected}
        selecting={props.selecting}
        onSelectedChanged={props.onSelectedChanged}
      />
    );
  }
);
