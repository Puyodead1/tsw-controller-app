import { useEffect, useMemo } from "react";
import { GetCabControlState } from "../../../wailsjs/go/main/App";
import useSWR from "swr";
import { useForm } from "react-hook-form";

export const CabDebuggerTab = () => {
  const { register, watch } = useForm<{ query: string }>({
    defaultValues: { query: "" },
  });
  const { data: cabControlState, mutate: refetchCabControlState } = useSWR(
    "cab-control-state",
    () => GetCabControlState(),
    { revalidateOnMount: true },
  );

  const query = watch("query");
  const sortedControls = useMemo(
    () =>
      cabControlState?.Controls.filter((c) =>
        [c.Identifier, c.PropertyName].some((t) =>
          t.toLowerCase().includes(query.toLowerCase()),
        ),
      ).sort((a, b) =>
        `${a.PropertyName}_${a.Identifier}`.localeCompare(
          `${b.PropertyName}_${b.Identifier}`,
        ),
      ),
    [cabControlState?.Controls, query],
  );

  useEffect(() => {
    let interval: ReturnType<typeof setInterval> | null = null;
    interval = setInterval(() => {
      refetchCabControlState();
    }, 100);
    return () => {
      if (interval) clearInterval(interval);
    };
  }, [refetchCabControlState]);

  return (
    <div>
      {!cabControlState?.Controls?.length && (
        <div className="py-12 text-center">
          <p className="text-base-content/50 text-sm">
            Waiting for cab state...
          </p>
        </div>
      )}
      {!!cabControlState?.Controls?.length && (
        <div className="p-4">
          <input
            className="input w-full"
            placeholder="Search for control(s)"
            {...register("query")}
          />
        </div>
      )}
      <ul className="list bg-base-100 rounded-box shadow-md">
        {sortedControls?.map((controlState) => (
          <li key={controlState.PropertyName} className="list-row">
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
                  <p>{controlState.CurrentValue.toFixed(4)}</p>
                </div>
                <div>
                  <p className="text-slate-400">Current Normalized Value</p>
                  <p>{controlState.CurrentNormalizedValue.toFixed(4)}</p>
                </div>
              </div>
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
};
