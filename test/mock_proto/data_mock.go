package mock_v1

import pb "otabe/pb"

var NationsRate = &pb.NationsRate{
	NationName: "Vietnamese",
	Rate: 4,
}
var Feedbacks = []*pb.Feedbacks{
	&pb.Feedbacks{
		Nation: "Vietnamese",
		Rate: 4,
		Comments: []*pb.Comments{
			&pb.Comments{
				UserId: 1,
				UserComment: "Oishii ne",
				Rate: 5,
			},
			&pb.Comments{
				UserId: 2,
				UserComment: "Oaaaa",
				Rate: 3,
			},
		},
	},
	&pb.Feedbacks{
		Nation: "Japanese",
		Rate: 5,
		Comments: []*pb.Comments{
			&pb.Comments{
				UserId: 3,
				UserComment: "Ngon ghe",
				Rate: 5,
			},
		},
	},
}
var Menu_items = []*pb.MenuItems{
	&pb.MenuItems{
		ItemName: "Banh mi nhan thit",
		Description: "With pork inside",
		Price: 5300,
		Feedbacks: Feedbacks,
	},
	&pb.MenuItems{
		ItemName: "Banh mi rau",
		Description: "Vegetable",
		Price: 3400,
	},
}
var RestaurantDetails = &pb.GetRestaurantResponse{
	Restaurant: &pb.Restaurant{
		Id: 1,
		Name: "HaNoi & Hanoi",
		Website: "hanoi.com",
		Phone: "09642540626",
		Description: "oishii",
		PostalCode: "1080023",
		Address: "Tokyo",
		Geo: &pb.Geo{
			Lat: 35.644597778,
			Long: 139.748714210,
		},
	},
	NationsRate: []*pb.NationsRate{NationsRate, &pb.NationsRate{NationName: "Japanese", Rate: 5}},
	Menus: []*pb.Menus{
		&pb.Menus{
			Name: "Main",
			MenuItems: Menu_items,
		},
		&pb.Menus{
			Name: "Haru",
			MenuItems: []*pb.MenuItems{
				&pb.MenuItems{
					ItemName: "Banh mi cuon",
					Description: "Okonomi yaki",
					Price: 6000,
				},
			},
		},
	},
}

var ListRestaurants = &pb.ListRestaurantsResponse{
	Data: []*pb.GetRestaurantResponse{
		RestaurantDetails,
	},
}
