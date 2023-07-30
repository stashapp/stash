import * as GQL from "src/core/generated-graphql";

interface IProps {
  scraper: GQL.Scraper;
  onHide: () => void;
  onSelectStudio: (
    studio: GQL.ScrapedStudioDataFragment,
    scraper: GQL.Scraper
  ) => void;
  name?: string;
}

const StudioScrapeModal: React.FC<IProps> = ({
  scraper,
  name,
  onHide,
  onSelectStudio,
}) => {
  return (
    <div>
      <div>StudioScrapeModal</div>
    </div>
  );
}

export default StudioScrapeModal;
