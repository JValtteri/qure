export async function generalRequest(url: string, method: string, payload: string): Promise<Response> {
    try {
        const response = await fetch(url, {
            method: method,
            body: payload,
            headers: {
                "Content-type": "application/json; charset=UTF-8"
            }
        });
        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`)
        }
        return response;
    } catch (error) {
        console.error(error);
        return new Response();
    }
}
