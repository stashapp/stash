import { PropsWithChildren } from "react";
import { LoadingIndicator } from "../LoadingIndicator";
import { FormattedMessage } from "react-intl";
import { PatchComponent } from "src/patch";

export const HeaderImage: React.FC<
  PropsWithChildren<{
    encodingImage: boolean;
  }>
> = PatchComponent("HeaderImage", ({ encodingImage, children }) => {
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
});
