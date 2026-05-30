export namespace main {
	
	export class Note {
	    id: number;
	    text: string;
	    createdAt: string;
	    updatedAt: string;
	    isLocked: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Note(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.text = source["text"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	        this.isLocked = source["isLocked"];
	    }
	}
	export class WindowSize {
	    width: number;
	    height: number;
	
	    static createFrom(source: any = {}) {
	        return new WindowSize(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.width = source["width"];
	        this.height = source["height"];
	    }
	}

}

