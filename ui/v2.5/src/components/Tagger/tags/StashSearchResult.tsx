import React, { useState } from "react";
import { Button } from "react-bootstrap";

import * as GQL from "src/core/generated-graphql";
import { useUpdateTag } from "../queries";
import TagModal from "./TagModal";
import { faTags } from "@fortawesome/free-solid-svg-icons";
import { useIntl } from "react-intl";
import { mergeTagStashIDs } from "../utils";

interface IStashSearchResultProps {
  tag: GQL.TagListDataFragment;
  stashboxTags: GQL.ScrapedSceneTagDataFragment[];
  endpoint: string;
  onTagTagged: (
    tag: Pick<GQL.TagListDataFragment, "id"> &
      Partial<Omit<GQL.TagListDataFragment, "id">>
  ) => void;
  excludedTagFields: string[];
}

const StashSearchResult: React.FC<IStashSearchResultProps> = ({
  tag,
  stashboxTags,
  onTagTagged,
  excludedTagFields,
  endpoint,
}) => {
  const intl = useIntl();

  const [modalTag, setModalTag] = useState<GQL.ScrapedSceneTagDataFragment>();
  const [saveState, setSaveState] = useState<string>("");
  const [error, setError] = useState<{ message?: string; details?: string }>(
    {}
  );

  const updateTag = useUpdateTag();

  const handleSave = async (input: GQL.TagCreateInput) => {
    setError({});
    setModalTag(undefined);
    setSaveState("Saving tag");

    const updateData: GQL.TagUpdateInput = {
      ...input,
      id: tag.id,
    };

    updateData.stash_ids = await mergeTagStashIDs(
      tag.id,
      input.stash_ids ?? []
    );

    const res = await updateTag(updateData);

    if (!res?.data?.tagUpdate) {
      setError({
        message: intl.formatMessage(
          { id: "tag_tagger.failed_to_save_tag" },
          { tag: input.name ?? tag.name }
        ),
        details:
          res?.errors?.[0]?.message === "UNIQUE constraint failed: tags.name"
            ? intl.formatMessage({
                id: "tag_tagger.name_already_exists",
              })
            : res?.errors?.[0]?.message ?? "",
      });
    } else {
      onTagTagged(tag);
    }
    setSaveState("");
  };

  const tags = stashboxTags.map((p) => (
    <Button
      className="StudioTagger-studio-search-item minimal col-6"
      variant="link"
      key={p.remote_site_id}
      onClick={() => setModalTag(p)}
    >
      <span>{p.name}</span>
    </Button>
  ));

  return (
    <>
      {modalTag && (
        <TagModal
          closeModal={() => setModalTag(undefined)}
          modalVisible={modalTag !== undefined}
          tag={modalTag}
          onSave={handleSave}
          icon={faTags}
          header="Update Tag"
          excludedTagFields={excludedTagFields}
          endpoint={endpoint}
        />
      )}
      <div className="StudioTagger-studio-search">{tags}</div>
      <div className="row no-gutters mt-2 align-items-center justify-content-end">
        {error.message && (
          <div className="text-right text-danger mt-1">
            <strong>
              <span className="mr-2">Error:</span>
              {error.message}
            </strong>
            <div>{error.details}</div>
          </div>
        )}
        {saveState && (
          <strong className="col-4 mt-1 mr-2 text-right">{saveState}</strong>
        )}
      </div>
    </>
  );
};

export default StashSearchResult;
