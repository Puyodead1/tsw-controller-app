import { Controller, UseFormReturn } from "react-hook-form";
import { main } from "../../../wailsjs/go/models";
import { useCallback } from "react";

type Props = {
  form: UseFormReturn<{
    profiles: Partial<Record<string, main.Interop_SelectedProfileInfo>>;
  }>;
  controller: main.Interop_GenericController;
  profiles: main.Interop_Profile[];
  onReloadConfiguration: () => void;
  onBrowseConfiguration: () => void;
  onCreateProfile: (controller: main.Interop_GenericController) => void;
  onSaveControllerProfileForSharing: (
    controller: main.Interop_GenericController,
  ) => void;
  onOpenProfileForController: (
    controller: main.Interop_GenericController,
  ) => void;
  onDeleteProfileForController: (
    controller: main.Interop_GenericController,
  ) => void;
};

const updatedAtFormatter = new Intl.DateTimeFormat(undefined, {
  dateStyle: "medium",
  timeStyle: "medium",
});

const unfocus = () => {
  if (document.activeElement && document.activeElement instanceof HTMLElement) {
    document.activeElement.blur();
  }
};

export function MainTabControllerProfileSelector({
  form,
  controller,
  profiles,
  onReloadConfiguration,
  onBrowseConfiguration,
  onCreateProfile,
  onSaveControllerProfileForSharing,
  onOpenProfileForController,
  onDeleteProfileForController,
}: Props) {
  const { watch, control } = form;
  const selectedProfile = watch(`profiles.${controller.GUID}`);
  const supportedProfiles = profiles?.filter(
    (profile) => !profile.UsbID || profile.UsbID === controller.UsbID,
  );
  const unsupportedProfiles = profiles?.filter(
    (profile) => profile.UsbID && profile.UsbID !== controller.UsbID,
  );

  const unfocusHandlerFactory = useCallback((func: () => void) => {
    return () => {
      func();
      setTimeout(unfocus, 0);
    };
  }, []);

  return (
    <div key={controller.GUID} className="flex flex-row gap-2">
      <Controller
        control={control}
        name={`profiles.${controller.GUID}`}
        render={({ field }) => (
          <div className="grow dropdown dropdown-start">
            <div tabIndex={0} role="button" className="select w-full">
              {selectedProfile?.Name ?? "Select profile"}
            </div>
            <div className="dropdown-content shadow-sm max-h-[50dvh] overflow-auto w-full">
              <ul className="menu w-full bg-base-300 rounded-box p-2">
                {supportedProfiles.map((profile) => (
                  <li key={profile.Name}>
                    <button
                      className="grid grid-cols-1 grid-flow-row auto-rows-max gap-0"
                      onClick={unfocusHandlerFactory(() => {
                        field.onChange({
                          Id: profile.Id,
                          Name: profile.Name,
                          IsAutoSelect: false,
                        });
                      })}
                    >
                      <span>{profile.Name}</span>
                      <span className="text-base-content/30 text-xs">
                        Last updated:{" "}
                        {updatedAtFormatter.format(
                          new Date(profile.Metadata.UpdatedAt),
                        )}
                      </span>
                      <span className="text-base-content/30 text-xs">
                        {profile.Metadata.Path}
                      </span>
                      {!!profile.Metadata.Warnings.length &&
                        profile.Metadata.Warnings.map((warning) => (
                          <div
                            key={warning}
                            role="alert"
                            className="alert alert-soft alert-warning my-2 p-2 text-xs"
                          >
                            {warning}
                          </div>
                        ))}
                    </button>
                  </li>
                ))}
                {unsupportedProfiles.map((profile) => (
                  <li key={profile.Name} className="menu-disabled">
                    <button className="grid grid-cols-1 grid-flow-row auto-rows-max gap-0">
                      <span>{profile.Name}</span>
                      <span className="text-base-content/30 text-xs">
                        Disabled for controller
                      </span>
                    </button>
                  </li>
                ))}
              </ul>
            </div>
          </div>
        )}
      />

      <div className="dropdown dropdown-end">
        <div tabIndex={0} role="button" className="btn">
          More
        </div>
        <ul
          tabIndex={-1}
          className="dropdown-content menu bg-base-100 rounded-box z-1 w-52 p-2 shadow-sm"
        >
          <li>
            <button onClick={unfocusHandlerFactory(onReloadConfiguration)}>
              Reload configuration
            </button>
          </li>
          <li>
            <button onClick={unfocusHandlerFactory(onBrowseConfiguration)}>
              Browse configuration
            </button>
          </li>
          <li>
            <button
              onClick={unfocusHandlerFactory(() => onCreateProfile(controller))}
            >
              Create new profile
            </button>
          </li>
          <li>
            <button
              disabled={!selectedProfile}
              onClick={unfocusHandlerFactory(() =>
                onSaveControllerProfileForSharing(controller),
              )}
              className="disabled:opacity-50 disabled:pointer-events-none"
            >
              Save profile for sharing
            </button>
          </li>
          <li>
            <button
              disabled={!selectedProfile}
              onClick={unfocusHandlerFactory(() =>
                onOpenProfileForController(controller),
              )}
              className="disabled:opacity-50 disabled:pointer-events-none"
            >
              Open profile in builder
            </button>
          </li>
          <li>
            <button
              disabled={!selectedProfile}
              onClick={unfocusHandlerFactory(() =>
                onDeleteProfileForController(controller),
              )}
              className="disabled:opacity-50 disabled:pointer-events-none"
            >
              Delete profile
            </button>
          </li>
        </ul>
      </div>
    </div>
  );
}
