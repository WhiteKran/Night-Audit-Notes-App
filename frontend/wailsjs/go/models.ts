export namespace main {
	
	export class Note {
	    id: number;
	    text: string;
	    createdAt: string;
	    updatedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new Note(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.text = source["text"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	    }
	}

}

