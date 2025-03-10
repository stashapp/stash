import React from "react";
import { FormattedNumber } from "react-intl";
import TextUtils from "src/utils/text";

export const FileSize: React.FC<{ size: number }> = ({ size: fileSize }) => {
  const { size, unit } = TextUtils.fileSize(fileSize);

  return (
    <>
      <FormattedNumber
        value={size}
        maximumFractionDigits={TextUtils.fileSizeFractionalDigits(unit)}
      />
      {` ${TextUtils.formatFileSizeUnit(unit)}`}
    </>
  );
};
