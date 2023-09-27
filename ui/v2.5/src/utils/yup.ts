import { IntlShape } from "react-intl";
import * as yup from "yup";

export function yupUniqueStringList(fieldName: string) {
  return yup
    .array(yup.string().required())
    .defined()
    .test({
      name: "unique",
      test: (value) => {
        const dupes = value
          .map((e, i, a) => {
            if (a.indexOf(e) !== i) {
              return String(i - 1);
            } else {
              return null;
            }
          })
          .filter((e) => e !== null) as string[];
        if (dupes.length === 0) return true;
        return new yup.ValidationError(dupes.join(" "), value, fieldName);
      },
    });
}

export function yupDateString(intl: IntlShape) {
  return yup
    .string()
    .ensure()
    .test({
      name: "date",
      test: (value) => {
        if (!value) return true;
        if (!value.match(/^\d{4}-\d{2}-\d{2}$/)) return false;
        if (Number.isNaN(Date.parse(value))) return false;
        return true;
      },
      message: intl.formatMessage({ id: "validation.date_invalid_form" }),
    });
}
