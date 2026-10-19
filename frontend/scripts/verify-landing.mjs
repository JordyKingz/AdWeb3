import assert from 'node:assert/strict'
import { existsSync, readdirSync, readFileSync } from 'node:fs'
import { fileURLToPath } from 'node:url'
import { join, relative } from 'node:path'

const root = fileURLToPath(new URL('../', import.meta.url))
const output = join(root, '.output/public')
function filesIn(directory) {
  return readdirSync(directory, { withFileTypes: true }).flatMap(entry => {
    const path = join(directory, entry.name)
    assert(!entry.isSymbolicLink(), `Unexpected output symlink: ${path}`)
    return entry.isDirectory() ? filesIn(path) : [path]
  })
}

assert(existsSync(join(output, 'index.html')), 'Generate the landing before verification')
const files = filesIn(output).map(path => relative(output, path))
const allowedHtml = new Set(['index.html', '200.html', '404.html'])
assert(files.filter(path => path.endsWith('.html')).every(path => allowedHtml.has(path)), 'Unexpected exported page')
assert(!files.some(path => /^(api|connect|dashboard|brands)(\/|\.|$)/.test(path)), 'Private app routes must not be exported')
assert(!files.some(path => /(^|\/)\.env(?:\.|$)/.test(path)), 'Environment files must not be exported')
assert(!existsSync(join(root, '.output/server')), 'Landing deployment must not include a server bundle')
assert(!existsSync(join(root, '.vercel/output/functions')), 'Landing deployment must not include server functions')

const html = readFileSync(join(output, 'index.html'), 'utf8')
assert(html.includes('Live campaigns and cash payouts are not available yet.'), 'Development notice is missing')
assert(html.includes('landingOnly:true'), 'Landing feature flag is disabled')
assert(!/href="\/(?:connect|dashboard|brands)(?:\/|"|\?)/.test(html), 'Landing links must not open unfinished app flows')
assert(html.includes('id="pilot"'), 'Pilot information target is missing')
assert(!html.includes('mailto:hello@example.com'), 'Placeholder contact address must not be published')
console.log(`Verified static landing: ${files.length} files, no exported app/API routes or server functions`)
