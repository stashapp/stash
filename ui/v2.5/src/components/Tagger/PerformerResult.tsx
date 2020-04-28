import React, { useEffect, useState } from 'react';
import { Button } from 'react-bootstrap';
import cx from 'classnames';

import { Icon, PerformerSelect } from 'src/components/Shared';
import * as GQL from 'src/core/generated-graphql';
import { ValidTypes } from 'src/components/Shared/Select';
import {
  SearchScene_searchScene_performers_performer as StashPerformer,
} from 'src/definitions-box/SearchScene';
import { sortImageURLs } from './utils';

import PerformerModal from './PerformerModal';

interface IIconProps {
  className?: string;
}

const SuccessIcon: React.FC<IIconProps> = ({ className }) => (
  <Icon icon="check" className={cx("success mr-4", className)} color="#0f9960" />
);
const FailIcon: React.FC<IIconProps> = ({ className }) => (
  <Icon icon="times" className={cx("secondary mr-4", className)} color="#394b59" />
);

interface IPerformerOperation {
  type: "Create"|"Existing"|"Update";
  data: StashPerformer|string;
}

interface IPerformerResultProps {
  performer: StashPerformer
  setPerformer: (data:IPerformerOperation) => void;
}

const PerformerResult: React.FC<IPerformerResultProps> = ({ performer, setPerformer }) => {
  const [selectedPerformer, setSelectedPerformer] = useState();
  const [selectedSource, setSelectedSource] = useState<'create'|'existing'|undefined>();
  const [modalVisible, showModal] = useState(false);
  const { data: stashData, loading: stashLoading } = GQL.useFindPerformersQuery({
    variables: {
      performer_filter: {
        stash_id: {
          value: performer.id,
          modifier: GQL.CriterionModifier.Equals
        }
      }
    }
  });
  const { loading } = GQL.useFindPerformersQuery({
    variables: {
      filter: {
        q: `"${performer.name}"`
      }
    },
    onCompleted: (data) => {
      const performerResult = data.findPerformers?.performers?.[0]?.id;
      if (performerResult) {
        setSelectedPerformer(performerResult);
        setSelectedSource('existing');
        setPerformer({
          type: 'Update',
          data: performerResult
        });
      }
    }
  });

  useEffect(() => {
    if(!stashData?.findPerformers.count)
      return;

    setPerformer({
      type: 'Existing',
      data: stashData.findPerformers.performers[0].id
    });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [stashData]);

  const handlePerformerSelect = (items: ValidTypes[]) => {
    if (items.length) {
      setSelectedSource('existing');
      setSelectedPerformer(items[0].id);
      setPerformer({
        type: 'Update',
        data: items[0].id
      });
    }
    else {
      setSelectedSource(undefined);
      setSelectedPerformer(null);
    }
  };

  const handlePerformerCreate = (imageIndex: number) => {
    const images = sortImageURLs(performer.images, 'portrait');
    const imageURLs = images.length ? [{
      url: images[imageIndex].url,
      id: images[imageIndex].id,
      width: null,
      height: null
    }] : [];
    setSelectedSource('create');
    setPerformer({
      type: 'Create',
      data: {
        ...performer,
        images: imageURLs
      }
    });
    showModal(false);
  };

  if(stashLoading || loading)
    return <div>Loading performer</div>;

  if((stashData?.findPerformers.count ?? 0) > 0) {
    return (
      <div className="row my-2">
        <span className="ml-auto">
          <SuccessIcon />Performer matched:
        </span>
        <b className="col-3 text-right">{ stashData!.findPerformers.performers[0].name }</b>
      </div>
    );
  }
  return (
    <div className="row align-items-center mt-2">
      <PerformerModal
        showModal={showModal}
        modalVisible={modalVisible}
        performer={performer}
        handlePerformerCreate={handlePerformerCreate}
      />
      <div className="entity-name">
        Performer:
        <b className="ml-2">{performer.name}</b>
      </div>
      <div>
        <Button variant="secondary" className="mr-1" onClick={() => showModal(true)}>Create</Button>
        { selectedSource === 'create'
          ? <SuccessIcon />
          : <FailIcon />
        }
      </div>
      <div className="select-existing">
        { selectedSource === 'existing'
          ? <SuccessIcon />
          : <FailIcon />
        }
      </div>
      <PerformerSelect
        ids={selectedPerformer ? [selectedPerformer] : []}
        onSelect={handlePerformerSelect}
        className="performer-select"
      />
    </div>
  );
}

export default PerformerResult;
