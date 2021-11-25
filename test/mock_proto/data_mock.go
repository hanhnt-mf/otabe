package mock_v1

import v1 "otabe/v1"

var NationsRate = &v1.NationsRate{
	NationName: "Vietnamese",
	Rate: 4,
}
var Feedbacks = []*v1.Feedbacks{
	&v1.Feedbacks{
		Nation: "Vietnamese",
		Rate: 4,
		Comments: []*v1.Comments{
			&v1.Comments{
				UserId: 1,
				UserComment: "Oishii ne",
				Rate: 5,
			},
			&v1.Comments{
				UserId: 2,
				UserComment: "Oaaaa",
				Rate: 3,
			},
		},
	},
	&v1.Feedbacks{
		Nation: "Japanese",
		Rate: 5,
		Comments: []*v1.Comments{
			&v1.Comments{
				UserId: 3,
				UserComment: "Ngon ghe",
				Rate: 5,
			},
		},
	},
}
var Menu_items = []*v1.MenuItems{
	&v1.MenuItems{
		ItemName: "Banh mi nhan thit",
		Description: "With pork inside",
		Price: 5300,
		Feedbacks: Feedbacks,
	},
	&v1.MenuItems{
		ItemName: "Banh mi rau",
		Description: "Vegetable",
		Price: 3400,
	},
}
var RestaurantDetails = &v1.GetRestaurantResponse{
	Restaurant: &v1.Restaurant{
		Id: 1,
		Name: "HaNoi & Hanoi",
		Website: "hanoi.com",
		Phone: "09642540626",
		Description: "oishii",
		PostalCode: "1080023",
		Address: "Tokyo",
		Geo: &v1.Geo{
			Lat: 35.644597778,
			Long: 139.748714210,
		},
	},
	NationsRate: []*v1.NationsRate{NationsRate, &v1.NationsRate{NationName: "Japanese", Rate: 5}},
	Menus: []*v1.Menus{
		&v1.Menus{
			Name: "Main",
			MenuItems: Menu_items,
		},
		&v1.Menus{
			Name: "Haru",
			MenuItems: []*v1.MenuItems{
				&v1.MenuItems{
					ItemName: "Banh mi cuon",
					Description: "Okonomi yaki",
					Price: 6000,
				},
			},
		},
	},
}

var ListRestaurants = &v1.ListRestaurantsResponse{
	Data: []*v1.GetRestaurantResponse{
		RestaurantDetails,
	},
}
