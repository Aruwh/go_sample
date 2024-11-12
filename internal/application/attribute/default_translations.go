package application

var defaultTranslations = `
{
  "attribute_bathroom": {
    "deDE": "Eigenes Badezimmer",
    "enGB": "Private bathroom",
    "frFR": "Salle de bains privative",
    "itIT": "Bagno privato"
  },
  "attribute_balcon": {
    "deDE": "Balkon",
    "enGB": "Balcony",
    "frFR": "Balcon",
    "itIT": "Balcone"
  },
  "attribute_kitchenette": {
    "deDE": "Kochnische",
    "enGB": "Kitchenette",
    "frFR": "Kitchenette",
    "itIT": "Angolo cottura"
  },
  "attribute_terrace": {
    "deDE": "Terrasse",
    "enGB": "Terrace",
    "frFR": "Terrasse",
    "itIT": "Terrazza"
  },
  "attribute_shower": {
    "deDE": "Dusche",
    "enGB": "Shower",
    "frFR": "Douche",
    "itIT": "Doccia"
  },
  "attribute_view": {
    "deDE": "Aussicht",
    "enGB": "View",
    "frFR": "Vue",
    "itIT": "Vista"
  },
  "attribute_refrigerator": {
    "deDE": "Kühlschrank",
    "enGB": "Refrigerator",
    "frFR": "Réfrigérateur",
    "itIT": "Frigorifero"
  },
  "attribute_bathtub": {
    "deDE": "Badewanne",
    "enGB": "Bathtub",
    "frFR": "Baignoire",
    "itIT": "Vasca da bagno"
  },
  "attribute_toilet": {
    "deDE": "WC",
    "enGB": "Toilet",
    "frFR": "Toilettes",
    "itIT": "WC"
  },
  "attribute_washingMachine": {
    "deDE": "Waschmaschine",
    "enGB": "Washing machine",
    "frFR": "Machine à laver",
    "itIT": "Lavatrice"
  },
  "attribute_tv": {
    "deDE": "TV",
    "enGB": "TV",
    "frFR": "Télévision",
    "itIT": "TV"
  },
  "attribute_kitchen": {
    "deDE": "Küche",
    "enGB": "Kitchen",
    "frFR": "Cuisine",
    "itIT": "Cucina"
  },
  "attribute_towels": {
    "deDE": "Handtücher",
    "enGB": "Towels",
    "frFR": "Serviettes",
    "itIT": "Asciugamani"
  },
  "attribute_soundproofing": {
    "deDE": "Schallisolierung",
    "enGB": "Soundproofing",
    "frFR": "Insonorisation",
    "itIT": "Isolamento acustico"
  },
  "attribute_toiletPaper": {
    "deDE": "Toilettenpapier",
    "enGB": "Toilet paper",
    "frFR": "Papier toilette",
    "itIT": "Carta igienica"
  },
  "attribute_MountainView": {
    "deDE": "Bergblick",
    "enGB": "Mountain view",
    "frFR": "Vue sur la montagne",
    "itIT": "Vista sulla montagna"
  },
  "attribute_accessibleByElevator": {
    "deDE": "Erreichbar mit dem Aufzug",
    "enGB": "Accessible by elevator",
    "frFR": "Accessible par ascenseur",
    "itIT": "Accessibile con l'ascensore"
  },
  "attribute_heating": {
    "deDE": "Heizung",
    "enGB": "Heating",
    "frFR": "Chauffage",
    "itIT": "Riscaldamento"
  },
  "attribute_sauna": {
    "deDE": "Sauna",
    "enGB": "Sauna",
    "frFR": "Sauna",
    "itIT": "Sauna"
  },
  "attribute_walkInShower": {
    "deDE": "Ebenerdige Dusche",
    "enGB": "Walk-in shower",
    "frFR": "Douche à l'italienne",
    "itIT": "Doccia a filo pavimento"
  },
  "attribute_coffeeTeaMaker": {
    "deDE": "Kaffee-/Teezubehör",
    "enGB": "Coffee/tea maker",
    "frFR": "Cafetière/bouilloire",
    "itIT": "Macchina per il caffè/te"
  },
  "attribute_bedding": {
    "deDE": "Bettwäsche",
    "enGB": "Bedding",
    "frFR": "Linge de lit",
    "itIT": "Biancheria da letto"
  },
  "attribute_outletNearTheBed": {
    "deDE": "Steckdose in Bettnähe",
    "enGB": "Outlet near the bed",
    "frFR": "Prise près du lit",
    "itIT": "Presa vicino al letto"
  },
  "attribute_bidet": {
    "deDE": "Bidet",
    "enGB": "Bidet",
    "frFR": "Bidet",
    "itIT": "Bidet"
  },
  "attribute_familyRoom": {
    "deDE": "Familienzimmer",
    "enGB": "Family room",
    "frFR": "Chambre familiale",
    "itIT": "Camera familiare"
  },
  "attribute_nonSmokingRoom": {
    "deDE": "Nichtraucherzimmer",
    "enGB": "Non-smoking room",
    "frFR": "Chambre non fumeurs",
    "itIT": "Camera non fumatori"
  },
  "attribute_parking": {
    "deDE": "Parkplatz",
    "enGB": "Parking",
    "frFR": "Parking",
    "itIT": "Parcheggio"
  },
  "attribute_wlan": {
    "deDE": "WLAN inklusive",
    "enGB": "Free Wi-Fi",
    "frFR": "Wi-Fi gratuit",
    "itIT": "Wi-Fi gratuito"
  },
  "attribute_restairant": {
    "deDE": "Restaurant",
    "enGB": "Restaurant",
    "frFR": "Restaurant",
    "itIT": "Ristorante"
  },
  "attribute_petsAllowed": {
    "deDE": "Haustiere erlaubt",
    "enGB": "Pets allowed",
    "frFR": "Animaux autorisés",
    "itIT": "Animali ammessi"
  },
  "attribute_roomService": {
    "deDE": "Zimmerservice",
    "enGB": "Room service",
    "frFR": "Service en chambre",
    "itIT": "Servizio in camera"
  },
  "attribute_pool": {
    "deDE": "Pool",
    "enGB": "Pool",
    "frFR": "Piscine",
    "itIT": "Piscina"
  },
  "attribute_fitnessCenter": {
    "deDE": "Fitnesscenter",
    "enGB": "Fitness center",
    "frFR": "Centre de fitness",
    "itIT": "Centro fitness"
  },
  "attribute_electricCarChargingStation": {
    "deDE": "Aufladestation für Elektro-Autos",
    "enGB": "Electric car charging station",
    "frFR": "Station de recharge pour voitures électriques",
    "itIT": "Stazione di ricarica per auto elettriche"
  },
  "attribute_facilitiesForDisavledGuests": {
    "deDE": "Einrichtungen für behinderte Gäste",
    "enGB": "Facilities for disabled guests",
    "frFR": "Installations pour les personnes handicapées",
    "itIT": "Strutture per ospiti disabili"
  }
}`
