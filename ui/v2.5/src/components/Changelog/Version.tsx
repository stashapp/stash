import React, { useState } from 'react';
import { Button, Collapse } from 'react-bootstrap';

interface IVersionProps {
  header: string;
  defaultOpen?: boolean
}


const Version: React.FC<IVersionProps> = ({header, defaultOpen = false, children}) => {
  const [open, setOpen] = useState(defaultOpen);
  return (
    <div>
      <h4><Button onClick={() => setOpen(!open)} variant="link">{header}</Button></h4>
      <Collapse in={open}>{children}</Collapse>
    </div>
  );
}

export default Version;
