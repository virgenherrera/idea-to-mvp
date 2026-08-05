'use strict';

const fs = require('node:fs');
const path = require('node:path');

// ---------------------------------------------------------------------------
// Constants
// ---------------------------------------------------------------------------

const SKIP_DIRS = new Set([
  'node_modules', '.git', 'dist', 'build', 'coverage',
  '.next', '.nuxt', '.output', '.svelte-kit', '.angular',
  '.turbo', '.cache', '.parcel-cache', '__pycache__',
  'target', 'vendor',
]);

const FRAMEWORK_MAP = {
  '@angular/core': 'Angular',
  'react': 'React',
  'react-dom': 'React',
  'vue': 'Vue',
  'svelte': 'Svelte',
  '@sveltejs/kit': 'SvelteKit',
  'next': 'Next.js',
  'nuxt': 'Nuxt',
  'remix': 'Remix',
  '@remix-run/node': 'Remix',
  '@remix-run/react': 'Remix',
  'astro': 'Astro',
  '@nestjs/core': 'NestJS',
  'express': 'Express',
  'fastify': 'Fastify',
  'koa': 'Koa',
  'hono': 'Hono',
};

const TEST_RUNNER_MAP = {
  'jest': 'jest',
  'vitest': 'vitest',
  'mocha': 'mocha',
  'playwright': 'playwright',
  '@playwright/test': 'playwright',
  'cypress': 'cypress',
  'ava': 'ava',
};

const LOCKFILE_MAP = {
  'pnpm-lock.yaml': 'pnpm',
  'yarn.lock': 'yarn',
  'package-lock.json': 'npm',
  'bun.lockb': 'bun',
  'bun.lock': 'bun',
};

const KEY_FILES = [
  'README.md', 'CONTRIBUTING.md', 'AGENTS.md', 'CLAUDE.md',
  'CURSOR.md', 'LICENSE', 'CHANGELOG.md',
];

const ECHO_STEPS = [
  { step: 'Setup', scripts: ['install', 'prepare', 'postinstall'] },
  { step: 'Build', scripts: ['build'] },
  { step: 'Static', scripts: ['lint', 'eslint', 'lint:fix', 'test:static'] },
  { step: 'Dynamic', scripts: ['test', 'test:unit'] },
  { step: 'E2E', scripts: ['test:e2e', 'e2e'] },
];

const TEST_CONFIG_PATTERNS = [
  'jest.config.js', 'jest.config.ts', 'jest.config.mjs', 'jest.config.cjs',
  'vitest.config.js', 'vitest.config.ts', 'vitest.config.mjs',
  'playwright.config.js', 'playwright.config.ts',
  'cypress.config.js', 'cypress.config.ts',
  '.mocharc.yml', '.mocharc.yaml', '.mocharc.json', '.mocharc.js',
];

const EXTENSION_CATEGORIES = {
  '.ts': 'TypeScript', '.tsx': 'TypeScript (JSX)', '.js': 'JavaScript',
  '.jsx': 'JavaScript (JSX)', '.vue': 'Vue', '.svelte': 'Svelte',
  '.go': 'Go', '.rs': 'Rust', '.py': 'Python', '.rb': 'Ruby',
  '.java': 'Java', '.kt': 'Kotlin', '.swift': 'Swift',
  '.html': 'HTML', '.css': 'CSS', '.scss': 'SCSS', '.less': 'Less',
  '.json': 'JSON', '.yaml': 'YAML', '.yml': 'YAML',
  '.md': 'Markdown', '.sh': 'Shell', '.sql': 'SQL',
};

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

/** Safely read and parse JSON from a file path. Returns null on failure. */
function readJson(filePath) {
  try {
    return JSON.parse(fs.readFileSync(filePath, 'utf-8'));
  } catch {
    return null;
  }
}

/** Check if a path exists (file or directory). */
function exists(p) {
  try {
    fs.accessSync(p);
    return true;
  } catch {
    return false;
  }
}

/** Check if path is a directory. */
function isDir(p) {
  try {
    return fs.statSync(p).isDirectory();
  } catch {
    return false;
  }
}

/** List directory entries safely. */
function readDir(p) {
  try {
    return fs.readdirSync(p, { withFileTypes: true });
  } catch {
    return [];
  }
}

/**
 * Walk a directory tree up to a max depth, skipping SKIP_DIRS.
 * Calls visitor(filePath, dirent) for every file.
 */
function walk(dir, visitor, maxDepth = Infinity, _depth = 0) {
  if (_depth > maxDepth) return;
  for (const entry of readDir(dir)) {
    const full = path.join(dir, entry.name);
    if (entry.isDirectory()) {
      if (!SKIP_DIRS.has(entry.name)) {
        walk(full, visitor, maxDepth, _depth + 1);
      }
    } else if (entry.isFile() || entry.isSymbolicLink()) {
      visitor(full, entry);
    }
  }
}

/**
 * Walk only top-level directories (max depth 2) and build a structure map.
 * Returns { dirStats: Map<relDir, { purpose, extensions: Map<ext, count> }> }
 */
function buildStructureMap(cwd) {
  const extensions = new Map(); // ext -> count (global)
  const dirStats = new Map();   // relDir -> { extensions: Map<ext, count> }

  walk(cwd, (filePath) => {
    const ext = path.extname(filePath).toLowerCase();
    if (!ext) return;

    // Global count
    extensions.set(ext, (extensions.get(ext) || 0) + 1);

    // Per top-level directory
    const rel = path.relative(cwd, filePath);
    const topDir = rel.split(path.sep)[0];
    if (!topDir || topDir === rel) {
      // Root-level file — skip dir stats
      return;
    }
    if (!dirStats.has(topDir)) {
      dirStats.set(topDir, { extensions: new Map() });
    }
    const stats = dirStats.get(topDir);
    stats.extensions.set(ext, (stats.extensions.get(ext) || 0) + 1);
  }, 6); // depth 6 to get meaningful file counts, but only report top-level dirs

  return { extensions, dirStats };
}

/** Guess the purpose of a directory by name. */
function guessPurpose(dirName) {
  const lower = dirName.toLowerCase();
  const map = {
    'src': 'Application source',
    'lib': 'Library source',
    'app': 'Application source',
    'apps': 'Application packages',
    'packages': 'Workspace packages',
    'test': 'Test files',
    'tests': 'Test files',
    '__tests__': 'Test files',
    'spec': 'Test specifications',
    'specs': 'Test specifications',
    'e2e': 'End-to-end tests',
    'docs': 'Documentation',
    'doc': 'Documentation',
    'scripts': 'Build/utility scripts',
    'config': 'Configuration',
    'public': 'Static assets',
    'static': 'Static assets',
    'assets': 'Assets',
    'styles': 'Stylesheets',
    'components': 'UI components',
    'pages': 'Page routes',
    'views': 'View templates',
    'api': 'API routes',
    'server': 'Server-side code',
    'client': 'Client-side code',
    'shared': 'Shared utilities',
    'common': 'Shared utilities',
    'utils': 'Utility functions',
    'helpers': 'Helper functions',
    'types': 'Type definitions',
    'interfaces': 'Interface definitions',
    'models': 'Data models',
    'services': 'Service layer',
    'controllers': 'Controllers',
    'middleware': 'Middleware',
    'guards': 'Route guards',
    'pipes': 'Data pipes',
    'migrations': 'Database migrations',
    'seeds': 'Database seeds',
    'fixtures': 'Test fixtures',
    'mocks': 'Test mocks',
    'stubs': 'Test stubs',
    '.github': 'GitHub configuration',
    '.circleci': 'CircleCI configuration',
    '.vscode': 'VS Code settings',
    'ci': 'CI configuration',
    'docker': 'Docker configuration',
    'deploy': 'Deployment configuration',
    'infra': 'Infrastructure',
    'terraform': 'Terraform IaC',
    'quality': 'Quality tools/config',
    'artifacts': 'Build/project artifacts',
    'templates': 'Templates',
    'openspec': 'OpenSpec artifacts',
    '.atl': 'Agent Team Lite config',
  };
  return map[lower] || '';
}

/** Format extension counts into a human-readable string. */
function formatExtCounts(extMap) {
  const sorted = [...extMap.entries()]
    .sort((a, b) => b[1] - a[1])
    .slice(0, 6);
  return sorted.map(([ext, count]) => `${count} ${ext}`).join(', ');
}

// ---------------------------------------------------------------------------
// Detection functions
// ---------------------------------------------------------------------------

function detectStack(cwd) {
  const result = {
    name: null,
    version: null,
    description: null,
    engines: null,
    frameworks: [],
    testRunners: [],
    packageManager: null,
    monorepo: null,
    scripts: {},
    allDeps: {},       // combined deps + devDeps for lookup
    topDeps: [],       // top deps by relevance
    hasPkgJson: false,
  };

  const pkg = readJson(path.join(cwd, 'package.json'));
  if (!pkg) return result;

  result.hasPkgJson = true;
  result.name = pkg.name || null;
  result.version = pkg.version || null;
  result.description = pkg.description || null;
  result.scripts = pkg.scripts || {};

  // Engines
  if (pkg.engines && Object.keys(pkg.engines).length > 0) {
    result.engines = pkg.engines;
  }

  // Merge all deps
  const deps = Object.assign({}, pkg.dependencies || {});
  const devDeps = Object.assign({}, pkg.devDependencies || {});
  result.allDeps = Object.assign({}, deps, devDeps);

  // Framework detection
  const seenFrameworks = new Set();
  for (const [depName, fwName] of Object.entries(FRAMEWORK_MAP)) {
    if (result.allDeps[depName] && !seenFrameworks.has(fwName)) {
      seenFrameworks.add(fwName);
      const ver = result.allDeps[depName];
      result.frameworks.push({ name: fwName, version: ver });
    }
  }

  // Test runner detection
  const seenRunners = new Set();
  for (const [depName, runnerName] of Object.entries(TEST_RUNNER_MAP)) {
    if (result.allDeps[depName] && !seenRunners.has(runnerName)) {
      seenRunners.add(runnerName);
      result.testRunners.push(runnerName);
    }
  }

  // Package manager from lockfile
  for (const [lockfile, manager] of Object.entries(LOCKFILE_MAP)) {
    if (exists(path.join(cwd, lockfile))) {
      result.packageManager = manager;
      break;
    }
  }

  // Monorepo detection
  const pnpmWorkspace = path.join(cwd, 'pnpm-workspace.yaml');
  if (exists(pnpmWorkspace)) {
    result.monorepo = { type: 'pnpm', packages: listWorkspacePackages(cwd, 'pnpm') };
  } else if (pkg.workspaces) {
    const wsGlobs = Array.isArray(pkg.workspaces)
      ? pkg.workspaces
      : (pkg.workspaces.packages || []);
    result.monorepo = { type: 'workspaces', packages: listWorkspacePackages(cwd, 'npm', wsGlobs) };
  }

  // Top dependencies (by relevance: frameworks first, then deps, skip devDeps for top list)
  const depEntries = Object.entries(deps).map(([n, v]) => ({ name: n, version: v, dev: false }));
  const devDepEntries = Object.entries(devDeps).map(([n, v]) => ({ name: n, version: v, dev: true }));
  const all = [...depEntries, ...devDepEntries];

  // Sort: frameworks first, then non-dev, then alphabetical
  all.sort((a, b) => {
    const aFw = FRAMEWORK_MAP[a.name] ? 0 : 1;
    const bFw = FRAMEWORK_MAP[b.name] ? 0 : 1;
    if (aFw !== bFw) return aFw - bFw;
    if (a.dev !== b.dev) return a.dev ? 1 : -1;
    return a.name.localeCompare(b.name);
  });
  result.topDeps = all.slice(0, 10).map(d => ({
    name: d.name,
    version: d.version,
    category: categorize(d.name, d.dev),
  }));

  return result;
}

/** Categorize a dependency. */
function categorize(name, isDev) {
  if (FRAMEWORK_MAP[name]) return 'Framework';
  if (TEST_RUNNER_MAP[name]) return 'Testing';
  if (/eslint|prettier|stylelint/.test(name)) return 'Linting';
  if (/webpack|vite|rollup|esbuild|tsup|unbuild|turbo/.test(name)) return 'Build';
  if (/typescript|ts-node|tsx/.test(name)) return 'TypeScript';
  if (isDev) return 'Dev Dependency';
  return 'Dependency';
}

/** List workspace packages from a monorepo. */
function listWorkspacePackages(cwd, type, globs) {
  const packages = [];

  if (type === 'pnpm') {
    // Read pnpm-workspace.yaml — simple parse without yaml dep
    try {
      const content = fs.readFileSync(path.join(cwd, 'pnpm-workspace.yaml'), 'utf-8');
      const wsGlobs = [];
      let inPackages = false;
      for (const line of content.split('\n')) {
        if (/^packages:/.test(line.trim())) {
          inPackages = true;
          continue;
        }
        if (inPackages) {
          const match = line.match(/^\s+-\s+['"]?([^'"]+)['"]?/);
          if (match) {
            wsGlobs.push(match[1]);
          } else if (/^\S/.test(line) && line.trim()) {
            inPackages = false;
          }
        }
      }
      globs = wsGlobs;
    } catch {
      return packages;
    }
  }

  if (!globs || globs.length === 0) return packages;

  // Expand globs (simple: only handle `dir/*` patterns)
  for (const glob of globs) {
    const baseDir = glob.replace(/\/?\*$/, '').replace(/\/?\*\*$/, '');
    const fullBase = path.join(cwd, baseDir);
    if (!isDir(fullBase)) continue;

    for (const entry of readDir(fullBase)) {
      if (!entry.isDirectory()) continue;
      const pkgJson = readJson(path.join(fullBase, entry.name, 'package.json'));
      if (pkgJson) {
        packages.push(pkgJson.name || entry.name);
      } else {
        packages.push(entry.name);
      }
    }
  }

  return packages;
}

function detectTests(cwd, stack) {
  const result = {
    configs: [],
    testFileCount: 0,
    coverageConfig: null,
    e2eRunner: null,
  };

  // Test config files
  for (const pattern of TEST_CONFIG_PATTERNS) {
    if (exists(path.join(cwd, pattern))) {
      result.configs.push(pattern);
    }
  }

  // Count test files
  walk(cwd, (filePath) => {
    const rel = path.relative(cwd, filePath);
    const base = path.basename(filePath);
    if (
      /\.(spec|test)\.(ts|tsx|js|jsx|mjs|cjs)$/.test(base) ||
      /^__tests__\//.test(rel) ||
      /\/__tests__\//.test(rel)
    ) {
      result.testFileCount++;
    }
  }, 8);

  // E2E runner
  if (stack.allDeps['@playwright/test'] || stack.allDeps['playwright']) {
    result.e2eRunner = 'playwright';
    const playwrightConfig = ['playwright.config.ts', 'playwright.config.js']
      .find(f => exists(path.join(cwd, f)));
    if (playwrightConfig) {
      result.e2eRunner = `playwright (${playwrightConfig} found)`;
    }
  } else if (stack.allDeps['cypress']) {
    result.e2eRunner = 'cypress';
    const cypressConfig = ['cypress.config.ts', 'cypress.config.js']
      .find(f => exists(path.join(cwd, f)));
    if (cypressConfig) {
      result.e2eRunner = `cypress (${cypressConfig} found)`;
    }
  }

  // Coverage config
  if (stack.allDeps['c8'] || stack.allDeps['@vitest/coverage-c8']) {
    result.coverageConfig = 'c8';
  } else if (stack.allDeps['nyc'] || stack.allDeps['@istanbuljs/nyc-config-typescript']) {
    result.coverageConfig = 'nyc (istanbul)';
  } else if (stack.allDeps['@vitest/coverage-v8']) {
    result.coverageConfig = 'v8';
  } else if (stack.allDeps['@vitest/coverage-istanbul']) {
    result.coverageConfig = 'istanbul (via vitest)';
  }

  // Check package.json for coverage-related config
  if (!result.coverageConfig) {
    const pkg = readJson(path.join(cwd, 'package.json'));
    if (pkg) {
      if (pkg.nyc) result.coverageConfig = 'nyc (istanbul)';
      if (pkg.c8) result.coverageConfig = 'c8';
      if (pkg.jest && pkg.jest.coverageDirectory) result.coverageConfig = 'jest built-in';
    }
  }

  return result;
}

function detectCICD(cwd) {
  const result = {
    github: [],
    gitlab: false,
    jenkins: false,
    dockerfile: false,
    dockerCompose: false,
    circleci: false,
  };

  // GitHub Actions
  const ghDir = path.join(cwd, '.github', 'workflows');
  if (isDir(ghDir)) {
    for (const entry of readDir(ghDir)) {
      if (entry.isFile() && /\.(yml|yaml)$/.test(entry.name)) {
        result.github.push(entry.name);
      }
    }
  }

  result.gitlab = exists(path.join(cwd, '.gitlab-ci.yml'));
  result.jenkins = exists(path.join(cwd, 'Jenkinsfile'));
  result.dockerfile = exists(path.join(cwd, 'Dockerfile'));
  result.dockerCompose = exists(path.join(cwd, 'docker-compose.yml'))
    || exists(path.join(cwd, 'docker-compose.yaml'));
  result.circleci = exists(path.join(cwd, '.circleci', 'config.yml'));

  return result;
}

function detectEchoReadiness(scripts) {
  return ECHO_STEPS.map(({ step, scripts: candidates }) => {
    const found = candidates.find(s => scripts[s]);
    return {
      step,
      script: found ? `\`${found}\`` : '--',
      available: !!found,
    };
  });
}

function detectKeyFiles(cwd) {
  const result = KEY_FILES.map(f => ({
    name: f,
    exists: exists(path.join(cwd, f)),
  }));

  // Check for docs/*.md at root of docs/
  const docsDir = path.join(cwd, 'docs');
  if (isDir(docsDir)) {
    for (const entry of readDir(docsDir)) {
      if (entry.isFile() && /\.md$/i.test(entry.name)) {
        const name = `docs/${entry.name}`;
        if (!result.find(r => r.name === name)) {
          result.push({ name, exists: true });
        }
      }
    }
  }

  // ADR detection
  const adrDirs = ['docs/adr', 'docs/adrs', 'adr', 'adrs', 'docs/architecture/decisions'];
  for (const dir of adrDirs) {
    const adrPath = path.join(cwd, dir);
    if (isDir(adrPath)) {
      let count = 0;
      for (const entry of readDir(adrPath)) {
        if (entry.isFile() && /\.md$/i.test(entry.name)) count++;
      }
      if (count > 0) {
        result.push({ name: `${dir}/ (${count} ADR files)`, exists: true });
      }
      break;
    }
  }

  return result;
}

// ---------------------------------------------------------------------------
// Audit functions
// ---------------------------------------------------------------------------

function scoreF1(keyFiles, docsCount) {
  const hasReadme = keyFiles.find(f => f.name === 'README.md' && f.exists);
  const docFiles = keyFiles.filter(f => f.exists && f.name !== 'README.md' && f.name !== 'LICENSE');
  const totalDocs = docFiles.length + docsCount;

  if (!hasReadme && totalDocs === 0) return { score: 0, reasoning: 'No documentation found' };
  if (hasReadme && totalDocs === 0) return { score: 0.5, reasoning: 'README exists, no additional documentation' };
  if (hasReadme && totalDocs < 3) return { score: 1, reasoning: `README + ${totalDocs} additional doc(s)` };
  return { score: 2, reasoning: `Comprehensive documentation (README + ${totalDocs} docs)` };
}

function scoreF2(stack, structure) {
  const hasFramework = stack.frameworks.length > 0;
  const hasPkgJson = stack.hasPkgJson;
  const hasLinter = Object.keys(stack.scripts).some(s => /lint|eslint|test:static/.test(s))
    || !!stack.allDeps['eslint'];
  const hasFormatter = Object.keys(stack.scripts).some(s => /format|prettier/.test(s))
    || !!stack.allDeps['prettier'];
  const hasConventionalDirs = ['src', 'lib', 'app', 'test', 'tests'].some(d =>
    structure.dirStats.has(d)
  );

  let points = 0;
  if (hasPkgJson) points += 0.5;
  if (hasFramework) points += 0.5;
  if (hasLinter || hasFormatter) points += 0.5;
  if (hasConventionalDirs) points += 0.5;

  if (points <= 0.5) return { score: 0, reasoning: 'No recognizable patterns or standards' };
  if (points <= 1) return { score: 1, reasoning: 'Some standards detected' };
  return { score: 2, reasoning: `Clear standards (${hasFramework ? stack.frameworks.map(f => f.name).join(', ') + ', ' : ''}conventional structure)` };
}

function scoreF3(stack, keyFiles, tests) {
  const hasReadme = keyFiles.find(f => f.name === 'README.md' && f.exists);
  const hasTests = tests.testFileCount > 0;
  const hasSpecs = keyFiles.some(f =>
    /spec|design|adr/i.test(f.name) && f.exists
  );

  if (!hasReadme) return { score: 0, reasoning: 'Purpose unclear (no README)' };
  if (hasReadme && !hasTests && !hasSpecs) return { score: 1, reasoning: 'Purpose clear, implementation coverage unclear' };
  return { score: 2, reasoning: 'Clear purpose and implementation (README + tests/specs)' };
}

function generateAuditSection(stack, structure, tests, cicd, echoReadiness, keyFiles) {
  const lines = [];

  // Count docs in docs/ dir
  const docsCount = keyFiles.filter(f => f.exists && f.name.startsWith('docs/')).length;

  const f1 = scoreF1(keyFiles, docsCount);
  const f2 = scoreF2(stack, structure);
  const f3 = scoreF3(stack, keyFiles, tests);
  const f4 = { score: 2, reasoning: 'Codebase exists (takeover fixed)' };
  const total = f1.score + f2.score + f3.score + f4.score;

  // Determine tier
  let tier;
  if (total >= 6) tier = 'Ligero';
  else if (total >= 3) tier = 'Estandar';
  else tier = 'Completo';

  lines.push('### Audit: fastForward Scoring');
  lines.push('');
  lines.push('| Factor | Score | Reasoning |');
  lines.push('|--------|-------|-----------|');
  lines.push(`| F1 (Artifacts) | ${f1.score} | ${f1.reasoning} |`);
  lines.push(`| F2 (Standards) | ${f2.score} | ${f2.reasoning} |`);
  lines.push(`| F3 (Ambiguity) | ${f3.score} | ${f3.reasoning} |`);
  lines.push(`| F4 (Reference) | ${f4.score} | ${f4.reasoning} |`);
  lines.push(`| **Total** | **${total}** | **Tier ${tier} probable** |`);
  lines.push('');

  // Artifact Equivalence
  lines.push('### Audit: Artifact Equivalence');
  lines.push('');
  lines.push('| Framework Artifact | Equivalent Found | Coverage |');
  lines.push('|-------------------|-----------------|----------|');

  const hasReadme = keyFiles.find(f => f.name === 'README.md' && f.exists);
  const hasDescription = stack.description && stack.description.length > 10;
  const ideaCoverage = hasReadme
    ? (hasDescription ? 'Partial -- has purpose' : 'Partial -- exists but minimal')
    : 'Missing';
  const ideaEquiv = hasReadme ? 'README.md' : 'Not found';
  lines.push(`| idea.md | ${ideaEquiv} | ${ideaCoverage} |`);

  const specEquiv = tests.testFileCount > 0
    ? `Test suite (${tests.testFileCount} test files)`
    : 'Not found';
  const specCoverage = tests.testFileCount > 5
    ? 'Partial -- tests exist, no formal acceptance criteria'
    : tests.testFileCount > 0
      ? 'Minimal -- few tests, no formal specs'
      : 'Missing';
  lines.push(`| spec.md | ${specEquiv} | ${specCoverage} |`);

  const hasAdr = keyFiles.some(f => /adr/i.test(f.name) && f.exists);
  const hasContrib = keyFiles.find(f => f.name === 'CONTRIBUTING.md' && f.exists);
  const designEquiv = hasAdr
    ? 'ADR files'
    : hasContrib
      ? 'CONTRIBUTING.md'
      : 'Not found';
  const designCoverage = hasAdr ? 'Partial -- ADRs exist' : hasContrib ? 'Minimal -- contributing guide only' : 'Missing';
  lines.push(`| design.md | ${designEquiv} | ${designCoverage} |`);

  const echoAvailable = echoReadiness.filter(e => e.available).length;
  const echoTotal = echoReadiness.length;
  const hasCi = cicd.github.length > 0 || cicd.gitlab || cicd.jenkins || cicd.circleci;
  const echoEquiv = hasCi
    ? `CI/CD + scripts (${echoAvailable}/${echoTotal} steps)`
    : `Scripts only (${echoAvailable}/${echoTotal} steps)`;
  const echoCoverage = echoAvailable === echoTotal && hasCi
    ? 'Complete'
    : echoAvailable >= 3
      ? 'Partial -- most steps covered'
      : 'Minimal -- few automation steps';
  lines.push(`| echo | ${echoEquiv} | ${echoCoverage} |`);
  lines.push('');

  // Gaps
  lines.push('### Audit: Gaps');
  lines.push('');

  const gaps = [];

  if (!hasReadme) gaps.push({ done: false, msg: 'No README.md (idea.md equivalent missing)' });
  else gaps.push({ done: true, msg: 'README.md exists (idea.md equivalent)' });

  if (tests.testFileCount === 0) gaps.push({ done: false, msg: 'No test files found (spec.md equivalent missing)' });
  else gaps.push({ done: true, msg: `Test suite exists (${tests.testFileCount} test files)` });

  if (!hasAdr) gaps.push({ done: false, msg: 'No ADRs found (design.md equivalent missing)' });
  else gaps.push({ done: true, msg: 'ADRs found (design.md equivalent)' });

  const hasLint = echoReadiness.find(e => e.step === 'Static' && e.available)
    || !!stack.allDeps['eslint'];
  if (!hasLint) gaps.push({ done: false, msg: 'No linter configured' });
  else gaps.push({ done: true, msg: 'Linter configured' });

  if (!hasCi) gaps.push({ done: false, msg: 'No CI/CD pipeline detected' });
  else gaps.push({ done: true, msg: 'CI/CD pipeline configured' });

  if (!tests.coverageConfig) gaps.push({ done: false, msg: 'No code coverage configured' });
  else gaps.push({ done: true, msg: `Code coverage configured (${tests.coverageConfig})` });

  const hasE2e = echoReadiness.find(e => e.step === 'E2E' && e.available);
  if (!hasE2e && !tests.e2eRunner) gaps.push({ done: false, msg: 'No E2E testing configured' });
  else gaps.push({ done: true, msg: 'E2E testing configured' });

  for (const gap of gaps) {
    lines.push(`- [${gap.done ? 'x' : ' '}] ${gap.msg}`);
  }
  lines.push('');

  // Recommended entry point
  lines.push('### Audit: Recommended Entry Point');
  lines.push('');

  if (total >= 6) {
    lines.push('**Scenario**: Codebase heredado, cambios planificados, bien documentado');
    lines.push('**Entry**: fastForward -- docs existentes cuentan como artefactos parciales');
    lines.push(`**Probable tier**: Ligero (score ${total})`);
  } else if (total >= 3) {
    lines.push('**Scenario**: Codebase heredado, cambios planificados, documentacion parcial');
    lines.push('**Entry**: planning -- generar artefactos faltantes, fastForward parcial posible');
    lines.push(`**Probable tier**: Estandar (score ${total})`);
  } else {
    lines.push('**Scenario**: Codebase heredado con baja certeza o sin documentacion');
    lines.push('**Entry**: planning completo -- generar todos los artefactos desde cero');
    lines.push(`**Probable tier**: Completo (score ${total})`);
  }

  return lines.join('\n');
}

// ---------------------------------------------------------------------------
// Main scan function
// ---------------------------------------------------------------------------

/**
 * Scan a codebase and return structured Markdown.
 * @param {string} cwd — the directory to scan
 * @param {{ audit?: boolean }} options
 * @returns {string} Markdown output
 */
function scan(cwd, options = {}) {
  const lines = [];

  try {
    if (!isDir(cwd)) {
      return `## Project Context (generated by virgil scan)\n\n> Error: "${cwd}" is not a valid directory.\n`;
    }

    const stack = detectStack(cwd);
    const structure = buildStructureMap(cwd);
    const tests = detectTests(cwd, stack);
    const cicd = detectCICD(cwd);
    const echoReadiness = detectEchoReadiness(stack.scripts);
    const keyFiles = detectKeyFiles(cwd);

    lines.push('## Project Context (generated by virgil scan)');
    lines.push('');

    // ---- Stack ----
    lines.push('### Stack');
    lines.push('');
    if (stack.hasPkgJson) {
      if (stack.name) lines.push(`- **Name**: ${stack.name}`);
      if (stack.version) lines.push(`- **Version**: ${stack.version}`);
      if (stack.frameworks.length > 0) {
        lines.push(`- **Framework**: ${stack.frameworks.map(f => `${f.name} ${f.version}`).join(', ')}`);
      }
      if (stack.testRunners.length > 0) {
        lines.push(`- **Test runner**: ${stack.testRunners.join(', ')}`);
      }
      if (stack.packageManager) {
        lines.push(`- **Package manager**: ${stack.packageManager}`);
      }
      if (stack.engines) {
        const engineStr = Object.entries(stack.engines).map(([k, v]) => `${k} ${v}`).join(', ');
        lines.push(`- **Engines**: ${engineStr}`);
      }
      if (stack.monorepo) {
        const pkgList = stack.monorepo.packages.length > 0
          ? ` (${stack.monorepo.packages.length} packages: ${stack.monorepo.packages.join(', ')})`
          : '';
        lines.push(`- **Monorepo**: yes${pkgList}`);
      }
    } else {
      lines.push('- **Package manager**: Not detected (no package.json)');
      // Try to detect language from file extensions
      const mainExts = [...structure.extensions.entries()]
        .sort((a, b) => b[1] - a[1])
        .slice(0, 3);
      if (mainExts.length > 0) {
        const langs = mainExts
          .map(([ext]) => EXTENSION_CATEGORIES[ext] || ext)
          .filter((v, i, a) => a.indexOf(v) === i);
        lines.push(`- **Primary languages**: ${langs.join(', ')}`);
      }
    }
    lines.push('');

    // ---- Structure ----
    lines.push('### Structure');
    lines.push('');
    if (structure.dirStats.size > 0) {
      lines.push('| Directory | Purpose | Files |');
      lines.push('|-----------|---------|-------|');
      const sorted = [...structure.dirStats.entries()].sort((a, b) => a[0].localeCompare(b[0]));
      for (const [dir, stats] of sorted) {
        const purpose = guessPurpose(dir);
        const files = formatExtCounts(stats.extensions);
        lines.push(`| ${dir}/ | ${purpose} | ${files} |`);
      }
    } else {
      lines.push('> No directories found (flat project).');
    }
    lines.push('');

    // ---- Dependencies ----
    if (stack.topDeps.length > 0) {
      lines.push('### Dependencies (top 10 by relevance)');
      lines.push('');
      lines.push('| Package | Version | Category |');
      lines.push('|---------|---------|----------|');
      for (const dep of stack.topDeps) {
        lines.push(`| ${dep.name} | ${dep.version} | ${dep.category} |`);
      }
      lines.push('');
    }

    // ---- Tests ----
    lines.push('### Tests');
    lines.push('');
    if (stack.testRunners.length > 0) {
      const mainRunner = stack.testRunners[0];
      const configFound = tests.configs.find(c => c.startsWith(mainRunner));
      lines.push(`- **Runner**: ${mainRunner}${configFound ? ` (${configFound} found)` : ''}`);
    } else if (tests.configs.length > 0) {
      lines.push(`- **Runner**: detected from config (${tests.configs[0]})`);
    } else {
      lines.push('- **Runner**: not detected');
    }
    lines.push(`- **Test files**: ${tests.testFileCount}${tests.testFileCount > 0 ? ' (*.spec.* / *.test.* pattern)' : ''}`);
    lines.push(`- **Coverage**: ${tests.coverageConfig || 'not configured'}`);
    if (tests.e2eRunner) {
      lines.push(`- **E2E**: ${tests.e2eRunner}`);
    }
    lines.push('');

    // ---- CI/CD ----
    lines.push('### CI/CD');
    lines.push('');
    const ciItems = [];
    if (cicd.github.length > 0) ciItems.push(`GitHub Actions: ${cicd.github.join(', ')}`);
    if (cicd.gitlab) ciItems.push('GitLab CI: .gitlab-ci.yml');
    if (cicd.jenkins) ciItems.push('Jenkins: Jenkinsfile');
    if (cicd.circleci) ciItems.push('CircleCI: .circleci/config.yml');
    if (cicd.dockerfile) ciItems.push('Docker: Dockerfile');
    if (cicd.dockerCompose) ciItems.push('Docker Compose: docker-compose.yml');

    if (ciItems.length > 0) {
      for (const item of ciItems) {
        lines.push(`- ${item}`);
      }
    } else {
      lines.push('- No CI/CD configuration detected');
    }
    lines.push('');

    // ---- Echo Readiness ----
    lines.push('### Echo Readiness');
    lines.push('');
    lines.push('| Step | Script | Status |');
    lines.push('|------|--------|--------|');
    for (const echo of echoReadiness) {
      const status = echo.available ? 'Available' : 'Missing';
      const mark = echo.available ? '✓' : '✗';
      lines.push(`| ${echo.step} | ${echo.script} | ${mark} ${status} |`);
    }
    lines.push('');

    // ---- Key Files ----
    lines.push('### Key Files');
    lines.push('');
    for (const file of keyFiles) {
      const mark = file.exists ? '✓' : '✗';
      lines.push(`- ${file.name} ${mark}`);
    }
    lines.push('');

    // ---- Audit ----
    if (options.audit) {
      lines.push(generateAuditSection(stack, structure, tests, cicd, echoReadiness, keyFiles));
      lines.push('');
    }

  } catch (err) {
    lines.push('## Project Context (generated by virgil scan)');
    lines.push('');
    lines.push(`> Scan error: ${err.message}`);
    lines.push('');
  }

  return lines.join('\n');
}

module.exports = { scan };
