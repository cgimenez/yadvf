<script>
  import MutationsTable from './lib/MutationsTable.svelte'
  import SearchForm from './lib/SearchForm.svelte'
  import GeoMap from './lib/GeoMap.svelte'

  let getGeo // bind from GeoMap component
  let map_visible = true

  function toggle_map() {
    map_visible = !map_visible
    document.getElementById('map').style.display = map_visible ? 'block' : 'none'
  }
</script>

<GeoMap bind:getGeo />
<div class="app">
  <button on:click={toggle_map}>{map_visible ? 'Masquer la carte' : 'Afficher la carte'}</button>
  <SearchForm />
  <MutationsTable {map_visible} {getGeo} />
</div>

<style>
</style>
