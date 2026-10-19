#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/../.."
python3 - <<'PY'
import subprocess,json,pathlib,datetime
base=pathlib.Path('/Users/homer/development/audits/capyfi-smart-contracts/audit/tools')
expected={'nemesis-auditor':'https://github.com/0xiehnnkta/nemesis-auditor.git','pashov-skills':'https://github.com/pashov/skills.git'}
record={}
for name,url in expected.items():
 p=base/name
 def g(*args):return subprocess.check_output(['git','-C',str(p),*args],text=True).strip()
 if g('status','--porcelain'):raise SystemExit(name+': local modifications preserved; update stopped')
 if g('remote','get-url','origin')!=url:raise SystemExit(name+': unexpected upstream')
 before=g('rev-parse','HEAD');g('fetch','origin');g('merge','--ff-only','@{upstream}')
 record[name]={'path':str(p),'origin':url,'before':before,'sha':g('rev-parse','HEAD'),'updated':True}
record['updated_at']=datetime.datetime.now(datetime.timezone.utc).isoformat()
p=pathlib.Path('audit/tools-lock.json');p.parent.mkdir(exist_ok=True);p.write_text(json.dumps(record,indent=2)+'\n')
PY
