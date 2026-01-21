import { writable } from 'svelte/store';

export type Tab = 'download' | 'history' | 'presets' | 'settings';

export const activeTab = writable<Tab>('download');
