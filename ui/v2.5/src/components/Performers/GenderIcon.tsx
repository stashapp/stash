import React from "react";
import {
  faVenus,
  faTransgenderAlt,
  faMars,
  faVenusMars,
} from "@fortawesome/free-solid-svg-icons";
import * as GQL from "src/core/generated-graphql";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { genderToString } from "src/utils/gender";
import { useIntl } from "react-intl";

interface IconProps {
  gender?: GQL.Maybe<GQL.GenderEnum>;
  className?: string;
}

const GenderIcon: React.FC<IconProps> = ({ gender, className }) => {
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
        title={intl.formatMessage({ id: "gender." + gender })}
        className={className}
        icon={icon}
      />
    );
  }
  return <FontAwesomeIcon className={className} icon={faVenusMars} />;
};

export default GenderIcon;
