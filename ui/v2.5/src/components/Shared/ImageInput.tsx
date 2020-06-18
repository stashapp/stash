import React from "react";
import { Button, Form } from "react-bootstrap";

interface IImageInput {
  isEditing: boolean;
  text?: string;
  onImageChange: (event: React.ChangeEvent<HTMLInputElement>) => void;
  acceptSVG?: boolean;
}

export const ImageInput: React.FC<IImageInput> = ({
  isEditing,
  text,
  onImageChange,
  acceptSVG = false,
}) => {
  if (!isEditing) return <div />;

  return (
    <Form.Label className="image-input">
      <Button variant="secondary">{text ?? "Browse for image..."}</Button>
      <Form.Control
        type="file"
        onChange={onImageChange}
        accept={`.jpg,.jpeg,.png${acceptSVG ? ",.svg" : ""}`}
      />
    </Form.Label>
  );
};
