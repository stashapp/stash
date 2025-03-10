import React from "react";
import { FormattedMessage } from "react-intl";
import cx from "classnames";
import { Icon } from "./Icon";
import { faCaretDown, faCaretUp } from "@fortawesome/free-solid-svg-icons";

export function maybeRenderShowMoreLess(
  height: number,
  limit: number,
  ref: React.MutableRefObject<HTMLDivElement | null>,
  setCollapsed: React.Dispatch<React.SetStateAction<boolean>>,
  collapsed: boolean
) {
  if (height < limit) {
    return;
  }
  return (
    <span
      className={`show-${collapsed ? "more" : "less"}`}
      onClick={() => {
        const container = ref.current;
        if (container == null) {
          return;
        }
        if (container.style.maxHeight) {
          container.style.maxHeight = "";
        } else {
          container.style.maxHeight = container.scrollHeight + "px";
        }
        setCollapsed(!collapsed);
      }}
    >
      {collapsed ? "Show more" : "Show less"}
      <Icon className="fa-solid" icon={collapsed ? faCaretDown : faCaretUp} />
    </span>
  );
}

interface IDetailItem {
  id?: string | null;
  label?: React.ReactNode;
  messageId?: string;
  heading?: React.ReactNode;
  value?: React.ReactNode;
  labelTitle?: string;
  title?: string;
  fullWidth?: boolean;
  showEmpty?: boolean;
}

export const DetailItem: React.FC<IDetailItem> = ({
  id,
  label,
  messageId,
  heading,
  value,
  labelTitle,
  title,
  fullWidth,
  showEmpty = false,
}) => {
  if (!id || (!showEmpty && (!value || value === "Na"))) {
    return <></>;
  }

  const message = label ?? <FormattedMessage id={messageId ?? id} />;

  // according to linter rule CSS classes shouldn't use underscores
  const sanitisedID = id.replace(/_/g, "-");


  return (
    <div className={cx(`detail-item ${sanitisedID}`, { "full-width": fullWidth })}>
      <span className={`detail-item-title ${sanitisedID}`} title={labelTitle}>
        {heading ? (
          heading
        ) : (
          <>
            {message}
            {fullWidth ? ":" : ""}
          </>
        )}
      </span>
      <span className={`detail-item-value ${sanitisedID}`} title={title}>
        {value}
      </span>
    </div>
  );
};
