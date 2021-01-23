document.addEventListener("DOMContentLoaded", function (event) {
   if (window.location.href.endsWith("graph.html")) {
      LoadGraph();
   } else {
      LoadMap("mainmap");
   }
});

function LoadGraph() {
   d3.json("/api/json/bla.json").then(function (data) {
      var margin = { top: 20, right: 20, bottom: 30, left: 50 };
      var width = 600 - margin.left - margin.right;
      var height = 400 - margin.top - margin.bottom;

      var x = d3.scaleTime().range([0, width]);
      var y = d3.scaleLinear().range([height, 0]);

      data.forEach(function (d) {
         d.DateTime = Date.parse(d.Date);
      });
      data.sort((a, b) => {
         return a.DateTime - b.DateTime;
      });
      var line = d3
         .line()
         .x(function (d) {
            return x(d.DateTime);
         })
         .y(function (d) {
            return y(d.Voltage);
         });

      var svg = d3
         .select("#graph")
         .append("svg")
         .attr("width", width + margin.left + margin.right)
         .attr("height", height + margin.top + margin.bottom)
         .append("g")
         .attr("transform", "translate(" + margin.left + "," + margin.top + ")");

      x.domain(
         d3.extent(data, function (d) {
            return d.DateTime;
         })
      );
      y.domain([3.5, 4.5]);

      svg.append("path")
         .data([data])
         .attr("fill", "none")
         .attr("stroke", "steelblue")
         .attr("stroke-width", "2px")
         .attr("d", line);
      svg.append("g")
         .attr("transform", "translate(0," + height + ")")
         .call(d3.axisBottom(x));
      svg.append("g").call(d3.axisLeft(y));
   });
}

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
            Voltage: 80,
         },
      ],
      Water: [
         {
            Id: 12,
            Lat: -30.8062,
            Long: 26.4178,

            Level: 60,
         },
         {
            Id: 12,
            Lat: -30.8042,
            Long: 26.4118,

            Level: 10,
         },
         {
            Id: 12,
            Lat: -30.8042,
            Long: 26.4168,
            Level: 85,
         },
      ],
   };

   var allSheeps = L.layerGroup();

   bla.Sheep.forEach(function (item, index) {
      var sheepIcon = L.divIcon({
         className: "", // pass empty so it doesn't do any styling.
         html: "<div class='marker-main sheep'></div>",
      });
      var x = L.marker([item.Lat, item.Long], { icon: sheepIcon });
      allSheeps.addLayer(x);
   });

   var allWaters = L.layerGroup();
   bla.Water.forEach(function (item, index) {
      var levelCss = "level-low";
      if (item.Level > 40) {
         levelCss = "level-medium";
      }
      if (item.Level > 70) {
         levelCss = "level-high";
      }
      var waterIcon = L.divIcon({
         className: "", // pass empty so it doesn't do any styling.
         html: `<div class='marker-main water'><div class='marker-bar-outer'><div class='marker-bar-level ${levelCss}' style='height:${item.Level}%'></div></div></div>`,
      });
      var x = L.marker([item.Lat, item.Long], { icon: waterIcon });
      allWaters.addLayer(x);
   });

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
      layers: [Esri_WorldImagery, allSheeps, allWaters],
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
            Waters: allWaters,
         }
      )

      .addTo(mainmap)
      .expand();

   var mapCentre = L.latLng(bla.Center.Lat, bla.Center.Long);
   mainmap.setView(mapCentre, 16);

   return mainmap;
}
