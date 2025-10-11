# CHANGELOG

## 1.0.0
**MAJOR RELEASE - BREAKING CHANGES**
- Complete runtime and UI overhaul
- Added visual Cab Debugger
- Added visual Calibration mode
- Added controller specific profile selection
- Added conditional assignments

## 0.2.5
- Fix calibration for calibrating multiple controllers.

## 0.2.4
- Fix calibration mode not exiting and writing files.

## 0.2.3
- Update the mixing of `null` values to act as free range zones instead of automatic interpolation zones. This makes for smoother actions between detents. Eg, the following steps value: `[0.0, null, 0.5, 0.6, null, 1.0]` - will snap to `0.5` and `0.6` but allow free range of motion between `0.0` and `0.5` and `0.6` and `1.0`.

## 0.2.2
- Improve performance overhead by reducing controller polling.

## 0.2.1

- Add support for mixing `null` values in the `steps` array for direct control assignments. This new feature allows automatic interpolation between step values without having to manually calculate the steps. Eg: `[0, null, null, 1]` will result in the following actual step list: `[0, 0.33, 0.66, 1]`. This is useful for levers where you want a combination of semi-free range and stepped values. (ie: North American suppression steps mixed with percentage based free range)

## 0.2.0

- Add "relative" option for direct control actions outside of direct control assignments. This allows relative value setting (ie: set the value to be -0.4 below the current value).

## 0.1.7

- Update usb_id check to be case insensitive
- Updated `Tick` and `None` checks using `Fname`(thanks to [@UE4SS](https://github.com/UE4SS))
