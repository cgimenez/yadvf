import { writable } from 'svelte/store';

export const mutations = writable([])
export const fetching = writable(false)
export const distincts = writable({
  Cultures: [],
  Cultures_speciales: [],
  Locaux: []
})