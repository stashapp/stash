import React from "react";
import {
  faVenus,
  faTransgenderAlt,
  faMars,
  faNonBinary,
} from "@fortawesome/free-solid-svg-icons";
import * as GQL from "src/core/generated-graphql";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { useIntl } from "react-intl";

interface IIconProps {
  gender?: GQL.Maybe<GQL.GenderEnum>;
  className?: string;
}

function genderIcon(gender: GQL.GenderEnum) {
  switch (gender) {
    case GQL.GenderEnum.Male:
      return faMars;
    case GQL.GenderEnum.Female:
      return faVenus;
    case GQL.GenderEnum.NonBinary:
      return faNonBinary;
    default:
      return faTransgenderAlt;
  }
}

const GenderIcon: React.FC<IIconProps> = ({ gender, className }) => {
  const intl = useIntl();
  if (gender) {
    const icon = genderIcon(gender);

    // new version of fontawesome doesn't seem to support titles on icons, so adding it
    // to a span instead
    return (
      <span title={intl.formatMessage({ id: "gender_types." + gender })}>
        <FontAwesomeIcon
          data-gender={gender}
          className={className}
          icon={icon}
        />
      </span>
    );
  }
  return null;
};

export default GenderIcon;
