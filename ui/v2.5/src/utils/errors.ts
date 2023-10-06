import { ApolloError } from "@apollo/client";

export const apolloError = (error: unknown) =>
  error instanceof ApolloError ? error.message : "";
