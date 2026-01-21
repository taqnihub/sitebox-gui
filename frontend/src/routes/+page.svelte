<script lang="ts">
	import DownloadForm from '$lib/components/DownloadForm.svelte';
	import ProgressPanel from '$lib/components/ProgressPanel.svelte';
	import LogOutput from '$lib/components/LogOutput.svelte';
	import HistoryList from '$lib/components/HistoryList.svelte';
	import PresetManager from '$lib/components/PresetManager.svelte';
	import SettingsPanel from '$lib/components/SettingsPanel.svelte';
	import { activeTab } from '$lib/stores/ui';

	const tabs = [
		{ id: 'download', label: 'Download' },
		{ id: 'history', label: 'History' },
		{ id: 'presets', label: 'Presets' },
		{ id: 'settings', label: 'Settings' }
	] as const;
</script>

<div class="app-container">
	<header class="header">
		<h1 class="title">SiteBox</h1>
		<nav class="tabs">
			{#each tabs as tab}
				<button
					class="tab"
					class:active={$activeTab === tab.id}
					onclick={() => activeTab.set(tab.id)}
				>
					{tab.label}
				</button>
			{/each}
		</nav>
	</header>

	<div class="content">
		{#if $activeTab === 'download'}
			<div class="download-view">
				<div class="download-main">
					<DownloadForm />
					<ProgressPanel />
				</div>
				<div class="download-log">
					<LogOutput />
				</div>
			</div>
		{:else if $activeTab === 'history'}
			<HistoryList />
		{:else if $activeTab === 'presets'}
			<PresetManager />
		{:else if $activeTab === 'settings'}
			<SettingsPanel />
		{/if}
	</div>
</div>

<style>
	.app-container {
		display: flex;
		flex-direction: column;
		height: 100%;
		background-color: var(--bg-primary);
	}

	.header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: var(--spacing-md) var(--spacing-lg);
		border-bottom: 1px solid var(--border-color);
		background-color: var(--bg-secondary);
		-webkit-app-region: drag;
	}

	.title {
		font-size: 20px;
		font-weight: 600;
		margin: 0;
		color: var(--text-primary);
	}

	.tabs {
		display: flex;
		gap: var(--spacing-xs);
		-webkit-app-region: no-drag;
	}

	.tab {
		padding: var(--spacing-sm) var(--spacing-md);
		border: none;
		border-radius: var(--radius-md);
		background: transparent;
		color: var(--text-secondary);
		font-size: 14px;
		font-weight: 500;
		cursor: pointer;
		transition: all var(--transition-fast);
	}

	.tab:hover {
		background-color: var(--bg-tertiary);
		color: var(--text-primary);
	}

	.tab.active {
		background-color: var(--accent-color);
		color: white;
	}

	.content {
		flex: 1;
		overflow: hidden;
		padding: var(--spacing-lg);
	}

	.download-view {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: var(--spacing-lg);
		height: 100%;
	}

	.download-main {
		display: flex;
		flex-direction: column;
		gap: var(--spacing-lg);
		overflow-y: auto;
	}

	.download-log {
		display: flex;
		flex-direction: column;
		overflow: hidden;
	}

	@media (max-width: 900px) {
		.download-view {
			grid-template-columns: 1fr;
			grid-template-rows: auto 1fr;
		}
	}
</style>
