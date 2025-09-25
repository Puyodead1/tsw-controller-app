import {
  CalibrationStateControl,
  UseCalibrationFormType,
} from "./useCalibrationForm";

type Props = {
  form: UseCalibrationFormType;
  index: number;
  field: CalibrationStateControl;
};

export const CalibrationModalFormControl = ({ form, index, field }: Props) => {
  return (
    <div
      key={`${field.kind}_${field.index}`}
      className="card card-sm shadow-sm"
    >
      <div className="card-body">
        <div>
          <div className="flex flex-col basis-full gap-2">
            <div className="flex justify-between items-center">
              <div>
                {field.kind} {field.index}
              </div>
              <div>
                <kbd className="kbd kbd-sm">{field.value}</kbd>
              </div>
            </div>
            {field.kind === "axis" && (
              <div>
                {field.invert && (
                  <progress
                    className="progress progress-primary w-full"
                    value={
                      Math.abs(field.min) +
                      field.max -
                      (field.value + Math.abs(field.min))
                    }
                    max={Math.abs(field.min) + field.max}
                  ></progress>
                )}
                {!field.invert && (
                  <progress
                    className="progress progress-primary w-full"
                    value={field.value + Math.abs(field.min)}
                    max={Math.abs(field.min) + field.max}
                  ></progress>
                )}
              </div>
            )}
            <div>
              <label className="input input-xs">
                Name
                <input
                  type="text"
                  className="grow"
                  {...form.register(`controls.${index}.name`, {
                    required: true,
                  })}
                />
              </label>
            </div>
            {field.kind === "axis" && (
              <>
                <div className="grid grid-cols-2 grid-flow-row auto-rows-max gap-2">
                  <label className="input input-xs">
                    Min
                    <input
                      type="number"
                      className="grow"
                      disabled={!field.override}
                      {...form.register(`controls.${index}.min`, {
                        valueAsNumber: true,
                        required: true,
                      })}
                    />
                  </label>
                  <label className="input input-xs">
                    Max
                    <input
                      type="number"
                      className="grow"
                      disabled={!field.override}
                      {...form.register(`controls.${index}.max`, {
                        valueAsNumber: true,
                        required: true,
                      })}
                    />
                  </label>
                  <label className="input input-xs">
                    Idle
                    <input
                      type="number"
                      className="grow"
                      disabled={!field.override}
                      {...form.register(`controls.${index}.idle`, {
                        valueAsNumber: true,
                        required: true,
                      })}
                    />
                  </label>
                  <label className="input input-xs">
                    Deadzone
                    <input
                      type="number"
                      className="grow"
                      disabled={!field.override}
                      {...form.register(`controls.${index}.deadzone`, {
                        valueAsNumber: true,
                        required: true,
                      })}
                    />
                  </label>
                </div>
                <div>
                  <label className="label">
                    <input
                      type="checkbox"
                      className="checkbox checkbox-xs"
                      {...form.register(`controls.${index}.invert`)}
                    />
                    Invert
                  </label>
                </div>
                <div>
                  <label className="label">
                    <input
                      type="checkbox"
                      className="checkbox checkbox-xs"
                      {...form.register(`controls.${index}.override`)}
                    />
                    Override calibration values
                  </label>
                </div>
              </>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};
