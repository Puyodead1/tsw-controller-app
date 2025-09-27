import { useEffect, useRef } from "react";
import { EventsOn } from "../../../wailsjs/runtime/runtime";
import { events } from "../../events";

export const LogsTab = () => {
  const logsRef = useRef<HTMLDivElement | null>(null)

  useEffect(() => {
    return EventsOn(events.log, (msg: string) => {
      if (logsRef.current) {
        const textNode = document.createTextNode(msg + '\n');
        logsRef.current.appendChild(textNode);
      }
    })
  }, [])

  return (
    <div>
      <div ref={logsRef} key="logs" className="whitespace-pre text-xs font-mono" />
    </div>
  );
};
