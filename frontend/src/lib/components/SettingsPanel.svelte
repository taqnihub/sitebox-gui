<script lang="ts">
	import { onMount } from 'svelte';
	import { settings, defaultSettings, type Theme } from '$lib/stores/settings';

	onMount(async () => {
		await loadSettings();
	});

	async function loadSettings() {
		try {
			// @ts-ignore - Wails runtime
			const saved = await window.go.gui.App.GetSettings();
			if (saved) {
				settings.set(saved);
			}
		} catch (err) {
			console.error('Failed to load settings:', err);
		}
	}

	async function saveSettings() {
		try {
			// @ts-ignore - Wails runtime
			await window.go.gui.App.SaveSettings($settings);
		} catch (err) {
			console.error('Failed to save settings:', err);
		}
	}

	async function selectDefaultOutput() {
		try {
			// @ts-ignore - Wails runtime
			const dir = await window.go.gui.App.SelectOutputDirectory();
			if (dir) {
				settings.update((s) => ({ ...s, defaultOutput: dir }));
				await saveSettings();
			}
		} catch (err) {
			console.error('Failed to select directory:', err);
		}
	}

	function handleThemeChange(theme: Theme) {
		settings.update((s) => ({ ...s, theme }));
		saveSettings();
	}

	function handleSettingChange() {
		saveSettings();
	}

	async function resetToDefaults() {
		if (!confirm('Are you sure you want to reset all settings to defaults?')) return;

		settings.set({ ...defaultSettings });
		await saveSettings();
	}
</script>

<div class="settings-panel">
	<h2 class="section-title">Settings</h2>

	<div class="settings-section card">
		<h3 class="section-subtitle">Appearance</h3>

		<div class="setting-item">
			<span class="setting-label">Theme</span>
			<div class="theme-options">
				<button
					class="theme-btn"
					class:active={$settings.theme === 'light'}
					onclick={() => handleThemeChange('light')}
				>
					Light
				</button>
				<button
					class="theme-btn"
					class:active={$settings.theme === 'dark'}
					onclick={() => handleThemeChange('dark')}
				>
					Dark
				</button>
				<button
					class="theme-btn"
					class:active={$settings.theme === 'system'}
					onclick={() => handleThemeChange('system')}
				>
					System
				</button>
			</div>
		</div>
	</div>

	<div class="settings-section card">
		<h3 class="section-subtitle">Default Download Options</h3>

		<div class="setting-item">
			<label for="defaultOutput" class="setting-label">Default Output Directory</label>
			<div class="input-with-button">
				<input
					type="text"
					id="defaultOutput"
					class="input"
					placeholder="Select default output directory"
					bind:value={$settings.defaultOutput}
					onblur={handleSettingChange}
				/>
				<button class="btn" onclick={selectDefaultOutput}>
					Browse
				</button>
			</div>
		</div>

		<div class="setting-row">
			<div class="setting-item">
				<label for="defaultConcurrent" class="setting-label">Concurrent Downloads</label>
				<input
					type="number"
					id="defaultConcurrent"
					class="input"
					min="1"
					max="20"
					bind:value={$settings.defaultConcurrent}
					onchange={handleSettingChange}
				/>
			</div>

			<div class="setting-item">
				<label for="defaultMaxDepth" class="setting-label">Max Depth</label>
				<input
					type="number"
					id="defaultMaxDepth"
					class="input"
					min="1"
					max="100"
					bind:value={$settings.defaultMaxDepth}
					onchange={handleSettingChange}
				/>
			</div>

			<div class="setting-item">
				<label for="defaultRetries" class="setting-label">Retries</label>
				<input
					type="number"
					id="defaultRetries"
					class="input"
					min="0"
					max="10"
					bind:value={$settings.defaultRetries}
					onchange={handleSettingChange}
				/>
			</div>
		</div>

		<div class="setting-item">
			<label class="checkbox-label">
				<input
					type="checkbox"
					bind:checked={$settings.includeImages}
					onchange={handleSettingChange}
				/>
				Include Images by Default
			</label>
		</div>
	</div>

	<div class="settings-section card">
		<h3 class="section-subtitle">Notifications</h3>

		<div class="setting-item">
			<label class="checkbox-label">
				<input
					type="checkbox"
					bind:checked={$settings.showNotifications}
					onchange={handleSettingChange}
				/>
				Show Desktop Notifications
			</label>
			<span class="setting-hint">Get notified when downloads complete or fail</span>
		</div>
	</div>

	<div class="settings-actions">
		<button class="btn btn-danger" onclick={resetToDefaults}>
			Reset to Defaults
		</button>
	</div>
</div>

<style>
	.settings-panel {
		display: flex;
		flex-direction: column;
		gap: var(--spacing-lg);
		max-width: 600px;
	}

	.section-title {
		font-size: 20px;
		font-weight: 600;
		margin: 0;
	}

	.settings-section {
		display: flex;
		flex-direction: column;
		gap: var(--spacing-md);
	}

	.section-subtitle {
		font-size: 16px;
		font-weight: 600;
		margin: 0;
		padding-bottom: var(--spacing-sm);
		border-bottom: 1px solid var(--border-color);
	}

	.setting-item {
		display: flex;
		flex-direction: column;
		gap: var(--spacing-xs);
	}

	.setting-row {
		display: flex;
		gap: var(--spacing-md);
	}

	.setting-row .setting-item {
		flex: 1;
	}

	.setting-label {
		font-size: 14px;
		font-weight: 500;
		color: var(--text-secondary);
	}

	.setting-hint {
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

	.theme-options {
		display: flex;
		gap: var(--spacing-xs);
	}

	.theme-btn {
		flex: 1;
		padding: var(--spacing-sm) var(--spacing-md);
		border: 1px solid var(--border-color);
		border-radius: var(--radius-md);
		background: var(--bg-primary);
		color: var(--text-secondary);
		font-size: 14px;
		cursor: pointer;
		transition: all var(--transition-fast);
	}

	.theme-btn:hover {
		background: var(--bg-tertiary);
		color: var(--text-primary);
	}

	.theme-btn.active {
		background: var(--accent-color);
		border-color: var(--accent-color);
		color: white;
	}

	.checkbox-label {
		display: flex;
		align-items: center;
		gap: var(--spacing-sm);
		font-size: 14px;
		cursor: pointer;
	}

	.checkbox-label input {
		width: 16px;
		height: 16px;
	}

	.settings-actions {
		display: flex;
		justify-content: flex-start;
	}
</style>
