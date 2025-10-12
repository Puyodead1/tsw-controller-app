import { UseFormReturn } from "react-hook-form";
import { main } from "../../../wailsjs/go/models";

type Props = {
  form: UseFormReturn<{
    profiles: Record<`${string}`, string>;
  }>;
  controller: main.Interop_GenericController;
  profiles: main.Interop_Profile[];
  onReloadConfiguration: () => void;
  onBrowseConfiguration: () => void;
  onCreateProfile: () => void;
  onSaveControllerProfileForSharing: (
    controller: main.Interop_GenericController,
  ) => void;
  onOpenProfileForController: (
    controller: main.Interop_GenericController,
  ) => void;
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
}: Props) {
  const { register, watch } = form;
  const selectedProfile = watch(`profiles.${controller.GUID}`);

  return (
    <div key={controller.GUID} className="flex flex-row gap-2">
      <fieldset className="fieldset grow">
        <legend className="fieldset-legend">
          Select profile for {controller.Name}
        </legend>
        <select
          className="select w-full"
          {...register(`profiles.${controller.GUID}`)}
        >
          <option selected value="">
            None
          </option>
          {profiles
            ?.filter(
              (profile) => !profile.UsbID || profile.UsbID === controller.UsbID,
            )
            .map((profile) => (
              <option key={profile.Name} value={profile.Name}>
                {profile.Name}
              </option>
            ))}
          {profiles
            ?.filter(
              (profile) => profile.UsbID && profile.UsbID !== controller.UsbID,
            )
            .map((profile) => (
              <option key={profile.Name} value={profile.Name} disabled>
                {profile.Name} (Not enabled for controller)
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
            <button onClick={onReloadConfiguration}>
              Reload configuration
            </button>
          </li>
          <li>
            <button onClick={onBrowseConfiguration}>
              Browse configuration
            </button>
          </li>
          <li>
            <button onClick={onCreateProfile}>Create new profile</button>
          </li>
          <li>
            <button
              disabled={!selectedProfile}
              onClick={() => onSaveControllerProfileForSharing(controller)}
              className="disabled:opacity-50 disabled:pointer-events-none"
            >
              Save profile for sharing
            </button>
          </li>
          <li>
            <button
              disabled={!selectedProfile}
              onClick={() => onOpenProfileForController(controller)}
              className="disabled:opacity-50 disabled:pointer-events-none"
            >
              Open profile in builder
            </button>
          </li>
        </ul>
      </div>
    </div>
  );
}
