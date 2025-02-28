import React from "react";

export const BackgroundImage: React.FC<{
  imagePath: string | undefined;
  show: boolean;
  alt?: string;
}> = ({ imagePath, show, alt }) => {
  if (imagePath && show) {
    const imageURL = new URL(imagePath);
    let isDefaultImage = imageURL.searchParams.get("default");
    if (!isDefaultImage) {
      return (
        <div className="background-image-container">
          <picture>
            <source src={imagePath} />
            <img className="background-image" src={imagePath} alt={alt} />
          </picture>
        </div>
      );
    }
  }

  return null;
};
