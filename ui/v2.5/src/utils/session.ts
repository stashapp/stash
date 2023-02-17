import Cookies from "universal-cookie";

const isLoggedIn = () => {
  return new Cookies().get("session") !== undefined;
};

const SessionUtils = {
  isLoggedIn,
};

export default SessionUtils;
