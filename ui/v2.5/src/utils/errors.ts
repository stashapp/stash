import { ApolloError } from "@apollo/client";

export const apolloError = (error: unknown) =>
  error instanceof ApolloError ? error.message : "";

export function errorToString(error: unknown) {
  let message;
  if (error instanceof Error) {
    message = error.message;
  }
  if (!message) {
    message = String(error);
  }
  if (!message) {
    message = "Unknown error";
  }

  return message;
}
