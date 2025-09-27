import useSWR from "swr";
import {
  GetProfiles,
  LoadConfiguration,
  SelectProfile,
  ClearProfile,
} from "../../../wailsjs/go/main/App";
import { useEffect } from "react";
import { EventsOn } from "../../../wailsjs/runtime/runtime";
import { events } from "../../events";
import { useForm } from "react-hook-form";

type FormValues = {
  profile: string;
};

export const MainTab = () => {
  const { data: profiles, mutate: refetchProfiles } = useSWR(
    "profiles",
    () => GetProfiles(),
    {
      revalidateOnMount: true,
    },
  );

  const { register, watch, getValues } = useForm<FormValues>({
    defaultValues: {
      profile: "",
    },
  });

  const handleReload = () => {
    LoadConfiguration();
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
      <fieldset className="fieldset">
        <legend className="fieldset-legend">Select profile</legend>
        <select className="select w-full" {...register("profile")}>
          <option selected value="">
            Auto-detect
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
      <button className="btn btn-sm" onClick={handleReload}>
        Reload Configurations
      </button>
    </div>
  );
};
