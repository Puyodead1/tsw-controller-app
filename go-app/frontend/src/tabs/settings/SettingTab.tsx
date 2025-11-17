import { DeepPartial, useForm } from "react-hook-form";
import {
  GetTSWAPIKeyLocation,
  SetTSWAPIKeyLocation,
} from "../../../wailsjs/go/main/App";
import { useEffect } from "react";
import { alert } from "../../utils/alert";
import debounce from "just-debounce";

type FormValues = {
  tswApiKeyLocation: string;
};

export const SettingsTab = () => {
  const { register, watch } = useForm<FormValues>({
    defaultValues: async () => ({
      tswApiKeyLocation: await GetTSWAPIKeyLocation(),
    }),
  });

  useEffect(() => {
    const watchFunc = debounce(async (values: DeepPartial<FormValues>) => {
      if (
        values.tswApiKeyLocation &&
        values.tswApiKeyLocation.endsWith("CommAPIKey.txt") &&
        values.tswApiKeyLocation !== (await GetTSWAPIKeyLocation())
      ) {
        SetTSWAPIKeyLocation(values.tswApiKeyLocation).then(() => {
          alert("Saved API key location", "success");
        });
      }
    }, 300);
    watch(watchFunc);
  }, [watch]);

  return (
    <div>
      <fieldset className="fieldset">
        <label className="fieldset-legend">TSW API Key Location</label>
        <input className="input w-full" {...register("tswApiKeyLocation")} />
        <p className="fieldset-label whitespace-normal">
          If the location has not been auto-detected you will need to enter it
          manually here. The API key is only requred for the "api_control"
          control mode.
        </p>
      </fieldset>
    </div>
  );
};
