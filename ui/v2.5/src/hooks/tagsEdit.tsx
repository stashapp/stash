import * as GQL from "src/core/generated-graphql";
import { useTagCreate } from "src/core/StashService";
import { useEffect, useState } from "react";
import { Tag, TagSelect, TagSelectProps } from "src/components/Tags/TagSelect";
import { useToast } from "src/hooks/Toast";
import { useIntl } from "react-intl";
import { Badge, Button } from "react-bootstrap";
import { Icon } from "src/components/Shared/Icon";
import { faPlus } from "@fortawesome/free-solid-svg-icons";
import { CollapseButton } from "src/components/Shared/CollapseButton";

export function useTagsEdit(
  srcTags: Tag[] | undefined,
  setFieldValue: (ids: string[]) => void
) {
  const intl = useIntl();
  const Toast = useToast();
  const [createTag] = useTagCreate();

  const [tags, setTags] = useState<Tag[]>([]);
  const [newTags, setNewTags] = useState<GQL.ScrapedTag[]>();

  function onSetTags(items: Tag[]) {
    setTags(items);
    setFieldValue(items.map((item) => item.id));
  }

  useEffect(() => {
    setTags(srcTags ?? []);
  }, [srcTags]);

  async function createNewTag(toCreate: GQL.ScrapedTag) {
    const tagInput: GQL.TagCreateInput = { name: toCreate.name ?? "" };
    try {
      const result = await createTag({
        variables: {
          input: tagInput,
        },
      });

      if (!result.data?.tagCreate) {
        Toast.error(new Error("Failed to create tag"));
        return;
      }

      // add the new tag to the new tags value
      onSetTags(
        tags.concat([
          {
            id: result.data.tagCreate.id,
            name: toCreate.name ?? "",
            aliases: [],
            stash_ids: result.data.tagCreate.stash_ids,
          },
        ])
      );

      // remove the tag from the list
      const newTagsClone = newTags!.concat();
      const pIndex = newTagsClone.indexOf(toCreate);
      newTagsClone.splice(pIndex, 1);

      setNewTags(newTagsClone);

      Toast.success(
        intl.formatMessage(
          { id: "toast.created_entity" },
          {
            entity: intl.formatMessage({ id: "tag" }).toLocaleLowerCase(),
            entity_name: toCreate.name,
          }
        )
      );
    } catch (e) {
      Toast.error(e);
    }
  }

  function updateTagsStateFromScraper(
    scrapedTags?: Pick<GQL.ScrapedTag, "name" | "stored_id">[]
  ) {
    if (!scrapedTags) {
      return;
    }

    // map tags to their ids and filter out those not found
    const idTags = scrapedTags.filter(
      (t) => t.stored_id !== undefined && t.stored_id !== null
    );
    const newNewTags = scrapedTags.filter((t) => !t.stored_id);
    onSetTags(
      idTags.map((p) => {
        return {
          id: p.stored_id!,
          name: p.name ?? "",
          aliases: [],
          stash_ids: [],
        };
      })
    );

    setNewTags(newNewTags);
  }

  function renderNewTags() {
    if (!newTags || newTags.length === 0) {
      return;
    }

    const ret = (
      <>
        {newTags.map((t) => (
          <Badge
            className="tag-item"
            variant="secondary"
            key={t.name}
            onClick={() => createNewTag(t)}
          >
            {t.name}
            <Button className="minimal ml-2">
              <Icon className="fa-fw" icon={faPlus} />
            </Button>
          </Badge>
        ))}
      </>
    );

    const minCollapseLength = 10;

    if (newTags.length >= minCollapseLength) {
      return (
        <CollapseButton text={`Missing (${newTags.length})`}>
          {ret}
        </CollapseButton>
      );
    }

    return ret;
  }

  function tagsControl(props?: TagSelectProps) {
    return (
      <>
        <TagSelect isMulti onSelect={onSetTags} values={tags} {...props} />
        {renderNewTags()}
      </>
    );
  }

  return {
    tags,
    onSetTags,
    tagsControl,
    updateTagsStateFromScraper,
  };
}
