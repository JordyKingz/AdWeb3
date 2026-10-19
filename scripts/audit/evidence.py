"""Shared source binding for audit orchestration and verification."""
import hashlib,pathlib,subprocess
AREAS=('contracts/src','contracts/script','internal','cmd','migrations','web','api')
def scoped_files(root):
 files={}
 for area in AREAS:
  p=root/area
  if p.exists():
   for f in p.rglob('*'):
    if f.is_file():files[str(f.relative_to(root))]=hashlib.sha256(f.read_bytes()).hexdigest()
 for n in ['go.mod','go.sum','contracts/foundry.toml']:
  f=root/n
  if f.is_file():files[n]=hashlib.sha256(f.read_bytes()).hexdigest()
 return files

def verify_source(root,manifest):
 files=manifest['files'];assert files,'missing audited source manifest'
 assert files==scoped_files(root),'current audited source differs'
 for path,digest in files.items():
  data=subprocess.check_output(['git','-C',str(root),'show',manifest['source_sha']+':'+path],stderr=subprocess.DEVNULL)
  assert hashlib.sha256(data).hexdigest()==digest,'source revision/content mismatch: '+path
