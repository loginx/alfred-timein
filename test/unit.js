#!/usr/bin/osascript -l JavaScript
// Unit tests for what the CLI can't pin down offline: how the lookup chain reacts to each kind of answer,
// and how Nominatim replies map onto it. Loads timein.js into a scope whose network edges are stubbed;
// lookups are swapped through its own LOOKUPS table. Run from the repo root.

const read = path => $.NSString.stringWithContentsOfFileEncodingError(path, $.NSUTF8StringEncoding, $()).js
const reply = name => JSON.parse(read(`test/fixtures/nominatim/${name}.json`)) // captured Nominatim replies

let nominatim, missing // per case: what Nominatim "returns", which class this "macOS" lacks
const t = load({
  env: () => undefined, // default lookup order, regardless of the caller's environment
  fetchJSON: () => nominatim,
  cls: name => (name === missing ? $() : $.NSClassFromString(name)),
})
const boom = () => { throw new Error('must not be consulted') }
const answer = value => () => value

// [case, query, lookups, cache before, resolved zone, cache after]
const CHAIN = [
  ['first answer wins and is cached', 'eiffel tower', { osm: answer('Europe/Paris'), mapkit: boom }, {}, 'Europe/Paris', { 'eiffel tower': 'Europe/Paris' }],
  ["no answer → next lookup", 'eiffel tower', { osm: answer(null), mapkit: answer('Europe/Paris') }, {}, 'Europe/Paris', { 'eiffel tower': 'Europe/Paris' }],
  ['NOT_FOUND stops the chain', 'nowhere', { osm: answer(t.NOT_FOUND), mapkit: boom }, {}, undefined, {}],
  ['zone this macOS lacks → next lookup', 'kyiv', { osm: answer('Mars/Olympus_Mons'), mapkit: answer('Europe/Kyiv') }, {}, 'Europe/Kyiv', { kyiv: 'Europe/Kyiv' }],
  ['nobody answers → nothing cached', 'nowhere', { osm: answer(null), mapkit: answer(null) }, {}, undefined, {}],
  ['typed IANA zone skips lookups', 'europe/paris', { osm: boom, mapkit: boom }, {}, 'Europe/Paris', {}],
  ['cache hit skips lookups', 'jfk', { osm: boom, mapkit: boom }, { jfk: 'America/New_York' }, 'America/New_York', { jfk: 'America/New_York' }],
  ['new answer joins the cache', 'jfk', { osm: answer('America/New_York'), mapkit: boom }, { lima: 'America/Lima' }, 'America/New_York', { lima: 'America/Lima', jfk: 'America/New_York' }],
]

// [case, Nominatim reply (null = unreachable or non-200), class this macOS lacks, osm() result]
const OSM = [
  ['no results → NOT_FOUND', reply('no-results'), null, t.NOT_FOUND],
  ['unreachable or non-200 → no answer', null, null, null],
  ['hit, but CLGeocoder removed → no answer (MapKit takes over)', reply('eiffel-tower'), 'CLGeocoder', null],
]

function run() {
  const dir = `${$.NSTemporaryDirectory().js}timein-unit-${Date.now()}`
  const failures = []
  const check = (name, test) => {
    try {
      test((what, got, want) => show(got) === show(want) || failures.push(`${name}: ${what} ${show(got)}, want ${show(want)}`))
    } catch (e) {
      failures.push(`${name}: threw ${e.message}`)
    }
  }

  CHAIN.forEach(([name, query, lookups, before, zone, after], i) => {
    const cache = `${dir}/${i}/cache.json`
    if (Object.keys(before).length) t.writeJSON(cache, before)
    Object.assign(t.LOOKUPS, lookups)
    check(name, expect => {
      expect('resolved', t.resolve(query, `${dir}/no-seed.json`, cache), zone)
      expect('cached', t.readJSON(cache), after)
    })
  })
  for (const [name, r, lacks, want] of OSM) {
    nominatim = r
    missing = lacks
    check(name, expect => expect('got', t.osm('eiffel tower'), want))
  }

  $.NSFileManager.defaultManager.removeItemAtPathError(dir, $())
  const total = CHAIN.length + OSM.length
  failures.forEach(f => console.log(`FAIL: ${f}`))
  console.log(`unit: ${total - failures.length}/${total} passed`)
  if (failures.length) $.exit(1)
}

// A fresh timein.js scope; `stubs` replace its top-level functions by name.
function load(stubs) {
  const src = read('timein.js').replace(/^#!.*\n/, '')
  const swaps = Object.keys(stubs).map(name => `${name} = stubs.${name}`).join('\n')
  return new Function('stubs', `${src}\n${swaps}\nreturn { resolve, osm, LOOKUPS, NOT_FOUND, readJSON, writeJSON }`)(stubs)
}

function show(value) { return typeof value === 'symbol' ? value.toString() : JSON.stringify(value) ?? String(value) }
