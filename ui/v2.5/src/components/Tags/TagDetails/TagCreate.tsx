import React, { useState } from "react";
import { useHistory, useLocation } from "react-router-dom";

import * as GQL from "src/core/generated-graphql";
import { useTagCreate } from "src/core/StashService";
import { ImageUtils } from "src/utils";
import { LoadingIndicator } from "src/components/Shared";
import { useToast } from "src/hooks";
import { tagRelationHook } from "src/core/tags";
import { TagEditPanel } from "./TagEditPanel";

const TagCreate: React.FC = () => {
  const history = useHistory();
  const Toast = useToast();

  function useQuery() {
    const { search } = useLocation();
    return React.useMemo(() => new URLSearchParams(search), [search]);
  }

  const query = useQuery();
  const nameQuery = query.get("name");

  // Editing tag state
  const [image, setImage] = useState<string | null>();

  const [createTag] = useTagCreate();

  function onImageLoad(imageData: string) {
    setImage(imageData);
  }

  const imageEncoding = ImageUtils.usePasteImage(onImageLoad, true);

  function getTagInput(
    input: Partial<GQL.TagCreateInput | GQL.TagUpdateInput>
  ) {
    const ret: Partial<GQL.TagCreateInput | GQL.TagUpdateInput> = {
      ...input,
      image,
    };

    return ret;
  }

  async function onSave(
    input: Partial<GQL.TagCreateInput | GQL.TagUpdateInput>
  ) {
    try {
      const oldRelations = {
        parents: [],
        children: [],
      };
      const result = await createTag({
        variables: {
          input: getTagInput(input) as GQL.TagCreateInput,
        },
      });
      if (result.data?.tagCreate?.id) {
        const created = result.data.tagCreate;
        tagRelationHook(created, oldRelations, {
          parents: created.parents,
          children: created.children,
        });
        history.push(`/tags/${result.data.tagCreate.id}`);
      }
    } catch (e) {
      Toast.error(e);
    }
  }

  function renderImage() {
    if (image) {
      return <img className="logo" alt="" src={image} />;
    }
  }

  return (
    <div className="row">
      <div className="tag-details col-md-8">
        <div className="text-center logo-container">
          {imageEncoding ? (
            <LoadingIndicator message="Encoding image..." />
          ) : (
            renderImage()
          )}
        </div>
        <TagEditPanel
          tag={{ name: nameQuery ?? "" }}
          onSubmit={onSave}
          onCancel={() => history.push("/tags")}
          onDelete={() => {}}
          setImage={setImage}
        />
      </div>
    </div>
  );
};

export default TagCreate;
