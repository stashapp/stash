import React, { useState, useContext } from "react";
import { ListFilterModel } from "src/models/list-filter/filter";
import { useLocalForage } from "./LocalForage";

export interface IPlaylist {
  query?: ListFilterModel;
  sceneIDs?: number[];
}

export const usePlaylist = () =>
  useLocalForage<IPlaylist>("playlist");

export default usePlaylist;
