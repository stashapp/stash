import React, { useEffect } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import * as yup from "yup";
import Mousetrap from "mousetrap";
import { Icon } from "src/components/Shared/Icon";
import { StudioSelect } from "src/components/Shared/Select";
import { DetailsEditNavbar } from "src/components/Shared/DetailsEditNavbar";
import { Button, Form, Col, Row } from "react-bootstrap";
import FormUtils from "src/utils/form";
import ImageUtils from "src/utils/image";
import { getStashIDs } from "src/utils/stashIds";
import { RatingSystem } from "src/components/Shared/Rating/RatingSystem";
import { useFormik } from "formik";
import { Prompt } from "react-router-dom";
import { StringListInput } from "../../Shared/StringListInput";
import { faTrashAlt } from "@fortawesome/free-solid-svg-icons";
import { useRatingKeybinds } from "src/hooks/keybinds";
import { ConfigurationContext } from "src/hooks/Config";
import isEqual from "lodash-es/isEqual";

interface IStudioEditPanel {
  studio: Partial<GQL.StudioDataFragment>;
  onSubmit: (studio: GQL.StudioCreateInput) => void;
  onCancel: () => void;
  onDelete: () => void;
  setImage: (image?: string | null) => void;
  setEncodingImage: (loading: boolean) => void;
}

export const StudioEditPanel: React.FC<IStudioEditPanel> = ({
  studio,
  onSubmit,
  onCancel,
  onDelete,
  setImage,
  setEncodingImage,
}) => {
  const intl = useIntl();

  const isNew = studio.id === undefined;
  const { configuration } = React.useContext(ConfigurationContext);

  const schema = yup.object({
    name: yup.string().required(),
    url: yup.string().ensure(),
    details: yup.string().ensure(),
    parent_id: yup.string().required().nullable(),
    rating100: yup.number().nullable().defined(),
    aliases: yup
      .array(yup.string().required())
      .defined()
      .test({
        name: "unique",
        test: (value, context) => {
          const aliases = [context.parent.name, ...value];
          const dupes = aliases
            .map((e, i, a) => {
              if (a.indexOf(e) !== i) {
                return String(i - 1);
              } else {
                return null;
              }
            })
            .filter((e) => e !== null) as string[];
          if (dupes.length === 0) return true;
          return new yup.ValidationError(dupes.join(" "), value, "aliases");
        },
      }),
    ignore_auto_tag: yup.boolean().defined(),
    stash_ids: yup.mixed<GQL.StashIdInput[]>().defined(),
    image: yup.string().nullable().optional(),
  });

  const initialValues = {
    name: studio.name ?? "",
    url: studio.url ?? "",
    details: studio.details ?? "",
    parent_id: studio.parent_studio?.id ?? null,
    rating100: studio.rating100 ?? null,
    aliases: studio.aliases ?? [],
    ignore_auto_tag: studio.ignore_auto_tag ?? false,
    stash_ids: getStashIDs(studio.stash_ids),
  };

  type InputValues = yup.InferType<typeof schema>;

  const formik = useFormik<InputValues>({
    initialValues,
    enableReinitialize: true,
    validationSchema: schema,
    onSubmit: (values) => onSubmit(values),
  });

  const encodingImage = ImageUtils.usePasteImage((imageData) =>
    formik.setFieldValue("image", imageData)
  );

  useEffect(() => {
    setImage(formik.values.image);
  }, [formik.values.image, setImage]);

  useEffect(() => {
    setEncodingImage(encodingImage);
  }, [setEncodingImage, encodingImage]);

  function setRating(v: number) {
    formik.setFieldValue("rating100", v);
  }

  useRatingKeybinds(
    true,
    configuration?.ui?.ratingSystemOptions?.type,
    setRating
  );

  function onCancelEditing() {
    setImage(undefined);
    onCancel?.();
  }

  // set up hotkeys
  useEffect(() => {
    Mousetrap.bind("s s", () => formik.handleSubmit());

    return () => {
      Mousetrap.unbind("s s");
    };
  });

  function onImageLoad(imageData: string | null) {
    formik.setFieldValue("image", imageData);
  }

  function onImageChange(event: React.FormEvent<HTMLInputElement>) {
    ImageUtils.onImageChange(event, onImageLoad);
  }

  const removeStashID = (stashID: GQL.StashIdInput) => {
    formik.setFieldValue(
      "stash_ids",
      (formik.values.stash_ids ?? []).filter(
        (s) =>
          !(s.endpoint === stashID.endpoint && s.stash_id === stashID.stash_id)
      )
    );
  };

  function renderStashIDs() {
    if (!formik.values.stash_ids?.length) {
      return;
    }

    return (
      <Row>
        <Form.Label column>StashIDs</Form.Label>
        <Col xs={9}>
          <ul className="pl-0">
            {formik.values.stash_ids.map((stashID) => {
              const base = stashID.endpoint.match(/https?:\/\/.*?\//)?.[0];
              const link = base ? (
                <a
                  href={`${base}studios/${stashID.stash_id}`}
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  {stashID.stash_id}
                </a>
              ) : (
                stashID.stash_id
              );
              return (
                <li key={stashID.stash_id} className="row no-gutters">
                  <Button
                    variant="danger"
                    className="mr-2 py-0"
                    title={intl.formatMessage(
                      { id: "actions.delete_entity" },
                      { entityType: intl.formatMessage({ id: "stash_id" }) }
                    )}
                    onClick={() => removeStashID(stashID)}
                  >
                    <Icon icon={faTrashAlt} />
                  </Button>
                  {link}
                </li>
              );
            })}
          </ul>
        </Col>
      </Row>
    );
  }

  const aliasErrors = Array.isArray(formik.errors.aliases)
    ? formik.errors.aliases[0]
    : formik.errors.aliases;
  const aliasErrorMsg = aliasErrors
    ? intl.formatMessage({ id: "validation.aliases_must_be_unique" })
    : undefined;
  const aliasErrorIdx = aliasErrors?.split(" ").map((e) => parseInt(e));

  return (
    <>
      <Prompt
        when={formik.dirty}
        message={(location, action) => {
          // Check if it's a redirect after studio creation
          if (action === "PUSH" && location.pathname.startsWith("/studios/"))
            return true;
          return intl.formatMessage({ id: "dialogs.unsaved_changes" });
        }}
      />

      <Form noValidate onSubmit={formik.handleSubmit} id="studio-edit">
        <Form.Group controlId="name" as={Row}>
          {FormUtils.renderLabel({
            title: intl.formatMessage({ id: "name" }),
          })}
          <Col xs={9}>
            <Form.Control
              className="text-input"
              {...formik.getFieldProps("name")}
              isInvalid={!!formik.errors.name}
            />
            <Form.Control.Feedback type="invalid">
              {formik.errors.name}
            </Form.Control.Feedback>
          </Col>
        </Form.Group>

        <Form.Group controlId="url" as={Row}>
          {FormUtils.renderLabel({
            title: intl.formatMessage({ id: "url" }),
          })}
          <Col xs={9}>
            <Form.Control
              className="text-input"
              {...formik.getFieldProps("url")}
              isInvalid={!!formik.errors.url}
            />
            <Form.Control.Feedback type="invalid">
              {formik.errors.url}
            </Form.Control.Feedback>
          </Col>
        </Form.Group>

        <Form.Group controlId="details" as={Row}>
          {FormUtils.renderLabel({
            title: intl.formatMessage({ id: "details" }),
          })}
          <Col xs={9}>
            <Form.Control
              as="textarea"
              className="text-input"
              {...formik.getFieldProps("details")}
              isInvalid={!!formik.errors.details}
            />
            <Form.Control.Feedback type="invalid">
              {formik.errors.details}
            </Form.Control.Feedback>
          </Col>
        </Form.Group>

        <Form.Group controlId="parent_studio" as={Row}>
          {FormUtils.renderLabel({
            title: intl.formatMessage({ id: "parent_studios" }),
          })}
          <Col xs={9}>
            <StudioSelect
              onSelect={(items) =>
                formik.setFieldValue(
                  "parent_id",
                  items.length > 0 ? items[0]?.id : null
                )
              }
              ids={formik.values.parent_id ? [formik.values.parent_id] : []}
              excludeIds={studio.id ? [studio.id] : []}
            />
          </Col>
        </Form.Group>

        <Form.Group controlId="rating" as={Row}>
          {FormUtils.renderLabel({
            title: intl.formatMessage({ id: "rating" }),
          })}
          <Col xs={9}>
            <RatingSystem
              value={formik.values.rating100 ?? undefined}
              onSetRating={(value) =>
                formik.setFieldValue("rating100", value ?? null)
              }
            />
          </Col>
        </Form.Group>

        {renderStashIDs()}

        <Form.Group controlId="aliases" as={Row}>
          <Form.Label column xs={3}>
            <FormattedMessage id="aliases" />
          </Form.Label>
          <Col xs={9}>
            <StringListInput
              value={formik.values.aliases ?? []}
              setValue={(value) => formik.setFieldValue("aliases", value)}
              errors={aliasErrorMsg}
              errorIdx={aliasErrorIdx}
            />
          </Col>
        </Form.Group>
      </Form>

      <hr />

      <Form.Group controlId="ignore-auto-tag" as={Row}>
        <Form.Label column xs={3}>
          <FormattedMessage id="ignore_auto_tag" />
        </Form.Label>
        <Col xs={9}>
          <Form.Check
            {...formik.getFieldProps({
              name: "ignore_auto_tag",
              type: "checkbox",
            })}
          />
        </Col>
      </Form.Group>

      <DetailsEditNavbar
        objectName={studio?.name ?? intl.formatMessage({ id: "studio" })}
        isNew={isNew}
        isEditing
        onToggleEdit={onCancelEditing}
        onSave={formik.handleSubmit}
        saveDisabled={(!isNew && !formik.dirty) || !isEqual(formik.errors, {})}
        onImageChange={onImageChange}
        onImageChangeURL={onImageLoad}
        onClearImage={() => onImageLoad(null)}
        onDelete={onDelete}
        acceptSVG
      />
    </>
  );
};
