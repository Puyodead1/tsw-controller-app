import useSWR from "swr";
import { GetControllers } from "../wailsjs/go/main/App";
import { EventsOn } from "../wailsjs/runtime/runtime";
import { useEffect } from "react";

function App() {
  const { data: controllers, mutate: mutateControllers } = useSWR(
    "controllers",
    () => GetControllers(),
  );

  useEffect(() => {
    EventsOn("joydevice_added_or_removed", () => {
      mutateControllers();
    });
  }, []);

  return (
    <div className="p-6">
      <fieldset className="fieldset">
        <legend className="fieldset-legend">Select profile</legend>
        <select className="select">
          <option disabled selected>
            Auto-detect
          </option>
        </select>
        <span className="label">Auto-detect only works for certain supported controllers</span>
      </fieldset>
    </div>
  );
}

export default App;
