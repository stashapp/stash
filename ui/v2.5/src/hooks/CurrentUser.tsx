import React from "react";
import * as GQL from "src/core/generated-graphql";

export type CurrentUser = NonNullable<GQL.MeQuery["me"]>;

export interface ICurrentUserContext {
  user: CurrentUser | null;
}

const CurrentUserContext = React.createContext<ICurrentUserContext>({
  user: null,
});

export const useCurrentUser = () => {
  return React.useContext(CurrentUserContext).user;
};

// Returns true if the current user is an administrator.
// Some fields, such as organized and ignore_auto_tag, may only be changed by administrators.
export const useIsAdmin = () => {
  const user = useCurrentUser();
  return user?.roles?.includes(GQL.RoleEnum.Admin) ?? false;
};

export const CurrentUserProvider: React.FC<
  React.PropsWithChildren<ICurrentUserContext>
> = ({ user, children }) => {
  return (
    <CurrentUserContext.Provider value={{ user }}>
      {children}
    </CurrentUserContext.Provider>
  );
};
