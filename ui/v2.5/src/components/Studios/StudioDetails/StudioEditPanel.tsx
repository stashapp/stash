import React, { useEffect } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import * as yup from "yup";
import Mousetrap from "mousetrap";
import { Icon, StudioSelect, DetailsEditNavbar } from "src/components/Shared";
import { Button, Form, Col, Row } from "react-bootstrap";
import { FormUtils, ImageUtils, getStashIDs } from "src/utils";
import { RatingStars } from "src/components/Scenes/SceneDetails/RatingStars";
import { useFormik } from "formik";
import { Prompt } from "react-router-dom";
import { StringListInput } from "../../Shared/StringListInput";

interface IStudioEditPanel {
  studio: Partial<GQL.StudioDataFragment>;
  onSubmit: (
    studio: Partial<GQL.StudioCreateInput | GQL.StudioUpdateInput>
  ) => void;
  onCancel: () => void;
  onDelete: () => void;
  onImageChange?: (image?: string | null) => void;
  onImageEncoding?: (loading?: boolean) => void;
}

export const StudioEditPanel: React.FC<IStudioEditPanel> = ({
  studio,
  onSubmit,
  onCancel,
  onDelete,
  onImageChange,
  onImageEncoding,
}) => {
  const intl = useIntl();

  const isNew = !studio || !studio.id;

  const imageEncoding = ImageUtils.usePasteImage(onImageLoad, true);

  const schema = yup.object({
    name: yup.string().required(),
    url: yup.string().optional().nullable(),
    details: yup.string().optional().nullable(),
    image: yup.string().optional().nullable(),
    rating: yup.number().optional().nullable(),
    parent_id: yup.string().optional().nullable(),
    stash_ids: yup.mixed<GQL.StashIdInput>().optional().nullable(),
    aliases: yup
      .array(yup.string().required())
      .optional()
      .test({
        name: "unique",
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        test: (value: any) => {
          return (value ?? []).length === new Set(value).size;
        },
        message: "aliases must be unique",
      }),
    ignore_auto_tag: yup.boolean().optional(),
  });

  const initialValues = {
    name: studio.name ?? "",
    url: studio.url ?? "",
    details: studio.details ?? "",
    image: undefined,
    rating: studio.rating ?? null,
    parent_id: studio.parent_studio?.id,
    stash_ids: studio.stash_ids ?? undefined,
    aliases: studio.aliases,
    ignore_auto_tag: studio.ignore_auto_tag ?? false,
  };

  type InputValues = typeof initialValues;

  const formik = useFormik({
    initialValues,
    validationSchema: schema,
    onSubmit: (values) => onSubmit(getStudioInput(values)),
  });

  function setRating(v: number) {
    formik.setFieldValue("rating", v);
  }

  function onImageLoad(imageData: string) {
    formik.setFieldValue("image", imageData);
  }

  function getStudioInput(values: InputValues) {
    const input: Partial<GQL.StudioCreateInput | GQL.StudioUpdateInput> = {
      ...values,
      stash_ids: getStashIDs(values.stash_ids),
    };

    if (studio && studio.id) {
      (input as GQL.StudioUpdateInput).id = studio.id;
    }
    return input;
  }

  // set up hotkeys
  useEffect(() => {
    Mousetrap.bind("s s", () => formik.handleSubmit());

    // numeric keypresses get caught by jwplayer, so blur the element
    // if the rating sequence is started
    Mousetrap.bind("r", () => {
      if (document.activeElement instanceof HTMLElement) {
        document.activeElement.blur();
      }

      Mousetrap.bind("0", () => setRating(NaN));
      Mousetrap.bind("1", () => setRating(1));
      Mousetrap.bind("2", () => setRating(2));
      Mousetrap.bind("3", () => setRating(3));
      Mousetrap.bind("4", () => setRating(4));
      Mousetrap.bind("5", () => setRating(5));

      setTimeout(() => {
        Mousetrap.unbind("0");
        Mousetrap.unbind("1");
        Mousetrap.unbind("2");
        Mousetrap.unbind("3");
        Mousetrap.unbind("4");
        Mousetrap.unbind("5");
      }, 1000);
    });

    return () => {
      Mousetrap.unbind("s s");

      Mousetrap.unbind("e");
    };
  });

  useEffect(() => {
    if (onImageChange) {
      onImageChange(formik.values.image);
    }
    return () => onImageChange?.();
  }, [formik.values.image, onImageChange]);

  useEffect(() => onImageEncoding?.(imageEncoding), [
    onImageEncoding,
    imageEncoding,
  ]);

  function onImageChangeHandler(event: React.FormEvent<HTMLInputElement>) {
    ImageUtils.onImageChange(event, onImageLoad);
  }

  function onImageChangeURL(url: string) {
    formik.setFieldValue("image", url);
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
                    <Icon icon="trash-alt" />
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
            <RatingStars
              value={formik.values.rating ?? undefined}
              onSetRating={(value) =>
                formik.setFieldValue("rating", value ?? null)
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
              errors={formik.errors.aliases}
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
        onToggleEdit={onCancel}
        onSave={() => formik.handleSubmit()}
        saveDisabled={!formik.dirty}
        onImageChange={onImageChangeHandler}
        onImageChangeURL={onImageChangeURL}
        onClearImage={() => {
          formik.setFieldValue("image", null);
        }}
        onDelete={onDelete}
        acceptSVG
      />
    </>
  );
};
