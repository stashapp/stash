const playerID = "main-jwplayer";
// eslint-disable-next-line @typescript-eslint/no-explicit-any
const getPlayer = () => (window as any).jwplayer(playerID);

// eslint-disable-next-line @typescript-eslint/no-explicit-any

export default {
  playerID,
  getPlayer,
};
