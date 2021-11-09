# Get Restaurant Details Function `restaurants/{restaurant_id}`
Get detail information about a restaurant

## Table referenced by this API
+ `restaurants`
+ `restaurants` reference table: `menus` 
+ `menus` reference table: `items`, `users` 

＋　others: `users`, `countries`, `feedbacks`, `items`

## Details
Return restaurant information in system with id : `restaurant_id`
### Input. request
#### Params
| Params        | Description                          | Type  | Required  | Example  |
| ------------- |:------------------------------------:| -----:| ---------:| --------:|
| restaurant_id | Numeric ID of the restaurant to get. | int   | True      |  0       |

#### Request conditions
+ Is logined: false
+ `restaurant_id` (required)
    + **filter**: match `restaurants`.id

#### Output conditions
+ `RestaurantSearchSortColumn`
    + None
+ `Paging`
    + None
+ Only return result with 1 restaurant - 1 `restaurant_id` exist on DB
    
### Output. response
+ Output data <table name.column name>
    + `restaurantContents`
        + <restaurants.id>
        + <restaurants.name>
        + <restaurants.website>
        + <restaurants.phone>
        + <restaurants.description>
        + `restaurantAddressContent`
            + <restaurants.postal_code>
            + <restaurants.prefecture>
            + <restaurants.city>
            + <restaurants.town> 
            + <restaurants.block>
        + `restaurantGeoContent`
            + <restaurants.lon>
            + <restaurants.lat>
        + `nationsRateContent`(multiple) : `menus` `items` `item_feedbacks` `users` reference
            + <users.country_name>
            + <rate> `calculate by get average from user's rate for every item in restaurant`
        + `menuContent` (multiple) : `menus` reference
            + <menus.name>
            + `menuItemsContent` (multiple) : `items` reference
                + <items.name>
                + <items.description>
                + <items.price> (`¥`)
                + `itemFeedbackContent` (multiple) : `item_feedbacks` `users` reference
                    + <nation> (get from `user.country_name`)
                    + <rate> (get average from user's rates by nation)
                    + `userCommentsContent` (multiple) 
                        + <item_feedbacks.user_id>
                        + <users.name>
                        + <item_feedbacks.comment>
                        + <item_feedbacks.rate>
    
    
