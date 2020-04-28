import {
  SearchScene_searchScene_performers_performer_urls as URL,
  SearchScene_searchScene_performers_performer_images as Image,
  SearchScene_searchScene_performers_performer_measurements as Measurements
} from 'src/definitions-box/SearchScene';
import { BreastTypeEnum, GenderEnum as StashGenderEnum} from 'src/definitions-box/globalTypes';
import { GenderEnum } from 'src/core/generated-graphql';

const CDN = 'https://cdn.stashdb.org';

export const sortImageURLs = (images: Image[], orientation: 'portrait'|'landscape') => (
  images.map((i) => ({
      url: i.id ? `${CDN}/${i.id.slice(0, 2)}/${i.id.slice(2, 4)}/${i.id}` : i.url,
      id: i.id,
      width: i.width ?? 1,
      height: i.height ?? 1,
      aspect: orientation === 'portrait' ? (i.height ?? 1) / (i.width ?? 1) > 1 : ((i.width ?? 1) / (i.height ?? 1)) > 1
  })).sort((a, b) => {
      if (a.aspect > b.aspect) return -1;
      if (a.aspect < b.aspect) return 1;
      if (orientation === 'portrait' && a.height > b.height) return -1;
      if (orientation === 'portrait' && a.height < b.height) return 1;
      if (orientation === 'landscape' && a.width > b.width) return -1;
      if (orientation === 'landscape' && a.width < b.width) return 1;
      return 0;
  })
)

export const getImage = (images: Image[], orientation: 'portrait'|'landscape') => (
  sortImageURLs(images, orientation)?.[0]?.url ?? ''
);

export const getUrlByType = (
    urls:(URL|null)[],
    type:string
) => (
  (urls && (urls.find((url) => url?.type === type) || {}).url) || ''
);


export const formatMeasurements = (measurements: Measurements) => (
  (measurements.cup_size && measurements.waist && measurements.hip) ?
    `${measurements.band_size}${measurements.cup_size}-${measurements.waist}-${measurements.hip}`
  : ''
);

export const formatBreastType = (type: BreastTypeEnum|null) => (
    type === BreastTypeEnum.FAKE ? "Yes" :
    type === BreastTypeEnum.NATURAL ? "No" : ""
);

export const formatGender = (type: StashGenderEnum|null) => (
  type === StashGenderEnum.FEMALE ? GenderEnum.Female :
  type === StashGenderEnum.MALE ? GenderEnum.Male :
  type === StashGenderEnum.TRANSGENDER_FEMALE ? GenderEnum.TransgenderFemale :
  type === StashGenderEnum.TRANSGENDER_MALE ? GenderEnum.TransgenderMale :
  type === StashGenderEnum.INTERSEX ? GenderEnum.Intersex : null
);
