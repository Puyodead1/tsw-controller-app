import { Controller, UseFormReturn } from "react-hook-form";
import { main } from "../../../wailsjs/go/models";
import {
  autoUpdate,
  flip,
  FloatingPortal,
  offset,
  size,
  useClick,
  useDismiss,
  useFloating,
  useInteractions,
} from "@floating-ui/react";
import { useState } from "react";

type Props = {
  form: UseFormReturn<{
    profiles: Record<`${string}`, string>;
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

  const [selectOpen, setSelectOpen] = useState(false);
  const { refs, floatingStyles, context } = useFloating<HTMLElement>({
    strategy: "fixed",
    placement: "bottom-start",
    open: selectOpen,
    onOpenChange: setSelectOpen,
    whileElementsMounted: autoUpdate,
    middleware: [
      offset(5),
      flip({ padding: 10 }),
      size({
        apply({ rects, elements, availableHeight }) {
          Object.assign(elements.floating.style, {
            maxHeight: `${availableHeight}px`,
            minWidth: `${rects.reference.width}px`,
          });
        },
        padding: 10,
      }),
    ],
  });

  const click = useClick(context);
  const dismiss = useDismiss(context);

  const { getReferenceProps, getFloatingProps } = useInteractions([
    click,
    dismiss,
  ]);

  return (
    <div key={controller.GUID} className="flex flex-row gap-2">
      <Controller
        control={control}
        name={`profiles.${controller.GUID}`}
        render={({ field }) => (
          <>
            <button
              ref={refs.setReference}
              className="select w-full"
              {...getReferenceProps()}
            >
              {selectedProfile ?? "Select profile"}
            </button>
            {!!selectOpen && (
              <FloatingPortal>
                <div
                  ref={refs.setFloating}
                  className="z-1"
                  style={floatingStyles}
                  {...getFloatingProps()}
                >
                  <div className="shadow-sm max-h-[50dvh] overflow-auto w-full">
                    <ul className="menu w-full bg-base-300 rounded-box p-2">
                      {supportedProfiles.map((profile) => (
                        <li key={profile.Name}>
                          <button
                            className="grid grid-cols-1 grid-flow-row auto-rows-max gap-0 dropdown-close"
                            onClick={() => {
                              field.onChange(profile.Name);
                              setSelectOpen(false);
                            }}
                          >
                            <span>{profile.Name}</span>
                            <span className="text-base-content/30 text-xs">
                              Last updated:{" "}
                              {updatedAtFormatter.format(
                                new Date(profile.Metadata.UpdatedAt),
                              )}
                            </span>
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
              </FloatingPortal>
            )}
          </>
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
            <button onClick={() => onCreateProfile(controller)}>
              Create new profile
            </button>
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
          <li>
            <button
              disabled={!selectedProfile}
              onClick={() => onDeleteProfileForController(controller)}
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
