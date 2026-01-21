import { writable } from 'svelte/store';
import type { DownloadConfig } from './download';

export interface Preset {
	id: string;
	name: string;
	createdAt: string;
	config: DownloadConfig;
}

export const presets = writable<Preset[]>([]);
