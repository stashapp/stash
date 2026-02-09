import { Button, Form, Col, Row } from "react-bootstrap";
import * as GQL from "src/core/generated-graphql";
import React, { useCallback, useEffect, useMemo, useState } from "react";
import { Icon } from "../Shared/Icon";
import { ModalComponent } from "src/components/Shared/Modal";
import * as FormUtils from "src/utils/form";
import { queryFindTagsByID, useTagsMerge } from "src/core/StashService";
import { FormattedMessage, useIntl } from "react-intl";
import { useToast } from "src/hooks/Toast";
import { faExchangeAlt, faSignInAlt } from "@fortawesome/free-solid-svg-icons";
import { Tag, TagSelect } from "./TagSelect";
import {
  CustomFieldScrapeResults,
  hasScrapedValues,
  ObjectListScrapeResult,
  ScrapeResult,
} from "../Shared/ScrapeDialog/scrapeResult";
import { sortStoredIdObjects } from "src/utils/data";
import ImageUtils from "src/utils/image";
import { uniq } from "lodash-es";
import { LoadingIndicator } from "../Shared/LoadingIndicator";
import {
  ScrapedCustomFieldRows,
  ScrapeDialogRow,
  ScrapedImageRow,
  ScrapedInputGroupRow,
  ScrapedStringListRow,
  ScrapedTextAreaRow,
} from "../Shared/ScrapeDialog/ScrapeDialogRow";
import { ScrapedTagsRow } from "../Shared/ScrapeDialog/ScrapedObjectsRow";
import { StringListSelect } from "../Shared/Select";
import { ScrapeDialog } from "../Shared/ScrapeDialog/ScrapeDialog";

interface IStashIDsField {
  values: GQL.StashId[];
}

const StashIDsField: React.FC<IStashIDsField> = ({ values }) => {
  return <StringListSelect value={values.map((v) => v.stash_id)} />;
};

interface ITagMergeDetailsProps {
  sources: GQL.TagDataFragment[];
  dest: GQL.TagDataFragment;
  onClose: (values?: GQL.TagUpdateInput) => void;
}

const TagMergeDetails: React.FC<ITagMergeDetailsProps> = ({
  sources,
  dest,
  onClose,
}) => {
  const intl = useIntl();

  const [loading, setLoading] = useState(true);

  const filterCandidates = useCallback(
    (t: { stored_id: string }) =>
      t.stored_id !== dest.id && sources.every((s) => s.id !== t.stored_id),
    [dest.id, sources]
  );

  const [name, setName] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(dest.name)
  );
  const [sortName, setSortName] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(dest.sort_name)
  );
  const [aliases, setAliases] = useState<ScrapeResult<string[]>>(
    new ScrapeResult<string[]>(dest.aliases)
  );
  const [description, setDescription] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(dest.description)
  );
  const [parentTags, setParentTags] = useState<
    ObjectListScrapeResult<GQL.ScrapedTag>
  >(
    new ObjectListScrapeResult<GQL.ScrapedTag>(
      sortStoredIdObjects(
        dest.parents.map(idToStoredID).filter(filterCandidates)
      )
    )
  );
  const [childTags, setChildTags] = useState<
    ObjectListScrapeResult<GQL.ScrapedTag>
  >(
    new ObjectListScrapeResult<GQL.ScrapedTag>(
      sortStoredIdObjects(
        dest.children.map(idToStoredID).filter(filterCandidates)
      )
    )
  );

  const [stashIDs, setStashIDs] = useState(new ScrapeResult<GQL.StashId[]>([]));

  const [image, setImage] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(dest.image_path)
  );

  const [customFields, setCustomFields] = useState<CustomFieldScrapeResults>(
    new Map()
  );

  function idToStoredID(o: { id: string; name: string }) {
    return {
      stored_id: o.id,
      name: o.name,
    };
  }

  // calculate the values for everything
  // uses the first set value for single value fields, and combines all
  useEffect(() => {
    async function loadImages() {
      const src = sources.find((s) => s.image_path);
      if (!dest.image_path || !src) return;

      setLoading(true);

      const destData = await ImageUtils.imageToDataURL(dest.image_path);
      const srcData = await ImageUtils.imageToDataURL(src.image_path!);

      // keep destination image by default
      const useNewValue = false;
      setImage(new ScrapeResult(destData, srcData, useNewValue));

      setLoading(false);
    }

    // append dest to all so that if dest has stash_ids with the same
    // endpoint, then it will be excluded first
    const all = sources.concat(dest);

    setName(
      new ScrapeResult(dest.name, sources.find((s) => s.name)?.name, !dest.name)
    );
    setSortName(
      new ScrapeResult(
        dest.sort_name,
        sources.find((s) => s.sort_name)?.sort_name,
        !dest.sort_name
      )
    );

    setDescription(
      new ScrapeResult(
        dest.description,
        sources.find((s) => s.description)?.description,
        !dest.description
      )
    );

    // default alias list should be the existing aliases, plus the names of all sources,
    // plus all source aliases, deduplicated
    const allAliases = uniq(
      dest.aliases.concat(
        sources.map((s) => s.name),
        sources.flatMap((s) => s.aliases)
      )
    );
    setAliases(new ScrapeResult(dest.aliases, allAliases, !!allAliases.length));

    // default parent/child tags should be the existing tags, plus all source parent/child tags, deduplicated
    const allParentTags = uniq(all.flatMap((s) => s.parents))
      .map(idToStoredID)
      .filter(filterCandidates); // exclude self and sources

    setParentTags(
      new ObjectListScrapeResult<GQL.ScrapedTag>(
        sortStoredIdObjects(dest.parents.map(idToStoredID)),
        sortStoredIdObjects(allParentTags),
        !!allParentTags.length
      )
    );

    const allChildTags = uniq(all.flatMap((s) => s.children))
      .map(idToStoredID)
      .filter(filterCandidates); // exclude self and sources

    setChildTags(
      new ObjectListScrapeResult<GQL.ScrapedTag>(
        sortStoredIdObjects(
          dest.children.map(idToStoredID).filter(filterCandidates)
        ),
        sortStoredIdObjects(allChildTags),
        !!allChildTags.length
      )
    );

    setStashIDs(
      new ScrapeResult(
        dest.stash_ids,
        all
          .map((s) => s.stash_ids)
          .flat()
          .filter((s, index, a) => {
            // remove entries with duplicate endpoints
            return index === a.findIndex((ss) => ss.endpoint === s.endpoint);
          })
      )
    );

    setImage(
      new ScrapeResult(
        dest.image_path,
        sources.find((s) => s.image_path)?.image_path,
        !dest.image_path
      )
    );

    const customFieldNames = new Set<string>(Object.keys(dest.custom_fields));

    for (const s of sources) {
      for (const n of Object.keys(s.custom_fields)) {
        customFieldNames.add(n);
      }
    }

    setCustomFields(
      new Map(
        Array.from(customFieldNames)
          .sort()
          .map((field) => {
            return [
              field,
              new ScrapeResult(
                dest.custom_fields?.[field],
                sources.find((s) => s.custom_fields?.[field])?.custom_fields?.[
                  field
                ],
                dest.custom_fields?.[field] === undefined
              ),
            ];
          })
      )
    );

    loadImages();
  }, [sources, dest, filterCandidates]);

  const hasCustomFieldValues = useMemo(() => {
    return hasScrapedValues(Array.from(customFields.values()));
  }, [customFields]);

  // ensure this is updated if fields are changed
  const hasValues = useMemo(() => {
    return (
      hasCustomFieldValues ||
      hasScrapedValues([
        name,
        sortName,
        aliases,
        description,
        parentTags,
        childTags,
        stashIDs,
        image,
      ])
    );
  }, [
    name,
    sortName,
    aliases,
    description,
    parentTags,
    childTags,
    stashIDs,
    image,
    hasCustomFieldValues,
  ]);

  function renderScrapeRows() {
    if (loading) {
      return (
        <div>
          <LoadingIndicator />
        </div>
      );
    }

    if (!hasValues) {
      return (
        <div>
          <FormattedMessage id="dialogs.merge.empty_results" />
        </div>
      );
    }

    return (
      <>
        <ScrapedInputGroupRow
          field="name"
          title={intl.formatMessage({ id: "name" })}
          result={name}
          onChange={(value) => setName(value)}
        />
        <ScrapedInputGroupRow
          field="sort_name"
          title={intl.formatMessage({ id: "sort_name" })}
          result={sortName}
          onChange={(value) => setSortName(value)}
        />
        <ScrapedStringListRow
          field="aliases"
          title={intl.formatMessage({ id: "aliases" })}
          result={aliases}
          onChange={(value) => setAliases(value)}
        />
        <ScrapedTagsRow
          field="parent_tags"
          title={intl.formatMessage({ id: "parent_tags" })}
          result={parentTags}
          onChange={(value) => setParentTags(value)}
        />
        <ScrapedTagsRow
          field="child_tags"
          title={intl.formatMessage({ id: "sub_tags" })}
          result={childTags}
          onChange={(value) => setChildTags(value)}
        />
        <ScrapedTextAreaRow
          field="description"
          title={intl.formatMessage({ id: "description" })}
          result={description}
          onChange={(value) => setDescription(value)}
        />
        <ScrapeDialogRow
          field="stash_ids"
          title={intl.formatMessage({ id: "stash_id" })}
          result={stashIDs}
          originalField={
            <StashIDsField values={stashIDs?.originalValue ?? []} />
          }
          newField={<StashIDsField values={stashIDs?.newValue ?? []} />}
          onChange={(value) => setStashIDs(value)}
        />
        <ScrapedImageRow
          field="image"
          title={intl.formatMessage({ id: "tag_image" })}
          className="performer-image"
          result={image}
          onChange={(value) => setImage(value)}
        />
        {hasCustomFieldValues && (
          <ScrapedCustomFieldRows
            results={customFields}
            onChange={(newCustomFields) => setCustomFields(newCustomFields)}
          />
        )}
      </>
    );
  }

  function createValues(): GQL.TagUpdateInput {
    // only set the cover image if it's different from the existing cover image
    const coverImage = image.useNewValue ? image.getNewValue() : undefined;

    return {
      id: dest.id,
      name: name.getNewValue(),
      sort_name: sortName.getNewValue(),
      aliases: aliases
        .getNewValue()
        ?.map((s) => s.trim())
        .filter((s) => s.length > 0),
      parent_ids: parentTags.getNewValue()?.map((t) => t.stored_id!),
      child_ids: childTags.getNewValue()?.map((t) => t.stored_id!),
      description: description.getNewValue(),
      stash_ids: stashIDs.getNewValue(),
      image: coverImage,
      custom_fields: {
        partial: Object.fromEntries(
          Array.from(customFields.entries()).flatMap(([field, v]) =>
            v.useNewValue ? [[field, v.getNewValue()]] : []
          )
        ),
      },
    };
  }

  const dialogTitle = intl.formatMessage({
    id: "actions.merge",
  });

  const destinationLabel = !hasValues
    ? ""
    : intl.formatMessage({ id: "dialogs.merge.destination" });
  const sourceLabel = !hasValues
    ? ""
    : intl.formatMessage({ id: "dialogs.merge.source" });

  return (
    <ScrapeDialog
      className="tag-merge-dialog"
      title={dialogTitle}
      existingLabel={destinationLabel}
      scrapedLabel={sourceLabel}
      onClose={(apply) => {
        if (!apply) {
          onClose();
        } else {
          onClose(createValues());
        }
      }}
    >
      {renderScrapeRows()}
    </ScrapeDialog>
  );
};

interface ITagMergeModalProps {
  show: boolean;
  onClose: (mergedID?: string) => void;
  tags: Tag[];
}

export const TagMergeModal: React.FC<ITagMergeModalProps> = ({
  show,
  onClose,
  tags,
}) => {
  const [src, setSrc] = useState<Tag[]>([]);
  const [dest, setDest] = useState<Tag | null>(null);

  const [loadedSources, setLoadedSources] = useState<GQL.TagDataFragment[]>([]);
  const [loadedDest, setLoadedDest] = useState<GQL.TagDataFragment>();

  const [secondStep, setSecondStep] = useState(false);

  const [running, setRunning] = useState(false);

  const [mergeTags] = useTagsMerge();

  const intl = useIntl();
  const Toast = useToast();

  const title = intl.formatMessage({
    id: "actions.merge",
  });

  useEffect(() => {
    if (tags.length > 0) {
      setDest(tags[0]);
      setSrc(tags.slice(1));
    }
  }, [tags]);

  async function loadTags() {
    try {
      const tagIDs = src.map((s) => s.id);
      tagIDs.push(dest!.id);
      const query = await queryFindTagsByID(tagIDs);
      const { tags: loadedTags } = query.data.findTags;

      setLoadedDest(loadedTags.find((s) => s.id === dest!.id));
      setLoadedSources(loadedTags.filter((s) => s.id !== dest!.id));
      setSecondStep(true);
    } catch (e) {
      Toast.error(e);
      return;
    }
  }

  async function onMerge(values: GQL.TagUpdateInput) {
    if (!dest) return;

    const source = src.map((s) => s.id);
    const destination = dest.id;

    try {
      setRunning(true);
      const result = await mergeTags({
        variables: {
          source,
          destination,
          values,
        },
      });
      if (result.data?.tagsMerge) {
        Toast.success(intl.formatMessage({ id: "toast.merged_tags" }));
        onClose(dest.id);
      }
    } catch (e) {
      Toast.error(e);
    } finally {
      setRunning(false);
    }
  }

  function canMerge() {
    return src.length > 0 && dest !== null;
  }

  function switchTags() {
    if (src.length && dest !== null) {
      const newDest = src[0];
      setSrc([...src.slice(1), dest]);
      setDest(newDest);
    }
  }

  if (secondStep && dest) {
    return (
      <TagMergeDetails
        sources={loadedSources}
        dest={loadedDest!}
        onClose={(values) => {
          setSecondStep(false);
          if (values) {
            onMerge(values);
          } else {
            onClose();
          }
        }}
      />
    );
  }

  return (
    <ModalComponent
      show={show}
      header={title}
      icon={faSignInAlt}
      accept={{
        text: intl.formatMessage({ id: "actions.merge" }),
        onClick: () => loadTags(),
      }}
      disabled={!canMerge()}
      cancel={{
        variant: "secondary",
        onClick: () => onClose(),
      }}
      isRunning={running}
    >
      <div className="form-container row px-3">
        <div className="col-12 col-lg-6 col-xl-12">
          <Form.Group controlId="source" as={Row}>
            {FormUtils.renderLabel({
              title: intl.formatMessage({ id: "dialogs.merge.source" }),
              labelProps: {
                column: true,
                sm: 3,
                xl: 12,
              },
            })}
            <Col sm={9} xl={12}>
              <TagSelect
                isMulti
                creatable={false}
                onSelect={(items) => setSrc(items)}
                values={src}
                menuPortalTarget={document.body}
              />
            </Col>
          </Form.Group>
          <Form.Group
            controlId="switch"
            as={Row}
            className="justify-content-center"
          >
            <Button
              variant="secondary"
              onClick={() => switchTags()}
              disabled={!src.length || !dest}
              title={intl.formatMessage({ id: "actions.swap" })}
            >
              <Icon className="fa-fw" icon={faExchangeAlt} />
            </Button>
          </Form.Group>
          <Form.Group controlId="destination" as={Row}>
            {FormUtils.renderLabel({
              title: intl.formatMessage({
                id: "dialogs.merge.destination",
              }),
              labelProps: {
                column: true,
                sm: 3,
                xl: 12,
              },
            })}
            <Col sm={9} xl={12}>
              <TagSelect
                isMulti={false}
                creatable={false}
                onSelect={(items) => setDest(items[0])}
                values={dest ? [dest] : undefined}
                menuPortalTarget={document.body}
              />
            </Col>
          </Form.Group>
        </div>
      </div>
    </ModalComponent>
  );
};
