# get filters by country + city
https://api.glovoapp.com/seo-content/locales/en/countries/kg/cities/BSK


# get stores by filter endpoint
https://api.glovoapp.com/v3/feeds/search?filter={filter_name}&categoryId=1
returns all stores that have a dish associated with the filter


# items endpoint
address_id has to be an address that is associated with the store
https://api.glovoapp.com/v3/stores/{glovo_store_id}/addresses/{glovo_address_id}/content?promoListViewWebVariation=CONTROL
