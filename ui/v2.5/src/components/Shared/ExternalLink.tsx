type IExternalLinkProps = JSX.IntrinsicElements["a"];

export const ExternalLink: React.FC<IExternalLinkProps> = (props) => {
  return <a target="_blank" rel="noopener noreferrer" {...props} />;
};
