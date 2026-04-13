import { Link } from "react-router-dom";
import { Studio } from "src/core/generated-graphql";
import { Icon } from "./Icon";
import { faVideo } from "@fortawesome/free-solid-svg-icons";

export const StudioLogo: React.FC<{
  studio: Pick<Studio, "id" | "image_path" | "name"> | undefined | null;
  showText?: boolean;
}> = ({ studio, showText = false }) => {
  if (!studio) return null;

  const hasLogo =
    !showText &&
    studio.image_path &&
    !studio.image_path.endsWith("default=true");

  return (
    <h1 className="text-center studio-logo">
      <Link to={`/studios/${studio.id}`}>
        {hasLogo ? (
          <img
            src={studio.image_path ?? ""}
            alt={`${studio.name} logo`}
            className="studio-logo"
          />
        ) : (
          <span className="studio-name">
            <Icon icon={faVideo} />
            {studio.name}
          </span>
        )}
      </Link>
    </h1>
  );
};
