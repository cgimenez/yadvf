<script>
  import { onMount } from 'svelte'
  import { mutations, distincts, fetching } from './store.js'
  import { setup, fetch_mutations, csv_export } from './Mutations.svelte'

  let cityCodes = []
  let form_data = {
    nature_culture: '',
    nature_culture_speciale: '',
    type_local: '',
    Surface_terrain_min: '',
    Surface_terrain_max: '',
    Surface_reelle_bati_min: '',
    Surface_reelle_bati_max: '',
    Nom_commune: '',
    Code_departement: '',
    date_debut: '',
    date_fin: ''
  }

  /*
  function submitable() {
    let r = false

    r = r || !fetch_pending
    r = r || cityCodes.length > 0
    let one_value = false
    for (const property in form_data) {
      if (form_data[property].length > 0) {
        one_value = true
        break
      }
    }
    r = r || one_value
    return r
  }
*/

  function addCityCode() {
    cityCodes = [...cityCodes, { kind: 'POS', code: '' }]
  }

  function removeCityCode(index) {
    cityCodes = [...cityCodes.slice(0, index), ...cityCodes.slice(index + 1, cityCodes.length)]
  }

  async function submitForm(e) {
    let action = e.target.action
    action = action.substring(action.lastIndexOf('/'), action.length) // remove protocol and host

    form_data['code_commune'] = cityCodes.filter((e) => e.kind == 'COM').map((e) => e['code'])
    form_data['code_postal'] = cityCodes.filter((e) => e.kind == 'POS').map((e) => e['code'])
    fetch_mutations(form_data, action)
  }

  onMount(() => {
    setup()
  })
</script>

<form action="/query" on:submit|preventDefault={submitForm}>
  <fieldset>
    <legend>Type de bien</legend>
    <select name="nature_culture" bind:value={form_data['nature_culture']}>
      {#each $distincts['Cultures'] as option}
        <option value={option}>
          {option}
        </option>
      {/each}
    </select>
    <select name="nature_culture_speciale" bind:value={form_data['nature_culture_speciale']}>
      {#each $distincts['Cultures_speciales'] as option}
        <option value={option}>
          {option}
        </option>
      {/each}
    </select>
    <select name="type_local" bind:value={form_data['type_local']}>
      {#each $distincts['Locaux'] as option}
        <option value={option}>
          {option}
        </option>
      {/each}
    </select>
  </fieldset>

  <fieldset>
    <legend>Superficies</legend>
    <label>Terrain min<input type="text" name="Surface_terrain_min" size="8" bind:value={form_data['Surface_terrain_min']} /></label>
    <label>max<input type="text" name="Surface_terrain_max" size="8" bind:value={form_data['Surface_terrain_max']} /></label><br />
    <label>Bâti min<input type="text" name="Surface_reelle_bati_min" size="8" bind:value={form_data['Surface_reelle_bati_min']} /></label>
    <label>max<input type="text" name="Surface_reelle_bati_max" size="8" bind:value={form_data['Surface_reelle_bati_max']} /></label>
  </fieldset>

  <fieldset>
    <legend>Codes communes et postaux</legend>
    {#each cityCodes as value, i}
      <select bind:value={value['kind']}>
        <option value="POS">Postal</option>
        <option value="COM">Commune (INSEE)</option>
      </select>
      <button on:click|preventDefault={() => removeCityCode(i)} class="small">-</button>
      <input bind:value={value['code']} />
      <br />
    {/each}
    <div>
      <button on:click|preventDefault={addCityCode}>Ajouter</button>
    </div>
  </fieldset>

  <fieldset>
    <legend>Autre</legend>
    <label>Nom commune<input type="text" name="Nom_commune" bind:value={form_data['Nom_commune']} /></label>
    <label>Département<input type="text" name="Code_departement" bind:value={form_data['Code_departement']} /></label>
  </fieldset>

  <fieldset>
    <legend>Période</legend>
    <label>Du<input type="text" name="date_debut" placeholder="YYYY-MM-DD" bind:value={form_data['date_debut']} /></label>
    <label>Au<input type="text" name="date_fin" placeholder="YYYY-MM-DD" bind:value={form_data['date_fin']} /></label>
  </fieldset>
  <div class="buttons">
    <button disabled={$fetching}>Rechercher</button>
    <button disabled={$mutations.length == 0} on:click={csv_export}>Export CSV</button>
  </div>
</form>

<style>
  .buttons {
    padding: 1rem;
  }

  button.small {
    padding: 0 0.3rem;
    background-color: red;
    border-radius: 0;
    border: none;
  }
</style>
