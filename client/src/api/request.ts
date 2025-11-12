export async function generalRequest(url: string, method: string, payload: string) {
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

        return await response.json();
    } catch (error) {
        console.error(error);
        return {};
    }
}
