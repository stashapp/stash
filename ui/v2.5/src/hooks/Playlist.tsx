import React, { useState, useContext } from "react";
import { ListFilterModel } from "src/models/list-filter/filter";

export interface IPlaylist {
  query?: ListFilterModel;
  sceneIDs?: [number];
}

interface IPlaylistContextData {
  playlist: IPlaylist;
  setPlaylist: (p: IPlaylist) => void;
}

const defaultPlaylistContextData: IPlaylistContextData = {
  playlist: {},
  setPlaylist: () => null,
};

const PlaylistContext = React.createContext<IPlaylistContextData>(
  defaultPlaylistContextData
);

export const PlaylistProvider: React.FC = ({ children }) => {
  const [playlist, setPlaylist] = useState<IPlaylist>({});

  return (
    <PlaylistContext.Provider
      value={{
        playlist,
        setPlaylist,
      }}
    >
      {children}
    </PlaylistContext.Provider>
  );
};

const usePlaylist = () => {
  return useContext(PlaylistContext);
};

export default usePlaylist;
