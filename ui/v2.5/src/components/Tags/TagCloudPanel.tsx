import React, { useMemo } from "react";
import { Link } from "react-router-dom";
import { FormattedMessage } from "react-intl";
import { useFindTags } from "src/core/StashService";
import { ListFilterModel } from "src/models/list-filter/filter";
import { FilterMode, SortDirectionEnum } from "src/core/generated-graphql";
import { RecommendationRow } from "../FrontPage/RecommendationRow";
import { PatchComponent } from "src/patch";

// Five English-language web-safe fonts cycled across tags
const TAG_CLOUD_FONTS = [
  "Georgia, serif",
  "'Courier New', Courier, monospace",
  "Arial, Helvetica, sans-serif",
  "'Trebuchet MS', Helvetica, sans-serif",
  "'Palatino Linotype', Palatino, serif",
];

const FONT_SIZES = [0.8, 0.9, 1.0, 1.1, 1.2, 1.35, 1.5];

// Stable per-tag style derived from the tag's numeric id.
// Uses a simple LCG-style mix so adjacent IDs produce visually different outputs.
function tagStyle(id: string): React.CSSProperties {
  const n = parseInt(id, 10) || 0;
  // cheap integer hash
  const h = ((n ^ (n >>> 16)) * 0x45d9f3b) >>> 0;

  const fontFamily = TAG_CLOUD_FONTS[h % TAG_CLOUD_FONTS.length];
  const fontSize = FONT_SIZES[(h >>> 3) % FONT_SIZES.length];
  const bold = (h >>> 6) % 3 === 0; // ~33% bold
  const italic = (h >>> 9) % 3 === 0; // ~33% italic
  const underline = (h >>> 12) % 4 === 0; // ~25% underline

  return {
    fontFamily,
    fontSize: `${fontSize}rem`,
    fontWeight: bold ? "bold" : "normal",
    fontStyle: italic ? "italic" : "normal",
    textDecoration: underline ? "underline" : "none",
  };
}

export const TagCloudPanel: React.FC = PatchComponent("TagCloudPanel", () => {
  const filter = useMemo(() => {
    const f = new ListFilterModel(FilterMode.Tags, undefined);
    f.itemsPerPage = 1000;
    f.currentPage = 1;
    f.sortBy = "name";
    f.sortDirection = SortDirectionEnum.Asc;
    return f;
  }, []);

  const result = useFindTags(filter);
  const tags = result.data?.findTags.tags ?? [];

  if (!result.loading && tags.length === 0) {
    return null;
  }

  return (
    <RecommendationRow
      className="tag-cloud-row"
      header={<FormattedMessage id="tag_cloud" />}
      link={
        <Link to="/tags">
          <FormattedMessage id="view_all" />
        </Link>
      }
    >
      <div className="tag-cloud">
        {result.loading
          ? null
          : tags.map((tag) => (
              <Link
                key={tag.id}
                to={`/tags/${tag.id}`}
                className="tag-cloud-link"
                style={tagStyle(tag.id)}
              >
                {tag.name}
              </Link>
            ))}
      </div>
    </RecommendationRow>
  );
});
