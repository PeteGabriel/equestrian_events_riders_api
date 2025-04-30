## Equestrian events API



### /competitions


```json
{
  "data": [
      {
        "type": "competitions",
        "id": "b332eddc-4cc5-4172-aeed-ea749052b6f7",
        "attributes": {
          "name": "61. Mannheim Maimarkt Turnier"
        },
        "relationships": {
          "events": {
            "data": [
              {
                "type": "events",
                "id": "63b38b3f-f24d-4525-bb4f-eb4170905fe7"
              },
              {
                "type": "events",
                "id": "efcc2936-27f6-456d-b49e-b1598d4cd3f7"
              }
            ]
          }
        }
      }
  ],
 "included": [
   {
     "type": "events",
     "id": "6039ad97-df79-4079-8c8d-43f24d22a60a",
     "attributes": {
       "competitors": [
         {
           "flag": "https://www.hippobase.com/flags/28x19/BEL.png",
           "country_code": "BEL",
           "country_name": "Belgium",
           "pairs": {
             "VERWIMP, Jorinde": [
               "Charmer"
             ]
           }
         },
         {
           "flag": "https://www.hippobase.com/flags/28x19/DEN.png",
           "country_code": "DEN",
           "country_name": "Denmark",
           "pairs": {
             "SKODBORG MERRALD, Nanna": [
               "Blue Hors Znickers"
             ]
           }
         }
       ]
     }
    }
  ] 
}

```