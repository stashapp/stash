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
        if (!value?.length) return true;

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

export function normalizeDateString(value?: string) {
  if (!value) return undefined;

  let year: number;
  let month: number | undefined;
  let day: number | undefined;

  const yearMatch = /^(\d{4})$/.exec(value);
  const yearMonthMatch = /^(\d{4})([-.])(\d{1,2})$/.exec(value);
  const fullDateMatch = /^(\d{4})([-.])(\d{1,2})\2(\d{1,2})$/.exec(value);
  const compactDateMatch = /^(\d{4})(\d{2})(\d{2})$/.exec(value);
  const shortDateMatch = /^(\d{2})([-.])(\d{1,2})\2(\d{1,2})$/.exec(value);

  if (yearMatch) {
    year = Number(yearMatch[1]);
  } else if (yearMonthMatch) {
    year = Number(yearMonthMatch[1]);
    month = Number(yearMonthMatch[3]);
  } else if (fullDateMatch) {
    year = Number(fullDateMatch[1]);
    month = Number(fullDateMatch[3]);
    day = Number(fullDateMatch[4]);
  } else if (compactDateMatch) {
    year = Number(compactDateMatch[1]);
    month = Number(compactDateMatch[2]);
    day = Number(compactDateMatch[3]);
  } else if (shortDateMatch) {
    const shortYear = Number(shortDateMatch[1]);
    year = shortYear <= 68 ? 2000 + shortYear : 1900 + shortYear;
    month = Number(shortDateMatch[3]);
    day = Number(shortDateMatch[4]);
  } else {
    return undefined;
  }

  if (year < 1 || year > 9999) return undefined;
  if (month === undefined) return `${year}`;
  if (month < 1 || month > 12) return undefined;
  if (day === undefined) return `${year}-${month.toString().padStart(2, "0")}`;

  const daysInMonth = [
    31,
    isLeapYear(year) ? 29 : 28,
    31,
    30,
    31,
    30,
    31,
    31,
    30,
    31,
    30,
    31,
  ][month - 1];
  if (day < 1 || day > daysInMonth) return undefined;

  return `${year}-${month.toString().padStart(2, "0")}-${day
    .toString()
    .padStart(2, "0")}`;
}

export function validateDateString(value?: string) {
  if (!value) return true;
  return normalizeDateString(value) !== undefined;
}

function isLeapYear(year: number) {
  return year % 4 === 0 && (year % 100 !== 0 || year % 400 === 0);
}

export function getDateError(
  value: string | undefined | null,
  intl: IntlShape
) {
  if (validateDateString(value ?? "")) return undefined;
  return (
    intl
      .formatMessage({ id: "validation.date_invalid_form" })
      // biome-ignore lint/suspicious/noTemplateCurlyInString: required for intl
      .replace("${path}", intl.formatMessage({ id: "date" }))
  );
}

export function yupDateString(intl: IntlShape) {
  return yup
    .string()
    .ensure()
    .test({
      name: "date",
      test(value) {
        return validateDateString(value);
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
  return async (values) => {
    try {
      await schema.validate(values, { abortEarly: false });
    } catch (err) {
      return yupToFormErrors(err);
    }
    return {};
  };
}
