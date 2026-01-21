<script lang="ts">
	import { onMount } from 'svelte';
	import { presets, type Preset } from '$lib/stores/presets';
	import { config } from '$lib/stores/download';
	import { activeTab } from '$lib/stores/ui';

	let newPresetName = $state('');
	let showSaveDialog = $state(false);

	onMount(async () => {
		await loadPresets();
	});

	async function loadPresets() {
		try {
			// @ts-ignore - Wails runtime
			const items = await window.go.gui.App.GetPresets();
			presets.set(items || []);
		} catch (err) {
			console.error('Failed to load presets:', err);
		}
	}

	async function savePreset() {
		if (!newPresetName.trim()) return;

		try {
			// @ts-ignore - Wails runtime
			await window.go.gui.App.SavePreset(newPresetName, $config);
			await loadPresets();
			newPresetName = '';
			showSaveDialog = false;
		} catch (err) {
			console.error('Failed to save preset:', err);
		}
	}

	async function loadPreset(preset: Preset) {
		config.set(preset.config);
		activeTab.set('download');
	}

	async function deletePreset(id: string) {
		if (!confirm('Are you sure you want to delete this preset?')) return;

		try {
			// @ts-ignore - Wails runtime
			await window.go.gui.App.DeletePreset(id);
			await loadPresets();
		} catch (err) {
			console.error('Failed to delete preset:', err);
		}
	}

	function formatDate(dateStr: string): string {
		const date = new Date(dateStr);
		return date.toLocaleDateString('en-US', {
			month: 'short',
			day: 'numeric',
			year: 'numeric'
		});
	}
</script>

<div class="preset-manager">
	<div class="preset-header">
		<h2 class="section-title">Presets</h2>
		<button class="btn btn-primary" onclick={() => (showSaveDialog = true)}>
			Save Current Settings
		</button>
	</div>

	{#if showSaveDialog}
		<div class="save-dialog card">
			<h3>Save Preset</h3>
			<div class="save-form">
				<input
					type="text"
					class="input"
					placeholder="Preset name"
					bind:value={newPresetName}
					onkeydown={(e) => e.key === 'Enter' && savePreset()}
				/>
				<div class="dialog-actions">
					<button class="btn" onclick={() => (showSaveDialog = false)}>
						Cancel
					</button>
					<button
						class="btn btn-primary"
						onclick={savePreset}
						disabled={!newPresetName.trim()}
					>
						Save
					</button>
				</div>
			</div>
		</div>
	{/if}

	{#if $presets.length === 0}
		<div class="empty-state card">
			<p>No presets saved yet</p>
			<p class="hint">Save your current download settings as a preset for quick access later</p>
		</div>
	{:else}
		<div class="preset-list">
			{#each $presets as preset (preset.id)}
				<div class="preset-item card">
					<div class="preset-info">
						<h3 class="preset-name">{preset.name}</h3>
						<span class="preset-date">Created {formatDate(preset.createdAt)}</span>
					</div>

					<div class="preset-details">
						{#if preset.config.domain}
							<div class="detail">
								<span class="detail-label">Domain:</span>
								<span>{preset.config.domain}</span>
							</div>
						{/if}
						<div class="detail">
							<span class="detail-label">Concurrent:</span>
							<span>{preset.config.concurrent}</span>
						</div>
						<div class="detail">
							<span class="detail-label">Max Depth:</span>
							<span>{preset.config.maxDepth}</span>
						</div>
						<div class="detail">
							<span class="detail-label">Images:</span>
							<span>{preset.config.includeImages ? 'Yes' : 'No'}</span>
						</div>
					</div>

					<div class="preset-actions">
						<button class="btn btn-primary" onclick={() => loadPreset(preset)}>
							Load
						</button>
						<button class="btn btn-danger" onclick={() => deletePreset(preset.id)}>
							Delete
						</button>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>

<style>
	.preset-manager {
		display: flex;
		flex-direction: column;
		height: 100%;
		overflow: hidden;
	}

	.preset-header {
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

	.save-dialog {
		margin-bottom: var(--spacing-lg);
	}

	.save-dialog h3 {
		font-size: 16px;
		font-weight: 600;
		margin: 0 0 var(--spacing-md) 0;
	}

	.save-form {
		display: flex;
		flex-direction: column;
		gap: var(--spacing-md);
	}

	.dialog-actions {
		display: flex;
		justify-content: flex-end;
		gap: var(--spacing-sm);
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

	.preset-list {
		display: flex;
		flex-direction: column;
		gap: var(--spacing-md);
		overflow-y: auto;
	}

	.preset-item {
		display: flex;
		flex-direction: column;
		gap: var(--spacing-md);
	}

	.preset-info {
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	.preset-name {
		font-size: 16px;
		font-weight: 600;
		margin: 0;
	}

	.preset-date {
		font-size: 13px;
		color: var(--text-muted);
	}

	.preset-details {
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

	.preset-actions {
		display: flex;
		gap: var(--spacing-sm);
	}
</style>
