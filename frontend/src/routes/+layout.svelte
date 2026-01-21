<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import { theme, settings } from '$lib/stores/settings';
	import {
		status,
		downloadedCount,
		errorCount,
		currentUrl,
		addLog,
		type ProgressEvent
	} from '$lib/stores/download';

	let { children } = $props();

	// Load settings from backend on app startup
	async function loadSettings() {
		try {
			// @ts-ignore - Wails runtime
			if (typeof window.go !== 'undefined') {
				// @ts-ignore
				const saved = await window.go.gui.App.GetSettings();
				if (saved) {
					settings.set(saved);
				}
			}
		} catch (err) {
			console.error('Failed to load settings:', err);
		}
	}

	onMount(() => {
		// Load settings on startup
		loadSettings();

		// Apply theme on mount
		const unsubscribeTheme = theme.subscribe((value) => {
			if (value === 'system') {
				const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
				document.documentElement.setAttribute('data-theme', prefersDark ? 'dark' : 'light');
			} else {
				document.documentElement.setAttribute('data-theme', value);
			}
		});

		// Listen for Wails events
		let unsubscribeProgress: (() => void) | null = null;

		// @ts-ignore - Wails runtime
		if (typeof window.runtime !== 'undefined') {
			// @ts-ignore
			unsubscribeProgress = window.runtime.EventsOn(
				'download:progress',
				(event: ProgressEvent) => {
					currentUrl.set(event.url);
					downloadedCount.set(event.downloaded);
					errorCount.set(event.errors);

					switch (event.type) {
						case 'start':
							status.set('downloading');
							addLog('info', event.message);
							break;
						case 'download':
							addLog('success', event.message, event.url);
							break;
						case 'error':
							addLog('error', event.message, event.url);
							break;
						case 'complete':
							status.set('completed');
							addLog('info', event.message);
							break;
					}
				}
			);
		}

		return () => {
			unsubscribeTheme();
			if (unsubscribeProgress) {
				unsubscribeProgress();
			}
		};
	});
</script>

<main>
	{@render children()}
</main>

<style>
	main {
		height: 100vh;
		display: flex;
		flex-direction: column;
	}
</style>
