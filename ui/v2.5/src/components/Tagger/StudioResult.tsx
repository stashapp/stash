import React, { useEffect, useState, Dispatch, SetStateAction } from 'react';
import { Button, ButtonGroup } from 'react-bootstrap';
import cx from 'classnames';

import { Icon, Modal, StudioSelect } from 'src/components/Shared';
import * as GQL from 'src/core/generated-graphql';
import {
  SearchScene_searchScene_studio as StashStudio
} from 'src/definitions-box/SearchScene';
import { getImage, getUrlByType } from './utils';
import { Operation } from './StashSearchResult';

interface IIconProps {
  className?: string;
}

const SuccessIcon: React.FC<IIconProps> = ({ className }) => (
  <Icon icon="check" className={cx("success mr-4", className)} color="#0f9960" />
);

interface IStudioOperation {
  type: Operation;
  data: StashStudio|string;
}

interface IStudioResultProps {
  studio: StashStudio|null;
  setStudio: Dispatch<SetStateAction<IStudioOperation|undefined>>;
}

const StudioResult: React.FC<IStudioResultProps> = ({ studio, setStudio }) => {
  const [selectedStudio, setSelectedStudio] = useState<string|null>();
  const [modalVisible, showModal] = useState(false);
  const [selectedSource, setSelectedSource] = useState<'create'|'existing'|'skip'|undefined>();
  const { data: stashData, loading: stashLoading } = GQL.useFindStudioByStashIdQuery({
    variables: {
      id: studio?.id ?? ''
    }
  })

  const handleStudioSelect = (id?: string) => {
    if (id) {
      setSelectedStudio(id);
      setSelectedSource('existing');
      setStudio({
        type: 'Update',
        data: id
      });
    }
    else {
      setSelectedSource(undefined);
      setSelectedStudio(null);
    }
  };


  const { loading } = GQL.useFindStudiosQuery({
    variables: {
      filter: {
        q: `"${studio?.name ?? ''}"`
      }
    },
    onCompleted: (data) => (
      handleStudioSelect(data.findStudios?.studios?.[0]?.id)
    )
  });

  useEffect(() => {
    if(!stashData?.findStudioByStashID)
      return;

    setStudio({
      type: 'Existing',
      data: stashData.findStudioByStashID.id
    });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [stashData]);

  const handleStudioCreate = () => {
    if(!studio)
      return;
    setSelectedSource('create');
    setStudio({
      type: 'Create',
      data: studio
    });
    showModal(false);
  };

  const handleStudioSkip = () => {
    setSelectedSource('skip');
    setStudio({
      type: 'Skip',
      data: ''
    });
  };

  if(loading || stashLoading)
    return <div>Loading studio</div>;

  if(stashData?.findStudioByStashID) {
    return (
      <div className="row my-2">
        <span className="ml-auto">
          <SuccessIcon />Studio matched:
        </span>
        <b className="col-3 text-right">{ stashData.findStudioByStashID.name }</b>
      </div>
    );
  }


  return (
    <div className="row align-items-center mt-2">
      <Modal
        show={modalVisible}
        accept={{ text: "Save", onClick: handleStudioCreate }}
        cancel={{ onClick: () => showModal(false), variant: "secondary" }}
      >
        <div className="row">
          <strong className="col-2">Name:</strong>
          <span className="col-10">{ studio?.name }</span>
        </div>
        <div className="row">
          <strong className="col-2">URL:</strong>
          <span className="col-10">{ getUrlByType(studio?.urls ?? [], 'HOME') }</span>
        </div>
        <div className="row">
          <strong className="col-2">Logo:</strong>
          <span className="col-10"><img src={ getImage(studio?.images ?? [], 'landscape') } alt='' /></span>
        </div>
      </Modal>

      <div className="entity-name">
        Studio:
        <b className="ml-2">{studio?.name}</b>
      </div>
      <ButtonGroup>
        <Button variant={selectedSource === 'create' ? 'primary' : "secondary"} onClick={() => showModal(true)}>Create</Button>
        <Button variant={selectedSource === 'skip' ? 'primary' : 'secondary'} onClick={() => handleStudioSkip()}>Skip</Button>
      <StudioSelect
        ids={selectedStudio ? [selectedStudio] : []}
        onSelect={(items) => handleStudioSelect(items.length ? items[0].id : undefined)}
          className={cx("studio-select", {'studio-select-active': selectedSource === 'existing' })}
      />
    </ButtonGroup>
  </div>
  );
}

export default StudioResult;
