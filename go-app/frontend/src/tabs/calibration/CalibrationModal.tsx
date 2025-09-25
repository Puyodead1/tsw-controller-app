import {
  MutableRefObject,
  useEffect,
  useImperativeHandle,
  useMemo,
  useRef,
  useState,
} from "react";
import { main } from "../../../wailsjs/go/models";
import {
  SaveCalibration,
  UnsubscribeRaw,
  SubscribeRaw,
} from "../../../wailsjs/go/main/App";
import { EventsOn } from "../../../wailsjs/runtime/runtime";
import { events } from "../../events";
import throttle from "just-throttle";
import { useForm } from "react-hook-form";
import { CalibrationModalForm } from "./CalibrationModalForm";

type Kind = "axis" | "button" | "hat";
type CalibrationState = {
  name: string;
  controls: Partial<
    Record<
      `${Kind}${number}`,
      {
        kind: Kind;
        index: number;
        value: number;
        name: string;
        min: number;
        max: number;
        idle: number;
        deadzone: number;
        invert: boolean;
        override: boolean;
      }
    >
  >;
};

type Props = {
  dialogRef: MutableRefObject<HTMLDialogElement | null>;
  controller: main.Interop_GenericController | null;
};

export const CalibrationModal = ({ dialogRef, controller }: Props) => {
  const ref = useRef<HTMLDialogElement | null>(null);

  const handleRef = (d: HTMLDialogElement | null) => {
    ref.current = d;
    dialogRef.current = d;
  };

  const handleClose = () => {
    ref.current?.close();
  };

  return (
    <dialog ref={handleRef} className="modal modal-s">
      <div className="modal-box w-11/12 max-w-5xl">
        {!!controller && (
          <CalibrationModalForm
            key={controller.GUID}
            controller={controller}
            onClose={handleClose}
          />
        )}
      </div>
    </dialog>
  );
};
