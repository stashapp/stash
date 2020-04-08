import Cookies from "universal-cookie";

const isLoggedIn = () => {
  return new Cookies().get("session") !== undefined;
};

export default {
  isLoggedIn,
};
