# Definition
API specification for Restaurants's functions
## Get Restaurant Details Function `restaurants/{restaurant_id}`
Get detail information about a restaurant

### Table referenced by this API
+ `restaurants`
+ `restaurants` reference table: `menus` 
+ `menus` reference table: `items`, `users` 

＋　others: `users`, `countries`, `feedbacks`, `items`

### Details
Return restaurant information in system with id : `restaurant_id`

### Params
| Params        | Description                          | Type  | Required  | Example  |
| ------------- |:------------------------------------:| -----:| ---------:| --------:|
| restaurant_id | Numeric ID of the restaurant to get. | int   | True      |  0       |

### Get conditions
+ `restaurant_id` (required)
    + **filter**: match restaurants.id

### Output conditions