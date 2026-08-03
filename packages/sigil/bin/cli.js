#!/usr/bin/env node

const fs = require('node:fs');
const path = require('node:path');

const TEMPLATE = path.join(__dirname, '..', 'templates', 'AGENTS.md');
const TARGET = path.join(process.cwd(), 'AGENTS.md');

const command = process.argv[2];

const commands = {
  init() {
    if (fs.existsSync(TARGET) && !process.argv.includes('--force')) {
      console.error('AGENTS.md already exists. Use --force to overwrite.');
      process.exit(1);
    }

    if (!fs.existsSync(TEMPLATE)) {
      console.error('Template not found. Package may be corrupted.');
      process.exit(1);
    }

    fs.copyFileSync(TEMPLATE, TARGET);

    const size = (fs.statSync(TARGET).size / 1024).toFixed(1);
    console.log(`AGENTS.md created (${size} KB)`);
    console.log('Your AI agent now has a methodology. Open it and start working.');
  },

  version() {
    const pkg = require('../package.json');
    console.log(pkg.version);
  },

  help() {
    console.log(`sigil v${require('../package.json').version}`);
    console.log('');
    console.log('Usage: sigil <command>');
    console.log('');
    console.log('Commands:');
    console.log('  init      Install AGENTS.md methodology into current project');
    console.log('  version   Show version');
    console.log('  help      Show this help');
    console.log('');
    console.log('Options:');
    console.log('  --force   Overwrite existing AGENTS.md');
  },
};

const handler = commands[command] || commands[command === '--version' || command === '-v' ? 'version' : 'help'];
handler();
