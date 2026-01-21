<script lang="ts">
	import { logs, clearLogs, type LogEntry } from '$lib/stores/download';
	import { onMount } from 'svelte';

	let logContainer: HTMLDivElement;
	let autoScroll = $state(true);

	function formatTime(date: Date): string {
		return date.toLocaleTimeString('en-US', {
			hour12: false,
			hour: '2-digit',
			minute: '2-digit',
			second: '2-digit'
		});
	}

	function getLogClass(type: LogEntry['type']): string {
		switch (type) {
			case 'success':
				return 'log-success';
			case 'error':
				return 'log-error';
			case 'warning':
				return 'log-warning';
			default:
				return 'log-info';
		}
	}

	function handleScroll() {
		if (!logContainer) return;
		const { scrollTop, scrollHeight, clientHeight } = logContainer;
		autoScroll = scrollHeight - scrollTop - clientHeight < 50;
	}

	$effect(() => {
		if ($logs.length && autoScroll && logContainer) {
			logContainer.scrollTop = logContainer.scrollHeight;
		}
	});
</script>

<div class="log-output card">
	<div class="log-header">
		<h2 class="section-title">Log</h2>
		<div class="log-actions">
			<label class="checkbox-label">
				<input type="checkbox" bind:checked={autoScroll} />
				Auto-scroll
			</label>
			<button class="btn btn-small" onclick={clearLogs}>
				Clear
			</button>
		</div>
	</div>

	<div class="log-container" bind:this={logContainer} onscroll={handleScroll}>
		{#if $logs.length === 0}
			<div class="log-empty">No log entries yet</div>
		{:else}
			{#each $logs as entry (entry.id)}
				<div class="log-entry {getLogClass(entry.type)}">
					<span class="log-time">{formatTime(entry.timestamp)}</span>
					<span class="log-message">{entry.message}</span>
				</div>
			{/each}
		{/if}
	</div>
</div>

<style>
	.log-output {
		display: flex;
		flex-direction: column;
		height: 100%;
		overflow: hidden;
	}

	.log-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: var(--spacing-md);
	}

	.section-title {
		font-size: 18px;
		font-weight: 600;
		margin: 0;
	}

	.log-actions {
		display: flex;
		align-items: center;
		gap: var(--spacing-md);
	}

	.checkbox-label {
		display: flex;
		align-items: center;
		gap: var(--spacing-xs);
		font-size: 13px;
		color: var(--text-secondary);
		cursor: pointer;
	}

	.checkbox-label input {
		width: 14px;
		height: 14px;
	}

	.btn-small {
		padding: var(--spacing-xs) var(--spacing-sm);
		font-size: 12px;
	}

	.log-container {
		flex: 1;
		overflow-y: auto;
		background-color: var(--bg-primary);
		border: 1px solid var(--border-color);
		border-radius: var(--radius-md);
		font-family: monospace;
		font-size: 12px;
		line-height: 1.5;
	}

	.log-empty {
		padding: var(--spacing-lg);
		text-align: center;
		color: var(--text-muted);
	}

	.log-entry {
		display: flex;
		gap: var(--spacing-sm);
		padding: var(--spacing-xs) var(--spacing-sm);
		border-bottom: 1px solid var(--border-color);
	}

	.log-entry:last-child {
		border-bottom: none;
	}

	.log-time {
		color: var(--text-muted);
		flex-shrink: 0;
	}

	.log-message {
		word-break: break-all;
	}

	.log-info {
		color: var(--text-primary);
	}

	.log-success .log-message {
		color: var(--success-color);
	}

	.log-error .log-message {
		color: var(--error-color);
	}

	.log-warning .log-message {
		color: var(--warning-color);
	}
</style>
