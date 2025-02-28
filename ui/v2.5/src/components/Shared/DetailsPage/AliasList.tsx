export const AliasList: React.FC<{ aliases: string[] | undefined }> = ({
  aliases,
}) => {
  if (!aliases?.length) {
    return null;
  }

  return (
    <div>
      <span className="alias-head">{aliases.join(", ")}</span>
    </div>
  );
};
