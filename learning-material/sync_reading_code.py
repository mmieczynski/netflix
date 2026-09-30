"""Extract the book's Go examples for maintainer compilation and testing."""
from pathlib import Path
import re

ROOT = Path(__file__).resolve().parent

def sync():
    parts = ['// Code generated from reading Markdown; DO NOT EDIT.',
             'package readingcode', 'import "sort"']
    for path in sorted((ROOT/'reading').glob('[0-9]*.md')):
        text = path.read_text(encoding='utf-8')
        section = text.split('## Go example:',1)[1]
        blocks = re.findall(r'```go\n(.*?)\n```',section,re.S)
        assert blocks, path
        parts.append('// Source: reading/'+path.name)
        parts.extend(blocks)
    target = ROOT/'readingcode'
    target.mkdir(exist_ok=True)
    (target/'examples.go').write_text('\n\n'.join(parts)+'\n',encoding='utf-8')

if __name__ == '__main__':
    sync()
