import { Button, Dropdown } from "react-bootstrap";
import { ExternalLink } from "./ExternalLink";
import TextUtils from "src/utils/text";
import { Icon } from "./Icon";
import { IconDefinition, faLink } from "@fortawesome/free-solid-svg-icons";

export const ExternalLinksButton: React.FC<{
  icon?: IconDefinition;
  urls: string[];
}> = ({ urls, icon = faLink }) => {
  if (!urls.length) {
    return null;
  }

  if (urls.length === 1) {
    return (
      <Button
        as={ExternalLink}
        href={TextUtils.sanitiseURL(urls[0])}
        className="minimal link external-links-button"
        title={urls[0]}
      >
        <Icon icon={icon} />
      </Button>
    );
  }

  return (
    <Dropdown className="external-links-button">
      <Dropdown.Toggle as={Button} className="minimal link">
        <Icon icon={icon} />
      </Dropdown.Toggle>

      <Dropdown.Menu>
        {urls.map((url) => (
          <Dropdown.Item
            key={url}
            as={ExternalLink}
            href={TextUtils.sanitiseURL(url)}
            title={url}
          >
            {url}
          </Dropdown.Item>
        ))}
      </Dropdown.Menu>
    </Dropdown>
  );
};
