"""Build a standalone course reader using only Python's standard library."""
from pathlib import Path
import re

ROOT = Path(__file__).resolve().parent
lessons = "\n".join(p.read_text(encoding="utf-8") for p in sorted((ROOT / "lessons").glob("*.html")))
style = (ROOT / "reader.css").read_text(encoding="utf-8")
script = (ROOT / "reader.js").read_text(encoding="utf-8")
assert sum(map(int, re.findall(r'data-minutes="(\d+)"', lessons))) == 840
assert len(re.findall(r'<section\b', lessons)) == 14
page = """<!doctype html>
<html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>Netflix L5 Go Interview Course</title><link rel="icon" href="data:,"><style>STYLE</style></head>
<body><a class="skip" href="#reading">Skip to lesson</a>
<aside><div class="brand">THE CODING INTERVIEW</div><div class="course-name">Reason it through.</div>
<p class="subtitle">Go · Netflix L5 · Live practice</p>
<label for="search">Find a topic</label><input id="search" type="search" placeholder="Try expiration or graphs">
<nav id="navigation" aria-label="Lessons"></nav>
<p id="progress" aria-live="polite"></p><p class="aside-note">14 hours of active study.<br>Read. Explain. Implement. Test.</p></aside>
<div class="workspace"><header><div><span class="eyebrow">YOUR PREPARATION WORKBOOK</span><p id="lesson-meta">Self-paced course</p></div><button id="print">Print all lessons</button></header>
<div class="toolbar" aria-label="Listening controls">
<label for="voice">Voice <select id="voice"><option value="">System default</option></select></label>
<label for="rate">Speed <select id="rate"><option value="0.85">0.85×</option><option value="1" selected>1×</option><option value="1.15">1.15×</option><option value="1.3">1.3×</option></select></label>
<button id="listen" class="primary">Listen</button><button id="pause">Pause</button><button id="continue">Continue</button><button id="stop">Stop</button>
<span id="audio-status" role="status">Narration pauses at questions.</span></div>
<div class="learning-actions"><button id="tutor">Copy AI tutor prompt</button><button id="export">Download lesson text</button><label><input type="checkbox" id="complete"> I can explain and apply this</label></div>
<div id="copy-fallback" hidden><label for="prompt-text">Copy this prompt into your AI conversation</label><textarea id="prompt-text" rows="8"></textarea></div>
<main id="reading" tabindex="-1">LESSONS</main>
<footer><button id="previous">← Previous lesson</button><button id="next" class="primary">Next lesson →</button><p>Original study material · Sources checked 30 September 2026</p></footer></div>
<noscript><p>JavaScript is off. All lessons remain readable; narration and navigation controls need JavaScript.</p></noscript>
<script>SCRIPT</script></body></html>"""
page = page.replace("STYLE", style).replace("LESSONS", lessons).replace("SCRIPT", script)
(ROOT / "course.html").write_text(page, encoding="utf-8")
print(f"Built course.html: {len(page):,} characters, 14 lessons, 840 active-study minutes")
