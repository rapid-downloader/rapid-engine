export namespace client {
	
	
	export class Cookie {
	    name: string;
	    value: string;
	    path: string;
	    domain: string;
	    // Go type: time
	    expirationDate: any;
	    secure: boolean;
	    httpOnly: boolean;
	    sameSite: string;
	
	    static createFrom(source: any = {}) {
	        return new Cookie(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.value = source["value"];
	        this.path = source["path"];
	        this.domain = source["domain"];
	        this.expirationDate = this.convertValues(source["expirationDate"], null);
	        this.secure = source["secure"];
	        this.httpOnly = source["httpOnly"];
	        this.sameSite = source["sameSite"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
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
	export class Request {
	    url: string;
	    provider: string;
	    client?: string;
	    contentType?: string;
	    userAgent?: string;
	    cookies?: Cookie[];
	
	    static createFrom(source: any = {}) {
	        return new Request(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.url = source["url"];
	        this.provider = source["provider"];
	        this.client = source["client"];
	        this.contentType = source["contentType"];
	        this.userAgent = source["userAgent"];
	        this.cookies = this.convertValues(source["cookies"], Cookie);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
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

