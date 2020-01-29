import React from "react";
import { Button, Form } from "react-bootstrap";

interface IImageInput {
  isEditing: boolean;
  onImageChange: (event: React.FormEvent<HTMLInputElement>) => void;
}

export const ImageInput: React.FC<IImageInput> = ({
  isEditing,
  onImageChange
}) => {
  if (!isEditing) return <div />;

  return (
    <Form.Label className="image-input">
      <Button variant="secondary">Browse for image...</Button>
      <Form.Control
        type="file"
        onChange={onImageChange}
        accept=".jpg,.jpeg,.png"
      />
    </Form.Label>
  );
};
