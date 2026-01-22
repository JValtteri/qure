export async function generalRequest(
    url: string,
    method: 'GET' | 'POST' | 'PUT' | 'DELETE',
    payload?: any
): Promise<Response> {
    let rq: RequestInit = {
        method: method,
        headers: {
            "Content-type": "application/json; charset=UTF-8"
        }
    }
    if (method != 'GET') {
        rq["body"] = payload ? JSON.stringify(payload) : undefined
    }
    const response = await fetch(url, rq);
    if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`)
    }
    return response;
}
