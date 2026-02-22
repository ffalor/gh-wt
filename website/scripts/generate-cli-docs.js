import { execSync } from 'child_process';
import { readdirSync, readFileSync, writeFileSync, rmSync } from 'fs';
import { join, dirname } from 'path';
import { fileURLToPath } from 'url';

const __dirname = dirname(fileURLToPath(import.meta.url));
const projectRoot = join(__dirname, '..', '..');
const cliDir = join(projectRoot, 'website/src/pages/docs/cli');

console.log('Project root:', projectRoot);
console.log('CLI dir:', cliDir);

// Change to project root
process.chdir(projectRoot);

// Run docgen
console.log('Generating CLI docs...');
execSync('go run ./internal/tools/docgen -out ./website/src/pages/docs/cli -format markdown -frontmatter', {
  stdio: 'inherit'
});

const files = readdirSync(cliDir);

// Remove .md files
files.filter(f => f.endsWith('.md')).forEach(f => {
  rmSync(join(cliDir, f));
  console.log(`Removed: ${f}`);
});

// Add layout to .mdx files and fix SEE ALSO links
files.filter(f => f.endsWith('.mdx')).forEach(f => {
  const filepath = join(cliDir, f);
  let content = readFileSync(filepath, 'utf8');
  
  if (!content.includes('layout:')) {
    content = content.replace(
      /^---\n/,
      '---\nlayout: "../../../layouts/CLILayout.astro"\n'
    );
    writeFileSync(filepath, content);
    console.log(`Updated: ${f}`);
  }
  
  // Fix SEE ALSO links - modify markdown links to include full path
  // [text](gh_wt.md) -> [text](/gh-wt/docs/cli/gh_wt)
  content = content.replace(/\]\(([^)]+)\.md\)/g, `](/gh-wt/docs/cli/$1)`);
  writeFileSync(filepath, content);
});

console.log('CLI docs generation complete!');
