import React from 'react';
import Version from './Version';

const Changelog: React.FC = () => {
  return (
    <div>
      <h4>Changelog:</h4>
      <Version
        header="v0.2.0 - WIP"
        defaultOpen
      />
      <Version header="v0.1.0 - 2020-02-24" />
    </div>
  );
}

export default Changelog;
