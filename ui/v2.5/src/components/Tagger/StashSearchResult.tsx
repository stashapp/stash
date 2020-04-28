import React, { useCallback, useState } from 'react';
import { blobToBase64 } from 'base64-blob';
import { loader } from 'graphql.macro';
import cx from 'classnames';
import { Button } from 'react-bootstrap';

import {
  SearchScene_searchScene as SearchResult,
  SearchScene_searchScene_performers_performer as StashPerformer,
  SearchScene_searchScene_studio as StashStudio
} from 'src/definitions-box/SearchScene';
import { FingerprintAlgorithm } from 'src/definitions-box/globalTypes';
import * as GQL from 'src/core/generated-graphql';
import {
  SubmitFingerprintVariables,
  SubmitFingerprint
} from 'src/definitions-box/SubmitFingerprint';
import { FindPerformersDocument, FindStudioByStashIdDocument, AllTagsForFilterDocument } from '../../core/generated-graphql';
import PerformerResult from './PerformerResult';
import StudioResult from './StudioResult';
import { formatGender, formatMeasurements, formatBreastType, getUrlByType, getImage } from './utils';
import { client } from './client';

const SubmitFingerprintMutation = loader('src/queries/submitFingerprint.gql');

const getDurationStatus = (scene: SearchResult, stashDuration: number|undefined|null) => {
  const fingerprintDuration = scene.fingerprints.map(f => f.duration)?.[0] ?? null;
  const sceneDuration = scene.duration || fingerprintDuration;
  if(!sceneDuration || !stashDuration) return '';
  const diff = Math.abs(sceneDuration - stashDuration);
  if(diff < 5) {
    return <div><b>Duration is a match</b></div>;
  }
  return <div>Duration off by {Math.floor(diff)}s</div>;
};

interface IStashSearchResultProps {
  scene: SearchResult;
  stashScene: Partial<GQL.Scene>;
  isActive: boolean;
  setActive: () => void;
  showMales: boolean;
  setScene: (scene: Partial<GQL.Scene>) => void;
  isFingerprintMatch?: boolean;
  setCoverImage: boolean;
}

const titleCase = (str?: string) => {
  if (!str) return '';
  return (str ?? '').split(' ')
    .map(w => w[0].toUpperCase() + w.substr(1).toLowerCase())
    .join(' ');
};

interface IPerformerOperation {
  type: "Create"|"Existing"|"Update";
  data: StashPerformer|string;
}

interface IStudioOperation {
  type: "Create"|"Existing"|"Update";
  data: StashStudio|string;
}

const StashSearchResult: React.FC<IStashSearchResultProps> = ({ scene, stashScene, isActive, setActive, showMales, setScene, setCoverImage }) => {
  const [studio, setStudio] = useState<IStudioOperation>();
  const [performers, setPerformers] = useState<Record<string, IPerformerOperation>>();

  const [createStudio] = GQL.useStudioCreateMutation();
  const [updateStudio] = GQL.useStudioUpdateMutation();
  const [updateScene] = GQL.useSceneUpdateMutation();
  const [createPerformer] = GQL.usePerformerCreateMutation();
  const [updatePerformer] = GQL.usePerformerUpdateMutation();
  const [createTag] = GQL.useTagCreateMutation();
  const { data: allTags } = GQL.useAllTagsForFilterQuery();

  const setPerformer = useCallback((performerData: IPerformerOperation, performerID: string) => (
    setPerformers({ ...performers, [performerID]: performerData })
  ), [performers]);

  const handleSave = async () => {
    if(!performers || !studio)
      return;

    let studioID:string;
    let performerIDs = [];

    if (studio.type === 'Update') {
      const studioUpdateResult = await updateStudio({
        variables: {
          id: studio.data as string,
          stash_id: scene.studio?.id ?? ''
        },
        update: (store, updatedStudio) => {
          if (!updatedStudio?.data?.studioUpdate)
            return;

          store.writeQuery({
            query: FindStudioByStashIdDocument,
            variables: {
              id: updatedStudio.data.studioUpdate.stash_id
            },
            data: {
              findStudioByStashID: updatedStudio.data.studioUpdate
            }
          });
        }
      });
      const id = studioUpdateResult.data?.studioUpdate?.id;
      if(studioUpdateResult.errors || !id)
        return;
      studioID = id;
    }
    else if(studio.type === 'Create') {
      const studioData = studio.data as StashStudio;
      const studioCreateResult = await createStudio({
        variables: {
          name: studioData.name,
          stash_id: studioData.id,
          ...(!!getUrlByType(studioData.urls, 'HOME') && { url: getUrlByType(studioData.urls, 'HOME') }),
          ...(!!getUrlByType(studioData.urls, 'PHOTO') && { image: getUrlByType(studioData.urls, 'PHOTO') })
        },
        update: (store, newStudio) => {
          if (!newStudio?.data?.studioCreate)
            return;

          store.writeQuery({
            query: FindStudioByStashIdDocument,
            variables: {
              id: newStudio.data.studioCreate.stash_id
            },
            data: {
              findStudioByStashID: newStudio.data.studioCreate
            }
          });
        }
      });

      const id = studioCreateResult.data?.studioCreate?.id;
      if(studioCreateResult.errors || !id)
        return;
      studioID = id;
    }
    else {
      studioID = studio.data as string;
    }

    performerIDs = await Promise.all(Object.keys(performers).map(async (performerID) => {
      const performer = performers[performerID];
      if (performer.type === 'Update') {
        const res = await updatePerformer({
          variables: {
            id: performer.data as string,
            stash_id: performerID
          },
          update: (store, updatedPerformer) => {
            if (!updatedPerformer?.data?.performerUpdate)
              return;

            store.writeQuery({
              query: FindPerformersDocument,
              variables: {
                performer_filter: {
                  stash_id: {
                    value: updatedPerformer.data.performerUpdate.stash_id,
                    modifier: GQL.CriterionModifier.Equals
                  }
                }
              },
              data: {
                findPerformers: {
                  performers: [updatedPerformer.data.performerUpdate],
                  count: 1,
                  __typename: "FindPerformersResultType"
                }
              }
            });
          }
        });

        if(res.errors)
          return;

        return res?.data?.performerUpdate?.id ?? null;
      }
      if(performer.type === 'Create') {
        const performerData = performer.data as StashPerformer;
        const imgurl = performerData.urls?.[0]?.url;
        let imgData = null;
        if(imgurl) {
          const img = await fetch(imgurl, {
            mode: 'cors',
            cache: 'no-store'
          });
          if(img.status === 200) {
            const blob = await img.blob();
            imgData = await blobToBase64(blob);
          }
        }

        const res = await createPerformer({
          variables: {
            name: performerData.name,
            gender: formatGender(performerData.gender),
            country: performerData.country,
            height: performerData.height?.toString(),
            ethnicity: titleCase(performerData.ethnicity ?? ''),
            birthdate: performerData.birthdate?.date ?? null,
            eye_color: titleCase(performerData.eye_color ?? ''),
            fake_tits: formatBreastType(performerData.breast_type),
            measurements: formatMeasurements(performerData.measurements),
            image: imgData,
            stash_id: performerID
          },
          update: (store, newPerformer) => {
            if (!newPerformer?.data?.performerCreate)
              return;

            store.writeQuery({
              query: FindPerformersDocument,
              variables: {
                performer_filter: {
                  stash_id: {
                    value: newPerformer.data.performerCreate.stash_id,
                    modifier: GQL.CriterionModifier.Equals
                  }
                }
              },
              data: {
                findPerformers: {
                  performers: [newPerformer.data.performerCreate],
                  count: 1,
                  __typename: "FindPerformersResultType"
                }
              }
            });
          }
        });

        if(res.errors)
          return;

        return res?.data?.performerCreate?.id ?? null;
      }
      return performer.data as string;
    }));

    if(studioID && !performerIDs.some(id => !id)) {
      const imgurl = getImage(scene.images, 'landscape');
      let imgData = null;
      if(imgurl && setCoverImage) {
        const img = await fetch(imgurl, {
          mode: 'cors',
          cache: 'no-store'
        });
        if(img.status === 200) {
          const blob = await img.blob();
          imgData = await blobToBase64(blob);
        }
      }

      const tagIDs:string[] = [];
      const tags = scene.tags ?? [];
      if (tags.length > 0) {
        const tagDict:Record<string, string> = (allTags?.allTagsSlim ?? []).reduce((dict, t) => ({ ...dict, [t.name.toLowerCase()]: t.id }), {});
        const newTags:string[] = [];
        tags.forEach(tag => {
          if (tagDict[tag.name.toLowerCase()])
            tagIDs.push(tagDict[tag.name.toLowerCase()]);
          else
            newTags.push(tag.name);
        });

        const createdTags = await Promise.all(newTags.map(tag => (
          createTag({
            variables: {
              name: tag
            },
            update: (store, _newTag) => {
              if (!_newTag.data?.tagCreate)
                return;

              store.writeQuery({
                query: AllTagsForFilterDocument,
                variables: {},
                data: {
                  allTagsSlim: [...(allTags?.allTagsSlim ?? []), _newTag.data.tagCreate],
                }
              });
            }
          })
        )));
        createdTags.forEach(createdTag => {
          if (createdTag?.data?.tagCreate?.id)
            tagIDs.push(createdTag.data.tagCreate.id);
        });
      }

      const sceneUpdateResult = await updateScene({
        variables: {
          id: stashScene.id ?? '',
          stash_id: scene.id,
          title: scene.title,
          details: scene.details,
          date: scene.date,
          performer_ids: performerIDs as string[],
          studio_id: studioID,
          cover_image: imgData,
          url: getUrlByType(scene.urls, 'STUDIO') ?? null,
          ...(tagIDs ? { tag_ids: tagIDs } : {})
        }
      });
      if(sceneUpdateResult.data?.sceneUpdate)
        setScene(sceneUpdateResult.data.sceneUpdate);

      if(stashScene.checksum && stashScene.file?.duration)
        client.mutate<SubmitFingerprint, SubmitFingerprintVariables>({
          mutation: SubmitFingerprintMutation,
          variables: {
            input: {
              scene_id: scene.id,
              fingerprint: {
                hash: stashScene.checksum,
                algorithm: FingerprintAlgorithm.MD5,
                duration: Math.floor(stashScene.file?.duration)
              }
            }
          }
        });
    }
  };

  const classname = cx('row mb-4 search-result', { 'selected-result': isActive });

  const sceneTitle = getUrlByType(scene.urls, 'STUDIO') ? (
    <a href={getUrlByType(scene.urls, 'STUDIO')} target="_blank" rel="noopener noreferrer" className="scene-link">{scene?.title}</a>
  ) : (<span>{scene?.title}</span>);

  return (
    // eslint-disable-next-line jsx-a11y/click-events-have-key-events, jsx-a11y/no-noninteractive-element-interactions
    <li className={classname} key={scene?.id} onClick={() => !isActive && setActive()}>
      <div className="col-6 row">
        <img src={getImage(scene?.images, 'landscape')} alt="" className="align-self-center scene-image" />
        <div className="d-flex flex-column justify-content-center scene-metadata">
          <h4 className="text-truncate">{ sceneTitle }</h4>
          <h5>{scene?.studio?.name} • {scene?.date}</h5>
          <div>Performers: {scene?.performers?.map(p => p.performer.name).join(', ')}</div>
          { getDurationStatus(scene, stashScene.file?.duration) }
        </div>
      </div>
      { isActive && (
        <div className="col-6">
          <StudioResult studio={scene.studio} setStudio={setStudio} />
          { scene.performers
            .filter(p => p.performer.gender !== 'MALE' || showMales)
            .map(performer => (
              <PerformerResult performer={performer.performer} setPerformer={(data:IPerformerOperation) => setPerformer(data, performer.performer.id)} key={`${scene.id}${performer.performer.id}`} />
            ))
          }
          <div className="row pr-3 mt-2">
            <Button className="col-1 offset-11" onClick={handleSave}>Save</Button>
          </div>
        </div>
      )}
    </li>
  );
};

export default StashSearchResult;
