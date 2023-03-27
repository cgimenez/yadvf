<script>
  import { onMount } from 'svelte'
  import maplibregl from 'maplibre-gl'

  let map
  let geoCache = {}

  onMount(() => {
    map = new maplibregl.Map({
      container: 'map',
      style: 'https://openmaptiles.geo.data.gouv.fr/styles/osm-bright/style.json', // Fond de carte
      zoom: 15, // Zoom
      center: [2.213749, 46.227638], // Centrage
      pitch: 0, // Inclinaison
      bearing: 0, // Rotation
      minZoom: 3, // Zoom min
      attributionControl: false
    })

    map.on('load', function () {})
  })

  export async function getGeo(mutation) {
    const section = mutation['Id_parcelle'].substring(8, 10)
    const numero = mutation['Id_parcelle'].substring(10, 14)
    const key = mutation['Code_commune'] + section + numero

    if (geoCache[key] === undefined) {
      const res = await fetch(`https://apicarto.ign.fr/api/cadastre/parcelle?code_insee=${mutation['Code_commune']}&section=${section}&numero=${numero}`)
      geoCache[key] = await res.json()
    }
    let geo_data = geoCache[key]

    map.setCenter([mutation['Longitude'], mutation['Latitude']])
    if (map.getSource('parcelles') === undefined) {
      map.addSource('parcelles', {
        type: 'geojson',
        data: { type: 'FeatureCollection', features: geo_data['features'] }
      })

      map.addLayer({
        id: 'parcelles',
        type: 'fill',
        source: 'parcelles',
        layout: {},
        paint: {
          'fill-color': '#088',
          'fill-opacity': 0.8
        }
      })

      map.addLayer({
        id: 'parcelles_outline',
        type: 'line',
        source: 'parcelles',
        layout: {},
        paint: {
          'line-color': '#FF0000',
          'line-opacity': 1,
          'line-width': 2
        }
      })
    } else {
      map.getSource('parcelles').setData({
        type: 'FeatureCollection',
        features: geo_data['features']
      })
    }
  }
</script>
