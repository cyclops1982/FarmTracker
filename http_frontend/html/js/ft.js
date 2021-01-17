document.addEventListener("DOMContentLoaded", function (event) {
   var mainmap = LoadMap("mainmap");
});

function LoadMap(mapId) {
   // Load data from JSON
   var bla = {
      Center: {
         Lat: -30.8,
         Long: 26.4,
      },
      Sheep: [
         {
            Id: 1,
            Lat: -30.8,
            Long: 26.417,
            CreatedOn: 2010 - 02 - 23,
         },
         {
            Id: 2,
            Lat: -30.8012,
            Long: 26.4188,
            CreatedOn: 2010 - 02 - 26,
         },
      ],
   };

   var allSheeps = L.layerGroup();

   bla.Sheep.forEach(function (item, index) {
      var x = L.marker([item.Lat, item.Long]);
      allSheeps.addLayer(x);
   });

   var objectsOnMap = {
      Sheep: allSheeps,
   };

   var StreetMaps = L.tileLayer("http://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
      attribution: '&copy; <a href="https://openstreetmap.org/copyright">OpenStreetMap contributors</a>',
      maxZoom: 19,
   });

   var OpenTopoMap = L.tileLayer("https://{s}.tile.opentopomap.org/{z}/{x}/{y}.png", {
      maxZoom: 17,
      attribution:
         '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors, <a href="http://viewfinderpanoramas.org">SRTM</a> | Map style: &copy; <a href="https://opentopomap.org">OpenTopoMap</a> (<a href="https://creativecommons.org/licenses/by-sa/3.0/">CC-BY-SA</a>)',
   });

   var Esri_WorldImagery = L.tileLayer(
      "https://server.arcgisonline.com/ArcGIS/rest/services/World_Imagery/MapServer/tile/{z}/{y}/{x}",
      {
         attribution:
            "Tiles &copy; Esri &mdash; Source: Esri, i-cubed, USDA, USGS, AEX, GeoEye, Getmapping, Aerogrid, IGN, IGP, UPR-EGP, and the GIS User Community",
      }
   );

   var mainmap = L.map(mapId, {
      layers: [Esri_WorldImagery, allSheeps],
   });

   L.control.scale().addTo(mainmap);

   var layersVar = L.control
      .layers(
         {
            Street: StreetMaps,
            Topology: OpenTopoMap,
            Satelite: Esri_WorldImagery,
         },
         {
            Sheeps: allSheeps,
         }
      )

      .addTo(mainmap)
      .expand();

   var mapCentre = L.latLng(bla.Center.Lat, bla.Center.Long);
   mainmap.setView(mapCentre, 10);

   return mainmap;
}
