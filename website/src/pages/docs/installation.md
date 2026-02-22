---
layout: ../../layouts/DocsLayout.astro
title: Installation
---

<article class="doc-page">
  <header class="doc-header">
    <h1><span class="neon-text">Installation</span></h1>
    <p class="doc-lead">Get gh-wt up and running in seconds.</p>
  </header>

  <section class="doc-section">
    <h2>Prerequisites</h2>
    <ul>
      <li><a href="https://cli.github.com">GitHub CLI</a> (<code>gh</code>) installed and authenticated</li>
      <li><a href="https://git-scm.com">Git</a> installed</li>
    </ul>
  </section>

  <section class="doc-section">
    <h2>Install</h2>
    <p>Install gh-wt directly from the GitHub repository:</p>
<pre is:raw><code>gh extension install ffalor/gh-wt</code></pre>
  </section>

  <section class="doc-section">
    <h2>Verify Installation</h2>
    <p>Confirm gh-wt is installed correctly:</p>
<pre is:raw><code>gh wt --help</code></pre>
    <p>You should see the help output with all available commands.</p>
  </section>

  <section class="doc-section">
    <h2>Upgrading</h2>
    <p>To upgrade to the latest version:</p>
<pre is:raw><code>gh extension upgrade gh-wt</code></pre>
  </section>

  <section class="doc-section">
    <h2>Uninstalling</h2>
    <p>To remove gh-wt:</p>
<pre is:raw><code>gh extension remove wt</code></pre>
  </section>

  <section class="doc-section">
    <h2>Install from Source</h2>
    <p>For development or manual installation, clone the repository and use the provided Taskfile:</p>
<pre is:raw><code>git clone https://github.com/ffalor/gh-wt.git
cd gh-wt
task install</code></pre>
    <p>This project uses Taskfile. Learn more at <a href="https://taskfile.dev">taskfile.dev</a>.</p>
<table class="config-table" is:raw>
  <thead>
    <tr>
      <th>Task</th>
      <th>Description</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>task install</code></td>
      <td>Build and install the extension</td>
    </tr>
    <tr>
      <td><code>task build</code></td>
      <td>Build the binary only</td>
    </tr>
    <tr>
      <td><code>task remove</code></td>
      <td>Remove the installed extension</td>
    </tr>
    <tr>
      <td><code>task clean</code></td>
      <td>Clean built files</td>
    </tr>
  </tbody>
</table>
  </section>

  <nav class="doc-nav">
    <a href="/gh-wt/docs/configuration" class="btn">
      Next: Configuration
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

  .doc-section ul {
    list-style: none;
    padding-left: 0;
  }

  .doc-section li {
    position: relative;
    padding-left: 1.5rem;
    margin-bottom: 0.75rem;
    color: var(--text-secondary);
  }

  .doc-section li::before {
    content: '›';
    position: absolute;
    left: 0;
    color: var(--color-cyan);
    font-weight: bold;
  }

  .doc-section pre {
    margin: 1.5rem 0;
  }

  .note {
    background: var(--bg-tertiary);
    border-left: 3px solid var(--color-amber);
    padding: 1rem;
    border-radius: 0 8px 8px 0;
    margin: 1.5rem 0;
    font-size: 0.9rem;
    color: var(--text-secondary);
  }

  .note strong {
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
