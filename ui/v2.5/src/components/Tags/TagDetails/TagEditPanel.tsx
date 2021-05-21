import React, { useEffect } from "react";
import * as GQL from "src/core/generated-graphql";
import * as yup from "yup";
import {
  DetailsEditNavbar,
} from "src/components/Shared";
import { useToast } from "src/hooks";
import { Form, Col, Row } from "react-bootstrap";
import { ImageUtils } from "src/utils";
import { useFormik } from "formik";
import { Prompt } from "react-router-dom";
import Mousetrap from "mousetrap";

interface ITagEditPanel {
  tag?: Partial<GQL.TagDataFragment>;
  onSubmit: (
    movie: Partial<GQL.TagCreateInput | GQL.TagUpdateInput>
  ) => void;
  onCancel: () => void;
  onDelete: () => void;
  setImage: (image?: string | null) => void;
}

export const TagEditPanel: React.FC<ITagEditPanel> = ({
  tag,
  onSubmit,
  onCancel,
  onDelete,
  setImage,
}) => {
  const Toast = useToast();

  const isNew = tag === undefined;

  const labelXS = 3;
  const labelXL = 3;
  const fieldXS = 9;
  const fieldXL = 9;

  const schema = yup.object({
    name: yup.string().required(),
    aliases: yup.array(yup.string().required()).required(),
  });

  const initialValues = {
    name: tag?.name,
    aliases: tag?.aliases,
  };

  type InputValues = typeof initialValues;

  const formik = useFormik({
    initialValues,
    validationSchema: schema,
    onSubmit: (values) => onSubmit(getTagInput(values)),
  });

  // set up hotkeys
  useEffect(() => {
    Mousetrap.bind("s s", () => formik.handleSubmit());

    return () => {
      Mousetrap.unbind("s s");
    };
  });

  function getTagInput(values: InputValues) {
    const input: Partial<GQL.TagCreateInput | GQL.TagUpdateInput> = {
      ...values,
    };

    if (tag && tag.id) {
      (input as GQL.TagUpdateInput).id = tag.id;
    }
    return input;
  }

  function onImageChange(event: React.FormEvent<HTMLInputElement>) {
    ImageUtils.onImageChange(event, setImage);
  }

  const isEditing = true;

  function renderTextField(field: string, title: string) {
    return (
      <Form.Group controlId={field} as={Row}>
        <Form.Label column xs={labelXS} xl={labelXL}>
          {title}
        </Form.Label>
        <Col xs={fieldXS} xl={fieldXL}>
          <Form.Control
            className="text-input"
            placeholder={title}
            {...formik.getFieldProps(field)}
            isInvalid={!!formik.getFieldMeta(field).error}
          />
        </Col>
      </Form.Group>
    );
  }

  // TODO: CSS class
  return (
    <div>
      {isNew && <h2>Add Tag</h2>}

      <Prompt
        when={formik.dirty}
        message="Unsaved changes. Are you sure you want to leave?"
      />

      <Form noValidate onSubmit={formik.handleSubmit} id="tag-edit">
        <Form.Group controlId="name" as={Row}>
          <Form.Label column xs={labelXS} xl={labelXL}>
            Name
          </Form.Label>
          <Col xs={fieldXS} xl={fieldXL}>
            <Form.Control
              className="text-input"
              placeholder="Name"
              {...formik.getFieldProps("name")}
              isInvalid={!!formik.errors.name}
            />
            <Form.Control.Feedback type="invalid">
              {formik.errors.name}
            </Form.Control.Feedback>
          </Col>
        </Form.Group>

        {/* {renderTextField("aliases", "Aliases")} */}

      </Form>

      <DetailsEditNavbar
        objectName={tag?.name ?? "tag"}
        isNew={isNew}
        isEditing={isEditing}
        onToggleEdit={onCancel}
        onSave={() => formik.handleSubmit()}
        onImageChange={onImageChange}
        onImageChangeURL={setImage}
        onClearImage={() => {
          setImage(null);
        }}
        onDelete={onDelete}
      />
    </div>
  );
};
