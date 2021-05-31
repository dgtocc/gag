//#region Base


var apibase="";

export function SetAPIBase(s:string){
	apibase=s;
}

export function GetAPIBase(): string{
	return apibase;
}

let REGEX_DATE = /^(\d{4})-(\d{2})-(\d{2})T(\d{2}):(\d{2}):(\d{2}(?:\.\d*)?)(Z|([+\-])(\d{2}):(\d{2}))$/

type HTMLMethod = "GET" | "POST" | "PUT" | "DELETE" | "HEAD" | "TRACE"

async function Invoke(path: string, method: HTMLMethod, body?: any): Promise<Response> {
	let jbody = undefined
	let init = {method: method, mode: "cors", credentials: "include", withCredentials: true}
	if (!!body) {
		let jbody = JSON.stringify(body)
		//@ts-ignore
		init.body = jbody
	}
	if (apibase.endsWith("/") && path.startsWith("/")) {
		path = path.substr(1, path.length)
	}
	let fpath = (apibase + path)
	//@ts-ignore
	let res = await fetch(fpath, init)

	return res
}
 
async function InvokeJSON(path: string, method: HTMLMethod, body?: any): Promise<any> {

	let txt = await InvokeTxt(path, method, body)
	if (txt == "") {
		txt = "{}"
	}
	let ret = JSON.parse(txt, (k: string, v: string) => {
		if (REGEX_DATE.exec(v)) {
			return new Date(v)
		}
		return v
	})

	return ret
}

async function InvokeTxt(path: string, method: HTMLMethod, body?: any): Promise<string> {
	//@ts-ignore
	let res = await Invoke(path, method, body)

	let txt = await res.text()

	if (res.status < 200 || res.status >= 400) {
		// webix.alert("API Error:" + res.status + "\n" + txt)
		console.error("API Error:" + res.status + "\n" + txt)
		let e = new Error(txt)
		throw e
	}

	return txt
}

async function InvokeOk(path: string, method: HTMLMethod, body?: any): Promise<boolean> {

	//@ts-ignore
	let res = await Invoke(path, method, body)

	let txt = await res.text()
	if (res.status >= 400) {
		console.error("API Error:" + res.status + "\n" + txt)
		return false
	}
	return true
}

//#endregion

//#region Types
export interface ASimpleReq {
	data ?: string
}

export interface ASimpleRes {
	data ?: string
}

export interface AComplexReq {
	city ?: string
	iscondo ?: boolean
	someweirdtest ?: string
	recursive ?: {[s:string]:AComplexReq}
	arrofpstr ?: string
	country ?: string
	housenumber ?: number
	when ?: Date
	some ?: crypto.Decrypter
}

//#endregion

//#region Methods
/**
ApiMethod01*/
export async function ApiMethod01(req:ASimpleReq):Promise<ASimpleRes>{
	return InvokeJSON("/someapi","POST",req)
}

/**
ApiMethod02*/
export async function ApiMethod02(req:AComplexReq):Promise<ASimpleRes>{
	return InvokeJSON("/someapi2","GET",req)
}

//#endregion
