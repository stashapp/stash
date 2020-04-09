/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

//==============================================================
// START Enums and Input Objects
//==============================================================

export enum BreastTypeEnum {
  FAKE = "FAKE",
  NA = "NA",
  NATURAL = "NATURAL",
}

export enum DateAccuracyEnum {
  DAY = "DAY",
  MONTH = "MONTH",
  YEAR = "YEAR",
}

export enum EthnicityEnum {
  ASIAN = "ASIAN",
  BLACK = "BLACK",
  CAUCASIAN = "CAUCASIAN",
  INDIAN = "INDIAN",
  LATIN = "LATIN",
  MIDDLE_EASTERN = "MIDDLE_EASTERN",
  MIXED = "MIXED",
  OTHER = "OTHER",
}

export enum EyeColorEnum {
  BLUE = "BLUE",
  BROWN = "BROWN",
  GREEN = "GREEN",
  GREY = "GREY",
  HAZEL = "HAZEL",
  RED = "RED",
}

export enum FingerprintAlgorithm {
  MD5 = "MD5",
}

export enum GenderEnum {
  FEMALE = "FEMALE",
  INTERSEX = "INTERSEX",
  MALE = "MALE",
  TRANSGENDER_FEMALE = "TRANSGENDER_FEMALE",
  TRANSGENDER_MALE = "TRANSGENDER_MALE",
}

export enum HairColorEnum {
  AUBURN = "AUBURN",
  BALD = "BALD",
  BLACK = "BLACK",
  BLONDE = "BLONDE",
  BRUNETTE = "BRUNETTE",
  GREY = "GREY",
  OTHER = "OTHER",
  RED = "RED",
  VARIOUS = "VARIOUS",
}

export interface FingerprintInput {
  hash: string;
  algorithm: FingerprintAlgorithm;
  duration: number;
}

export interface FingerprintQueryInput {
  hash: string;
  algorithm: FingerprintAlgorithm;
}

export interface FingerprintSubmission {
  scene_id: string;
  fingerprint: FingerprintInput;
}

//==============================================================
// END Enums and Input Objects
//==============================================================
