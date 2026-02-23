---
layout: ../../layouts/DocsLayout.astro
title: Configuration
---

<article class="doc-page">
  <header class="doc-header">
    <h1><span class="neon-text">Configuration</span></h1>
    <p class="doc-lead">Customize gh-wt to fit your workflow.</p>
  </header>

  <section class="doc-section">
    <h2>Configuration File</h2>
    <p>gh-wt reads configuration from <code>~/.config/gh-wt/config.yaml</code>.</p>

<pre is:raw><code>worktree_dir: "~/github/worktree"

actions:
  - name: tmux
    cmds:
      - tmux new-session -d -s {{.BranchName}}
      - tmux send-keys -t {{.BranchName}} "cd {{.WorktreePath}}" C-m</code></pre>
  </section>

  <section class="doc-section">
    <h2>Configuration Options</h2>
<table class="config-table" is:raw>
  <thead>
    <tr>
      <th>Option</th>
      <th>Type</th>
      <th>Description</th>
      <th>Default</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>worktree_dir</code></td>
      <td>string</td>
      <td>Directory where worktrees are created</td>
      <td><code>~/github/worktree</code></td>
    </tr>
    <tr>
      <td><code>actions</code></td>
      <td>array</td>
      <td>List of post-creation actions</td>
      <td><code>[]</code></td>
    </tr>
  </tbody>
</table>
  </section>

  <section class="doc-section">
    <h2>Environment Variables</h2>
    <p>Configuration can also be set via environment variables:</p>
<table class="config-table" is:raw>
  <thead>
    <tr>
      <th>Variable</th>
      <th>Description</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>GH_WT_WORKTREE_DIR</code></td>
      <td>Override worktree directory</td>
    </tr>
    <tr>
      <td><code>GH_WT_VERBOSE</code></td>
      <td>Enable verbose output</td>
    </tr>
    <tr>
      <td><code>GH_WT_NO_COLOR</code></td>
      <td>Disable color output</td>
    </tr>
  </tbody>
</table>
  </section>

  <section class="doc-section">
    <h2>Priority Order</h2>
    <p>Configuration priority (highest to lowest):</p>
    <ol>
      <li>Command-line flags</li>
      <li>Environment variables (prefix: <code>GH_WT_</code>)</li>
      <li>Configuration file</li>
      <li>Default values</li>
    </ol>
  </section>

  <nav class="doc-nav">
    <a href="/gh-wt/docs/actions" class="btn">
      Next: Actions
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="M5 12h14M12 5l7 7-7 7"/>
      </svg>
    </a>
  </nav>
</article>

<style>
  .doc-page {
    max-width: 800px;
    margin: 0 auto;
  }

  .doc-header {
    margin-bottom: 3rem;
    padding-bottom: 2rem;
    border-bottom: 1px solid var(--border-color);
  }

  .doc-header h1 {
    margin-bottom: 0.5rem;
  }

  .doc-lead {
    font-size: 1.2rem;
    color: var(--text-secondary);
  }

  .doc-section {
    margin-bottom: 3rem;
  }

  .doc-section h2 {
    font-family: var(--font-display);
    font-size: 1.5rem;
    margin-bottom: 1rem;
    color: var(--color-amber);
  }

  .doc-section p {
    color: var(--text-secondary);
    margin-bottom: 1rem;
  }

  .doc-section ol {
    list-style: decimal;
    padding-left: 1.5rem;
    color: var(--text-secondary);
  }

  .doc-section ol li {
    margin-bottom: 0.5rem;
  }

  .doc-section pre {
    margin: 1.5rem 0;
  }

  .config-table {
    width: 100%;
    border-collapse: collapse;
    margin: 1.5rem 0;
    font-size: 0.9rem;
  }

  .config-table th,
  .config-table td {
    padding: 0.75rem 1rem;
    text-align: left;
    border: 1px solid var(--border-color);
  }

  .config-table th {
    background: var(--bg-tertiary);
    color: var(--color-cyan);
    font-family: var(--font-display);
    font-weight: 600;
  }

  .config-table td {
    background: var(--bg-secondary);
    color: var(--text-secondary);
  }

  .config-table code {
    color: var(--color-amber);
  }

  .doc-nav {
    margin-top: 3rem;
    padding-top: 2rem;
    border-top: 1px solid var(--border-color);
    display: flex;
    justify-content: flex-end;
  }
</style>
