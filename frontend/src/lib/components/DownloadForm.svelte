<script lang="ts">
	import { onMount } from 'svelte';
	import {
		config,
		status,
		isDownloading,
		canStart,
		extractDomain,
		resetProgress,
		addLog
	} from '$lib/stores/download';
	import { settings } from '$lib/stores/settings';

	let showAdvanced = $state(false);
	let blacklistInput = $state('');
	let extensionsInput = $state('');
	let imageExtsInput = $state('');

	// Update domain when URL changes
	$effect(() => {
		if ($config.url) {
			const domain = extractDomain($config.url);
			if (domain && !$config.domain) {
				config.update((c) => ({ ...c, domain }));
			}
		}
	});

	// Initialize from settings
	onMount(() => {
		const unsubscribe = settings.subscribe((s) => {
			config.update((c) => ({
				...c,
				output: c.output || s.defaultOutput,
				concurrent: s.defaultConcurrent,
				maxDepth: s.defaultMaxDepth,
				retries: s.defaultRetries,
				includeImages: s.includeImages
			}));
		});
		return unsubscribe;
	});

	// Sync array inputs
	$effect(() => {
		blacklistInput = $config.blacklist.join(', ');
	});

	$effect(() => {
		extensionsInput = $config.extensions.join(', ');
	});

	$effect(() => {
		imageExtsInput = $config.imageExtensions.join(', ');
	});

	function updateBlacklist() {
		const items = blacklistInput
			.split(',')
			.map((s) => s.trim())
			.filter((s) => s);
		config.update((c) => ({ ...c, blacklist: items }));
	}

	function updateExtensions() {
		const items = extensionsInput
			.split(',')
			.map((s) => s.trim())
			.filter((s) => s);
		config.update((c) => ({ ...c, extensions: items }));
	}

	function updateImageExts() {
		const items = imageExtsInput
			.split(',')
			.map((s) => s.trim())
			.filter((s) => s);
		config.update((c) => ({ ...c, imageExtensions: items }));
	}

	async function selectOutputDirectory() {
		try {
			// @ts-ignore - Wails runtime
			const dir = await window.go.gui.App.SelectOutputDirectory();
			if (dir) {
				config.update((c) => ({ ...c, output: dir }));
			}
		} catch (err) {
			console.error('Failed to select directory:', err);
		}
	}

	async function startDownload() {
		if (!$canStart) return;

		resetProgress();
		status.set('downloading');
		addLog('info', `Starting download from ${$config.url}`);

		try {
			// @ts-ignore - Wails runtime
			await window.go.gui.App.StartDownload($config);
		} catch (err) {
			status.set('error');
			addLog('error', `Download failed: ${err}`);
		}
	}

	async function stopDownload() {
		try {
			// @ts-ignore - Wails runtime
			await window.go.gui.App.StopDownload();
			status.set('stopped');
			addLog('warning', 'Download stopped by user');
		} catch (err) {
			console.error('Failed to stop download:', err);
		}
	}
</script>

<div class="download-form card">
	<h2 class="section-title">Download Website</h2>

	<div class="form-group">
		<label for="url">URL</label>
		<input
			type="url"
			id="url"
			class="input"
			placeholder="https://example.com/docs"
			bind:value={$config.url}
			disabled={$isDownloading}
		/>
	</div>

	<div class="form-group">
		<label for="domain">Domain Filter</label>
		<input
			type="text"
			id="domain"
			class="input"
			placeholder="https://example.com"
			bind:value={$config.domain}
			disabled={$isDownloading}
		/>
		<span class="hint">Only download URLs under this domain</span>
	</div>

	<div class="form-group">
		<label for="output">Output Directory</label>
		<div class="input-with-button">
			<input
				type="text"
				id="output"
				class="input"
				placeholder="/path/to/output"
				bind:value={$config.output}
				disabled={$isDownloading}
			/>
			<button class="btn" onclick={selectOutputDirectory} disabled={$isDownloading}>
				Browse
			</button>
		</div>
	</div>

	<div class="form-row">
		<label class="checkbox-label">
			<input
				type="checkbox"
				bind:checked={$config.includeImages}
				disabled={$isDownloading}
			/>
			Include Images
		</label>
	</div>

	<button
		class="btn toggle-advanced"
		onclick={() => (showAdvanced = !showAdvanced)}
	>
		{showAdvanced ? 'Hide' : 'Show'} Advanced Options
	</button>

	{#if showAdvanced}
		<div class="advanced-options">
			<div class="form-row">
				<div class="form-group half">
					<label for="concurrent">Concurrent Downloads</label>
					<input
						type="number"
						id="concurrent"
						class="input"
						min="1"
						max="20"
						bind:value={$config.concurrent}
						disabled={$isDownloading}
					/>
				</div>
				<div class="form-group half">
					<label for="maxDepth">Max Depth</label>
					<input
						type="number"
						id="maxDepth"
						class="input"
						min="1"
						max="100"
						bind:value={$config.maxDepth}
						disabled={$isDownloading}
					/>
				</div>
			</div>

			<div class="form-group">
				<label for="retries">Retries</label>
				<input
					type="number"
					id="retries"
					class="input"
					min="0"
					max="10"
					bind:value={$config.retries}
					disabled={$isDownloading}
				/>
			</div>

			<div class="form-group">
				<label for="blacklist">Blacklist (comma-separated)</label>
				<input
					type="text"
					id="blacklist"
					class="input"
					placeholder="https://github, /admin"
					bind:value={blacklistInput}
					onblur={updateBlacklist}
					disabled={$isDownloading}
				/>
			</div>

			<div class="form-group">
				<label for="extensions">File Extensions (comma-separated)</label>
				<input
					type="text"
					id="extensions"
					class="input"
					placeholder=".js, .css, .woff2"
					bind:value={extensionsInput}
					onblur={updateExtensions}
					disabled={$isDownloading}
				/>
			</div>

			<div class="form-group">
				<label for="imageExts">Image Extensions (comma-separated)</label>
				<input
					type="text"
					id="imageExts"
					class="input"
					placeholder=".png, .jpg, .svg"
					bind:value={imageExtsInput}
					onblur={updateImageExts}
					disabled={$isDownloading}
				/>
			</div>
		</div>
	{/if}

	<div class="actions">
		{#if $isDownloading}
			<button class="btn btn-danger" onclick={stopDownload}>
				Stop Download
			</button>
		{:else}
			<button class="btn btn-primary" onclick={startDownload} disabled={!$canStart}>
				Start Download
			</button>
		{/if}
	</div>
</div>

<style>
	.download-form {
		display: flex;
		flex-direction: column;
		gap: var(--spacing-md);
	}

	.section-title {
		font-size: 18px;
		font-weight: 600;
		margin: 0 0 var(--spacing-sm) 0;
	}

	.form-group {
		display: flex;
		flex-direction: column;
		gap: var(--spacing-xs);
	}

	.form-group.half {
		flex: 1;
	}

	.form-row {
		display: flex;
		gap: var(--spacing-md);
	}

	label {
		font-size: 14px;
		font-weight: 500;
		color: var(--text-secondary);
	}

	.hint {
		font-size: 12px;
		color: var(--text-muted);
	}

	.input-with-button {
		display: flex;
		gap: var(--spacing-sm);
	}

	.input-with-button .input {
		flex: 1;
	}

	.checkbox-label {
		display: flex;
		align-items: center;
		gap: var(--spacing-sm);
		cursor: pointer;
	}

	.checkbox-label input {
		width: 16px;
		height: 16px;
	}

	.toggle-advanced {
		align-self: flex-start;
		background: transparent;
		border: none;
		color: var(--accent-color);
		padding: 0;
		font-size: 14px;
	}

	.toggle-advanced:hover {
		text-decoration: underline;
		background: transparent;
	}

	.advanced-options {
		display: flex;
		flex-direction: column;
		gap: var(--spacing-md);
		padding: var(--spacing-md);
		background-color: var(--bg-tertiary);
		border-radius: var(--radius-md);
	}

	.actions {
		display: flex;
		justify-content: flex-end;
		margin-top: var(--spacing-sm);
	}

	.actions .btn {
		min-width: 140px;
	}
</style>
