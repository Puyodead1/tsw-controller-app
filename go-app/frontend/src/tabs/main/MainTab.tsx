import useSWR from "swr";
import { lt as semverLt } from "semver";
import {
  GetProfiles,
  LoadConfiguration,
  SelectProfile,
  ClearProfile,
  GetSelectedProfile,
  InstallTrainSimWorldMod,
  OpenConfigDirectory,
  GetLastInstalledModVersion,
  SetLastInstalledModVersion,
  GetVersion,
} from "../../../wailsjs/go/main/App";
import { useEffect } from "react";
import { BrowserOpenURL, EventsOn } from "../../../wailsjs/runtime/runtime";
import { events } from "../../events";
import { useForm } from "react-hook-form";

type FormValues = {
  profile: string;
};

export const MainTab = () => {
  const { data: versionInfo, mutate: refetchVersionInfo } = useSWR(
    "version-info",
    () =>
      Promise.all([GetVersion(), GetLastInstalledModVersion()]).then(
        ([version, lastInstalledModVersion]) => ({
          version,
          lastInstalledModVersion,
        }),
      ),
    { revalidateOnMount: true },
  );
  const { data: profiles, mutate: refetchProfiles } = useSWR(
    "profiles",
    () => GetProfiles(),
    {
      revalidateOnMount: true,
    },
  );

  const { register, watch, getValues } = useForm<FormValues>({
    defaultValues: async () => ({
      profile: await GetSelectedProfile(),
    }),
  });

  const handleReloadConfiguration = () => {
    LoadConfiguration().then(() => {
      /* re-select profile after reloading */
      const selectedProfile = getValues('profile')
      if (selectedProfile.length) {
        SelectProfile(selectedProfile)
      }
    });
  };

  const handleBrowseConfig = () => {
    OpenConfigDirectory();
  };

  const openInWindow = (url: string) => {
    BrowserOpenURL(url);
  };

  const handleInstall = () => {
    InstallTrainSimWorldMod()
      .then(() => refetchVersionInfo())
      .catch((err) => alert(String(err)));
  };

  const handleIgnoreModInstallWarning = () => {
    if (versionInfo) {
      SetLastInstalledModVersion(versionInfo.version).then(() =>
        refetchVersionInfo(),
      );
    }
  };

  useEffect(() => {
    watch(() => {
      const profile = getValues("profile");
      if (profile.length) SelectProfile(profile);
      else ClearProfile();
    });
  }, []);

  useEffect(() => {
    return EventsOn(events.profiles_updated, () => {
      refetchProfiles();
    });
  }, []);

  return (
    <div className="grid grid-cols-1 grid-flow-row auto-rows-max gap-2">
      <div className="flex flex-row gap-2">
        <fieldset className="fieldset grow">
          <legend className="fieldset-legend">Select profile</legend>
          <select className="select w-full" {...register("profile")}>
            <option selected value="">
              None
            </option>
            {profiles?.map((profile) => (
              <option key={profile.Name} value={profile.Name}>
                {profile.Name}
              </option>
            ))}
          </select>
          <span className="label">
            Auto-detect only works for certain supported controllers
          </span>
        </fieldset>
        <div className="dropdown dropdown-end mt-[34px]">
          <div tabIndex={0} role="button" className="btn">
            More
          </div>
          <ul
            tabIndex={0}
            className="dropdown-content menu bg-base-100 rounded-box z-1 w-52 p-2 shadow-sm"
          >
            <li>
              <button onClick={handleReloadConfiguration}>
                Reload configuration
              </button>
            </li>
            <li>
              <button onClick={handleBrowseConfig}>Browse configuration</button>
            </li>
          </ul>
        </div>
      </div>
      {/* steam://controllerconfig/2967990/3576092503 */}
      <button className="btn btn-sm w-full" onClick={handleInstall}>
        Install/Reinstall Train Sim World mod
      </button>
      {!versionInfo?.lastInstalledModVersion && (
        <div role="alert" className="alert alert-soft alert-warning">
          <span>
            It looks like you have not installed the Train Sim World mod yet,
            make sure you install the mod first.
          </span>
          <div>
            <button
              className="btn btn-sm"
              onClick={handleIgnoreModInstallWarning}
            >
              Ignore
            </button>
          </div>
        </div>
      )}
      {versionInfo?.lastInstalledModVersion &&
        semverLt(versionInfo.lastInstalledModVersion, versionInfo.version) && (
          <div role="alert" className="alert alert-soft alert-warning">
            <span>
              It looks like the app has updated since the last time you
              installed the mod, make sure to reinstall the updated mod version
              before starting the game.
            </span>
            <div>
              <button
                className="btn btn-sm"
                onClick={handleIgnoreModInstallWarning}
              >
                Ignore
              </button>
            </div>
          </div>
        )}
      <div role="alert" className="alert">
        <span>
          For this app to correctly work you will need to make sure Train Sim
          World is not able to process the controller input. You can achieve
          this by configuring your controller in using{" "}
          <button
            className="inline link"
            onClick={() =>
              openInWindow(
                "https://github.com/LiamMartens/tsw-controller-app/blob/main/STEAM_INPUT_SETUP.md",
              )
            }
          >
            Steam Input
          </button>{" "}
          and applying the following "Disabled Controller" layout preset for the
          game. Alternatively, you can also a software like{" "}
          <button
            className="inline link"
            onClick={() =>
              openInWindow("https://ds4-windows.com/download/hidhide/")
            }
          >
            HidHide
          </button>{" "}
          to hide the controller from the game altogether
        </span>
      </div>
    </div>
  );
};
