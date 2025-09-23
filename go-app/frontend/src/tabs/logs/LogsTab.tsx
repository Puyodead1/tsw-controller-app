import { SaveLogsAsFile } from "../../../wailsjs/go/main/App";
import { useStore } from "../../store";

export const LogsTab = () => {
  const { logs } = useStore();

  const handleDownload = () => {
    SaveLogsAsFile();
  };

  return (
    <div>
      <p className="whitespace-pre font-mono text-xs">{logs.join("")}</p>
      <div className="sticky pt-2 bottom-2 flex justify-end">
        <button className="btn btn-sm" onClick={handleDownload}>Save logs</button>
      </div>
    </div>
  );
};
