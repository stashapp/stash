import React from "react";
import * as GQL from "src/core/generated-graphql";

interface IContext {
  configuration?: GQL.ConfigDataFragment;
  loading?: boolean;
  isTouch: boolean;
}

export const ConfigurationContext = React.createContext<IContext>({
  isTouch: false,
});

export const ConfigurationProvider: React.FC<IContext> = ({
  loading,
  configuration,
  isTouch,
  children,
}) => {
  return (
    <ConfigurationContext.Provider
      value={{
        configuration,
        loading,
        isTouch,
      }}
    >
      {children}
    </ConfigurationContext.Provider>
  );
};
