import {getContext, setContext} from "svelte";

type Code = {
    key: string
    content: string
}

export class API {
    baseURL: string

    constructor(baseURL?: string) {
        this.baseURL = baseURL || 'http://localhost:8080'
    }

    static provide = (endpoint?:string) => {
        const api = new API(endpoint)
        setContext('api', api)

        return api
    }

    static use = (): API => {
        return getContext('api')
    }

    buildUrl(path: string) {
        return new URL(path, this.baseURL)
    }

    async request(
        method: string,
        path: string,
        payload?: string,
        headers?: { [key: string]: string }
    ) {
        try {
            const url = this.buildUrl(path)
            const response = await fetch(url, {
                method: method,
                headers: {
                    'Content-Type': 'text/plain',
                    ...headers,
                },
                body: payload,
                credentials: 'include',
            })
            if (!response.ok) {
                throw new Error('Network response was not OK')
            }
            if (response.headers.get("content-type") === "application/json"){
                return response.json()
            }
            return response.text()
        } catch (err) {
            return err
        }
    }

    async listCodes(): Promise<Code[]> {
        const response = await this.request('GET', "/codes")
        if (response.ok) {
            return response as Code[]
        }
        return
    }

    async getContent(key: string): Promise<string> {
        try {
            const url = this.buildUrl(`/codes/content?key=${key}`)
            const response = await fetch(url, {
                method: 'GET',
                credentials: 'include',
            })
            if (response.ok) {
                return response.text()
            }
        } catch(err){
            throw new Error(err)
        }
    }
    async deleteContent(key: string) {
        try {
            const url = this.buildUrl(`/codes/content?key=${key}`)
            const response = await fetch(url, {
                method: 'DELETE',
                credentials: 'include',
            })
            if (response.ok) {
                return
            }
        } catch(err) {
            throw new Error(err)
        }
    }
    async putCode(key: string, content: string) {
        return this.request('POST', `/codes?key=${key}`, content)
    }
}