import { FormikErrors, yupToFormErrors } from "formik";
import { IntlShape } from "react-intl";
import * as yup from "yup";

// equivalent to yup.array(yup.string().required())
// except that error messages will be e.g.
// 'urls must not be blank' instead of
// 'urls["0"] is a required field'
export function yupRequiredStringArray(intl: IntlShape) {
  return yup
    .array(
      // we enforce that each string in the array is "required" in the outer test function
      // so cast to avoid having to add a redundant `.required()` here
      yup.string() as yup.StringSchema<string>
    )
    .test({
      name: "blank",
      test(value) {
        if (!value || !value.length) return true;

        const blanks: number[] = [];
        for (let i = 0; i < value.length; i++) {
          const s = value[i];
          if (!s) {
            blanks.push(i);
          }
        }
        if (blanks.length === 0) return true;

        // each error message is identical
        const msg = yup.ValidationError.formatError(
          intl.formatMessage({ id: "validation.blank" }),
          {
            label: this.schema.spec.label,
            path: this.path,
          }
        );

        // return multiple errors, one for each blank string
        const errors = blanks.map(
          (i) =>
            new yup.ValidationError(
              msg,
              value[i],
              // the path to this "sub-error": e.g. 'urls["0"]'
              `${this.path}["${i}"]`,
              "blank"
            )
        );

        return new yup.ValidationError(errors, value, this.path, "blank");
      },
    });
}

export function yupUniqueStringList(intl: IntlShape) {
  return yupRequiredStringArray(intl)
    .defined()
    .test({
      name: "unique",
      test(value) {
        const values: string[] = [];
        const dupes: number[] = [];
        for (let i = 0; i < value.length; i++) {
          const s = value[i];
          if (values.includes(s)) {
            dupes.push(i);
          } else {
            values.push(s);
          }
        }
        if (dupes.length === 0) return true;

        const msg = yup.ValidationError.formatError(
          intl.formatMessage({ id: "validation.unique" }),
          {
            label: this.schema.spec.label,
            path: this.path,
          }
        );
        const errors = dupes.map(
          (i) =>
            new yup.ValidationError(
              msg,
              value[i],
              `${this.path}["${i}"]`,
              "unique"
            )
        );
        return new yup.ValidationError(errors, value, this.path, "unique");
      },
    });
}

export function yupUniqueAliases(intl: IntlShape, nameField: string) {
  return yupRequiredStringArray(intl)
    .defined()
    .test({
      name: "unique",
      test(value) {
        const aliases = [this.parent[nameField].toLowerCase()];
        const dupes: number[] = [];
        for (let i = 0; i < value.length; i++) {
          const s = value[i].toLowerCase();
          if (aliases.includes(s)) {
            dupes.push(i);
          } else {
            aliases.push(s);
          }
        }
        if (dupes.length === 0) return true;

        const msg = yup.ValidationError.formatError(
          intl.formatMessage({ id: "validation.unique" }),
          {
            label: this.schema.spec.label,
            path: this.path,
          }
        );
        const errors = dupes.map(
          (i) =>
            new yup.ValidationError(
              msg,
              value[i],
              `${this.path}["${i}"]`,
              "unique"
            )
        );
        return new yup.ValidationError(errors, value, this.path, "unique");
      },
    });
}

export function yupDateString(intl: IntlShape) {
  return yup
    .string()
    .ensure()
    .test({
      name: "date",
      test(value) {
        if (!value) return true;
        // Allow YYYY, YYYY-MM, or YYYY-MM-DD formats
        if (!value.match(/^\d{4}(-\d{2}(-\d{2})?)?$/)) return false;
        // Validate the date components
        const parts = value.split("-");
        const year = parseInt(parts[0], 10);
        if (year < 1 || year > 9999) return false;
        if (parts.length >= 2) {
          const month = parseInt(parts[1], 10);
          if (month < 1 || month > 12) return false;
        }
        if (parts.length === 3) {
          const day = parseInt(parts[2], 10);
          if (day < 1 || day > 31) return false;
          // Full date - validate it parses correctly
          if (Number.isNaN(Date.parse(value))) return false;
        }
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
