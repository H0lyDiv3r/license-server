export namespace domain {
	
	export class LicenseClaims {
	    license_id: number;
	    user_id: number;
	    machine_id: string;
	    status: string;
	    iss?: string;
	    sub?: string;
	    aud?: string[];
	    // Go type: jwt
	    exp?: any;
	    // Go type: jwt
	    nbf?: any;
	    // Go type: jwt
	    iat?: any;
	    jti?: string;
	
	    static createFrom(source: any = {}) {
	        return new LicenseClaims(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.license_id = source["license_id"];
	        this.user_id = source["user_id"];
	        this.machine_id = source["machine_id"];
	        this.status = source["status"];
	        this.iss = source["iss"];
	        this.sub = source["sub"];
	        this.aud = source["aud"];
	        this.exp = this.convertValues(source["exp"], null);
	        this.nbf = this.convertValues(source["nbf"], null);
	        this.iat = this.convertValues(source["iat"], null);
	        this.jti = source["jti"];
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

