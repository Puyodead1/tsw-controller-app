export namespace main {
	
	export class Interop_GenericController {
	    Name: string;
	    IsConfigured: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Interop_GenericController(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.IsConfigured = source["IsConfigured"];
	    }
	}

}

