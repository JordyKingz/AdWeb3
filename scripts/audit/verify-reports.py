#!/usr/bin/env python3
"""Check evidence completeness; an independent reviewer still judges the findings."""
import json, pathlib, re, sys
def verify(path):
 root=path.parent; m=json.loads(path.read_text())
 assert re.fullmatch(r'[0-9a-f]{40}',m['source_sha']), 'invalid source SHA'
 assert m.get('source_clean') and m.get('source_unchanged'), 'unbound/modified source'
 assert m.get('review_status')=='closed', 'review not closed'
 for name in ['nemesis','pashov']:
  t=m['tools'][name];assert t.get('updated') and t.get('complete'),name+' incomplete'
  assert re.fullmatch(r'[0-9a-f]{40}',t['sha']),name+' unpinned'
  report=(root/t['report']).resolve();assert report.is_relative_to(root.resolve()),'report escapes audit evidence'
  assert report.is_file() and report.stat().st_size>0,name+' missing report'
 for f in m['findings']:
  assert f.get('status') in ('fixed','resolved','note'),f.get('id','unknown')+' unresolved'
  if f.get('classification') in ('MUST-FIX','SHOULD-FIX'):
   assert f.get('regression') and f.get('retest') and f.get('reviewer_resolution'),f['id']+' missing closure evidence'
if __name__=='__main__':
 try:verify(pathlib.Path(sys.argv[1]));print('Audit evidence complete; see independent review disposition.')
 except (AssertionError,KeyError,ValueError,OSError,IndexError) as e:print(str(e),file=sys.stderr);sys.exit(1)
