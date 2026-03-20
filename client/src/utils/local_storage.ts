
export function setLocalStorage(name: string, value: string) {
    localStorage.setItem(name, value);
}

export function clearLocalStorage(name: string) {
    localStorage.removeItem(name);
}

export function getLocalStorage(name: string): string {
    const storeItem = localStorage.getItem(name)
    if (!storeItem) {
        return "";
    }
    return storeItem;
}
