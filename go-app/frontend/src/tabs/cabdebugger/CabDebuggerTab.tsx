import { useEffect } from "react";
import { EventsOn } from "../../../wailsjs/runtime/runtime";
import {
  GetSyncControlState,
  ResetSyncControlState,
} from "../../../wailsjs/go/main/App";
import { events } from "../../events";
import useSWR from "swr";

export const CabDebuggerTab = () => {
  const { data: syncControlState, mutate: refetchSyncControlState } = useSWR(
    "sync-control-state",
    () =>
      GetSyncControlState().then((r) =>
        r.sort((a, b) => a.Identifier.localeCompare(b.Identifier)),
      ),
  );

  const handleReset = () => {
    ResetSyncControlState().then(() => {
      refetchSyncControlState();
    });
  };

  useEffect(() => {
    return EventsOn(events.synccontrolstate, () => {
      refetchSyncControlState();
    });
  }, []);

  return (
    <div>
      <ul className="list bg-base-100 rounded-box shadow-md">
        {syncControlState?.map((controlState) => (
          <li key={controlState.Identifier} className="list-row">
            <div className="flex flex-col gap-2">
              <div className="grid grid-cols-2">
                <div>
                  <p className="text-slate-400">Sync Control Name</p>
                  <p>{controlState.Identifier}</p>
                </div>
                <div>
                  <p className="text-slate-400">Direct Control Name</p>
                  <p>{controlState.PropertyName}</p>
                </div>
              </div>
              <div className="grid grid-cols-2">
                <div>
                  <p className="text-slate-400">Current Value</p>
                  <p>{controlState.CurrentValue}</p>
                </div>
                <div>
                  <p className="text-slate-400">Current Normalized Value</p>
                  <p>{controlState.CurrentNormalizedValue}</p>
                </div>
              </div>
            </div>
          </li>
        ))}
      </ul>
      {!!syncControlState?.length && (
        <div className="sticky bottom-0 left-0 right-0 py-2 bg-[var(--root-bg,var(--color-base-100))] border-t border-t-base-100">
          <button className="btn btn-primary btn-xs" onClick={handleReset}>
            Reset State
          </button>
        </div>
      )}
    </div>
  );
};
