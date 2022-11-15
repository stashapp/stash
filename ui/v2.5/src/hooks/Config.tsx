import React from "react";
import * as GQL from "src/core/generated-graphql";

export interface IContext {
  configuration?: GQL.ConfigDataFragment;
  loading?: boolean;
}

export const ConfigurationContext = React.createContext<IContext>({});

export const ConfigurationProvider: React.FC<IContext> = ({
  loading,
  configuration,
  children,
}) => {
  return (
    <ConfigurationContext.Provider
      value={{
        configuration,
        loading,
      }}
    >
      {children}
    </ConfigurationContext.Provider>
  );
};
