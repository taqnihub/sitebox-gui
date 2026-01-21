import { writable } from 'svelte/store';

export type Theme = 'light' | 'dark' | 'system';

export interface Settings {
	theme: Theme;
	defaultOutput: string;
	defaultConcurrent: number;
	defaultMaxDepth: number;
	defaultRetries: number;
	includeImages: boolean;
	showNotifications: boolean;
}

export const defaultSettings: Settings = {
	theme: 'system',
	defaultOutput: '',
	defaultConcurrent: 5,
	defaultMaxDepth: 50,
	defaultRetries: 3,
	includeImages: false,
	showNotifications: true
};

export const settings = writable<Settings>(defaultSettings);
export const theme = writable<Theme>('system');

// Sync theme with settings
settings.subscribe((s) => {
	theme.set(s.theme);
});
