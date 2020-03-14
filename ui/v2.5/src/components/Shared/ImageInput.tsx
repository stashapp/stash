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
  onImageChange
}) => {
  if (!isEditing) return <div />;

  text = text ?? "Browse for image...";

  return (
    <Form.Label className="image-input">
      <Button variant="secondary">{text}</Button>
      <Form.Control
        type="file"
        onChange={onImageChange}
        accept=".jpg,.jpeg,.png"
      />
    </Form.Label>
  );
};
