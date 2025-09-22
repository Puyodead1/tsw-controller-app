import useSWR from "swr";
import { GetControllers } from "../wailsjs/go/main/App";
import { EventsOn } from "../wailsjs/runtime/runtime"
import { useEffect } from "react";

function App() {
  const { data: controllers, mutate: mutateControllers } = useSWR('controllers', () => GetControllers())



useEffect(() => {
  EventsOn("joydevice_added_or_removed", () => {
    mutateControllers()
  });
}, [])

  return (
    <div className="p-6">
      <select className="select w-full">
        <option disabled selected>
          Select a controller
        </option>
        {controllers?.map((c) => (
          <option key={c.Name}>
            {c.Name}
          </option>
        ))}
      </select>
    </div>
  );
}

export default App;
