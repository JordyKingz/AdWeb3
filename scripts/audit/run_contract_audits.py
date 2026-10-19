#!/usr/bin/env python3
"""Invoke the actual skill workflows on immutable, isolated snapshots."""
import concurrent.futures,datetime,hashlib,json,pathlib,shutil,subprocess,sys,tarfile
root=pathlib.Path.cwd();phase=int(sys.argv[1]);choice=sys.argv[2] if len(sys.argv)>2 else 'both'
sha=subprocess.check_output(['git','rev-parse','HEAD'],text=True).strip()
if subprocess.check_output(['git','status','--porcelain','--untracked-files=normal'],text=True).strip():raise SystemExit('Commit the candidate before auditing; dirty source is not a reviewed SHA.')
lock=json.loads((root/'audit/tools-lock.json').read_text());stamp=datetime.datetime.now(datetime.timezone.utc).strftime('%Y%m%dT%H%M%SZ');dest=root/'audit/reports'/sha/stamp;dest.mkdir(parents=True)
manifest={'source_sha':sha,'source_clean':True,'source_unchanged':False,'phase':phase,'tools':{},'review_status':'pending','findings':[]}
archive=root/'audit/work'/f'{sha}.tar';archive.parent.mkdir(parents=True,exist_ok=True)
subprocess.run(['git','archive','--format=tar','--output',str(archive),sha],check=True)
def digest(path):
 return {str(p.relative_to(path)):hashlib.sha256(p.read_bytes()).hexdigest() for p in (path/'contracts/src').rglob('*.sol')}
def run(name,entry,key):
 work=root/'audit/work'/stamp/name;work.mkdir(parents=True)
 with tarfile.open(archive) as a:a.extractall(work,filter='data')
 source=pathlib.Path(lock[key]['path']);target=work/'.claude/skills';target.mkdir(parents=True,exist_ok=True)
 if name=='nemesis':
  for n in ['nemesis-auditor','feynman-auditor','state-inconsistency-auditor']:shutil.copytree(source/'.claude/skills'/n,target/n,dirs_exist_ok=True)
 else:shutil.copytree(source/'solidity-auditor',target/'solidity-auditor',dirs_exist_ok=True)
 before=digest(work)
 prompt=f'''Execute the actual skill at {entry} for the AdWeb3 Phase {phase} internal contract audit. Read its full instructions and required references. Target production contracts/src and deployment scripts, with interfaces, tests, relevant Go chain/ledger/campaign code and docs/features/adweb3-platform/feature-spec.md as context. This is a testnet escrow/distributor, not Capyfi. Do not implement fixes or change production source. You may write audit reports, metadata and isolated evidence harnesses outside production source. Use the current runtime model; do not ask a model/pass-count question. {'Run the complete Nemesis workflow: full Feynman and full state passes, alternate until convergence or max six passes, then verified consolidated report.' if name=='nemesis' else 'Explicit skill flags: --loop 3 --file-output. Run all twelve specialties each pass, follow skill report assembly, preserve the native full-report and memory. Report failed agents or passes accurately.'} Do not deploy or send transactions. Preserve raw native skill output files. Finish with paths, completion/coverage limitations, and findings summary. Existing unrelated .claude skills are background, not a direction to implement product code.'''
 (dest/f'{name}-prompt.txt').write_text(prompt)
 args=['claude','-p',prompt,'--permission-mode','acceptEdits','--allowedTools','Read,Write,Edit,Glob,Grep,Bash,Agent','--output-format','json']
 with (dest/f'{name}-runtime.json').open('w') as out,(dest/f'{name}-stderr.log').open('w') as err:result=subprocess.run(args,cwd=work,stdout=out,stderr=err)
 unchanged=before==digest(work)
 native=work/'.audit/findings/nemesis-verified.md' if name=='nemesis' else next(iter(sorted((work/'.solidity-auditor/runs').glob('*/full-report.md'))),None)
 data={'sha':lock[key]['sha'],'updated':lock[key]['updated'],'complete':False,'runtime_exit':result.returncode,'source_unchanged':unchanged,'report':name+'.md','workspace':str(work)}
 if native and native.is_file():
  shutil.copy2(native,dest/(name+'.md'))
  evidence=work/('.audit' if name=='nemesis' else '.solidity-auditor');shutil.copytree(evidence,dest/(name+'-native'),dirs_exist_ok=True)
  data['report_present']=True
 # Human/independent reviewer must validate native pass coverage before complete=true.
 return name,data
names=[('nemesis','.claude/skills/nemesis-auditor/SKILL.md','nemesis-auditor'),('pashov','.claude/skills/solidity-auditor/SKILL.md','pashov-skills')]
selected=[x for x in names if choice in ('both',x[0])]
with concurrent.futures.ThreadPoolExecutor(max_workers=2) as ex:
 for name,data in ex.map(lambda x:run(*x),selected):manifest['tools'][name]=data;print(name,'finished, review required',flush=True)
manifest['source_unchanged']=all(v.get('source_unchanged') for v in manifest['tools'].values())
(dest/'manifest.json').write_text(json.dumps(manifest,indent=2)+'\n');print(dest,flush=True)
if not manifest['source_unchanged'] or any(x.get('runtime_exit')!=0 or not x.get('report_present') for x in manifest['tools'].values()):sys.exit(1)
