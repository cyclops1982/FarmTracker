document.addEventListener("DOMContentLoaded", function (event) {
   var mainmap = L.map("mainmap").setView([51.505, -0.09], 13);
   /*
  L.tileLayer(
    "https://api.mapbox.com/styles/v1/{id}/tiles/{z}/{x}/{y}?access_token={accessToken}",
    {
      attribution:
        '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors.',
      maxZoom: 18,
      id: "openstreetmap",
      tileSize: 256,
      zoomOffset: -1,
    }
  ).addTo(mymap);*/

   var layer1 = L.tileLayer("http://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
      attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
      subdomains: ["a", "b", "c"],
   }).addTo(mainmap);

   layers: [grayscale, cities];
});
