<script lang="ts">
	import { status, downloadedCount, errorCount, currentUrl, isDownloading } from '$lib/stores/download';

	function getStatusText(s: typeof $status) {
		switch (s) {
			case 'idle':
				return 'Ready';
			case 'downloading':
				return 'Downloading...';
			case 'completed':
				return 'Completed';
			case 'error':
				return 'Error';
			case 'stopped':
				return 'Stopped';
			default:
				return 'Unknown';
		}
	}

	function getStatusClass(s: typeof $status) {
		switch (s) {
			case 'downloading':
				return 'status-active';
			case 'completed':
				return 'status-success';
			case 'error':
				return 'status-error';
			case 'stopped':
				return 'status-warning';
			default:
				return 'status-idle';
		}
	}
</script>

<div class="progress-panel card">
	<h2 class="section-title">Progress</h2>

	<div class="status-row">
		<span class="status-label">Status:</span>
		<span class="status-value {getStatusClass($status)}">
			{getStatusText($status)}
			{#if $isDownloading}
				<span class="spinner"></span>
			{/if}
		</span>
	</div>

	<div class="stats">
		<div class="stat">
			<span class="stat-value success">{$downloadedCount}</span>
			<span class="stat-label">Downloaded</span>
		</div>
		<div class="stat">
			<span class="stat-value error">{$errorCount}</span>
			<span class="stat-label">Errors</span>
		</div>
	</div>

	{#if $currentUrl}
		<div class="current-url">
			<span class="current-label">Current:</span>
			<span class="current-value" title={$currentUrl}>
				{$currentUrl}
			</span>
		</div>
	{/if}
</div>

<style>
	.progress-panel {
		display: flex;
		flex-direction: column;
		gap: var(--spacing-md);
	}

	.section-title {
		font-size: 18px;
		font-weight: 600;
		margin: 0;
	}

	.status-row {
		display: flex;
		align-items: center;
		gap: var(--spacing-sm);
	}

	.status-label {
		font-size: 14px;
		color: var(--text-secondary);
	}

	.status-value {
		display: flex;
		align-items: center;
		gap: var(--spacing-sm);
		font-size: 14px;
		font-weight: 600;
		padding: var(--spacing-xs) var(--spacing-sm);
		border-radius: var(--radius-sm);
	}

	.status-idle {
		background-color: var(--bg-tertiary);
		color: var(--text-secondary);
	}

	.status-active {
		background-color: rgba(0, 102, 204, 0.1);
		color: var(--accent-color);
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

	.spinner {
		width: 12px;
		height: 12px;
		border: 2px solid currentColor;
		border-right-color: transparent;
		border-radius: 50%;
		animation: spin 0.75s linear infinite;
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}

	.stats {
		display: flex;
		gap: var(--spacing-lg);
	}

	.stat {
		display: flex;
		flex-direction: column;
		align-items: center;
		padding: var(--spacing-md);
		background-color: var(--bg-tertiary);
		border-radius: var(--radius-md);
		min-width: 100px;
	}

	.stat-value {
		font-size: 28px;
		font-weight: 700;
	}

	.stat-value.success {
		color: var(--success-color);
	}

	.stat-value.error {
		color: var(--error-color);
	}

	.stat-label {
		font-size: 12px;
		color: var(--text-secondary);
		text-transform: uppercase;
		letter-spacing: 0.5px;
	}

	.current-url {
		display: flex;
		flex-direction: column;
		gap: var(--spacing-xs);
		padding: var(--spacing-sm);
		background-color: var(--bg-tertiary);
		border-radius: var(--radius-md);
		overflow: hidden;
	}

	.current-label {
		font-size: 12px;
		color: var(--text-secondary);
		text-transform: uppercase;
	}

	.current-value {
		font-size: 13px;
		font-family: monospace;
		color: var(--text-primary);
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}
</style>
