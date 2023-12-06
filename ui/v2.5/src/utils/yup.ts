import { FormikErrors, yupToFormErrors } from "formik";
import { IntlShape } from "react-intl";
import * as yup from "yup";

export function yupUniqueStringList(fieldName: string) {
  return yup
    .array(yup.string().required())
    .defined()
    .test({
      name: "unique",
      test: (value) => {
        const values: string[] = [];
        const dupes: number[] = [];
        for (let i = 0; i < value.length; i++) {
          const a = value[i];
          if (values.includes(a)) {
            dupes.push(i);
          } else {
            values.push(a);
          }
        }
        if (dupes.length === 0) return true;
        return new yup.ValidationError(dupes.join(" "), value, fieldName);
      },
    });
}

export function yupUniqueAliases(fieldName: string, nameField: string) {
  return yup
    .array(yup.string().required())
    .defined()
    .test({
      name: "unique",
      test: (value, context) => {
        const aliases = [context.parent[nameField].toLowerCase()];
        const dupes: number[] = [];
        for (let i = 0; i < value.length; i++) {
          const a = value[i].toLowerCase();
          if (aliases.includes(a)) {
            dupes.push(i);
          } else {
            aliases.push(a);
          }
        }
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

type StringEnum<T extends string> = {
  [k: string]: T;
};

// Use yupInputEnum to validate a string enum from a <select>.
// If "" is not a value in the enum, a "" input will be transformed to null.
export function yupInputEnum<T extends string>(e: StringEnum<T>) {
  const enumValues = Object.values(e);
  const schema = yup.string<T>().oneOf(enumValues);
  if (enumValues.includes("" as T)) {
    return schema;
  } else {
    return schema.transform((v, o) => (o === "" ? null : v));
  }
}

// Use yupInputNumber to validate a number from an <input type="number">.
// A "" input will be transformed to null.
export function yupInputNumber() {
  return yup.number().transform((v, o) => (o === "" ? null : v));
}

// Formik converts "" into undefined when validating with a yup schema,
// which prevents transformations from running.
// Interfacing with yup ourselves avoids this.
// https://github.com/jaredpalmer/formik/pull/2902#issuecomment-922492137
export function yupFormikValidate<T>(
  schema: yup.AnySchema
): (values: T) => Promise<FormikErrors<T>> {
  return async function (values) {
    try {
      await schema.validate(values, { abortEarly: false });
    } catch (err) {
      return yupToFormErrors(err);
    }
    return {};
  };
}
