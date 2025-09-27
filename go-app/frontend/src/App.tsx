import { useForm } from "react-hook-form";
import { MainTab } from "./tabs/main";
import { CalibrationTab } from "./tabs/calibration";
import { LogsTab } from "./tabs/logs";

const App = () => {
  const tabsForm = useForm<{ tab: "main" | "calibration" | "logs" }>({
    defaultValues: { tab: "main" },
  });
  const tab = tabsForm.watch("tab");

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
        {tab === 'logs' && <LogsTab />}
      </div>
    </div>
  );
};

export default App;
