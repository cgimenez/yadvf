<script context="module">
  import { get } from 'svelte/store'
  import { mutations, fetching, distincts } from './store.js'

  // --------------------------------------------------------------------------
  // Exported functions
  // --------------------------------------------------------------------------

  export async function fetch_mutations(form_data, action) {
    fetching.set(true)
    mutations.set([])

    try {
      var response = await fetch(build_url(action), {
        method: 'POST',
        body: JSON.stringify(form_data)
      })

      const response_data = await response.json()
      if (response_data != null) {
        mutations.set(response_data)
      }
    } finally {
      fetching.set(false)
    }
  }

  export function csv_export() {
    const table_data = []
    const row_data = []

    for (const key of Object.keys(get(mutations)[0])) {
      row_data.push(key)
    }
    table_data.push(row_data.join(','))

    for (const mutation of get(mutations)) {
      const row_data = []
      for (const property in mutation) {
        row_data.push(mutation[property])
      }
      table_data.push(row_data.join(','))
    }

    const csv_data = table_data.join('\n')

    const a = document.createElement('a')
    a.href = URL.createObjectURL(new Blob([csv_data], { type: 'text/csv' }))
    a.setAttribute('download', 'data.csv')
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
  }

  export function setup() {
    const port = import.meta.env.VITE_BACKEND_HTTP_PORT
    backend_url = port === undefined ? `${window.location.origin}` : `${window.location.protocol}//${window.location.hostname}:${port}`
    setMutationsTypes()
  }

  // --------------------------------------------------------------------------
  // Private
  // --------------------------------------------------------------------------

  let backend_url

  async function setMutationsTypes() {
    const res = await fetch(build_url('/types'))
    const d = await res.json()
    d['Cultures'].unshift('')
    d['Cultures'].sort()
    d['Cultures_speciales'].unshift('')
    d['Cultures_speciales'].sort()
    d['Locaux'].unshift('')
    d['Locaux'].sort()
    distincts.set(d)
  }

  function build_url(path) {
    return `${backend_url}${path}`
  }
</script>
