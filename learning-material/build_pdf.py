"""Export the reading edition to PDF and plain text.

Maintainer command: python learning-material/build_pdf.py
No browser reader, remote service, or HTML source files are involved.
"""
from pathlib import Path
import argparse
import html
import json
import re
import math
from sync_reading_code import sync as sync_reading_code

from reportlab.lib import colors
from reportlab.lib.styles import ParagraphStyle
from reportlab.pdfbase import pdfmetrics
from reportlab.pdfbase.ttfonts import TTFont
from reportlab.platypus import (
    SimpleDocTemplate, Paragraph, Spacer, PageBreak, Preformatted,
    LongTable, TableStyle, KeepTogether, Flowable,
)
from reportlab.platypus.tableofcontents import TableOfContents
from pypdf import PdfReader

ROOT = Path(__file__).resolve().parent
REPO = ROOT.parent
PAGE = (432, 648)  # 6 x 9 inch book format, readable on a small screen.
MARGIN = 39
WIDTH = PAGE[0] - 2 * MARGIN
FIGURES = json.loads((ROOT/'reading'/'figures.json').read_text(encoding='utf-8'))


def register_fonts():
    windows = Path('C:/Windows/Fonts')
    linux = Path('/usr/share/fonts/truetype/dejavu')
    if all((windows / name).exists() for name in ['georgia.ttf', 'georgiab.ttf', 'georgiai.ttf', 'georgiaz.ttf']):
        names = ['georgia.ttf', 'georgiab.ttf', 'georgiai.ttf', 'georgiaz.ttf']
        font_dir = windows
        mono = windows / 'consola.ttf'
    elif all((linux / name).exists() for name in ['DejaVuSerif.ttf', 'DejaVuSerif-Bold.ttf', 'DejaVuSerif-Italic.ttf', 'DejaVuSerif-BoldItalic.ttf']):
        names = ['DejaVuSerif.ttf', 'DejaVuSerif-Bold.ttf', 'DejaVuSerif-Italic.ttf', 'DejaVuSerif-BoldItalic.ttf']
        font_dir = linux
        mono = linux / 'DejaVuSansMono.ttf'
    else:
        import reportlab
        font_dir = Path(reportlab.__file__).parent / 'fonts'
        names = ['Vera.ttf', 'VeraBd.ttf', 'VeraIt.ttf', 'VeraBI.ttf']
        mono = linux / 'DejaVuSansMono.ttf'
    for suffix, name in zip(['', '-Bold', '-Italic', '-BoldItalic'], names):
        pdfmetrics.registerFont(TTFont('Book' + suffix, str(font_dir / name)))
    pdfmetrics.registerFontFamily('Book', normal='Book', bold='Book-Bold', italic='Book-Italic', boldItalic='Book-BoldItalic')
    if mono and mono.exists():
        pdfmetrics.registerFont(TTFont('Code', str(mono)))
        return 'Code'
    return 'Courier'


MONO = register_fonts()
INK = colors.HexColor('#202020')
styles = {
    'body': ParagraphStyle('Body', fontName='Book', fontSize=10.7, leading=15.2,
                           spaceAfter=8, textColor=INK, allowWidows=0, allowOrphans=0),
    'h1': ParagraphStyle('Chapter', fontName='Book-Bold', fontSize=21, leading=26,
                         spaceBefore=8, spaceAfter=19, keepWithNext=True),
    'h2': ParagraphStyle('Section', fontName='Book-Bold', fontSize=14, leading=18,
                         spaceBefore=16, spaceAfter=8, keepWithNext=True),
    'h3': ParagraphStyle('Subsection', fontName='Book-Bold', fontSize=11.5, leading=16,
                         spaceBefore=12, spaceAfter=6, keepWithNext=True),
    'answer': ParagraphStyle('Answer', fontName='Book-Italic', fontSize=9.3, leading=13,
                             spaceBefore=5, spaceAfter=5, keepWithNext=True, textColor=colors.HexColor('#454545')),
    'question': ParagraphStyle('Question', fontName='Book', fontSize=10.7, leading=15.2,
                               spaceBefore=16, spaceAfter=16, borderPadding=8, keepWithNext=True,
                               borderWidth=0.5, borderColor=colors.HexColor('#b0b0b0'),
                               backColor=colors.HexColor('#f4f4f4'), allowWidows=0, allowOrphans=0),
    'code': ParagraphStyle('Code', fontName=MONO, fontSize=9, leading=12,
                           spaceBefore=9, spaceAfter=12, leftIndent=7),
    'cell': ParagraphStyle('Cell', fontName='Book', fontSize=8.2, leading=11.2),
    'caption': ParagraphStyle('Caption', fontName='Book-Italic', fontSize=8.5, leading=11.5,
                              spaceBefore=7, spaceAfter=14),
    'toc': ParagraphStyle('ContentsEntry', fontName='Book', fontSize=10.5, leading=15,
                          spaceBefore=8, leftIndent=0, firstLineIndent=0, rightIndent=18),
    'cover': ParagraphStyle('Cover', fontName='Book-Bold', fontSize=28, leading=34, spaceAfter=25),
    'subtitle': ParagraphStyle('Subtitle', fontName='Book', fontSize=15, leading=22, spaceAfter=20),
}


class BookFigure(Flowable):
    """Small monochrome vector diagrams; coordinates are in page points."""
    def __init__(self, key):
        super().__init__()
        self.spec = FIGURES[key]
        self.width = WIDTH
        self.height = self.spec['height'] + 35

    def draw(self):
        c = self.canv
        c.setFillColor(colors.HexColor('#182b42'))
        c.roundRect(0,self.spec['height']+10,WIDTH,25,5,fill=1,stroke=0)
        c.setFillColor(colors.white)
        c.setFont('Book-Bold',8.5)
        assert pdfmetrics.stringWidth(self.spec['title'],'Book-Bold',8.5) < WIDTH-14
        c.drawString(7,self.spec['height']+18,self.spec['title'])
        nodes = {n['id']: n for n in self.spec['nodes']}
        c.setStrokeColor(colors.HexColor('#555555'))
        c.setFillColor(colors.HexColor('#555555'))
        c.setLineWidth(.8)
        for start, end in self.spec['edges']:
            a, b = nodes[start], nodes[end]
            ax, ay = a['x']+a['w']/2, a['y']+a['h']/2
            bx, by = b['x']+b['w']/2, b['y']+b['h']/2
            dx, dy = bx-ax, by-ay
            def border(node, cx, cy, dx, dy):
                tx = node['w']/2/abs(dx) if dx else float('inf')
                ty = node['h']/2/abs(dy) if dy else float('inf')
                t = min(tx,ty)
                return cx+t*dx, cy+t*dy
            x1,y1 = border(a,ax,ay,dx,dy)
            x2,y2 = border(b,bx,by,-dx,-dy)
            c.line(x1,y1,x2,y2)
            angle = math.atan2(y2-y1,x2-x1)
            for offset in [-.45,.45]:
                c.line(x2,y2,x2-5*math.cos(angle+offset),y2-5*math.sin(angle+offset))
        for n in nodes.values():
            c.setFillColor(colors.HexColor('#dde5ed') if n.get('shade') else colors.HexColor('#f8f9fb'))
            c.roundRect(n['x'],n['y'],n['w'],n['h'],5,fill=1,stroke=1)
            c.setFillColor(INK)
            c.setFont('Book',9)
            lines = n['text'].split('\n')
            for i,line in enumerate(lines):
                assert pdfmetrics.stringWidth(line,'Book',9) <= n['w']-8, line
                c.drawCentredString(n['x']+n['w']/2,n['y']+n['h']/2+(len(lines)-1)*6-i*12-3,line)
        c.setFont('Book',8)
        for x,y,label in self.spec['labels']:
            assert x+pdfmetrics.stringWidth(label,'Book',8) <= WIDTH, label
            c.setFillColor(colors.white)
            c.rect(x-2,y-2,pdfmetrics.stringWidth(label,'Book',8)+4,11,fill=1,stroke=0)
            c.setFillColor(INK)
            c.drawString(x,y,label)


class CodeBlock(Flowable):
    """Monospace code with measured widths and labelled page continuations."""
    line_height = 12
    padding = 30

    def __init__(self, lines, continued=False, allow_split=False):
        super().__init__()
        self.lines = lines
        self.continued = continued
        self.allow_split = allow_split
        self.width = WIDTH
        self.height = len(lines) * self.line_height + self.padding
        self.spaceBefore = 9
        self.spaceAfter = 12

    def split(self, availWidth, availHeight):
        # Move a short complete example intact to the next page when possible.
        if self.height <= 550 and not self.allow_split:
            return []
        count = int((availHeight - self.padding) // self.line_height)
        if count < 6:
            return []
        return [CodeBlock(self.lines[:count], self.continued, True),
                CodeBlock(self.lines[count:], True, True)]

    def draw(self):
        c = self.canv
        c.setFillColor(colors.HexColor('#f3f5f7'))
        c.roundRect(0, 0, WIDTH, self.height, 4, fill=1, stroke=0)
        c.setStrokeColor(colors.HexColor('#182b42'))
        c.setLineWidth(2)
        c.line(0, 0, 0, self.height)
        c.setFillColor(colors.HexColor('#526477'))
        c.setFont('Book-Bold', 7)
        c.drawString(8, self.height-11,
                     'GO (continued)' if self.continued else 'GO')
        c.setFont(MONO, 9)
        c.setFillColor(INK)
        for i, line in enumerate(self.lines):
            c.drawString(8, self.height-25-i*self.line_height, line)


class SvgCanvas:
    """Render the same diagram primitives as SVG for Markdown/GitHub readers."""
    def __init__(self, height):
        self.height = height
        self.parts = []
        self.stroke = '#555555'
        self.fill = '#202020'
        self.line_width = .8
        self.font_size = 8.2

    def setStrokeColor(self, value): self.stroke = value.hexval().replace('0x','#')
    def setFillColor(self, value): self.fill = value.hexval().replace('0x','#')
    def setLineWidth(self, value): self.line_width = value
    def setFont(self, name, size):
        self.font_size = size
        self.font_weight = 'bold' if 'Bold' in name else 'normal'

    def line(self, x1,y1,x2,y2):
        self.parts.append(f'<line x1="{x1}" y1="{self.height-y1}" x2="{x2}" y2="{self.height-y2}" stroke="{self.stroke}" stroke-width="{self.line_width}"/>')

    def rect(self,x,y,w,h,fill=0,stroke=1):
        self.parts.append(f'<rect x="{x}" y="{self.height-y-h}" width="{w}" height="{h}" fill="{self.fill if fill else "none"}" stroke="{self.stroke if stroke else "none"}" stroke-width="{self.line_width}"/>')

    def roundRect(self,x,y,w,h,r,fill=0,stroke=1):
        self.parts.append(f'<rect x="{x}" y="{self.height-y-h}" width="{w}" height="{h}" rx="{r}" fill="{self.fill if fill else "none"}" stroke="{self.stroke if stroke else "none"}" stroke-width="{self.line_width}"/>')

    def drawString(self,x,y,text): self._text(x,y,text,'start')
    def drawCentredString(self,x,y,text): self._text(x,y,text,'middle')
    def _text(self,x,y,text,anchor):
        self.parts.append(f'<text x="{x}" y="{self.height-y}" text-anchor="{anchor}" fill="{self.fill}" font-family="Georgia, serif" font-weight="{self.font_weight}" font-size="{self.font_size}">{html.escape(text)}</text>')


def export_svg_figures():
    target = ROOT/'reading'/'figures'
    target.mkdir(exist_ok=True)
    for key in FIGURES:
        figure = BookFigure(key)
        canvas = SvgCanvas(figure.height)
        figure.canv = canvas
        figure.draw()
        svg = f'<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 {WIDTH} {figure.height}" role="img"><title>{html.escape(key)}</title><rect width="100%" height="100%" fill="white"/>' + ''.join(canvas.parts) + '</svg>\n'
        (target/(key+'.svg')).write_text(svg,encoding='utf-8')


def normalize(text):
    return re.sub('[\u2010-\u2015\u2212]', '-', text).replace('\u00a0', ' ')


def inline(text, code_size=9):
    result = []
    for token in re.split(r'(\[[^\]]+\]\([^)]+\)|`[^`]+`|\*\*.+?\*\*)', normalize(text)):
        if token.startswith('[') and '](' in token:
            match = re.fullmatch(r'\[([^\]]+)\]\(([^)]+)\)', token)
            label, url = match.groups()
            if url.startswith(('https://','http://')):
                result.append(f'<link href="{html.escape(url, quote=True)}" color="#333333"><u>{html.escape(label)}</u></link>')
            else:
                result.append(html.escape(label))
        elif token.startswith('`') and token.endswith('`'):
            result.append(f'<font name="{MONO}" size="{code_size}">{html.escape(token[1:-1])}</font>')
        elif token.startswith('**') and token.endswith('**'):
            result.append('<b>' + html.escape(token[2:-2]) + '</b>')
        else:
            result.append(html.escape(token))
    return ''.join(result)


class BookDoc(SimpleDocTemplate):
    def afterFlowable(self, flowable):
        if not isinstance(flowable, Paragraph) or not hasattr(flowable, 'bookmark_key'):
            return
        title = flowable.getPlainText()
        self.canv.bookmarkPage(flowable.bookmark_key)
        self.canv.addOutlineEntry(title, flowable.bookmark_key, flowable.outline_level, False)
        if flowable.outline_level == 0:
            self.notify('TOCEntry', (0, title, self.page, flowable.bookmark_key))


def page_decor(canvas, doc):
    if doc.page == 1:
        return
    canvas.saveState()
    canvas.setFont('Book', 7.4)
    canvas.setFillColor(colors.HexColor('#555555'))
    canvas.drawString(MARGIN, PAGE[1] - 24, 'NETFLIX L5  |  GO INTERVIEW PREPARATION')
    canvas.setStrokeColor(colors.HexColor('#bbbbbb'))
    canvas.setLineWidth(0.35)
    canvas.line(MARGIN, PAGE[1] - 30, PAGE[0] - MARGIN, PAGE[1] - 30)
    canvas.drawCentredString(PAGE[0] / 2, 20, str(doc.page))
    canvas.restoreState()


def markdown_flowables(text, chapter_number):
    lines = text.splitlines()
    out = []
    i = 0
    heading_count = 0
    while i < len(lines):
        line = lines[i].strip()
        if not line:
            i += 1
            continue
        figure = re.fullmatch(r'!\[(.*)\]\(([^)]+)\)', line)
        if figure:
            caption, image_path = figure.groups()
            key = Path(image_path).stem
            out.append(KeepTogether([BookFigure(key), Paragraph(inline(caption), styles['caption'])]))
            i += 1
            continue
        if line.startswith('```'):
            code = []
            i += 1
            while i < len(lines) and not lines[i].startswith('```'):
                # Use compact indentation and preserve gofmt's field alignment.
                code_line = normalize(lines[i])
                indent = len(code_line) - len(code_line.lstrip('\t'))
                code_line = ' ' * (4 * indent) + code_line[indent:].expandtabs(8)
                assert pdfmetrics.stringWidth(code_line,MONO,9) <= WIDTH-14, code_line
                code.append(code_line)
                i += 1
            out.append(CodeBlock(code))
            i += 1
            continue
        if line.startswith('|'):
            rows = []
            while i < len(lines) and lines[i].strip().startswith('|'):
                raw = lines[i].strip().strip('|')
                cells = [x.strip() for x in raw.split('|')]
                if not all(re.fullmatch(r'[-: ]+', x) for x in cells):
                    rows.append([Paragraph(inline(x, code_size=7.6), styles['cell']) for x in cells])
                i += 1
            table = LongTable(rows, colWidths=[WIDTH / len(rows[0])] * len(rows[0]), repeatRows=1, hAlign='LEFT')
            table.setStyle(TableStyle([
                ('VALIGN', (0,0), (-1,-1), 'TOP'),
                ('BACKGROUND', (0,0), (-1,0), colors.HexColor('#eeeeee')),
                ('LINEBELOW', (0,0), (-1,-1), .35, colors.HexColor('#bbbbbb')),
                ('TOPPADDING', (0,0), (-1,-1), 7), ('BOTTOMPADDING', (0,0), (-1,-1), 7),
                ('LEFTPADDING', (0,0), (-1,-1), 6), ('RIGHTPADDING', (0,0), (-1,-1), 6),
            ]))
            if len(rows) <= 8:
                out.append(KeepTogether([table, Spacer(1,12)]))
            else:
                out.extend([table, Spacer(1,12)])
            continue
        h = re.match(r'^(#{1,4}) (.*)', line)
        if h:
            level, title = len(h.group(1)), h.group(2)
            style = styles['answer' if level == 4 else 'h'+str(level)]
            para = Paragraph(inline(title), style)
            if level <= 2:
                para.bookmark_key = f'chapter-{chapter_number}-heading-{heading_count}'
                para.outline_level = level-1
                heading_count += 1
            out.append(para)
            i += 1
            continue
        if line.startswith('- '):
            parts = [line[2:]]
            i += 1
            while i < len(lines) and lines[i].strip() and not lines[i].startswith(('- ','#','|','```')):
                parts.append(lines[i].strip())
                i += 1
            out.append(Paragraph(inline(' '.join(parts)), styles['body'], bulletText='-'))
            continue
        parts = [line]
        i += 1
        while i < len(lines) and lines[i].strip() and not lines[i].startswith(('#','|','```','- ','![')):
            parts.append(lines[i].strip())
            i += 1
        body = ' '.join(parts)
        is_question = body.startswith(('**Round ', '**Think aloud.'))
        out.append(Paragraph(inline(body), styles['question' if is_question else 'body']))
    return out


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--output', type=Path, default=REPO/'output'/'pdf'/'netflix-go-interview-book.pdf')
    args = parser.parse_args()
    code_blocks = sync_reading_code()
    export_svg_figures()
    index = json.loads((ROOT/'reading-index.json').read_text(encoding='utf-8'))
    files = [ROOT/'BOOK-INTRO.md'] + [ROOT/item['path'] for item in index] + [ROOT/'SOURCES.md']
    chapters = [p.read_text(encoding='utf-8') for p in files]
    forbidden = ['Practice for ', 'set a timer', 'Copy AI tutor prompt', 'course.html', 'build.py',
                 '**Round ', '#### Worked answer', '**Think aloud.**', '**Optional spoken walkthrough:**']
    for path, text in zip(files, chapters):
        for phrase in forbidden:
            assert phrase not in text, f'Obsolete instruction {phrase!r} in {path}'
    assert len(index) == 12
    for path, text in zip(files[1:13], chapters[1:13]):
        assert text.count('## Worked example:') >= 2, path
        assert '## Putting the program together:' in text, path
        assert '## Go example:' in text, path
        assert text.count('```go') >= 3, path
        assert '![' in text, path
    args.output.parent.mkdir(parents=True, exist_ok=True)
    doc = BookDoc(str(args.output), pagesize=PAGE, rightMargin=MARGIN, leftMargin=MARGIN,
                  topMargin=43, bottomMargin=37, title='Learning to Solve Coding Problems',
                  author='Netflix L5 Go Interview Preparation',
                  subject='Self-contained interview learning away from a computer')
    story = [Spacer(1,80), Paragraph('Learning to solve<br/>coding problems', styles['cover']),
             Paragraph('Netflix L5 interview preparation in Go', styles['subtitle']),
             Paragraph('The reading edition<br/>Worked examples and visual explanations', styles['body']),
             Spacer(1,28), Paragraph('Data structures. Algorithms. Caches. Rate limiting.<br/>Recommendations. Worked program designs.', styles['body']),
             Spacer(1,28), Paragraph('Learn how state changes, why an approach works, and when a different structure is needed. No computer required.', styles['body']),
             PageBreak(), Paragraph('Contents', styles['h1'])]
    toc = TableOfContents()
    toc.levelStyles = [styles['toc']]
    story.append(toc)
    for n, chapter in enumerate(chapters):
        story.append(PageBreak())
        flows = markdown_flowables(chapter,n)
        # Avoid leaving a chapter's final two or three lines alone on a new page.
        if len(flows) >= 2 and all(isinstance(x, Paragraph) for x in flows[-2:]):
            flows = flows[:-2] + [KeepTogether(flows[-2:])]
        story.extend(flows)
    doc.multiBuild(story, onFirstPage=page_decor, onLaterPages=page_decor)
    whole = '\n\n'.join(chapters)
    plain = re.sub(r'!\[(.*)\]\(([^)]+)\)', r'Figure: \1', whole)
    plain = re.sub(r'\[([^\]]+)\]\(([^)]+)\)', r'\1 (\2)', plain)
    plain = re.sub(r'^#{1,4} ', '', plain, flags=re.M).replace('**','').replace('```go','').replace('```','').replace('`','')
    (ROOT/'book.txt').write_text(plain,encoding='utf-8')
    reader = PdfReader(str(args.output))
    extracted = '\n'.join(p.extract_text() or '' for p in reader.pages)
    assert 'PLACEHOLDER' not in extracted.upper()
    assert '](figures/' not in extracted, 'Unrendered Markdown image'
    for key, spec in FIGURES.items():
        assert spec['title'] in extracted, f'Figure not rendered: {key}'
    assert len(reader.pages) > 30
    for n in range(1,13):
        assert f'Chapter {n}' in extracted, f'Missing chapter {n}'
    words = len(re.findall(r'\b[\w]+(?:[\x27-][\w]+)*\b',whole))
    report = {'edition':'reading','chapters':12,'worked_examples':whole.count('## Worked example:'),
              'vector_figures':len(FIGURES),'go_example_sections':whole.count('## Go example:'),
              'go_code_blocks':code_blocks,'state_tables':sum(1 for line in whole.splitlines() if line.startswith('| ---')),
              'problem_briefs':whole.count('**Problem.**')+whole.count('**Brief.**'),
              'conversation_questions_in_separate_sessions':132,
              'worked_program_designs':12,'complete_case_studies':2,
              'source_words_approx':words,'pdf_pages':len(reader.pages),'pdf_bytes':args.output.stat().st_size,
              'source_files':[str(p.relative_to(REPO)).replace('\\','/') for p in files]}
    (ROOT/'book-manifest.json').write_text(json.dumps(report,indent=2)+'\n',encoding='utf-8')
    print(json.dumps(report,indent=2))


if __name__ == '__main__':
    main()
