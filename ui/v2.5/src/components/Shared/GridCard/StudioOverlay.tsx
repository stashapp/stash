import React, { useMemo } from "react";
import { Link } from "react-router-dom";
import { useConfigurationContext } from "src/hooks/Config";

interface IStudio {
  id: string;
  name: string;
  image_path?: string | null;
}

export const StudioOverlay: React.FC<{
  studio: IStudio | null | undefined;
}> = ({ studio }) => {
  const { configuration } = useConfigurationContext();

  const configValue = configuration?.interface.showStudioAsText;

  const showStudioAsText = useMemo(() => {
    if (configValue || !studio?.image_path) {
      return true;
    }

    // If the studio has a default image, show the studio name as text
    const studioImageURL = new URL(studio.image_path);
    if (studioImageURL.searchParams.get("default") === "true") {
      return true;
    }

    return false;
  }, [configValue, studio?.image_path]);

  if (!studio) return <></>;

  return (
    // this class name is incorrect
    <div className="studio-overlay">
      <Link to={`/studios/${studio.id}`}>
        {showStudioAsText ? (
          studio.name
        ) : (
          <img
            className="image-thumbnail"
            loading="lazy"
            alt={studio.name}
            src={studio.image_path ?? ""}
          />
        )}
      </Link>
    </div>
  );
};
