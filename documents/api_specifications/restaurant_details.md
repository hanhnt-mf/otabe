# Get Restaurant Details
Get detail information about a restaurant

## Table referenced by this API
+ `restaurants`
+ `restaurants` reference table: `menus` 
+ `menus` reference table: `items`, `users` 
＋　others: `item_feedbacks`

## Details
Return restaurant information in system with id : `restaurant_id`
### Input. request
#### Params
| Params        | Description                          | Type  | Required  | Validation  | Example  |
| ------------- |:------------------------------------:| -----:| ---------:| -----------:| --------:|
| restaurant_id | Numeric ID of the restaurant to get. | int   | True      |             |   0      |

#### Search Logic
+ Need to login: false
+ `restaurant_id` (required)
    + **filter**: match `restaurants`.id

#### Output conditions
+ `RestaurantSearchSortColumn`
    + None
+ `Paging`
    + None
+ Only return result with 1 restaurant - 1 `restaurant_id` exist on DB
    
### Output. response
+ [Restaurant details response](https://github.com/hanhnt-mf/otabe/blob/master/documents/api_specifications/restaurant_details_response.md)
    
### Error Cases (INVALID_ARGUMENT)
+ Restaurant cannot be found if `restaurant_id` param doesn't exist in DB . Return {} (empty object)
+ Error showed if `restaurant_id`'s not type `int`
