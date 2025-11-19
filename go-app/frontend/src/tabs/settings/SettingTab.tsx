import { DeepPartial, useForm } from "react-hook-form";
import {
  GetPreferredControlMode,
  GetTSWAPIKeyLocation,
  SetPreferredControlMode,
  SetTSWAPIKeyLocation,
} from "../../../wailsjs/go/main/App";
import { useEffect } from "react";
import { alert } from "../../utils/alert";
import debounce from "just-debounce";
import { BrowserOpenURL } from "../../../wailsjs/runtime/runtime";

type FormValues = {
  tswApiKeyLocation: string;
  preferredControlMode: "direct_control" | "sync_control" | "api_control";
};

export const SettingsTab = () => {
  const { register, formState, reset, handleSubmit } = useForm<FormValues>({
    defaultValues: async () => ({
      tswApiKeyLocation: await GetTSWAPIKeyLocation(),
      preferredControlMode:
        (await GetPreferredControlMode()) as FormValues["preferredControlMode"],
    }),
  });

  const handleOpenForumLink = () => {
    BrowserOpenURL(
      "https://forums.dovetailgames.com/threads/train-sim-world-api-support.94488/",
    );
  };

  const handleSubmitSuccess = async (values: FormValues) => {
    const promises: Promise<void>[] = [];

    if (
      values.tswApiKeyLocation &&
      values.tswApiKeyLocation.endsWith("CommAPIKey.txt") &&
      values.tswApiKeyLocation !== (await GetTSWAPIKeyLocation())
    ) {
      promises.push(SetTSWAPIKeyLocation(values.tswApiKeyLocation));
    }

    if (
      values.preferredControlMode &&
      values.preferredControlMode !== (await GetPreferredControlMode())
    ) {
      promises.push(SetPreferredControlMode(values.preferredControlMode));
    }

    if (promises.length) {
      Promise.all(promises).then(() => {
        reset(values);
        alert("Saved settings", "success");
      });
    }
  };

  return (
    <form
      className="grid grid-cols-1 grid-flow-row auto-rows-max gap-2"
      onSubmit={handleSubmit(handleSubmitSuccess)}
    >
      <fieldset className="fieldset">
        <label htmlFor="preferred-control-mode" className="fieldset-legend">
          Preferred Control Mode
        </label>
        <select className="select w-full" {...register("preferredControlMode")}>
          <option value="direct_control">Direct Control</option>
          <option value="sync_control">Sync Control</option>
          <option value="api_control">API Control</option>
        </select>
        <p className="fieldset-label whitespace-normal">
          Sets which control mode to prefer if multiple are defined
        </p>
      </fieldset>
      <fieldset className="fieldset">
        <label className="fieldset-legend">TSW API Key Location</label>
        <input className="input w-full" {...register("tswApiKeyLocation")} />
        <p className="fieldset-label whitespace-normal">
          If the location has not been auto-detected you will need to enter it
          manually here. The API key is only requred for the "api_control"
          control mode.
        </p>
      </fieldset>
      <div role="alert" className="alert">
        <span>
          <strong>TSW API Notice</strong>
          <br />
          The API connection only works if -HTTPAPI is enabled in Train Sim
          World. You can find instructions in the linked PDF on the{" "}
          <button type="button" className="link" onClick={handleOpenForumLink}>
            online forum
          </button>
          .
        </span>
      </div>
      <div className="flex justify-end">
        <button
          type="submit"
          className="btn btn-primary"
          disabled={
            formState.disabled || !formState.isDirty || formState.isSubmitting
          }
        >
          Save
        </button>
      </div>
    </form>
  );
};
