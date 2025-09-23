import { useForm } from "react-hook-form";
import { LogsTab } from "./tabs/logs/LogsTab";
import { MainTab } from "./tabs/main";
import { OnFrontendReady } from "../wailsjs/go/main/App";
import { CalibrationTab } from "./tabs/calibration";

OnFrontendReady();

const App = () => {
  const tabsForm = useForm<{ tab: "main" | "calibration" | "logs" }>({
    defaultValues: { tab: "main" },
  });
  const tab = tabsForm.watch("tab");
  console.log(tab)

  return (
    <div className="p-2">
      <div className="sticky top-2 tabs tabs-box">
        <input
          type="radio"
          className="tab"
          aria-label="Main"
          value="main"
          {...tabsForm.register("tab", { value: "main" })}
        />
        <input
          type="radio"
          className="tab"
          aria-label="Calibration"
          value="calibration"
          {...tabsForm.register("tab", { value: 'calibration' })}
        />
        <input
          type="radio"
          className="tab"
          aria-label="Logs"
          value="logs"
          {...tabsForm.register("tab", { value: 'logs' })}
        />
      </div>

      <div className="p-2">
        {tab === "main" && <MainTab />}
        {tab === 'calibration' && <CalibrationTab />}
        {tab === "logs" && <LogsTab />}
      </div>
    </div>
  );
};

export default App;
