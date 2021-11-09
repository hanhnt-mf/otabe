### Output. response for restaurant details
+ Output data `<table name.column name>`
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
            + <users.id>
            + `rate`   (calculate by get average from user's rate for every item in restaurant)
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
        + `usersContents`
            + <user.id>
            + <user.name>
            + <user.nation>
