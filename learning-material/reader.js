'use strict';
const sections = Array.from(document.querySelectorAll('main section'));
const $ = id => document.getElementById(id);
const synth = window.speechSynthesis;
let selected = 0, voices = [], queue = [], cursor = 0, generation = 0;
let waiting = false, active = false, utterance = null, highlighted = null;
let completed = {};
try { completed = JSON.parse(localStorage.getItem('netflix-go-progress-v1') || '{}') || {}; } catch (_) {}
const tutorInstructions = `You are my interactive tutor for a Netflix L5 live coding interview in Go with another senior engineer. This is preparation, not assistance during a real interview. Use the lesson below as the source. Do not claim its exercises are actual Netflix questions.

Begin by asking whether I want learn mode, retrieval practice, or mock interview. Ask one question at a time and wait for my spoken reply. In learn mode, teach in short spoken passages, then ask me to predict or explain. Do not read Go punctuation aloud: explain state changes and invariants in words. When code is needed, ask me to switch to typing and inspect what I actually provide. Never pretend to have run code.

Do not reveal the hidden coach notes or solutions before I attempt the question. If I struggle, first ask for a tiny example, then offer one conceptual hint, then a stronger hint only if needed. Require me to justify a data structure by its operations, state an invariant, trace a boundary case, and give honest time and space complexity. Challenge weak claims with small counterexamples. Let me correct myself before explaining.

After success, change one assumption and ask whether my approach still works. Revisit a previous mistake after a few turns. In mock mode, give only the candidate brief, answer clarifications, and introduce follow-ups after meaningful progress without teaching unless I ask for help. Ask me to set a timer rather than pretending to measure elapsed time. Finish with demonstrated strengths, specific mistakes, and one transfer exercise. Keep an error log of assumption, counterexample, invariant, and transfer. Treat lesson hints and interviewer notes as private coaching guidance until appropriate.

LESSON CONTENT FOLLOWS:\n\n`;

function lessonText(section) {
  const copy = section.cloneNode(true);
  copy.querySelectorAll('h1,h2,h3,p,li,div,summary,pre,tr').forEach(el => {
    el.prepend(document.createTextNode('\n'));
    el.append(document.createTextNode('\n'));
  });
  copy.querySelectorAll('td,th').forEach(el => el.append(document.createTextNode(' | ')));
  return copy.textContent.replace(/[ \t]+/g, ' ').replace(/\n[ \t]+/g,'\n').replace(/\n{3,}/g,'\n\n').trim();
}

function renderNav() {
  const term = $('search').value.toLowerCase().trim();
  const nav = $('navigation'); nav.replaceChildren();
  sections.forEach((s, i) => {
    if (term && !s.textContent.toLowerCase().includes(term)) return;
    const a = document.createElement('a');
    a.href = '#' + s.id;
    a.className = i === selected ? 'active' : '';
    if (i === selected) a.setAttribute('aria-current', 'page');
    a.textContent = (completed[s.id] ? '✓ ' : '') + (i < 13 ? i + '. ' : '') + s.dataset.title;
    const sub = document.createElement('span'); sub.className = 'nav-minutes';
    sub.textContent = Number(s.dataset.minutes) ? s.dataset.minutes + ' minutes · active study' : 'Reference';
    a.append(sub); a.addEventListener('click', e => { e.preventDefault(); show(i, true); }); nav.append(a);
  });
  if (!nav.children.length) { const p = document.createElement('p'); p.className='no-results'; p.textContent='No matching lessons.'; nav.append(p); }
  const core = sections.filter(s => Number(s.dataset.minutes));
  const done = core.filter(s => completed[s.id]);
  $('progress').textContent = done.length + ' of ' + core.length + ' core lessons marked ready · ' + done.reduce((n,s)=>n+Number(s.dataset.minutes),0) + ' / 840 study minutes';
}

function clearHighlight() { if (highlighted) highlighted.classList.remove('speaking'); highlighted = null; }
function stopAudio(message = 'Narration stopped.') {
  generation++; active = false; waiting = false; queue = []; cursor = 0;
  if (synth) synth.cancel(); utterance = null; clearHighlight();
  $('audio-status').textContent = message;
}

function show(i, scroll = false) {
  stopAudio('Narration pauses at questions.'); selected = i;
  sections.forEach((s,j) => { s.hidden = i !== j; });
  $('complete').checked = Boolean(completed[sections[i].id]);
  $('lesson-meta').textContent = sections[i].dataset.title + (Number(sections[i].dataset.minutes) ? ' · ' + sections[i].dataset.minutes + ' minutes' : ' · Research and scope');
  $('previous').disabled = i === 0; $('next').disabled = i === sections.length - 1;
  $('copy-fallback').hidden = true;
  try { history.replaceState(null, '', '#' + sections[i].id); } catch (_) {}
  renderNav();
  if (scroll) { $('reading').scrollIntoView({behavior:'smooth',block:'start'}); $('reading').focus({preventScroll:true}); }
}

function buildQueue(section) {
  const result = [];
  section.querySelectorAll('h1,h2,h3,p,li,.checkpoint').forEach(el => {
    if (el.closest('details,pre,table') || (el.closest('.checkpoint') && !el.classList.contains('checkpoint')) || (el.parentElement.closest('li') && el.tagName !== 'LI')) return;
    const text = el.textContent.trim().replace(/\s+/g, ' ');
    if (!text) return;
    // Short utterances avoid long browser speech queues. Preserve a checkpoint stop
    // only after its last chunk, so its question is always fully read.
    const sentences = text.match(/[^.!?]+[.!?]+(?:\s|$)|[^.!?]+$/g) || [text];
    let chunk = '';
    const chunks = [];
    sentences.forEach(sentence => {
      if (chunk.length + sentence.length > 300 && chunk) { chunks.push(chunk); chunk=''; }
      if (sentence.length > 350) {
        if (chunk) { chunks.push(chunk); chunk=''; }
        const words = sentence.split(/\s+/);
        words.forEach(word => { if ((chunk+word).length > 280) { chunks.push(chunk); chunk=''; } chunk += word+' '; });
      } else chunk += sentence;
    });
    if (chunk.trim()) chunks.push(chunk.trim());
    chunks.forEach((text,j) => result.push({text,el,checkpoint:el.classList.contains('checkpoint') && j===chunks.length-1}));
  });
  return result;
}

function speakNext(token) {
  if (token !== generation || !active) return;
  if (cursor >= queue.length) { stopAudio('Lesson narration complete. Now do the timed practice.'); return; }
  const part = queue[cursor]; clearHighlight(); highlighted=part.el; highlighted.classList.add('speaking');
  highlighted.scrollIntoView({behavior:'smooth',block:'center'});
  utterance = new SpeechSynthesisUtterance(part.text);
  const chosen = voices.find(v => v.voiceURI === $('voice').value);
  if (chosen) utterance.voice = chosen;
  utterance.rate = Number($('rate').value);
  utterance.onend = () => {
    if (token !== generation) return;
    cursor++;
    if (part.checkpoint) { waiting = true; $('audio-status').textContent='Your turn. Answer aloud, then press Continue. Answers remain hidden.'; }
    else speakNext(token);
  };
  utterance.onerror = e => {
    if (token !== generation) return;
    stopAudio('Speech unavailable ('+e.error+'). Try another voice or use the AI tutor prompt.');
  };
  $('audio-status').textContent='Listening · code, tables, and answer panels are skipped.';
  synth.speak(utterance);
}

function loadVoices() {
  if (!synth) return;
  const previous = $('voice').value; voices=synth.getVoices(); $('voice').replaceChildren(new Option('System default',''));
  voices.forEach(v => $('voice').add(new Option(v.name+' · '+v.lang, v.voiceURI)));
  if (voices.some(v => v.voiceURI === previous)) $('voice').value=previous;
}

$('listen').addEventListener('click', () => {
  if (!synth) { $('audio-status').textContent='This browser has no speech synthesis. Use the AI tutor prompt.'; return; }
  stopAudio(); queue=buildQueue(sections[selected]); active=true; speakNext(generation);
});
$('pause').addEventListener('click', () => { if (synth && active && !waiting) { synth.pause(); $('audio-status').textContent='Paused. Press Continue to resume.'; } });
$('continue').addEventListener('click', () => {
  if (!synth || !active) return;
  if (waiting) { waiting=false; speakNext(generation); }
  else if (synth.paused) { synth.resume(); $('audio-status').textContent='Listening.'; }
});
$('stop').addEventListener('click', () => stopAudio());
$('search').addEventListener('input', renderNav);
$('previous').addEventListener('click', () => { if (selected>0) show(selected-1,true); });
$('next').addEventListener('click', () => { if (selected<sections.length-1) show(selected+1,true); });
$('complete').addEventListener('change', () => {
  completed[sections[selected].id]=$('complete').checked;
  try { localStorage.setItem('netflix-go-progress-v1',JSON.stringify(completed)); } catch (_) { $('audio-status').textContent='Progress saved for this session only; browser storage is unavailable.'; }
  renderNav();
});
$('tutor').addEventListener('click', async () => {
  const prompt = tutorInstructions + lessonText(sections[selected]);
  try { await navigator.clipboard.writeText(prompt); $('audio-status').textContent='Copied. Paste into a voice-capable AI conversation to begin.'; }
  catch (_) { $('copy-fallback').hidden=false; $('prompt-text').value=prompt; $('prompt-text').focus(); $('prompt-text').select(); $('audio-status').textContent='Automatic copy is unavailable. Copy the selected prompt manually.'; }
});
$('export').addEventListener('click', () => {
  const url=URL.createObjectURL(new Blob([lessonText(sections[selected])],{type:'text/plain;charset=utf-8'}));
  const a=document.createElement('a'); a.href=url; a.download=sections[selected].id+'-lesson.txt'; a.click(); setTimeout(()=>URL.revokeObjectURL(url),1000);
});
let beforePrint = [];
window.addEventListener('beforeprint', () => { beforePrint=Array.from(document.querySelectorAll('details')).map(el=>[el,el.open]); beforePrint.forEach(([el])=>el.open=true); });
window.addEventListener('afterprint', () => { beforePrint.forEach(([el,open])=>el.open=open); });
$('print').addEventListener('click', () => window.print());
window.addEventListener('hashchange', () => { const i=sections.findIndex(s=>'#'+s.id===location.hash); if(i>=0 && i!==selected) show(i,true); });
window.addEventListener('pagehide', () => { if(synth) synth.cancel(); });
if (synth) { loadVoices(); synth.addEventListener('voiceschanged',loadVoices); }
else { ['listen','pause','continue','stop'].forEach(id=>$(id).disabled=true); }
const initial = sections.findIndex(s => '#'+s.id === location.hash);
show(initial>=0 ? initial : 0);
