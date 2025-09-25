import { useEffect, useRef } from "react";
import { useFieldArray, useForm } from "react-hook-form";
import { EventsOn } from "../../../wailsjs/runtime/runtime";
import { events } from "../../events";
import { main } from "../../../wailsjs/go/models";

export type Kind = "axis" | "button" | "hat";
export type CalibrationStateControl = {
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
export type CalibrationState = {
  name: string;
  controls: CalibrationStateControl[]
};

export type UseCalibrationFormType = ReturnType<typeof useCalibrationForm>

const EMPTY_CONTROL_STATE: Omit<CalibrationStateControl, 'kind' | 'index'> = {
  deadzone: 0,
  invert: false,
  name: "",
  value: 0,
  /* default to MAX_SAFE_INTEGER so any value is < min */
  min: Number.MAX_SAFE_INTEGER,
  idle: Number.MAX_SAFE_INTEGER,
  /* default to MIN_SAFE_INTEGER so any value is > max */
  max: Number.MIN_SAFE_INTEGER,
  override: false,
}

export const useCalibrationForm = () => {
  const form = useForm<CalibrationState>({
    defaultValues: {
      name: "",
      controls: [],
    },
  });

  useEffect(() => {
    return EventsOn(events.rawevent, (data: main.Interop_RawEvent) => {
      const controls = form.getValues('controls')
      const existingIndex = form.getValues('controls').findIndex((c) => (
        c.kind === data.Kind && c.index === data.Index
      ))

      const controlState: CalibrationStateControl = existingIndex === -1
        ? { ...EMPTY_CONTROL_STATE, kind: data.Kind as Kind, index: data.Index }
        : { ...controls[existingIndex] }

      if (!controlState.override) {
        Object.assign(controlState, {
          value: data.Value,
          min: Math.min(controlState.min, data.Value),
          max: Math.max(controlState.max, data.Value),
          idle: Math.min(controlState.min, data.Value),
        } as Partial<CalibrationStateControl>)
      } else {
        controlState.value = data.Value
      }

      if (existingIndex === -1) {
        form.setValue('controls', [...form.getValues('controls'), controlState], {
          shouldDirty: true,
          shouldTouch: true,
        })
      } else {
        form.setValue(`controls.${existingIndex}`, controlState, {
          shouldDirty: true,
          shouldTouch: true,
        })
      }
    });
  }, []);

  return form
}