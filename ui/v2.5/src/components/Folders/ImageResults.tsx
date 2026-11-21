import type React from "react";
import { useState } from "react";
import { FormattedMessage } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import { ImageCardGrid } from "../Images/ImageCardGrid";
import { Pagination, PaginationIndex } from "../List/Pagination";
import { ErrorMessage } from "../Shared/ErrorMessage";
import { LoadingIndicator } from "../Shared/LoadingIndicator";

export const ImageResults: React.FC<{
  folderId: string;
  imagesPerPage: number;
}> = ({ folderId, imagesPerPage }) => {
  const [page, setPage] = useState(1);

  const { loading, error, data } = GQL.useFindImagesQuery({
    variables: {
      filter: {
        page,
        per_page: imagesPerPage,
        sort: "title",
        direction: GQL.SortDirectionEnum.Asc,
      },
      image_filter: {
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
    return <ErrorMessage error="Query unexpectedly has no data" />;
  }

  if (data.findImages.count === 0) {
    return (
      <p className="text-muted">
        <FormattedMessage id="folder_browser.no_images" />
      </p>
    );
  }

  return (
    <>
      <div className="pagination-index-container">
        <Pagination
          currentPage={page}
          itemsPerPage={imagesPerPage}
          totalItems={data.findImages.count}
          onChangePage={setPage}
        />
        <PaginationIndex
          loading={loading}
          currentPage={page}
          itemsPerPage={imagesPerPage}
          totalItems={data.findImages.count}
        />
      </div>

      <ImageCardGrid
        images={data.findImages.images}
        selectedIds={new Set()}
        zoomIndex={0}
        onSelectChange={() => {}}
        onPreview={() => {}}
      />
    </>
  );
};
