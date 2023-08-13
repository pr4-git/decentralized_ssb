export namespace dba {
	
	export class Peer {
	    id: string;
	    pubkey: number[];
	    alias: string;
	
	    static createFrom(source: any = {}) {
	        return new Peer(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.pubkey = source["pubkey"];
	        this.alias = source["alias"];
	    }
	}
	export class Post {
	    id: string;
	    content: string;
	    hash: number[];
	    signature: number[];
	    author: number[];
	    parent: number[];
	
	    static createFrom(source: any = {}) {
	        return new Post(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.content = source["content"];
	        this.hash = source["hash"];
	        this.signature = source["signature"];
	        this.author = source["author"];
	        this.parent = source["parent"];
	    }
	}

}

