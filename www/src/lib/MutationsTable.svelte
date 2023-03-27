<script>
  import { mutations } from './store.js'

  export let getGeo // function from GeoMap component binded to self in App
  export let map_visible = true // from App

  let selected_mutation_id = null

  const sel_color = '#B1F76F'
  const color1 = '#ffe033'
  const color2 = '#ffc433'
  let mut_color = color1

  const intl_date_format = new Intl.DateTimeFormat('fr-FR')

  function is_same_mutation(index) {
    if (index == 0) return false
    return mutations_ids[index] == mutations_ids[index - 1]
  }

  function mutation_color(same, selected) {
    if (selected) return sel_color
    if (!same) {
      mut_color = mut_color == color2 ? color1 : color2
    }
    return mut_color
  }

  function intl_date(d) {
    return intl_date_format.format(Date.parse(d))
  }

  function mutation_type(mv) {
    let a = [mv['Type_local'].toLowerCase(), mv['Nature_culture'].toLowerCase(), mv['Nature_culture_speciale'].toLowerCase()]
    return a.filter(Boolean).join(' - ') // filter(Boolean) is a nice trick to remove null and empty strings
  }

  function mutation_click(mv) {
    selected_mutation_id = mv['Id_parcelle']
    getGeo(mv)
  }

  $: mutations_ids = Object.values($mutations).map((mv) => mv['Date_mutation'] + mv['Id_parcelle'])
</script>

{#if Object.values($mutations).length > 0}
  <table class="mutations">
    <thead>
      <tr>
        <th />
        <th>Date</th>
        <th style="width: 7rem">Commune</th>
        <th>Num</th>
        <th>Voie</th>
        <th>INSEE</th>
        <th>Postal</th>
        <th>Type</th>
        <th style="width: 4rem">Terrain m²</th>
        <th style="width: 4rem">Bâti m²</th>
        <th>Pièces</th>
        <th>Prix</th>
      </tr>
    </thead>
    <tbody>
      {#each $mutations as mv, i}
        <tr style:background-color={mutation_color(is_same_mutation(i), selected_mutation_id == mv['Id_parcelle'])}>
          <td>
            <a href="/" class:disabled={!map_visible} title={mv['Id_parcelle']} on:click|preventDefault={() => mutation_click(mv)}>voir</a>
          </td>
          <td>{intl_date(mv['Date_mutation'])}</td>
          <td>{mv['Nom_commune']}</td>
          <td>{mv['Adresse_numero']}</td>
          <td>{mv['Adresse_nom_voie']}</td>
          <td>{mv['Code_commune']}</td>
          <td>{mv['Code_postal']}</td>
          <td>{mutation_type(mv)}</td>
          <td>{mv['Surface_terrain']}</td>
          <td>{mv['Surface_relle_bati']}</td>
          <td>{mv['Nombre_pieces_principales']}</td>
          <td>{Math.trunc(parseFloat(mv['Valeur_fonciere']))}</td>
        </tr>
      {/each}
    </tbody>
  </table>
{/if}

<style>
  .mutations {
    width: 100%;
    font-size: 9pt;
    border-collapse: collapse;
  }
  .mutations td {
    border: 1px solid #eee;
    word-wrap: break-word;
    line-height: normal;
    padding: 0.2rem;
  }
  .mutations a {
    font-weight: bolder;
    color: black;
  }
  .mutations a.disabled {
    pointer-events: none;
    color: #777;
  }
</style>
