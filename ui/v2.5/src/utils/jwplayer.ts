const playerID = "main-jwplayer";
const getPlayer = () => (window as any).jwplayer(playerID);

export default {
  playerID,
  getPlayer
};
