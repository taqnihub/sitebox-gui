import { writable } from 'svelte/store';
import type { DownloadConfig } from './download';

export interface HistoryEntry {
	id: string;
	url: string;
	domain: string;
	output: string;
	startedAt: string;
	completedAt?: string;
	downloaded: number;
	errors: number;
	status: 'completed' | 'error' | 'stopped';
	config: DownloadConfig;
}

export const history = writable<HistoryEntry[]>([]);
