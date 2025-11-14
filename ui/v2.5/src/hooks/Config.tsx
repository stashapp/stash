import React from "react";
import * as GQL from "src/core/generated-graphql";

export interface IContext {
  configuration: GQL.ConfigDataFragment;
}

export const ConfigurationContext = React.createContext<IContext | null>(null);

export const useConfigurationContext = () => {
  const context = React.useContext(ConfigurationContext);

  if (context === null) {
    throw new Error(
      "useConfigurationContext must be used within a ConfigurationProvider"
    );
  }

  return context;
};

export const useConfigurationContextOptional = () => {
  return React.useContext(ConfigurationContext);
};

export const ConfigurationProvider: React.FC<IContext> = ({
  configuration,
  children,
}) => {
  return (
    <ConfigurationContext.Provider
      value={{
        configuration,
      }}
    >
      {children}
    </ConfigurationContext.Provider>
  );
};
