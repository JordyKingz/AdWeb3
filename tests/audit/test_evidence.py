import json, pathlib, subprocess, tempfile, unittest
VERIFY=pathlib.Path(__file__).resolve().parents[2]/'scripts/audit/verify-reports.py'
class EvidenceTests(unittest.TestCase):
 def run_verifier(self, value):
  with tempfile.TemporaryDirectory() as tmp:
   p=pathlib.Path(tmp);(p/'nemesis.md').write_text('original verified report');(p/'pashov.md').write_text('original report');(p/'manifest.json').write_text(json.dumps(value))
   return subprocess.run(['python3',str(VERIFY),str(p/'manifest.json')],capture_output=True).returncode
 def manifest(self):return dict(source_sha='a'*40,source_clean=True,source_unchanged=True,phase=1,tools={k:dict(sha='b'*40,updated=True,complete=True,report=k+'.md') for k in ['nemesis','pashov']},review_status='closed',findings=[])
 def test_complete_evidence_accepted(self):self.assertEqual(self.run_verifier(self.manifest()),0)
 def test_failed_tool_does_not_pass(self):
  m=self.manifest();m['tools']['pashov']['complete']=False;self.assertNotEqual(self.run_verifier(m),0)
 def test_unresolved_findings_cannot_pass(self):
  m=self.manifest();m['findings']=[dict(id='F1',classification='MUST-FIX',status='open')];self.assertNotEqual(self.run_verifier(m),0)
 def test_unreviewed_scan_cannot_pass(self):
  m=self.manifest();m['review_status']='pending';self.assertNotEqual(self.run_verifier(m),0)
if __name__=='__main__':unittest.main()
