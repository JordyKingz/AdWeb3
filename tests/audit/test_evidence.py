import hashlib,json,pathlib,subprocess,tempfile,unittest
VERIFY=pathlib.Path(__file__).resolve().parents[2]/'scripts/audit/verify-reports.py'
class EvidenceTests(unittest.TestCase):
 def setUp(self):
  self.tmp=tempfile.TemporaryDirectory();self.addCleanup(self.tmp.cleanup);self.root=pathlib.Path(self.tmp.name)
  self.source=self.root/'source';self.source.mkdir();(self.source/'contracts/src').mkdir(parents=True);(self.source/'contracts/src/A.sol').write_text('contract A {}')
  (self.source/'audit').mkdir();(self.source/'audit/tools-lock.json').write_text(json.dumps({k:{'sha':'b'*40} for k in ['nemesis-auditor','pashov-skills']}))
  subprocess.run(['git','init','-q',str(self.source)],check=True);subprocess.run(['git','-C',str(self.source),'add','.'],check=True)
  subprocess.run(['git','-C',str(self.source),'-c','user.name=Test','-c','user.email=test@example.invalid','commit','-qm','fixture'],check=True)
  self.sha=subprocess.check_output(['git','-C',str(self.source),'rev-parse','HEAD'],text=True).strip()
 def manifest(self):return dict(source_sha=self.sha,files={'contracts/src/A.sol':hashlib.sha256(b'contract A {}').hexdigest()},source_clean=True,source_unchanged=True,phase=1,tools={k:dict(sha='b'*40,updated=True,complete=True,report=k+'.md',runtime_exit=0,source_unchanged=True) for k in ['nemesis','pashov']},review_status='closed',findings=[])
 def run_verifier(self,m):
  for k in ['nemesis','pashov']:(self.root/(k+'.md')).write_text('original report')
  p=self.root/'manifest.json';p.write_text(json.dumps(m));return subprocess.run(['python3',str(VERIFY),str(p),str(self.source)],capture_output=True).returncode
 def test_complete_evidence_accepted(self):self.assertEqual(self.run_verifier(self.manifest()),0)
 def test_incomplete_tool_fails(self):
  m=self.manifest();m['tools']['pashov']['complete']=False;self.assertNotEqual(self.run_verifier(m),0)
 def test_failed_runtime_fails_even_if_complete(self):
  m=self.manifest();m['tools']['pashov']['runtime_exit']=1;self.assertNotEqual(self.run_verifier(m),0)
 def test_tool_source_mutation_fails(self):
  m=self.manifest();m['tools']['nemesis']['source_unchanged']=False;self.assertNotEqual(self.run_verifier(m),0)
 def test_stale_source_fails(self):
  m=self.manifest();(self.source/'contracts/src/A.sol').write_text('contract Changed {}');self.assertNotEqual(self.run_verifier(m),0)
 def test_unbound_source_revision_fails(self):
  m=self.manifest();m['source_sha']='c'*40;self.assertNotEqual(self.run_verifier(m),0)
 def test_tool_lock_mismatch_fails(self):
  m=self.manifest();m['tools']['pashov']['sha']='c'*40;self.assertNotEqual(self.run_verifier(m),0)
 def test_unresolved_findings_cannot_pass(self):
  m=self.manifest();m['findings']=[dict(id='F1',classification='MUST-FIX',status='open')];self.assertNotEqual(self.run_verifier(m),0)
 def test_unreviewed_scan_cannot_pass(self):
  m=self.manifest();m['review_status']='pending';self.assertNotEqual(self.run_verifier(m),0)
if __name__=='__main__':unittest.main()
