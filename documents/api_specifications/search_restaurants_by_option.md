# Search restaurants API
Search Restaurants By Optional Fields
Return list of restaurants matching with conditions.

## Table referenced by this API
+ `restaurants`
+ `restaurants` reference table: `menus` 
+ `menus` reference table: `items`, `users` 
+  others: `item_feedbacks`

## Details
Return list of restaurants in system matching with conditions :
<table name.column name> 
  
  + nation          <`users.nation`>
  + restaurant name <`restaurants.name`>
  + item name       <`items.name`>
  + prefecture      <`restaurants.prefecture`>
  + distance        (using current `lon` & `lat` of user and search by distance (radius). Returns list of restaurants sorted by closest)

### Input. request
#### Params
| Params          | Description                          | Type   | Required  | Validation   | Example      |
| --------------- |:------------------------------------:| ------:| ---------:| ------------:| ------------:|
| restaurant_name | Search by restaurant name.           | string | False     | | Hanoi & Hanoi   |
| nation | Search by user feedback's nation.             | string | False     | |Vietnamese   |
| item_name       | Search by item name.                 | string | False     | |Banh mi      |
| prefecture      | Search by prefecture name.           | string | False     |  |Tokyo       |
| location        | Search by restaurant location.       | object | False     | identified as a Japan location & be required if one of them not null |   |
| location.lat    | Latitude of search area.             | number | False     | |35.64479921  |
| location.lon    | Longtitude of search area.           | number | False     | |139.74933603 |
| location.distance| Search area radius (in meters).      | number | False     |  |200         |
| paging.page            | Page through results.                | int    | False     |  |2           |
| paging.page_limit      | Number of results on 1 page.         | int    | False     |  |10          |
| sorted_by       | Results sorted by                    | string | False     | in [`restaurant_name`, `nation`, `distance`, `rate`], default value: `rate` & `restaurant_name` | distance    | 
| is_menu       | Results include menu or not     | boolean | False     |  | true    | 

#### Search logic
+ Need to login: false
+ Parameters:
  + `restaurant_name` 
     + **filter**: match `restaurants.name`
  + `nation` 
     + **filter**: match `users.nation`
     + meaning: search by user's nation who feedbacked item of a restaurant 
  + `item_name` 
     + **filter**: match `items.name`
  + `prefecture`
     + **filter**: match `restaurants.prefecture` 
  + `is_menu`
     + **validation**: `false` for API `search restaurants nearby`
#### Output conditions
+ `RestaurantsSearchSortColumn`
    + Results should be sorted based on value of `sorted_by` param.
+ `Paging`
    + Number of results should be matched with condition of `page` & `page_limit` params.
    
### Output. response
+ Output data `<table name.column name>`
  + `paginationContents`
    + `totalResults`
    + `totalPages`
    + `page`
    + `pageLimit`
  + `restaurantsContents` (multiple) `restaurants` reference
    + [Details Response for 1 restaurant](https://github.com/hanhnt-mf/otabe/blob/master/documents/api_specifications/restaurant_details_response.md)
#### Error case
  + Restaurant cannot be found if not match with conditions . Return [] (empty array)
  + Error showed if params's type validated not right
  
  
