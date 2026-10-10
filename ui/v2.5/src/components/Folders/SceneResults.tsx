import type React from "react";
import { useState } from "react";
import { FormattedMessage } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import SceneQueue from "src/models/sceneQueue";
import { Pagination, PaginationIndex } from "../List/Pagination";
import { SceneCardGrid } from "../Scenes/SceneCardGrid";
import { ErrorMessage } from "../Shared/ErrorMessage";
import { LoadingIndicator } from "../Shared/LoadingIndicator";

export const SceneResults: React.FC<{
  folderId: string;
  scenesPerPage: number;
}> = ({ folderId, scenesPerPage }) => {
  const [page, setPage] = useState(1);

  const { loading, error, data } = GQL.useFindScenesQuery({
    variables: {
      filter: {
        page,
        per_page: scenesPerPage,
        sort: "title",
        direction: GQL.SortDirectionEnum.Asc,
      },
      scene_filter: {
        files_filter: {
          parent_folder: {
            value: [folderId],
            modifier: GQL.CriterionModifier.Includes,
            depth: 0,
          },
        },
      },
    },
  });

  if (loading) {
    return <LoadingIndicator />;
  }

  if (error) {
    return <ErrorMessage error={error.message} />;
  }

  if (data === undefined) {
    // should never happen
    return <ErrorMessage error="Query data unexpectedly empty." />;
  }

  if (data.findScenes.count === 0) {
    return (
      <p className="text-muted">
        <FormattedMessage id="folder_browser.no_scenes" />
      </p>
    );
  }

  return (
    <>
      <div className="pagination-index-container">
        <Pagination
          currentPage={page}
          itemsPerPage={scenesPerPage}
          totalItems={data.findScenes.count}
          onChangePage={setPage}
        />
        <PaginationIndex
          loading={loading}
          currentPage={page}
          itemsPerPage={scenesPerPage}
          totalItems={data.findScenes.count}
        />
      </div>

      <SceneCardGrid
        scenes={data.findScenes.scenes}
        queue={SceneQueue.fromSceneIDList(
          data.findScenes.scenes.map((scene) => scene.id)
        )}
        selectedIds={new Set()}
        zoomIndex={0}
        onSelectChange={() => {}}
      />
    </>
  );
};
