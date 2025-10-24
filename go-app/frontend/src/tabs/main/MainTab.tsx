import useSWR from "swr";
import { lt as semverLt } from "semver";
import {
  GetProfiles,
  LoadConfiguration,
  SelectProfile,
  ClearProfile,
  GetSelectedProfiles,
  InstallTrainSimWorldMod,
  OpenConfigDirectory,
  GetLastInstalledModVersion,
  SetLastInstalledModVersion,
  GetVersion,
  GetControllers,
  OpenProfileBuilder,
  OpenNewProfileBuilder,
  SaveProfileForSharing,
  ImportProfile,
} from "../../../wailsjs/go/main/App";
import { useEffect } from "react";
import { BrowserOpenURL, EventsOn } from "../../../wailsjs/runtime/runtime";
import { events } from "../../events";
import { useForm } from "react-hook-form";
import { MainTabControllerProfileSelector } from "./MainTabControllerProfileSelecor";
import { main } from "../../../wailsjs/go/models";

type FormValues = {
  profiles: Record<`${string}`, string>;
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
    { revalidateOnMount: true },
  );
  const { data: controllers, mutate: refetchControllers } = useSWR(
    "controllers",
    () => GetControllers(),
    { revalidateOnMount: true },
  );

  const form = useForm<FormValues>({
    defaultValues: async () => ({
      profiles: await GetSelectedProfiles(),
    }),
  });
  const { register, watch, getValues } = form;

  const handleReloadConfiguration = () => {
    LoadConfiguration().then(() => {
      /* re-select profile after reloading */
      const profiles = getValues("profiles");
      for (const guid in profiles) {
        if (profiles[guid]) {
          SelectProfile(guid, profiles[guid]);
        }
      }
    });
  };

  const handleBrowseConfig = () => {
    OpenConfigDirectory();
  };

  const handleCreateProfile = (controller: main.Interop_GenericController) => {
    OpenNewProfileBuilder(controller.UsbID);
  };

  const handleOpenProfile = (controller: main.Interop_GenericController) => {
    const profile = getValues(`profiles.${controller.GUID}`);
    if (profile) OpenProfileBuilder(profile);
  };

  const handleSaveProfileForSharing = (
    controller: main.Interop_GenericController,
  ) => {
    const profile = getValues(`profiles.${controller.GUID}`);
    if (profile) SaveProfileForSharing(controller.GUID, profile);
  };

  const openInWindow = (url: string) => {
    BrowserOpenURL(url);
  };

  const handleInstall = () => {
    InstallTrainSimWorldMod()
      .then(() => refetchVersionInfo())
      .catch((err) => alert(String(err)));
  };

  const handleImportProfile = () => {
    ImportProfile()
      .then(() => LoadConfiguration())
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
      const profiles = getValues("profiles");
      for (const guid in profiles) {
        if (profiles[guid]) SelectProfile(guid, profiles[guid]);
        else ClearProfile(guid);
      }
    });
  }, []);

  useEffect(() => {
    return EventsOn(events.profiles_updated, () => {
      refetchProfiles();
    });
  }, []);

  useEffect(() => {
    return EventsOn(events.joydevices_updated, () => {
      refetchControllers();
    });
  }, []);

  return (
    <div className="grid grid-cols-1 grid-flow-row auto-rows-max gap-2">
      {controllers?.map((c) => (
        <div key={c.GUID}>
          <MainTabControllerProfileSelector
            controller={c}
            profiles={profiles ?? []}
            form={form}
            onBrowseConfiguration={handleBrowseConfig}
            onCreateProfile={handleCreateProfile}
            onReloadConfiguration={handleReloadConfiguration}
            onSaveControllerProfileForSharing={handleSaveProfileForSharing}
            onOpenProfileForController={handleOpenProfile}
          />
        </div>
      ))}
      {/* steam://controllerconfig/2967990/3576092503 */}
      <div className="flex gap-2">
        <button className="btn btn-sm grow" onClick={handleInstall}>
          Install/Reinstall Train Sim World mod
        </button>
        <button className="btn btn-sm grow" onClick={handleImportProfile}>
          Import profile (.tswprofile)
        </button>
      </div>
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
