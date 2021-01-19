DELETE FROM Location;
DELETE FROM BatteryStatus;
DELETE FROM Device;
DELETE FROM Area;
DELETE FROM Org;


INSERT INTO Org(Name) VALUES('Farm1');
INSERT INTO Device(Name, OrgId, DeviceEUI) VALUES('My First Device',(SELECT Id FROM Org WHERE Name = 'Farm1'),'0004a30b00eb5e28');





INSERT INTO Area(OrgId, Name, GeoJSON) VALUES((SELECT Id FROM Org WHERE Name = 'Farm1'), 'Area 1', '
{
  "type": "FeatureCollection",
  "features": [
    {
      "type": "Feature",
      "properties": {},
      "geometry": {
        "type": "Polygon",
        "coordinates": [
          [
            [
              26.405768394470215,
              -30.798511042705595
            ],
            [
              26.42538070678711,
              -30.798511042705595
            ],
            [
              26.42538070678711,
              -30.794087365234297
            ],
            [
              26.405768394470215,
              -30.794087365234297
            ],
            [
              26.405768394470215,
              -30.798511042705595
            ]
          ]
        ]
      }
    },
    {
      "type": "Feature",
      "properties": {},
      "geometry": {
        "type": "Polygon",
        "coordinates": [
          [
            [
              26.406755447387695,
              -30.811596563200276
            ],
            [
              26.428213119506836,
              -30.811596563200276
            ],
            [
              26.428213119506836,
              -30.805035595376328
            ],
            [
              26.406755447387695,
              -30.805035595376328
            ],
            [
              26.406755447387695,
              -30.811596563200276
            ]
          ]
        ]
      }
    },
    {
      "type": "Feature",
      "properties": {},
      "geometry": {
        "type": "Polygon",
        "coordinates": [
          [
            [
              26.42207622528076,
              -30.79976438097815
            ],
            [
              26.414308547973633,
              -30.800980840731068
            ],
            [
              26.416325569152832,
              -30.804114075252052
            ],
            [
              26.42207622528076,
              -30.79976438097815
            ]
          ]
        ]
      }
    }
  ]
}');


INSERT INTO Area(OrgId, Name, GeoJSON) VALUES((SELECT Id FROM Org WHERE Name = 'Farm1'), 'Home 1', '{
  "type": "FeatureCollection",
  "features": [
    {
      "type": "Feature",
      "properties": {},
      "geometry": {
        "type": "Polygon",
        "coordinates": [
          [
            [
              -0.4086184501647949,
              51.386405665079955
            ],
            [
              -0.4081892967224121,
              51.386405665079955
            ],
            [
              -0.4081892967224121,
              51.38700825531624
            ],
            [
              -0.4086184501647949,
              51.38700825531624
            ],
            [
              -0.4086184501647949,
              51.386405665079955
            ]
          ]
        ]
      }
    }
  ]
}');
