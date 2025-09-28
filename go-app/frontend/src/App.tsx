import { useForm } from "react-hook-form";
import { MainTab } from "./tabs/main";
import { CalibrationTab } from "./tabs/calibration";
import { LogsTab } from "./tabs/logs";
import { CabDebuggerTab } from "./tabs/cabdebugger";
import { SelfUpdateBanner } from "./SelfUpdateBanner";

const App = () => {
  const tabsForm = useForm<{
    tab: "main" | "calibration" | "cab_debugger" | "logs";
  }>({
    defaultValues: { tab: "main" },
  });
  const tab = tabsForm.watch("tab");

  return (
    <div className="p-2">
      <SelfUpdateBanner />

      <div className="sticky top-2 tabs tabs-box z-10">
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
          aria-label="Cab Debugger"
          value="cab_debugger"
          {...tabsForm.register("tab", { value: "cab_debugger" })}
        />
        <input
          type="radio"
          className="tab"
          aria-label="Calibration"
          value="calibration"
          {...tabsForm.register("tab", { value: "calibration" })}
        />
        <input
          type="radio"
          className="tab"
          aria-label="Logs"
          value="logs"
          {...tabsForm.register("tab", { value: "logs" })}
        />
      </div>

      <div className="p-2">
        {tab === "main" && <MainTab />}
        {tab === 'cab_debugger' && <CabDebuggerTab />}
        {tab === "calibration" && <CalibrationTab />}
        {tab === "logs" && <LogsTab />}
      </div>
    </div>
  );
};

export default App;
