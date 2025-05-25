import React, { useMemo } from "react";
import { Link } from "react-router-dom";
import { ConfigurationContext } from "src/hooks/Config";

interface IStudio {
  id: string;
  name: string;
  image_path?: string | null;
}

// Single studio overlay component (for backward compatibility)
export const StudioOverlay: React.FC<{
  studio: IStudio | null | undefined;
}> = ({ studio }) => {
  const { configuration } = React.useContext(ConfigurationContext);

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

// Multiple studios overlay component
export const StudiosOverlay: React.FC<{
  studios: IStudio[] | null | undefined;
}> = ({ studios }) => {
  const { configuration } = React.useContext(ConfigurationContext);

  const configValue = configuration?.interface.showStudioAsText;

  const shouldShowStudioAsText = (studio: IStudio) => {
    if (configValue || !studio?.image_path) {
      return true;
    }

    // If the studio has a default image, show the studio name as text
    const studioImageURL = new URL(studio.image_path);
    if (studioImageURL.searchParams.get("default") === "true") {
      return true;
    }

    return false;
  };

  if (!studios || studios.length === 0) return <></>;

  return (
    <div className="studios-overlay">
      {studios.map((studio) => (
        <div key={studio.id} className="studio-overlay">
          <Link to={`/studios/${studio.id}`}>
            {shouldShowStudioAsText(studio) ? (
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
      ))}
    </div>
  );
};
