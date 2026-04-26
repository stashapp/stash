import React from "react";
import { FormattedMessage } from "react-intl";
import { Link } from "react-router-dom";
import { ExternalLink } from "src/components/Shared/ExternalLink";
import { Icon } from "src/components/Shared/Icon";
import { TruncatedText } from "src/components/Shared/TruncatedText";
import { FontAwesomeIconProps } from "@fortawesome/react-fontawesome";

interface ITextField {
  id?: string;
  name?: string;
  abbr?: string | null;
  value?: string | null;
  truncate?: boolean | null;
}

export const TextField: React.FC<ITextField> = ({
  id,
  name,
  value,
  abbr,
  truncate,
  children,
}) => {
  if (!value && !children) {
    return null;
  }

  const message = (
    <>{id ? <FormattedMessage id={id} defaultMessage={name} /> : name}:</>
  );

  return (
    <>
      <dt>{abbr ? <abbr title={abbr}>{message}</abbr> : message}</dt>
      <dd>
        {value ? truncate ? <TruncatedText text={value} /> : value : children}
      </dd>
    </>
  );
};

interface IURLField {
  id?: string;
  name?: string;
  abbr?: string | null;
  value?: string | null;
  url?: string | null;
  truncate?: boolean | null;
  target?: string;
  // an internal link (uses <Link to={url}>)
  internal?: boolean;
}

export const URLField: React.FC<IURLField> = ({
  id,
  name,
  value,
  url,
  abbr,
  truncate,
  target = "_blank",
  internal,
}) => {
  if (!value) {
    return null;
  }

  const message = (
    <>{id ? <FormattedMessage id={id} defaultMessage={name} /> : name}:</>
  );

  function maybeRenderUrl() {
    if (!url) return;

    const children = truncate ? <TruncatedText text={value} /> : value;

    if (internal) {
      return (
        <Link to={url} target={target}>
          {children}
        </Link>
      );
    } else {
      return (
        <ExternalLink href={url} target={target}>
          {children}
        </ExternalLink>
      );
    }
  }

  return (
    <>
      <dt>{abbr ? <abbr title={abbr}>{message}</abbr> : message}</dt>
      <dd>{maybeRenderUrl()}</dd>
    </>
  );
};

interface IURLsField {
  id?: string;
  name?: string;
  abbr?: string | null;
  urls?: string[] | null;
  truncate?: boolean | null;
  target?: string;
  // an internal link (uses <Link to={url}>)
  internal?: boolean;
}

export const URLsField: React.FC<IURLsField> = ({
  id,
  name,
  urls,
  abbr,
  truncate,
  target = "_blank",
  internal,
}) => {
  if (!urls || !urls.length) {
    return null;
  }

  const message = (
    <>{id ? <FormattedMessage id={id} defaultMessage={name} /> : name}:</>
  );

  const renderUrls = () => {
    return urls.map((url, i) => {
      if (!url) return;

      const children = truncate ? <TruncatedText text={url} /> : url;

      if (internal) {
        return (
          <Link key={i} to={url} target={target}>
            {children}
          </Link>
        );
      } else {
        return (
          <ExternalLink key={i} href={url} target={target}>
            {children}
          </ExternalLink>
        );
      }
    });
  };

  return (
    <>
      <dt>{abbr ? <abbr title={abbr}>{message}</abbr> : message}</dt>
      <dd>
        <dl>{renderUrls()}</dl>
      </dd>
    </>
  );
};

interface IIconField {
  id?: string;
  name?: string;
  abbr?: string | null;
  icon?: FontAwesomeIconProps["icon"];
  tooltip?: string;
  value?: string | null;
  url?: string | null;
  stashEndpoint?: string | null;
  truncate?: boolean | null;
  target?: string;
}

export const IconField: React.FC<IIconField> = ({
  id,
  name,
  abbr,
  icon,
  tooltip,
  value,
  url,
  stashEndpoint,
  truncate,
  target = "_blank",
}) => {
  if (!value && !url && !stashEndpoint) {
    return null;
  }

  const message = (
    <>{id ? <FormattedMessage id={id} defaultMessage={name} /> : name}:</>
  );

  const displayValue = value ?? url ?? stashEndpoint ?? "";

  function renderValue() {
    if (stashEndpoint) {
      return (
        <span className="stash-id-pill" data-endpoint={stashEndpoint}>
          <span>{displayValue}</span>
        </span>
      );
    }

    if (url) {
      const children = truncate ? (
        <TruncatedText
          text={displayValue}
          className="d-inline-block align-middle"
        />
      ) : (
        displayValue
      );
      return (
        <ExternalLink href={url} target={target} title={url} className="flex-grow-1">
          {children}
        </ExternalLink>
      );
    }

    return displayValue;
  }

  return (
    <>
      <dt>{abbr ? <abbr title={abbr}>{message}</abbr> : message}</dt>
      <dd>
        <span
          title={tooltip}
          className={
            url ? "d-inline-flex align-items-center w-100" : "d-inline-flex align-items-center"
          }
          style={url ? { minWidth: 0 } : undefined}
        >
          {icon && <Icon icon={icon} className={url ? "mr-1 flex-shrink-0" : "mr-1"} />}
          {renderValue()}
        </span>
      </dd>
    </>
  );
};
