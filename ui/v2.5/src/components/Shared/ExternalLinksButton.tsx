import { Button, Dropdown } from "react-bootstrap";
import { ExternalLink } from "./ExternalLink";
import TextUtils from "src/utils/text";
import { Icon } from "./Icon";
import { IconDefinition, faLink } from "@fortawesome/free-solid-svg-icons";
import { useMemo } from "react";
import { faInstagram, faTwitter } from "@fortawesome/free-brands-svg-icons";
import ReactDOM from "react-dom";

export const ExternalLinksButton: React.FC<{
  icon?: IconDefinition;
  urls: string[];
  className?: string;
}> = ({ urls, icon = faLink, className = "" }) => {
  if (!urls.length) {
    return null;
  }

  const Menu = () =>
    ReactDOM.createPortal(
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
      </Dropdown.Menu>,
      document.body
    );

  return (
    <Dropdown className="external-links-button">
      <Dropdown.Toggle as={Button} className={`minimal link ${className}`}>
        <Icon icon={icon} />
      </Dropdown.Toggle>
      <Menu />
    </Dropdown>
  );
};

export const ExternalLinkButtons: React.FC<{ urls: string[] | undefined }> = ({
  urls,
}) => {
  const urlSpecs = useMemo(() => {
    if (!urls?.length) {
      return [];
    }

    const twitter = urls.filter((u) =>
      u.match(/https?:\/\/(?:www\.)?(?:twitter|x).com\//)
    );
    const instagram = urls.filter((u) =>
      u.match(/https?:\/\/(?:www\.)?instagram.com\//)
    );
    const others = urls.filter(
      (u) => !twitter.includes(u) && !instagram.includes(u)
    );

    return [
      { icon: faLink, className: "", urls: others },
      { icon: faTwitter, className: "twitter", urls: twitter },
      { icon: faInstagram, className: "instagram", urls: instagram },
    ];
  }, [urls]);

  return (
    <>
      {urlSpecs.map((spec, i) => (
        <ExternalLinksButton key={i} {...spec} />
      ))}
    </>
  );
};
