import useSWR from "swr";
import { lt } from "semver";
import {
  GetLatestReleaseVersion,
  GetVersion,
  UpateApp,
} from "../wailsjs/go/main/App";

export const SelfUpdateBanner = () => {
  const { data: versionInfo, mutate: refetchVersionInfo } = useSWR(
    "version-info-update-banner",
    () =>
      Promise.all([GetVersion(), GetLatestReleaseVersion()]).then(
        ([version, latestReleaseVersion]) => ({
          version,
          latestReleaseVersion,
        }),
      ),
    { revalidateOnMount: true },
  );

  const handleUpdate = () => {
    UpateApp().catch((err) => alert(err));
  };

  if (
    versionInfo?.version &&
    lt(versionInfo.version, versionInfo.latestReleaseVersion)
  ) {
    return (
      <div className="flex flex-row gap-2 items-center p-2">
        <div className="inline-grid *:[grid-area:1/1]">
          <div className="status status-info"></div>
          <div className="status status-info"></div>
        </div>{" "}
        <p className="text-xs">
          A new version is available{" "}
          <button className="link" onClick={handleUpdate}>Update now</button>
        </p>
      </div>
    );
  }

  return null;
};
