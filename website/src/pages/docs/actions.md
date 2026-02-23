---
layout: ../../layouts/DocsLayout.astro
title: Actions
---

<article class="doc-page">
  <header class="doc-header">
    <h1><span class="neon-text">Actions</span></h1>
    <p class="doc-lead">Automate your workflow with post-creation actions.</p>
  </header>

  <section class="doc-section">
    <h2>What are Actions?</h2>
    <p>Actions are commands that run automatically after a worktree is created. They can be used to open tmux sessions, launch editors, or run any other command.</p>
  </section>

  <section class="doc-section">
    <h2>Configuration</h2>
    <p>Define actions in your configuration file:</p>

    <h3>Action Schema</h3>
    <p>Each action supports the following fields:</p>
<table class="config-table" is:raw>
  <thead>
    <tr>
      <th>Field</th>
      <th>Type</th>
      <th>Required</th>
      <th>Description</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>name</code></td>
      <td>string</td>
      <td>Yes</td>
      <td>Unique name for the action</td>
    </tr>
    <tr>
      <td><code>dir</code></td>
      <td>string</td>
      <td>No</td>
      <td>Working directory to run commands in. Supports template variables. Defaults to the worktree path.</td>
    </tr>
    <tr>
      <td><code>cmds</code></td>
      <td>[]string</td>
      <td>Yes</td>
      <td>List of commands to execute</td>
    </tr>
  </tbody>
</table></pre>
  </section>

  <section class="doc-section">
    <h2>Template Variables</h2>
    <p>Actions support the following template variables:</p>
<table class="config-table" is:raw>
  <thead>
    <tr>
      <th>Variable</th>
      <th>Description</th>
      <th>Example</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>{{.WorktreePath}}</code></td>
      <td>Path to the worktree</td>
      <td><code>~/github/worktree/pr_123</code></td>
    </tr>
    <tr>
      <td><code>{{.WorktreeName}}</code></td>
      <td>Name of the worktree directory</td>
      <td><code>pr_123</code></td>
    </tr>
    <tr>
      <td><code>{{.BranchName}}</code></td>
      <td>Name of the branch</td>
      <td><code>pr_123</code></td>
    </tr>
    <tr>
      <td><code>{{.Action}}</code></td>
      <td>Name of the action being run</td>
      <td><code>tmux</code></td>
    </tr>
    <tr>
      <td><code>{{.CLI_ARGS}}</code></td>
      <td>Arguments passed after <code>--</code></td>
      <td><code>"fix bug"</code></td>
    </tr>
    <tr>
      <td><code>{{.OS}}</code></td>
      <td>Operating system</td>
      <td><code>linux</code></td>
    </tr>
    <tr>
      <td><code>{{.ARCH}}</code></td>
      <td>System architecture</td>
      <td><code>amd64</code></td>
    </tr>
    <tr>
      <td><code>{{.ROOT_DIR}}</code></td>
      <td>Git root directory</td>
      <td><code>~/projects/my-repo</code></td>
    </tr>
  </tbody>
</table>
  </section>

  <section class="doc-section">
    <h2>Using Actions</h2>
    <h3>Actions on Worktree Creation</h3>
    <p>Pass the <code>-a</code> or <code>--action</code> flag when creating a worktree:</p>
    <pre is:raw><code>gh wt add https://github.com/owner/repo/pull/123 -a tmux</code></pre>
    <p>You can also pass arguments to the action:</p>
    <pre is:raw><code>gh wt add my-branch -a editor -- --debug</code></pre>
    <p>The arguments after <code>--</code> will be available in <code>{{.CLI_ARGS}}</code>.</p>
    <h3>Using Actions After Creation</h3>
    <p>Run actions on existing worktrees using the <code>run</code> command:</p>
    <h4>Run a Named Action</h4>
    <pre is:raw><code>gh wt run pr_123 tmux</code></pre>
    <h4>Run with Arguments</h4>
    <p>Pass arguments to the action using <code>--</code>:</p>
    <pre is:raw><code>gh wt run pr_123 editor -- --debug</code></pre>
    <h4>Run a Command</h4>
    <p>Run any command directly in a worktree:</p>
    <pre is:raw><code>gh wt run pr_123 -- ls -la</code></pre>
  </section>

  <section class="doc-section">
    <h2>Example Actions</h2>
    <h3>VS Code</h3>
    <p>Opens the worktree directory in Visual Studio Code.</p>
<pre is:raw><code>actions:
  - name: code
    dir: {{.WorktreePath}}
    cmds:
      - code .</code></pre>
    <h3>tmux with Vim</h3>
    <p>Creates a new tmux session named after the branch and opens Vim in the worktree directory.</p>
<pre is:raw><code>actions:
  - name: tmux-vim
    dir: {{.WorktreePath}}
    cmds:
      - tmux new-session -d -s {{.BranchName}}
      - tmux send-keys -t {{.BranchName}} "vim ." C-m</code></pre>
    <h3>Shell in Worktree</h3>
    <p>Opens a new shell (using your default shell) in the worktree directory.</p>
<pre is:raw><code>actions:
  - name: shell
    dir: {{.WorktreePath}}
    cmds:
      - $SHELL</code></pre>
  </section>

  <nav class="doc-nav">
    <a href="/gh-wt/docs/cli/gh_wt" class="btn">
      Next: CLI Reference
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

  .doc-section h3 {
    font-family: var(--font-display);
    font-size: 1.1rem;
    margin: 1.5rem 0 1rem;
    color: var(--color-cyan);
  }

  .doc-section p {
    color: var(--text-secondary);
    margin-bottom: 1rem;
  }

  .doc-section pre {
    margin: 1rem 0 1.5rem;
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
