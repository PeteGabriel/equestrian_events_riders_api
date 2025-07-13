## Equestrian events API

This API provides information about equestrian competitions and events.

### Running with Docker

A Dockerfile is provided to containerize the application. Follow these steps to build and run the application using Docker:

1. Build the Docker image:
   ```bash
   docker build -t equestrian-events-api .
   ```

2. Run the container:
   ```bash
   docker run -p 8080:8080 equestrian-events-api
   ```

3. The API will be available at http://localhost:8080

### API Endpoints



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



# competitions/{event_id}

```
{
    "data": {
        "type": "competitions",
        "id": "778",
        "attributes": {
            "entry_list_url": "http://www.hippobase.com/EventInfo/Entries/CompetitorHorse.aspx?EventID=778",
            "homepage_url": "http://www.csicasasnovas.com",
            "name": "A Coruna - CH-EU-S/CSI3"
        },
        "relationships": {
            "events": {
                "data": [
                    {
                        "type": "events",
                        "id": "e237b4df-90b2-4adb-b196-5ad348a47ea9"
                    },
                    {
                        "type": "events",
                        "id": "51bd163d-b09c-43c1-90e0-0121ca0e4bd8"
                    }
                ]
            }
        }
    },
    "included": [
        {
            "type": "events",
            "id": "e237b4df-90b2-4adb-b196-5ad348a47ea9",
            "attributes": {
                "competitors": [ // 91 items
                  {
                      "rider": "KÜHNER, Max",
                      "horses": [
                          "249 Count On Me 19 (3)",
                          "56 EIC COOLEY JUMP THE Q (CH-EU)"
                      ]
                  }
                ],
                "date": "2025-07-13 16:20:58",
                "name": "Entries CH-EU-S",
                "total_of_athletes": 91,
                "total_of_horses": 121,
                "total_of_nations": 23
            }
        },
        {
            "type": "events",
            "id": "51bd163d-b09c-43c1-90e0-0121ca0e4bd8",
            "attributes": [
              ...
            ],
                "date": "2025-07-13 16:20:58",
                "name": "Entries CSI3*",
                "total_of_athletes": 18,
                "total_of_horses": 33,
                "total_of_nations": 4
            }
        }
    ]
}
```
