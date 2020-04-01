import React from "react";
import { Button, Form } from "react-bootstrap";

interface IImageInput {
  isEditing: boolean;
  text?: string;
  onImageChange: (event: React.FormEvent<HTMLInputElement>) => void;
}

export const ImageInput: React.FC<IImageInput> = ({
  isEditing,
  text,
  onImageChange,
}) => {
  if (!isEditing) return <div />;

  return (
    <Form.Label className="image-input ml-2">
      <Button variant="secondary">{text ?? "Browse for image..."}</Button>
      <Form.Control
        type="file"
        onChange={onImageChange}
        accept=".jpg,.jpeg,.png"
      />
    </Form.Label>
  );
};
