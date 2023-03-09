import React, { useMemo, useState } from "react";
import { useHistory, useLocation } from "react-router-dom";

import * as GQL from "src/core/generated-graphql";
import { useTagCreate } from "src/core/StashService";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { useToast } from "src/hooks/Toast";
import { tagRelationHook } from "src/core/tags";
import { TagEditPanel } from "./TagEditPanel";

const TagCreate: React.FC = () => {
  const history = useHistory();
  const location = useLocation();
  const Toast = useToast();

  const query = useMemo(() => new URLSearchParams(location.search), [location]);
  const tag = {
    name: query.get("q") ?? undefined,
  };

  // Editing tag state
  const [image, setImage] = useState<string | null>();
  const [encodingImage, setEncodingImage] = useState<boolean>(false);

  const [createTag] = useTagCreate();

  async function onSave(input: GQL.TagCreateInput) {
    try {
      const oldRelations = {
        parents: [],
        children: [],
      };
      const result = await createTag({
        variables: { input },
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
          {encodingImage ? (
            <LoadingIndicator message="Encoding image..." />
          ) : (
            renderImage()
          )}
        </div>
        <TagEditPanel
          tag={tag}
          onSubmit={onSave}
          onCancel={() => history.push("/tags")}
          onDelete={() => {}}
          setImage={setImage}
          setEncodingImage={setEncodingImage}
        />
      </div>
    </div>
  );
};

export default TagCreate;
