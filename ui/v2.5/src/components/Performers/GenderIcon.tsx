import React from "react";
import {
  faVenus,
  faTransgenderAlt,
  faMars,
} from "@fortawesome/free-solid-svg-icons";
import * as GQL from "src/core/generated-graphql";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { useIntl } from "react-intl";

interface IIconProps {
  gender?: GQL.Maybe<GQL.GenderEnum>;
  className?: string;
}

const GenderIcon: React.FC<IIconProps> = ({ gender, className }) => {
  const intl = useIntl();
  if (gender) {
    const icon =
      gender === GQL.GenderEnum.Male
        ? faMars
        : gender === GQL.GenderEnum.Female
        ? faVenus
        : faTransgenderAlt;
    return (
      <FontAwesomeIcon
        title={intl.formatMessage({ id: "gender_types." + gender })}
        className={className}
        icon={icon}
      />
    );
  }
  return null;
};

export default GenderIcon;
