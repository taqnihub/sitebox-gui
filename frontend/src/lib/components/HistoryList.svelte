<script lang="ts">
	import { onMount } from 'svelte';
	import { history, type HistoryEntry } from '$lib/stores/history';
	import { config } from '$lib/stores/download';
	import { activeTab } from '$lib/stores/ui';

	onMount(async () => {
		await loadHistory();
	});

	async function loadHistory() {
		try {
			// @ts-ignore - Wails runtime
			const entries = await window.go.gui.App.GetHistory();
			history.set(entries || []);
		} catch (err) {
			console.error('Failed to load history:', err);
		}
	}

	async function clearHistory() {
		if (!confirm('Are you sure you want to clear all history?')) return;

		try {
			// @ts-ignore - Wails runtime
			await window.go.gui.App.ClearHistory();
			history.set([]);
		} catch (err) {
			console.error('Failed to clear history:', err);
		}
	}

	async function rerunDownload(entry: HistoryEntry) {
		config.set(entry.config);
		activeTab.set('download');
	}

	async function openOutputFolder(entry: HistoryEntry) {
		try {
			// @ts-ignore - Wails runtime
			await window.go.gui.App.OpenFolder(entry.output);
		} catch (err) {
			console.error('Failed to open folder:', err);
		}
	}

	function formatDate(dateStr: string): string {
		const date = new Date(dateStr);
		return date.toLocaleDateString('en-US', {
			month: 'short',
			day: 'numeric',
			year: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}

	function getStatusClass(status: HistoryEntry['status']): string {
		switch (status) {
			case 'completed':
				return 'status-success';
			case 'error':
				return 'status-error';
			case 'stopped':
				return 'status-warning';
			default:
				return '';
		}
	}
</script>

<div class="history-list">
	<div class="history-header">
		<h2 class="section-title">Download History</h2>
		{#if $history.length > 0}
			<button class="btn btn-danger" onclick={clearHistory}>
				Clear History
			</button>
		{/if}
	</div>

	{#if $history.length === 0}
		<div class="empty-state card">
			<p>No download history yet</p>
			<p class="hint">Completed downloads will appear here</p>
		</div>
	{:else}
		<div class="history-items">
			{#each $history as entry (entry.id)}
				<div class="history-item card">
					<div class="history-item-header">
						<span class="history-url" title={entry.url}>{entry.url}</span>
						<span class="history-status {getStatusClass(entry.status)}">
							{entry.status}
						</span>
					</div>

					<div class="history-item-details">
						<div class="detail">
							<span class="detail-label">Started:</span>
							<span>{formatDate(entry.startedAt)}</span>
						</div>
						{#if entry.completedAt}
							<div class="detail">
								<span class="detail-label">Completed:</span>
								<span>{formatDate(entry.completedAt)}</span>
							</div>
						{/if}
						<div class="detail">
							<span class="detail-label">Downloaded:</span>
							<span class="success">{entry.downloaded} files</span>
						</div>
						{#if entry.errors > 0}
							<div class="detail">
								<span class="detail-label">Errors:</span>
								<span class="error">{entry.errors}</span>
							</div>
						{/if}
					</div>

					<div class="history-item-output">
						<span class="detail-label">Output:</span>
						<span class="output-path">{entry.output}</span>
					</div>

					<div class="history-item-actions">
						<button class="btn" onclick={() => rerunDownload(entry)}>
							Re-run
						</button>
						<button class="btn" onclick={() => openOutputFolder(entry)}>
							Open Folder
						</button>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>

<style>
	.history-list {
		display: flex;
		flex-direction: column;
		height: 100%;
		overflow: hidden;
	}

	.history-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: var(--spacing-lg);
	}

	.section-title {
		font-size: 20px;
		font-weight: 600;
		margin: 0;
	}

	.empty-state {
		text-align: center;
		padding: var(--spacing-xl);
	}

	.empty-state p {
		margin: 0;
		color: var(--text-secondary);
	}

	.empty-state .hint {
		font-size: 14px;
		color: var(--text-muted);
		margin-top: var(--spacing-sm);
	}

	.history-items {
		display: flex;
		flex-direction: column;
		gap: var(--spacing-md);
		overflow-y: auto;
	}

	.history-item {
		display: flex;
		flex-direction: column;
		gap: var(--spacing-md);
	}

	.history-item-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		gap: var(--spacing-md);
	}

	.history-url {
		font-weight: 600;
		font-size: 14px;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.history-status {
		flex-shrink: 0;
		padding: var(--spacing-xs) var(--spacing-sm);
		border-radius: var(--radius-sm);
		font-size: 12px;
		font-weight: 500;
		text-transform: capitalize;
	}

	.status-success {
		background-color: rgba(40, 167, 69, 0.1);
		color: var(--success-color);
	}

	.status-error {
		background-color: rgba(220, 53, 69, 0.1);
		color: var(--error-color);
	}

	.status-warning {
		background-color: rgba(255, 193, 7, 0.1);
		color: var(--warning-color);
	}

	.history-item-details {
		display: flex;
		flex-wrap: wrap;
		gap: var(--spacing-md) var(--spacing-lg);
	}

	.detail {
		display: flex;
		gap: var(--spacing-xs);
		font-size: 13px;
	}

	.detail-label {
		color: var(--text-secondary);
	}

	.success {
		color: var(--success-color);
	}

	.error {
		color: var(--error-color);
	}

	.history-item-output {
		display: flex;
		flex-direction: column;
		gap: var(--spacing-xs);
		font-size: 13px;
	}

	.output-path {
		font-family: monospace;
		color: var(--text-muted);
		word-break: break-all;
	}

	.history-item-actions {
		display: flex;
		gap: var(--spacing-sm);
	}
</style>
