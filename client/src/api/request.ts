export async function generalRequest(
    url: string,
    method: 'GET' | 'POST' | 'PUT' | 'DELETE',
    payload?: any
): Promise<Response> {
    const response = await fetch(url, {
        method: method,
        body: payload ? JSON.stringify(payload) : undefined,
        headers: {
            "Content-type": "application/json; charset=UTF-8"
        }
    });
    if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`)
    }
    return response;
}
