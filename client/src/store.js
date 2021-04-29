import { writable } from 'svelte/store';

export const state = writable();

if (sessionStorage.GetItem("uuid")) {
    export const uuid = sessionStorage.GetItem("uuid")
}
