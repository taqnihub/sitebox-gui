import { writable, derived } from 'svelte/store';

export interface DownloadConfig {
	url: string;
	domain: string;
	output: string;
	maxDepth: number;
	concurrent: number;
	retries: number;
	includeImages: boolean;
	blacklist: string[];
	extensions: string[];
	imageExtensions: string[];
}

export interface ProgressEvent {
	type: 'start' | 'download' | 'error' | 'complete';
	url: string;
	filePath: string;
	downloaded: number;
	errors: number;
	message: string;
}

export interface LogEntry {
	id: string;
	timestamp: Date;
	type: 'info' | 'success' | 'error' | 'warning';
	message: string;
	url?: string;
}

export type DownloadStatus = 'idle' | 'downloading' | 'completed' | 'error' | 'stopped';

export const defaultConfig: DownloadConfig = {
	url: '',
	domain: '',
	output: '',
	maxDepth: 50,
	concurrent: 5,
	retries: 3,
	includeImages: false,
	blacklist: ['https://github'],
	extensions: ['.js', '.css', '.woff2', '.woff', '.ttf'],
	imageExtensions: ['.png', '.svg', '.jpg', '.jpeg', '.gif', '.webp']
};

// Stores
export const config = writable<DownloadConfig>({ ...defaultConfig });
export const status = writable<DownloadStatus>('idle');
export const downloadedCount = writable<number>(0);
export const errorCount = writable<number>(0);
export const currentUrl = writable<string>('');
export const logs = writable<LogEntry[]>([]);

// Derived stores
export const isDownloading = derived(status, ($status) => $status === 'downloading');
export const canStart = derived(
	[config, status],
	([$config, $status]) => $config.url !== '' && $config.output !== '' && $status !== 'downloading'
);

// Helper functions
export function addLog(type: LogEntry['type'], message: string, url?: string) {
	const entry: LogEntry = {
		id: crypto.randomUUID(),
		timestamp: new Date(),
		type,
		message,
		url
	};
	logs.update((l) => [...l, entry]);
}

export function clearLogs() {
	logs.set([]);
}

export function resetProgress() {
	downloadedCount.set(0);
	errorCount.set(0);
	currentUrl.set('');
	clearLogs();
}

export function extractDomain(url: string): string {
	try {
		const parsed = new URL(url);
		return `${parsed.protocol}//${parsed.host}`;
	} catch {
		return '';
	}
}
