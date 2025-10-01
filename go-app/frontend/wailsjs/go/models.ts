export namespace config {
	
	export class Config_Controller_Profile_Control_Assignment_Action_Keys {
	    keys: string;
	    press_time?: number;
	    wait_time?: number;
	
	    static createFrom(source: any = {}) {
	        return new Config_Controller_Profile_Control_Assignment_Action_Keys(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.keys = source["keys"];
	        this.press_time = source["press_time"];
	        this.wait_time = source["wait_time"];
	    }
	}
	export class Config_Controller_Profile_Control_Assignment_DirectOrSyncControl_InputValue {
	    min: number;
	    max: number;
	    step?: number;
	    steps?: number[];
	    invert?: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Config_Controller_Profile_Control_Assignment_DirectOrSyncControl_InputValue(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.min = source["min"];
	        this.max = source["max"];
	        this.step = source["step"];
	        this.steps = source["steps"];
	        this.invert = source["invert"];
	    }
	}
	export class Config_Controller_Profile_Control_Assignment_SyncControl {
	    type: string;
	    identifier: string;
	    input_value: Config_Controller_Profile_Control_Assignment_DirectOrSyncControl_InputValue;
	    action_increase: Config_Controller_Profile_Control_Assignment_Action_Keys;
	    action_decrease: Config_Controller_Profile_Control_Assignment_Action_Keys;
	
	    static createFrom(source: any = {}) {
	        return new Config_Controller_Profile_Control_Assignment_SyncControl(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.identifier = source["identifier"];
	        this.input_value = this.convertValues(source["input_value"], Config_Controller_Profile_Control_Assignment_DirectOrSyncControl_InputValue);
	        this.action_increase = this.convertValues(source["action_increase"], Config_Controller_Profile_Control_Assignment_Action_Keys);
	        this.action_decrease = this.convertValues(source["action_decrease"], Config_Controller_Profile_Control_Assignment_Action_Keys);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Config_Controller_SDLMap_Control {
	    kind: string;
	    index: number;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new Config_Controller_SDLMap_Control(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.kind = source["kind"];
	        this.index = source["index"];
	        this.name = source["name"];
	    }
	}
	export class Config_Controller_SDLMap {
	    name: string;
	    usb_id: string;
	    data: Config_Controller_SDLMap_Control[];
	
	    static createFrom(source: any = {}) {
	        return new Config_Controller_SDLMap(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.usb_id = source["usb_id"];
	        this.data = this.convertValues(source["data"], Config_Controller_SDLMap_Control);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace main {
	
	export class Interop_ControllerCalibration_Control {
	    Kind: string;
	    Index: number;
	    Name: string;
	    Min: number;
	    Max: number;
	    Idle: number;
	    Deadzone: number;
	    Invert: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Interop_ControllerCalibration_Control(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Kind = source["Kind"];
	        this.Index = source["Index"];
	        this.Name = source["Name"];
	        this.Min = source["Min"];
	        this.Max = source["Max"];
	        this.Idle = source["Idle"];
	        this.Deadzone = source["Deadzone"];
	        this.Invert = source["Invert"];
	    }
	}
	export class Interop_ControllerCalibration {
	    Name: string;
	    UsbId: string;
	    Controls: Interop_ControllerCalibration_Control[];
	
	    static createFrom(source: any = {}) {
	        return new Interop_ControllerCalibration(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.UsbId = source["UsbId"];
	        this.Controls = this.convertValues(source["Controls"], Interop_ControllerCalibration_Control);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class Interop_ControllerConfiguration {
	    Calibration: Interop_ControllerCalibration;
	    SDLMapping: config.Config_Controller_SDLMap;
	
	    static createFrom(source: any = {}) {
	        return new Interop_ControllerConfiguration(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Calibration = this.convertValues(source["Calibration"], Interop_ControllerCalibration);
	        this.SDLMapping = this.convertValues(source["SDLMapping"], config.Config_Controller_SDLMap);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Interop_GenericController {
	    GUID: string;
	    UsbID: string;
	    Name: string;
	    IsConfigured: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Interop_GenericController(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.GUID = source["GUID"];
	        this.UsbID = source["UsbID"];
	        this.Name = source["Name"];
	        this.IsConfigured = source["IsConfigured"];
	    }
	}
	export class Interop_Profile {
	    Name: string;
	
	    static createFrom(source: any = {}) {
	        return new Interop_Profile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	    }
	}
	export class Interop_RawEvent {
	    GUID: string;
	    UsbID: string;
	    Kind: string;
	    Index: number;
	    Value: number;
	    Timestamp: number;
	
	    static createFrom(source: any = {}) {
	        return new Interop_RawEvent(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.GUID = source["GUID"];
	        this.UsbID = source["UsbID"];
	        this.Kind = source["Kind"];
	        this.Index = source["Index"];
	        this.Value = source["Value"];
	        this.Timestamp = source["Timestamp"];
	    }
	}

}

export namespace profile_runner {
	
	export class SyncController_ControlState {
	    Identifier: string;
	    PropertyName: string;
	    CurrentValue: number;
	    CurrentNormalizedValue: number;
	    TargetValue: number;
	    Moving: number;
	    ControlProfile: config.Config_Controller_Profile_Control_Assignment_SyncControl;
	
	    static createFrom(source: any = {}) {
	        return new SyncController_ControlState(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Identifier = source["Identifier"];
	        this.PropertyName = source["PropertyName"];
	        this.CurrentValue = source["CurrentValue"];
	        this.CurrentNormalizedValue = source["CurrentNormalizedValue"];
	        this.TargetValue = source["TargetValue"];
	        this.Moving = source["Moving"];
	        this.ControlProfile = this.convertValues(source["ControlProfile"], config.Config_Controller_Profile_Control_Assignment_SyncControl);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

