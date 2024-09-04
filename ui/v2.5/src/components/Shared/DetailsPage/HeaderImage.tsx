import { PropsWithChildren } from "react";
import { LoadingIndicator } from "../LoadingIndicator";
import { FormattedMessage } from "react-intl";

export const HeaderImage: React.FC<
  PropsWithChildren<{
    encodingImage: boolean;
  }>
> = ({ encodingImage, children }) => {
  return (
    <div className="detail-header-image">
      {encodingImage ? (
        <LoadingIndicator
          message={<FormattedMessage id="actions.encoding_image" />}
        />
      ) : (
        children
      )}
    </div>
  );
};
