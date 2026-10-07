"""Extract every reading-edition Go block for compilation and behavior checks.

All fences contain declarations, not loose statement fragments. Standard-library
imports are supplied here because examples are printed without package boilerplate.
"""
from pathlib import Path
import re

ROOT = Path(__file__).resolve().parent


def sync():
    parts = ['// Code generated from reading Markdown; DO NOT EDIT.\n'
             'package readingcode\n\n'
             'import (\n\t"container/heap"\n\t"sort"\n\t"sync"\n)']
    blocks_total = 0
    for path in sorted((ROOT / 'reading').glob('[0-9]*.md')):
        text = path.read_text(encoding='utf-8')
        blocks = re.findall(r'```go\n(.*?)\n```', text, re.S)
        assert blocks, f'No Go examples in {path}'
        parts.append('// Source: reading/' + path.name)
        parts.extend(blocks)
        blocks_total += len(blocks)
    target = ROOT / 'readingcode'
    target.mkdir(exist_ok=True)
    (target / 'examples.go').write_text('\n\n'.join(parts) + '\n', encoding='utf-8')
    return blocks_total


if __name__ == '__main__':
    print(f'Extracted {sync()} Go blocks')
