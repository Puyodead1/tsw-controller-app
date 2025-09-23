import { useEffect, useState } from "react";
import { main } from "../../../wailsjs/go/models";
import useSWR from "swr";
import { GetControllers } from "../../../wailsjs/go/main/App";

export const CalibrationTab = () => {
  const { data: controllers } = useSWR("controllers", () => GetControllers(), {
    revalidateOnMount: true,
  });

  return (
    <div>
      <ul className="list bg-base-100 rounded-box shadow-md">
        {controllers?.map((c) => (
          <li key={c.Name} className="list-row">
            <div className="list-col-grow">
              <div>{c.Name}</div>
            </div>
            <div>
              {c.IsConfigured && (
                <div className="tooltip" data-tip="Re-configure">
                  <button className="btn btn-success btn-soft btn-xs">
                    Configured
                  </button>
                </div>
              )}
              {!c.IsConfigured && (
                <div className="tooltip" data-tip="Configure now">
                  <button className="btn btn-error btn-soft btn-xs">
                    Unconfigured
                  </button>
                </div>
              )}
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
};
