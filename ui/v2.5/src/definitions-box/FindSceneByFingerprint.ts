/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

import { FingerprintQueryInput, GenderEnum, DateAccuracyEnum, EthnicityEnum, EyeColorEnum, HairColorEnum, BreastTypeEnum, FingerprintAlgorithm } from "./globalTypes";

// ====================================================
// GraphQL query operation: FindSceneByFingerprint
// ====================================================

export interface FindSceneByFingerprint_findSceneByFingerprint_urls {
  url: string;
  type: string;
  image_id: string | null;
  width: number | null;
  height: number | null;
}

export interface FindSceneByFingerprint_findSceneByFingerprint_studio_urls {
  url: string;
  type: string;
  image_id: string | null;
  width: number | null;
  height: number | null;
}

export interface FindSceneByFingerprint_findSceneByFingerprint_studio {
  name: string;
  id: string;
  urls: (FindSceneByFingerprint_findSceneByFingerprint_studio_urls | null)[];
}

export interface FindSceneByFingerprint_findSceneByFingerprint_tags {
  name: string;
  id: string;
}

export interface FindSceneByFingerprint_findSceneByFingerprint_performers_performer_urls {
  url: string;
  type: string;
  image_id: string | null;
  width: number | null;
  height: number | null;
}

export interface FindSceneByFingerprint_findSceneByFingerprint_performers_performer_birthdate {
  date: any;
  accuracy: DateAccuracyEnum;
}

export interface FindSceneByFingerprint_findSceneByFingerprint_performers_performer_measurements {
  band_size: number | null;
  cup_size: string | null;
  waist: number | null;
  hip: number | null;
}

export interface FindSceneByFingerprint_findSceneByFingerprint_performers_performer_tattoos {
  location: string;
  description: string | null;
}

export interface FindSceneByFingerprint_findSceneByFingerprint_performers_performer_piercings {
  location: string;
  description: string | null;
}

export interface FindSceneByFingerprint_findSceneByFingerprint_performers_performer {
  id: string;
  name: string;
  disambiguation: string | null;
  aliases: string[];
  gender: GenderEnum | null;
  urls: FindSceneByFingerprint_findSceneByFingerprint_performers_performer_urls[];
  birthdate: FindSceneByFingerprint_findSceneByFingerprint_performers_performer_birthdate | null;
  ethnicity: EthnicityEnum | null;
  country: string | null;
  eye_color: EyeColorEnum | null;
  hair_color: HairColorEnum | null;
  /**
   * Height in cm
   */
  height: number | null;
  measurements: FindSceneByFingerprint_findSceneByFingerprint_performers_performer_measurements;
  breast_type: BreastTypeEnum | null;
  career_start_year: number | null;
  career_end_year: number | null;
  tattoos: FindSceneByFingerprint_findSceneByFingerprint_performers_performer_tattoos[] | null;
  piercings: FindSceneByFingerprint_findSceneByFingerprint_performers_performer_piercings[] | null;
}

export interface FindSceneByFingerprint_findSceneByFingerprint_performers {
  /**
   * Performing as alias
   */
  as: string | null;
  performer: FindSceneByFingerprint_findSceneByFingerprint_performers_performer;
}

export interface FindSceneByFingerprint_findSceneByFingerprint_fingerprints {
  algorithm: FingerprintAlgorithm;
  hash: string;
  duration: number;
}

export interface FindSceneByFingerprint_findSceneByFingerprint {
  id: string;
  title: string | null;
  details: string | null;
  duration: number | null;
  date: any | null;
  urls: FindSceneByFingerprint_findSceneByFingerprint_urls[];
  studio: FindSceneByFingerprint_findSceneByFingerprint_studio | null;
  tags: FindSceneByFingerprint_findSceneByFingerprint_tags[];
  performers: FindSceneByFingerprint_findSceneByFingerprint_performers[];
  fingerprints: FindSceneByFingerprint_findSceneByFingerprint_fingerprints[];
}

export interface FindSceneByFingerprint {
  /**
   * Finds a scene by an algorithm-specific checksum
   */
  findSceneByFingerprint: FindSceneByFingerprint_findSceneByFingerprint[];
}

export interface FindSceneByFingerprintVariables {
  fingerprint: FingerprintQueryInput;
}
