import React from "react";
import { PatchComponent } from "src/patch";

export const SweatDrops: React.FC = PatchComponent("SweatDrops", () => (
  <span>
    <svg
      xmlns="http://www.w3.org/2000/svg"
      xmlnsXlink="http://www.w3.org/1999/xlink"
      aria-hidden="true"
      focusable="false"
      width="1em"
      height="1em"
      style={{ transform: "rotate(360deg)" }}
      preserveAspectRatio="xMidYMid meet"
      viewBox="0 0 36 36"
    >
      <path
        fill="currentColor"
        d="M22.855.758L7.875 7.024l12.537 9.733c2.633 2.224 6.377 2.937 9.77 1.518c4.826-2.018 7.096-7.576 5.072-12.413C33.232 1.024 27.68-1.261 22.855.758zm-9.962 17.924L2.05 10.284L.137 23.529a7.993 7.993 0 0 0 2.958 7.803a8.001 8.001 0 0 0 9.798-12.65zm15.339 7.015l-8.156-4.69l-.033 9.223c-.088 2 .904 3.98 2.75 5.041a5.462 5.462 0 0 0 7.479-2.051c1.499-2.644.589-6.013-2.04-7.523z"
      />
      <rect x="0" y="0" width="36" height="36" fill="rgba(0, 0, 0, 0)" />
    </svg>
  </span>
));
